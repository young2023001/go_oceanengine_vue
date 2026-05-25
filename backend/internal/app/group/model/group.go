package model

import (
	"time"

	"gorm.io/gorm"
)

type GroupType string

const (
	GroupTypeFranchisee GroupType = "franchisee"
	GroupTypeRegion     GroupType = "region"
	GroupTypeCustom     GroupType = "custom"
)

type AccountGroup struct {
	ID        uint64         `gorm:"primaryKey" json:"id"`
	TenantID  uint64         `gorm:"not null;index:idx_tenant_type" json:"tenant_id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Type      GroupType      `gorm:"size:20;not null;default:custom;index:idx_tenant_type" json:"type"`
	ParentID  uint64         `gorm:"default:0" json:"parent_id"`
	Sort      int            `gorm:"default:0" json:"sort"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (AccountGroup) TableName() string {
	return "account_group"
}

type AccountGroupMember struct {
	AccountID uint64 `gorm:"primaryKey" json:"account_id"`
	GroupID   uint64 `gorm:"primaryKey;index" json:"group_id"`
}

func (AccountGroupMember) TableName() string {
	return "account_group_member"
}
