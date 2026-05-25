package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
	"oceanengine-backend/internal/app/analytics/model"
)

type ExportRepository struct {
	db *gorm.DB
}

func NewExportRepository(db *gorm.DB) *ExportRepository {
	return &ExportRepository{db: db}
}

func (r *ExportRepository) Create(ctx context.Context, task *model.ExportTask) error {
	return r.db.WithContext(ctx).Create(task).Error
}

func (r *ExportRepository) GetByID(ctx context.Context, tenantID, id uint64) (*model.ExportTask, error) {
	var task model.ExportTask
	err := r.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).First(&task).Error
	return &task, err
}

func (r *ExportRepository) List(ctx context.Context, tenantID uint64) ([]model.ExportTask, error) {
	var tasks []model.ExportTask
	err := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).Order("id DESC").Limit(20).Find(&tasks).Error
	return tasks, err
}

func (r *ExportRepository) UpdateStatus(ctx context.Context, id uint64, status, filePath string, rowCount int) error {
	return r.db.WithContext(ctx).Model(&model.ExportTask{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":    status,
		"file_path": filePath,
		"row_count": rowCount,
	}).Error
}

func (r *ExportRepository) CleanExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("expire_at < ?", time.Now()).Delete(&model.ExportTask{}).Error
}
