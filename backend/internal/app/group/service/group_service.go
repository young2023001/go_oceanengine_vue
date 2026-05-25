package service

import (
	"context"

	"oceanengine-backend/internal/app/group/model"
	"oceanengine-backend/internal/app/group/repository"
)

type CreateGroupRequest struct {
	Name     string `json:"name" binding:"required"`
	Type     string `json:"type" binding:"required,oneof=franchisee region custom"`
	ParentID uint64 `json:"parent_id"`
	Sort     int    `json:"sort"`
}

type UpdateGroupRequest struct {
	Name string `json:"name"`
	Sort *int   `json:"sort"`
}

type MembersRequest struct {
	AccountIDs []uint64 `json:"account_ids" binding:"required,min=1"`
}

type GroupService struct {
	repo *repository.GroupRepository
}

func NewGroupService(repo *repository.GroupRepository) *GroupService {
	return &GroupService{repo: repo}
}

func (s *GroupService) Create(ctx context.Context, tenantID uint64, req *CreateGroupRequest) (*model.AccountGroup, error) {
	group := &model.AccountGroup{
		TenantID: tenantID,
		Name:     req.Name,
		Type:     model.GroupType(req.Type),
		ParentID: req.ParentID,
		Sort:     req.Sort,
	}
	if err := s.repo.Create(ctx, group); err != nil {
		return nil, err
	}
	return group, nil
}

func (s *GroupService) List(ctx context.Context, tenantID uint64, groupType string) ([]model.AccountGroup, error) {
	return s.repo.List(ctx, tenantID, groupType)
}

func (s *GroupService) Update(ctx context.Context, tenantID, id uint64, req *UpdateGroupRequest) error {
	group, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if req.Name != "" {
		group.Name = req.Name
	}
	if req.Sort != nil {
		group.Sort = *req.Sort
	}
	return s.repo.Update(ctx, group)
}

func (s *GroupService) Delete(ctx context.Context, tenantID, id uint64) error {
	return s.repo.Delete(ctx, tenantID, id)
}

func (s *GroupService) AddMembers(ctx context.Context, tenantID, groupID uint64, req *MembersRequest) error {
	// Verify group belongs to tenant
	if _, err := s.repo.GetByID(ctx, tenantID, groupID); err != nil {
		return err
	}
	return s.repo.AddMembers(ctx, groupID, req.AccountIDs)
}

func (s *GroupService) RemoveMembers(ctx context.Context, tenantID, groupID uint64, req *MembersRequest) error {
	// Verify group belongs to tenant
	if _, err := s.repo.GetByID(ctx, tenantID, groupID); err != nil {
		return err
	}
	return s.repo.RemoveMembers(ctx, groupID, req.AccountIDs)
}
