package service

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"oceanengine-backend/internal/app/tenant/model"
)

// DBTokenStore 基于数据库的 TokenStore 实现
type DBTokenStore struct {
	db *gorm.DB
}

// NewDBTokenStore 创建 DBTokenStore
func NewDBTokenStore(db *gorm.DB) *DBTokenStore {
	return &DBTokenStore{db: db}
}

// GetToken 从数据库获取租户的 access token
func (s *DBTokenStore) GetToken(ctx context.Context, tenantID uint64) (string, error) {
	var tenant model.Tenant
	err := s.db.WithContext(ctx).Select("access_token").First(&tenant, tenantID).Error
	if err != nil {
		return "", fmt.Errorf("get token for tenant %d: %w", tenantID, err)
	}
	if tenant.AccessToken == "" {
		return "", fmt.Errorf("tenant %d has no access token", tenantID)
	}
	return tenant.AccessToken, nil
}
