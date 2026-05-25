package model

import (
	"time"

	"gorm.io/gorm"
)

type TaskStatus string

const (
	TaskStatusPending        TaskStatus = "pending"
	TaskStatusRunning        TaskStatus = "running"
	TaskStatusCompleted      TaskStatus = "completed"
	TaskStatusPartialSuccess TaskStatus = "partial_success"
	TaskStatusFailed         TaskStatus = "failed"
	TaskStatusCancelled      TaskStatus = "cancelled"
)

type ItemStatus string

const (
	ItemStatusPending   ItemStatus = "pending"
	ItemStatusRunning   ItemStatus = "running"
	ItemStatusSuccess   ItemStatus = "success"
	ItemStatusFailed    ItemStatus = "failed"
	ItemStatusSkipped   ItemStatus = "skipped"
	ItemStatusCancelled ItemStatus = "cancelled"
)

type BatchTask struct {
	ID           uint64         `gorm:"primaryKey" json:"id"`
	TenantID     uint64         `gorm:"not null;index:idx_tenant_status" json:"tenant_id"`
	TaskType     string         `gorm:"size:50;not null" json:"task_type"`
	Title        string         `gorm:"size:200;not null" json:"title"`
	Config       string         `gorm:"type:json;not null" json:"config"`
	Status       TaskStatus     `gorm:"size:30;default:pending;index:idx_tenant_status" json:"status"`
	TotalCount   int            `gorm:"default:0" json:"total_count"`
	SuccessCount int            `gorm:"default:0" json:"success_count"`
	FailedCount  int            `gorm:"default:0" json:"failed_count"`
	CreatedBy    uint64         `gorm:"default:0" json:"created_by"`
	StartedAt    *time.Time     `json:"started_at"`
	CompletedAt  *time.Time     `json:"completed_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (BatchTask) TableName() string {
	return "batch_task"
}

type BatchTaskItem struct {
	ID             uint64     `gorm:"primaryKey" json:"id"`
	TenantID       uint64     `gorm:"not null;index" json:"tenant_id"`
	TaskID         uint64     `gorm:"not null;index" json:"task_id"`
	AccountID      uint64     `gorm:"not null" json:"account_id"`
	LocalAccountID uint64     `gorm:"not null" json:"local_account_id"`
	ConfigOverride string     `gorm:"type:json" json:"config_override"`
	Status         ItemStatus `gorm:"size:20;default:pending;index" json:"status"`
	ResultID       uint64     `gorm:"default:0" json:"result_id"`
	ErrorMsg       string     `gorm:"size:500" json:"error_msg"`
	ExecutedAt     *time.Time `json:"executed_at"`
}

func (BatchTaskItem) TableName() string {
	return "batch_task_item"
}
