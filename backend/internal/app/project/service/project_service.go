package service

import (
	"context"

	"oceanengine-backend/internal/app/project/model"
	"oceanengine-backend/internal/app/project/repository"
)

type ListProjectsRequest struct {
	Status   string `form:"status"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}

type ListPromotionsRequest struct {
	ProjectID uint64 `form:"project_id"`
	Page      int    `form:"page"`
	PageSize  int    `form:"page_size"`
}

type ProjectService struct {
	repo *repository.ProjectRepository
}

func NewProjectService(repo *repository.ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) ListProjects(ctx context.Context, tenantID uint64, accountIDs []uint64, req ListProjectsRequest) ([]model.LocalProject, int64, error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	return s.repo.ListProjects(ctx, tenantID, accountIDs, req.Status, page, pageSize)
}

func (s *ProjectService) GetProject(ctx context.Context, tenantID, id uint64) (*model.LocalProject, error) {
	return s.repo.GetByID(ctx, tenantID, id)
}

func (s *ProjectService) UpdateStatus(ctx context.Context, tenantID, id uint64, status string) error {
	return s.repo.UpdateStatus(ctx, tenantID, id, status)
}

func (s *ProjectService) ListPromotions(ctx context.Context, tenantID uint64, accountIDs []uint64, req ListPromotionsRequest) ([]model.LocalPromotion, int64, error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	return s.repo.ListPromotions(ctx, tenantID, req.ProjectID, accountIDs, page, pageSize)
}
