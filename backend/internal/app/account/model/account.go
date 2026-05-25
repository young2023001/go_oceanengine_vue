package model

import (
	"time"

	"gorm.io/gorm"
)

type LocalAccount struct {
	ID             uint64         `gorm:"primaryKey" json:"id"`
	TenantID       uint64         `gorm:"not null;uniqueIndex:uk_tenant_account" json:"tenant_id"`
	LocalAccountID uint64         `gorm:"not null;uniqueIndex:uk_tenant_account" json:"local_account_id"`
	Name           string         `gorm:"size:200;not null" json:"name"`
	Status         string         `gorm:"size:50;default:active" json:"status"`
	Balance        float64        `gorm:"type:decimal(12,2);default:0" json:"balance"`
	LastSyncAt     *time.Time     `json:"last_sync_at"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (LocalAccount) TableName() string {
	return "local_account"
}

type Store struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	TenantID  uint64    `gorm:"not null;index" json:"tenant_id"`
	AccountID uint64    `gorm:"not null;index" json:"account_id"`
	PoiID     uint64    `gorm:"default:0" json:"poi_id"`
	Name      string    `gorm:"size:200;not null" json:"name"`
	City      string    `gorm:"size:50" json:"city"`
	District  string    `gorm:"size:50" json:"district"`
	Address   string    `gorm:"size:500" json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Store) TableName() string {
	return "store"
}
