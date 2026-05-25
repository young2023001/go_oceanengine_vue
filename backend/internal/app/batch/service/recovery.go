package service

import (
	"context"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"oceanengine-backend/internal/app/batch/model"
)

type TaskRecovery struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewTaskRecovery(db *gorm.DB, logger *zap.Logger) *TaskRecovery {
	return &TaskRecovery{db: db, logger: logger}
}

func (r *TaskRecovery) RecoverOnStartup(ctx context.Context) {
	var tasks []model.BatchTask
	r.db.Where("status = ?", model.TaskStatusRunning).Find(&tasks)

	r.logger.Info("recovering batch tasks", zap.Int("count", len(tasks)))

	for _, task := range tasks {
		r.db.Model(&model.BatchTaskItem{}).
			Where("task_id = ? AND status = ?", task.ID, model.ItemStatusRunning).
			Update("status", model.ItemStatusPending)
	}
}
