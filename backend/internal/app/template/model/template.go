package model

import (
	"time"

	"gorm.io/gorm"
)

type ProjectTemplate struct {
	ID        uint64         `gorm:"primaryKey" json:"id"`
	TenantID  uint64         `gorm:"not null;index" json:"tenant_id"`
	Name      string         `gorm:"size:200;not null" json:"name"`
	Config    string         `gorm:"type:json;not null" json:"config"`
	UseCount  int            `gorm:"default:0" json:"use_count"`
	CreatedBy uint64         `gorm:"default:0" json:"created_by"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ProjectTemplate) TableName() string {
	return "project_template"
}

type PromotionTemplate struct {
	ID        uint64         `gorm:"primaryKey" json:"id"`
	TenantID  uint64         `gorm:"not null;index" json:"tenant_id"`
	Name      string         `gorm:"size:200;not null" json:"name"`
	Config    string         `gorm:"type:json;not null" json:"config"`
	UseCount  int            `gorm:"default:0" json:"use_count"`
	CreatedBy uint64         `gorm:"default:0" json:"created_by"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (PromotionTemplate) TableName() string {
	return "promotion_template"
}
