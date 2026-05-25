package service

import (
	"context"
	"encoding/json"
	"errors"

	"oceanengine-backend/internal/app/batch/model"
	"oceanengine-backend/internal/app/batch/repository"
)

type CreateBatchRequest struct {
	TaskType   string                 `json:"task_type"`
	Title      string                 `json:"title"`
	Config     map[string]interface{} `json:"config"`
	AccountIDs []uint64               `json:"account_ids"`
	Overrides  map[uint64]interface{} `json:"overrides"`
}

type BatchService struct {
	repo *repository.BatchRepository
}

func NewBatchService(repo *repository.BatchRepository) *BatchService {
	return &BatchService{repo: repo}
}

func (s *BatchService) CreateTask(ctx context.Context, tenantID, userID uint64, localAccountMap map[uint64]uint64, req CreateBatchRequest) (*model.BatchTask, error) {
	configJSON, err := json.Marshal(req.Config)
	if err != nil {
		return nil, err
	}

	task := &model.BatchTask{
		TenantID:   tenantID,
		TaskType:   req.TaskType,
		Title:      req.Title,
		Config:     string(configJSON),
		Status:     model.TaskStatusPending,
		TotalCount: len(req.AccountIDs),
		CreatedBy:  userID,
	}

	items := make([]model.BatchTaskItem, 0, len(req.AccountIDs))
	for _, accountID := range req.AccountIDs {
		localID := localAccountMap[accountID]
		item := model.BatchTaskItem{
			AccountID:      accountID,
			LocalAccountID: localID,
			Status:         model.ItemStatusPending,
		}
		if override, ok := req.Overrides[accountID]; ok {
			overrideJSON, _ := json.Marshal(override)
			item.ConfigOverride = string(overrideJSON)
		}
		items = append(items, item)
	}

	if err := s.repo.CreateTask(ctx, task, items); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *BatchService) GetTask(ctx context.Context, tenantID, taskID uint64) (*model.BatchTask, error) {
	return s.repo.GetTask(ctx, tenantID, taskID)
}

func (s *BatchService) ListTasks(ctx context.Context, tenantID uint64, page, pageSize int) ([]model.BatchTask, int64, error) {
	return s.repo.ListTasks(ctx, tenantID, page, pageSize)
}

func (s *BatchService) CancelTask(ctx context.Context, tenantID, taskID uint64) error {
	task, err := s.repo.GetTask(ctx, tenantID, taskID)
	if err != nil {
		return err
	}
	if task.Status != model.TaskStatusPending && task.Status != model.TaskStatusRunning {
		return errors.New("只能取消待执行或执行中的任务")
	}
	if _, err := s.repo.CancelPendingItems(ctx, taskID); err != nil {
		return err
	}
	return s.repo.UpdateTaskStatus(ctx, taskID, model.TaskStatusCancelled)
}

func (s *BatchService) RetryTask(ctx context.Context, tenantID, taskID uint64) error {
	task, err := s.repo.GetTask(ctx, tenantID, taskID)
	if err != nil {
		return err
	}
	if task.Status != model.TaskStatusFailed && task.Status != model.TaskStatusPartialSuccess {
		return errors.New("只能重试失败或部分成功的任务")
	}
	if _, err := s.repo.RetryFailedItems(ctx, taskID); err != nil {
		return err
	}
	return s.repo.UpdateTaskStatus(ctx, taskID, model.TaskStatusPending)
}
