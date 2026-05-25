package repository

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"oceanengine-backend/internal/app/group/model"
)

type GroupRepository struct {
	db *gorm.DB
}

func NewGroupRepository(db *gorm.DB) *GroupRepository {
	return &GroupRepository{db: db}
}

func (r *GroupRepository) Create(ctx context.Context, group *model.AccountGroup) error {
	return r.db.WithContext(ctx).Create(group).Error
}

func (r *GroupRepository) GetByID(ctx context.Context, tenantID, id uint64) (*model.AccountGroup, error) {
	var group model.AccountGroup
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&group).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *GroupRepository) List(ctx context.Context, tenantID uint64, groupType string) ([]model.AccountGroup, error) {
	var groups []model.AccountGroup
	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if groupType != "" {
		query = query.Where("type = ?", groupType)
	}
	err := query.Order("sort ASC").Find(&groups).Error
	return groups, err
}

func (r *GroupRepository) Update(ctx context.Context, group *model.AccountGroup) error {
	return r.db.WithContext(ctx).Save(group).Error
}

func (r *GroupRepository) Delete(ctx context.Context, tenantID, id uint64) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&model.AccountGroup{}).Error
}

func (r *GroupRepository) AddMembers(ctx context.Context, groupID uint64, accountIDs []uint64) error {
	members := make([]model.AccountGroupMember, 0, len(accountIDs))
	for _, aid := range accountIDs {
		members = append(members, model.AccountGroupMember{
			GroupID:   groupID,
			AccountID: aid,
		})
	}
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&members).Error
}

func (r *GroupRepository) RemoveMembers(ctx context.Context, groupID uint64, accountIDs []uint64) error {
	return r.db.WithContext(ctx).
		Where("group_id = ? AND account_id IN ?", groupID, accountIDs).
		Delete(&model.AccountGroupMember{}).Error
}

func (r *GroupRepository) ListMembers(ctx context.Context, groupID uint64) ([]uint64, error) {
	var ids []uint64
	err := r.db.WithContext(ctx).
		Model(&model.AccountGroupMember{}).
		Where("group_id = ?", groupID).
		Pluck("account_id", &ids).Error
	return ids, err
}
