package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
	"oceanengine-backend/internal/app/project/model"
)

type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) ListProjects(ctx context.Context, tenantID uint64, accountIDs []uint64, status string, page, pageSize int) ([]model.LocalProject, int64, error) {
	var projects []model.LocalProject
	var total int64

	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if len(accountIDs) > 0 {
		query = query.Where("account_id IN ?", accountIDs)
	}
	if status != "" {
		query = query.Where("status_first = ?", status)
	}

	if err := query.Model(&model.LocalProject{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&projects).Error; err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}

func (r *ProjectRepository) GetByID(ctx context.Context, tenantID, id uint64) (*model.LocalProject, error) {
	var project model.LocalProject
	if err := r.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).First(&project).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepository) Upsert(ctx context.Context, project *model.LocalProject) error {
	return r.db.WithContext(ctx).Save(project).Error
}

func (r *ProjectRepository) UpdateStatus(ctx context.Context, tenantID, id uint64, status string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&model.LocalProject{}).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Updates(map[string]interface{}{
			"status_first":     status,
			"local_updated_at": &now,
		}).Error
}

func (r *ProjectRepository) ListPromotions(ctx context.Context, tenantID, projectID uint64, accountIDs []uint64, page, pageSize int) ([]model.LocalPromotion, int64, error) {
	var promotions []model.LocalPromotion
	var total int64

	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if projectID > 0 {
		query = query.Where("project_id = ?", projectID)
	}
	if len(accountIDs) > 0 {
		query = query.Where("account_id IN ?", accountIDs)
	}

	if err := query.Model(&model.LocalPromotion{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&promotions).Error; err != nil {
		return nil, 0, err
	}

	return promotions, total, nil
}

func (r *ProjectRepository) UpsertPromotion(ctx context.Context, promo *model.LocalPromotion) error {
	return r.db.WithContext(ctx).Save(promo).Error
}
