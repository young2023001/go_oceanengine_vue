package model

import (
	"time"
)

type ReportDaily struct {
	ID               uint64    `gorm:"primaryKey" json:"id"`
	TenantID         uint64    `gorm:"not null;uniqueIndex:uk_report" json:"tenant_id"`
	AccountID        uint64    `gorm:"not null;uniqueIndex:uk_report;index:idx_account_date" json:"account_id"`
	ProjectID        uint64    `gorm:"default:0;uniqueIndex:uk_report" json:"project_id"`
	PromotionID      uint64    `gorm:"default:0;uniqueIndex:uk_report" json:"promotion_id"`
	StatDate         time.Time `gorm:"type:date;not null;uniqueIndex:uk_report;index:idx_account_date" json:"stat_date"`
	Cost             float64   `gorm:"type:decimal(12,2);default:0" json:"cost"`
	ShowCnt          int64     `gorm:"default:0" json:"show_cnt"`
	ClickCnt         int64     `gorm:"default:0" json:"click_cnt"`
	Ctr              float64   `gorm:"type:decimal(8,4);default:0" json:"ctr"`
	Cpc              float64   `gorm:"type:decimal(8,2);default:0" json:"cpc"`
	Cpm              float64   `gorm:"type:decimal(8,2);default:0" json:"cpm"`
	ConvertCnt       int64     `gorm:"default:0" json:"convert_cnt"`
	ConversionRate   float64   `gorm:"type:decimal(8,4);default:0" json:"conversion_rate"`
	ConversionCost   float64   `gorm:"type:decimal(10,2);default:0" json:"conversion_cost"`
	FormCnt          int64     `gorm:"default:0" json:"form_cnt"`
	PhoneConfirmCnt  int64     `gorm:"default:0" json:"phone_confirm_cnt"`
	PhoneConnectCnt  int64     `gorm:"default:0" json:"phone_connect_cnt"`
	MessageActionCnt int64     `gorm:"default:0" json:"message_action_cnt"`
	DyLike           int64     `gorm:"default:0" json:"dy_like"`
	DyComment        int64     `gorm:"default:0" json:"dy_comment"`
	DyShare          int64     `gorm:"default:0" json:"dy_share"`
	DyFollow         int64     `gorm:"default:0" json:"dy_follow"`
	TotalPlay        int64     `gorm:"default:0" json:"total_play"`
	PlayOverRate     float64   `gorm:"type:decimal(8,4);default:0" json:"play_over_rate"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (ReportDaily) TableName() string {
	return "report_daily"
}

type ExportTask struct {
	ID        uint64     `gorm:"primaryKey" json:"id"`
	TenantID  uint64     `gorm:"not null;index" json:"tenant_id"`
	Title     string     `gorm:"size:200" json:"title"`
	Status    string     `gorm:"size:20;default:pending" json:"status"`
	FilePath  string     `gorm:"size:500" json:"file_path"`
	RowCount  int        `gorm:"default:0" json:"row_count"`
	ExpireAt  *time.Time `json:"expire_at"`
	CreatedBy uint64     `gorm:"default:0" json:"created_by"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (ExportTask) TableName() string {
	return "export_task"
}
