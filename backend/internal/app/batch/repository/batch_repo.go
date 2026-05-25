package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"oceanengine-backend/internal/app/batch/model"
)

// BatchRepository 批量任务数据仓库
type BatchRepository struct {
	db *gorm.DB
}

// NewBatchRepository 创建 BatchRepository
func NewBatchRepository(db *gorm.DB) *BatchRepository {
	return &BatchRepository{db: db}
}

// CreateTask 在事务中创建任务和子项
func (r *BatchRepository) CreateTask(ctx context.Context, task *model.BatchTask, items []model.BatchTaskItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(task).Error; err != nil {
			return err
		}
		for i := range items {
			items[i].TaskID = task.ID
			items[i].TenantID = task.TenantID
		}
		if len(items) > 0 {
			if err := tx.CreateInBatches(items, 100).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetTask 根据租户ID和任务ID获取任务
func (r *BatchRepository) GetTask(ctx context.Context, tenantID, taskID uint64) (*model.BatchTask, error) {
	var task model.BatchTask
	err := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", taskID, tenantID).
		First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// GetTaskByID 根据任务ID获取任务（不限租户）
func (r *BatchRepository) GetTaskByID(ctx context.Context, taskID uint64) (*model.BatchTask, error) {
	var task model.BatchTask
	err := r.db.WithContext(ctx).
		Where("id = ?", taskID).
		First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// ListTasks 分页获取租户的任务列表
func (r *BatchRepository) ListTasks(ctx context.Context, tenantID uint64, page, pageSize int) ([]model.BatchTask, int64, error) {
	var tasks []model.BatchTask
	var total int64

	query := r.db.WithContext(ctx).Model(&model.BatchTask{}).Where("tenant_id = ?", tenantID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

// ListRunningTasks 获取所有运行中的任务
func (r *BatchRepository) ListRunningTasks(ctx context.Context) ([]model.BatchTask, error) {
	var tasks []model.BatchTask
	err := r.db.WithContext(ctx).
		Where("status = ?", model.TaskStatusRunning).
		Find(&tasks).Error
	return tasks, err
}

// GetPendingItems 获取指定任务的待处理子项
func (r *BatchRepository) GetPendingItems(ctx context.Context, taskID uint64, limit int) ([]model.BatchTaskItem, error) {
	var items []model.BatchTaskItem
	err := r.db.WithContext(ctx).
		Where("task_id = ? AND status = ?", taskID, model.ItemStatusPending).
		Order("id ASC").
		Limit(limit).
		Find(&items).Error
	return items, err
}

// UpdateItemStatus 更新子项状态
func (r *BatchRepository) UpdateItemStatus(ctx context.Context, itemID uint64, status model.ItemStatus, resultID uint64, errMsg string) error {
	now := time.Now()
	updates := map[string]interface{}{
		"status":      status,
		"result_id":   resultID,
		"error_msg":   errMsg,
		"executed_at": &now,
	}
	return r.db.WithContext(ctx).
		Model(&model.BatchTaskItem{}).
		Where("id = ?", itemID).
		Updates(updates).Error
}

// UpdateTaskStatus 更新任务状态，running 时设 started_at，完成时设 completed_at
func (r *BatchRepository) UpdateTaskStatus(ctx context.Context, taskID uint64, status model.TaskStatus) error {
	updates := map[string]interface{}{
		"status": status,
	}

	now := time.Now()
	if status == model.TaskStatusRunning {
		updates["started_at"] = &now
	}
	if status == model.TaskStatusCompleted || status == model.TaskStatusPartialSuccess || status == model.TaskStatusFailed || status == model.TaskStatusCancelled {
		updates["completed_at"] = &now
	}

	return r.db.WithContext(ctx).
		Model(&model.BatchTask{}).
		Where("id = ?", taskID).
		Updates(updates).Error
}

// IncrSuccess 增加成功计数
func (r *BatchRepository) IncrSuccess(ctx context.Context, taskID uint64) error {
	return r.db.WithContext(ctx).
		Model(&model.BatchTask{}).
		Where("id = ?", taskID).
		UpdateColumn("success_count", gorm.Expr("success_count + 1")).Error
}

// IncrFailed 增加失败计数
func (r *BatchRepository) IncrFailed(ctx context.Context, taskID uint64) error {
	return r.db.WithContext(ctx).
		Model(&model.BatchTask{}).
		Where("id = ?", taskID).
		UpdateColumn("failed_count", gorm.Expr("failed_count + 1")).Error
}

// CancelPendingItems 取消任务中所有待处理子项
func (r *BatchRepository) CancelPendingItems(ctx context.Context, taskID uint64) (int64, error) {
	result := r.db.WithContext(ctx).
		Model(&model.BatchTaskItem{}).
		Where("task_id = ? AND status = ?", taskID, model.ItemStatusPending).
		Update("status", model.ItemStatusCancelled)
	return result.RowsAffected, result.Error
}

// RetryFailedItems 将失败子项重置为待处理
func (r *BatchRepository) RetryFailedItems(ctx context.Context, taskID uint64) (int64, error) {
	result := r.db.WithContext(ctx).
		Model(&model.BatchTaskItem{}).
		Where("task_id = ? AND status = ?", taskID, model.ItemStatusFailed).
		Updates(map[string]interface{}{
			"status":    model.ItemStatusPending,
			"error_msg": "",
			"result_id": 0,
		})
	return result.RowsAffected, result.Error
}
