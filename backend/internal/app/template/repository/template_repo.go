package repository

import (
	"context"

	"gorm.io/gorm"

	"oceanengine-backend/internal/app/template/model"
)

// TemplateRepository 模板数据仓库
type TemplateRepository struct {
	db *gorm.DB
}

// NewTemplateRepository 创建模板仓库实例
func NewTemplateRepository(db *gorm.DB) *TemplateRepository {
	return &TemplateRepository{db: db}
}

// CreateProject 创建项目模板
func (r *TemplateRepository) CreateProject(ctx context.Context, t *model.ProjectTemplate) error {
	return r.db.WithContext(ctx).Create(t).Error
}

// ListProjects 获取项目模板列表，按使用次数降序排列
func (r *TemplateRepository) ListProjects(ctx context.Context, tenantID uint64) ([]model.ProjectTemplate, error) {
	var list []model.ProjectTemplate
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Order("use_count DESC").
		Find(&list).Error
	return list, err
}

// GetProject 获取单个项目模板
func (r *TemplateRepository) GetProject(ctx context.Context, tenantID, id uint64) (*model.ProjectTemplate, error) {
	var t model.ProjectTemplate
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// UpdateProject 更新项目模板
func (r *TemplateRepository) UpdateProject(ctx context.Context, t *model.ProjectTemplate) error {
	return r.db.WithContext(ctx).Save(t).Error
}

// DeleteProject 软删除项目模板
func (r *TemplateRepository) DeleteProject(ctx context.Context, tenantID, id uint64) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&model.ProjectTemplate{}).Error
}

// IncrUseCount 增加使用次数
func (r *TemplateRepository) IncrUseCount(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).
		Model(&model.ProjectTemplate{}).
		Where("id = ?", id).
		UpdateColumn("use_count", gorm.Expr("use_count + 1")).Error
}

// CreatePromotion 创建推广模板
func (r *TemplateRepository) CreatePromotion(ctx context.Context, t *model.PromotionTemplate) error {
	return r.db.WithContext(ctx).Create(t).Error
}

// ListPromotions 获取推广模板列表
func (r *TemplateRepository) ListPromotions(ctx context.Context, tenantID uint64) ([]model.PromotionTemplate, error) {
	var list []model.PromotionTemplate
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Order("use_count DESC").
		Find(&list).Error
	return list, err
}

// DeletePromotion 软删除推广模板
func (r *TemplateRepository) DeletePromotion(ctx context.Context, tenantID, id uint64) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&model.PromotionTemplate{}).Error
}
