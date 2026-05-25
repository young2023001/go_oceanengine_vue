package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"

	"oceanengine-backend/internal/app/batch/model"
	"oceanengine-backend/internal/app/batch/repository"
	oceanengine "oceanengine-backend/pkg/oceanengine"
	"oceanengine-backend/pkg/ratelimiter"
	"oceanengine-backend/pkg/retry"
)

// TokenStore 提供租户 access token 的接口
type TokenStore interface {
	GetToken(ctx context.Context, tenantID uint64) (string, error)
}

// Worker 批量任务执行引擎
type Worker struct {
	repo       *repository.BatchRepository
	sdk        *oceanengine.Client
	limiter    *ratelimiter.TenantRateLimiter
	tokenStore TokenStore
	logger     *zap.Logger
}

// NewWorker 创建 Worker
func NewWorker(
	repo *repository.BatchRepository,
	sdk *oceanengine.Client,
	limiter *ratelimiter.TenantRateLimiter,
	tokenStore TokenStore,
	logger *zap.Logger,
) *Worker {
	return &Worker{
		repo:       repo,
		sdk:        sdk,
		limiter:    limiter,
		tokenStore: tokenStore,
		logger:     logger,
	}
}

// Start 启动 worker，每 2 秒 tick 一次
func (w *Worker) Start(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	w.logger.Info("batch worker started")

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("batch worker stopped")
			return
		case <-ticker.C:
			w.processNext(ctx)
		}
	}
}

// processNext 处理所有运行中的任务
func (w *Worker) processNext(ctx context.Context) {
	tasks, err := w.repo.ListRunningTasks(ctx)
	if err != nil {
		w.logger.Error("list running tasks failed", zap.Error(err))
		return
	}

	for i := range tasks {
		w.processTask(ctx, &tasks[i])
	}
}

// processTask 获取 pending items 并逐条执行
func (w *Worker) processTask(ctx context.Context, task *model.BatchTask) {
	items, err := w.repo.GetPendingItems(ctx, task.ID, 10)
	if err != nil {
		w.logger.Error("get pending items failed", zap.Uint64("task_id", task.ID), zap.Error(err))
		return
	}

	if len(items) == 0 {
		w.finalizeTask(ctx, task)
		return
	}

	for i := range items {
		if err := w.limiter.Wait(ctx, task.TenantID); err != nil {
			w.logger.Error("rate limiter wait failed", zap.Error(err))
			return
		}
		w.executeItem(ctx, task, &items[i])
	}
}

// executeItem 根据 task_type 分发执行
func (w *Worker) executeItem(ctx context.Context, task *model.BatchTask, item *model.BatchTaskItem) {
	var resultID uint64
	var execErr error

	switch task.TaskType {
	case "create_project":
		resultID, execErr = w.execCreateProject(ctx, task, item)
	case "update_budget":
		execErr = w.execUpdateBudget(ctx, task, item)
	case "update_bid":
		execErr = w.execUpdateBid(ctx, task, item)
	case "update_status":
		execErr = w.execUpdateStatus(ctx, task, item)
	default:
		execErr = fmt.Errorf("unknown task type: %s", task.TaskType)
	}

	if execErr != nil {
		_ = w.repo.UpdateItemStatus(ctx, item.ID, model.ItemStatusFailed, 0, execErr.Error())
		_ = w.repo.IncrFailed(ctx, task.ID)
		w.logger.Warn("execute item failed",
			zap.Uint64("task_id", task.ID),
			zap.Uint64("item_id", item.ID),
			zap.Error(execErr),
		)
	} else {
		_ = w.repo.UpdateItemStatus(ctx, item.ID, model.ItemStatusSuccess, resultID, "")
		_ = w.repo.IncrSuccess(ctx, task.ID)
	}
}

// execCreateProject 创建广告项目
func (w *Worker) execCreateProject(ctx context.Context, task *model.BatchTask, item *model.BatchTaskItem) (uint64, error) {
	accessToken, err := w.tokenStore.GetToken(ctx, task.TenantID)
	if err != nil {
		return 0, err
	}

	config := mergeConfig(task.Config, item.ConfigOverride)

	var req map[string]interface{}
	if err := json.Unmarshal([]byte(config), &req); err != nil {
		return 0, err
	}
	req["advertiser_id"] = item.AccountID

	var result struct {
		Data struct {
			ProjectID uint64 `json:"project_id"`
		} `json:"data"`
	}
	err = retry.Do(ctx, retry.DefaultConfig(), func() error {
		return w.sdk.PostWithToken(ctx, accessToken, "/2/campaign/create/", req, &result)
	})
	if err != nil {
		return 0, err
	}

	return result.Data.ProjectID, nil
}

// execUpdateBudget 更新预算
func (w *Worker) execUpdateBudget(ctx context.Context, task *model.BatchTask, item *model.BatchTaskItem) error {
	accessToken, err := w.tokenStore.GetToken(ctx, task.TenantID)
	if err != nil {
		return err
	}

	config := mergeConfig(task.Config, item.ConfigOverride)

	var req map[string]interface{}
	if err := json.Unmarshal([]byte(config), &req); err != nil {
		return err
	}
	req["advertiser_id"] = item.AccountID

	return retry.Do(ctx, retry.DefaultConfig(), func() error {
		return w.sdk.PostWithToken(ctx, accessToken, "/2/campaign/update/", req, nil)
	})
}

// execUpdateBid 更新出价
func (w *Worker) execUpdateBid(ctx context.Context, task *model.BatchTask, item *model.BatchTaskItem) error {
	accessToken, err := w.tokenStore.GetToken(ctx, task.TenantID)
	if err != nil {
		return err
	}

	config := mergeConfig(task.Config, item.ConfigOverride)

	var req map[string]interface{}
	if err := json.Unmarshal([]byte(config), &req); err != nil {
		return err
	}
	req["advertiser_id"] = item.AccountID

	return retry.Do(ctx, retry.DefaultConfig(), func() error {
		return w.sdk.PostWithToken(ctx, accessToken, "/2/ad/update/", req, nil)
	})
}

// execUpdateStatus 更新广告状态
func (w *Worker) execUpdateStatus(ctx context.Context, task *model.BatchTask, item *model.BatchTaskItem) error {
	accessToken, err := w.tokenStore.GetToken(ctx, task.TenantID)
	if err != nil {
		return err
	}

	config := mergeConfig(task.Config, item.ConfigOverride)

	var req map[string]interface{}
	if err := json.Unmarshal([]byte(config), &req); err != nil {
		return err
	}
	req["advertiser_id"] = item.AccountID

	return retry.Do(ctx, retry.DefaultConfig(), func() error {
		return w.sdk.PostWithToken(ctx, accessToken, "/2/ad/update/status/", req, nil)
	})
}

// finalizeTask 根据 success/failed 计数决定最终状态
func (w *Worker) finalizeTask(ctx context.Context, task *model.BatchTask) {
	// 重新获取最新的任务数据
	latest, err := w.repo.GetTaskByID(ctx, task.ID)
	if err != nil {
		w.logger.Error("get task for finalize failed", zap.Uint64("task_id", task.ID), zap.Error(err))
		return
	}

	var finalStatus model.TaskStatus
	switch {
	case latest.SuccessCount == latest.TotalCount:
		finalStatus = model.TaskStatusCompleted
	case latest.FailedCount == latest.TotalCount:
		finalStatus = model.TaskStatusFailed
	default:
		finalStatus = model.TaskStatusPartialSuccess
	}

	if err := w.repo.UpdateTaskStatus(ctx, task.ID, finalStatus); err != nil {
		w.logger.Error("finalize task failed", zap.Uint64("task_id", task.ID), zap.Error(err))
	} else {
		w.logger.Info("task finalized",
			zap.Uint64("task_id", task.ID),
			zap.String("status", string(finalStatus)),
		)
	}
}

// mergeConfig 合并任务配置和子项覆盖配置
func mergeConfig(taskConfig, itemOverride string) string {
	if itemOverride == "" || itemOverride == "{}" || itemOverride == "null" {
		return taskConfig
	}

	var base map[string]interface{}
	if err := json.Unmarshal([]byte(taskConfig), &base); err != nil {
		return taskConfig
	}

	var override map[string]interface{}
	if err := json.Unmarshal([]byte(itemOverride), &override); err != nil {
		return taskConfig
	}

	for k, v := range override {
		base[k] = v
	}

	merged, err := json.Marshal(base)
	if err != nil {
		return taskConfig
	}
	return string(merged)
}
