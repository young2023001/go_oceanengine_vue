package model

import (
	"time"

	"gorm.io/gorm"
)

type LocalProject struct {
	ID             uint64         `gorm:"primaryKey" json:"id"`
	TenantID       uint64         `gorm:"not null;uniqueIndex:uk_tenant_project" json:"tenant_id"`
	AccountID      uint64         `gorm:"not null;index" json:"account_id"`
	ProjectID      uint64         `gorm:"not null;uniqueIndex:uk_tenant_project" json:"project_id"`
	Name           string         `gorm:"size:200" json:"name"`
	MarketingGoal  string         `gorm:"size:50" json:"marketing_goal"`
	DeliveryScene  string         `gorm:"size:50" json:"delivery_scene"`
	AdType         string         `gorm:"size:50" json:"ad_type"`
	ExternalAction string         `gorm:"size:50" json:"external_action"`
	BidType        string         `gorm:"size:20" json:"bid_type"`
	Bid            int            `gorm:"default:0" json:"bid"`
	BudgetMode     string         `gorm:"size:30" json:"budget_mode"`
	Budget         int            `gorm:"default:0" json:"budget"`
	StatusFirst    string         `gorm:"size:50;index" json:"status_first"`
	StatusSecond   string         `gorm:"size:50" json:"status_second"`
	StartTime      string         `gorm:"size:20" json:"start_time"`
	EndTime        string         `gorm:"size:20" json:"end_time"`
	PoiID          uint64         `gorm:"default:0" json:"poi_id"`
	PoiName        string         `gorm:"size:200" json:"poi_name"`
	ProductID      uint64         `gorm:"default:0" json:"product_id"`
	SyncedAt       *time.Time     `json:"synced_at"`
	LocalUpdatedAt *time.Time     `json:"local_updated_at"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (LocalProject) TableName() string {
	return "local_project"
}

type LocalPromotion struct {
	ID           uint64         `gorm:"primaryKey" json:"id"`
	TenantID     uint64         `gorm:"not null;uniqueIndex:uk_tenant_promotion" json:"tenant_id"`
	AccountID    uint64         `gorm:"not null;index" json:"account_id"`
	ProjectID    uint64         `gorm:"not null;index" json:"project_id"`
	PromotionID  uint64         `gorm:"not null;uniqueIndex:uk_tenant_promotion" json:"promotion_id"`
	Name         string         `gorm:"size:200" json:"name"`
	StatusFirst  string         `gorm:"size:50" json:"status_first"`
	StatusSecond string         `gorm:"size:50" json:"status_second"`
	AwemeID      string         `gorm:"size:100" json:"aweme_id"`
	SyncedAt     *time.Time     `json:"synced_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (LocalPromotion) TableName() string {
	return "local_promotion"
}
