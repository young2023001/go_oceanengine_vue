package model

import (
	"time"

	"gorm.io/gorm"
)

type TokenStatus string

const (
	TokenStatusActive   TokenStatus = "active"
	TokenStatusExpiring TokenStatus = "expiring"
	TokenStatusExpired  TokenStatus = "expired"
	TokenStatusFailed   TokenStatus = "failed"
)

type Tenant struct {
	ID             uint64         `gorm:"primaryKey" json:"id"`
	Name           string         `gorm:"size:200;not null" json:"name"`
	OAuthAppID     string         `gorm:"column:oauth_app_id;size:100;not null" json:"oauth_app_id"`
	OAuthSecret    string         `gorm:"column:oauth_secret;size:200;not null" json:"-"`
	AccessToken    string         `gorm:"size:500" json:"-"`
	RefreshToken   string         `gorm:"size:500" json:"-"`
	TokenExpireAt  *time.Time     `json:"-"`
	TokenStatus    TokenStatus    `gorm:"size:20;default:active" json:"token_status"`
	OrganizationID uint64         `gorm:"default:0" json:"organization_id"`
	Status         int8           `gorm:"default:1;index" json:"status"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Tenant) TableName() string {
	return "tenant"
}
