package service

import (
	"context"
	"encoding/json"

	"oceanengine-backend/internal/app/template/model"
	"oceanengine-backend/internal/app/template/repository"
)

// CreateTemplateRequest 创建模板请求
type CreateTemplateRequest struct {
	Name   string                 `json:"name" binding:"required"`
	Config map[string]interface{} `json:"config"`
}

// TemplateService 模板服务
type TemplateService struct {
	repo *repository.TemplateRepository
}

// NewTemplateService 创建模板服务实例
func NewTemplateService(repo *repository.TemplateRepository) *TemplateService {
	return &TemplateService{repo: repo}
}

// CreateProject 创建项目模板
func (s *TemplateService) CreateProject(ctx context.Context, tenantID, userID uint64, req *CreateTemplateRequest) error {
	configJSON, err := json.Marshal(req.Config)
	if err != nil {
		return err
	}

	t := &model.ProjectTemplate{
		TenantID:  tenantID,
		Name:      req.Name,
		Config:    string(configJSON),
		CreatedBy: userID,
	}
	return s.repo.CreateProject(ctx, t)
}

// ListProjects 获取项目模板列表
func (s *TemplateService) ListProjects(ctx context.Context, tenantID uint64) ([]model.ProjectTemplate, error) {
	return s.repo.ListProjects(ctx, tenantID)
}

// GetProject 获取单个项目模板
func (s *TemplateService) GetProject(ctx context.Context, tenantID, id uint64) (*model.ProjectTemplate, error) {
	return s.repo.GetProject(ctx, tenantID, id)
}

// UpdateProject 更新项目模板
func (s *TemplateService) UpdateProject(ctx context.Context, t *model.ProjectTemplate) error {
	return s.repo.UpdateProject(ctx, t)
}

// DeleteProject 删除项目模板
func (s *TemplateService) DeleteProject(ctx context.Context, tenantID, id uint64) error {
	return s.repo.DeleteProject(ctx, tenantID, id)
}

// CreatePromotion 创建推广模板
func (s *TemplateService) CreatePromotion(ctx context.Context, tenantID, userID uint64, req *CreateTemplateRequest) error {
	configJSON, err := json.Marshal(req.Config)
	if err != nil {
		return err
	}

	t := &model.PromotionTemplate{
		TenantID:  tenantID,
		Name:      req.Name,
		Config:    string(configJSON),
		CreatedBy: userID,
	}
	return s.repo.CreatePromotion(ctx, t)
}

// ListPromotions 获取推广模板列表
func (s *TemplateService) ListPromotions(ctx context.Context, tenantID uint64) ([]model.PromotionTemplate, error) {
	return s.repo.ListPromotions(ctx, tenantID)
}

// DeletePromotion 删除推广模板
func (s *TemplateService) DeletePromotion(ctx context.Context, tenantID, id uint64) error {
	return s.repo.DeletePromotion(ctx, tenantID, id)
}
