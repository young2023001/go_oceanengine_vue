package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
	"oceanengine-backend/internal/app/tenant/model"
)

type TenantRepository struct {
	db *gorm.DB
}

func NewTenantRepository(db *gorm.DB) *TenantRepository {
	return &TenantRepository{db: db}
}

func (r *TenantRepository) Create(ctx context.Context, tenant *model.Tenant) error {
	return r.db.WithContext(ctx).Create(tenant).Error
}

func (r *TenantRepository) GetByID(ctx context.Context, id uint64) (*model.Tenant, error) {
	var tenant model.Tenant
	err := r.db.WithContext(ctx).First(&tenant, id).Error
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *TenantRepository) List(ctx context.Context) ([]model.Tenant, error) {
	var tenants []model.Tenant
	err := r.db.WithContext(ctx).Where("status = 1").Find(&tenants).Error
	return tenants, err
}

func (r *TenantRepository) Update(ctx context.Context, tenant *model.Tenant) error {
	return r.db.WithContext(ctx).Save(tenant).Error
}

func (r *TenantRepository) UpdateToken(ctx context.Context, id uint64, accessToken, refreshToken string, expireAt *time.Time, status model.TokenStatus) error {
	return r.db.WithContext(ctx).Model(&model.Tenant{}).Where("id = ?", id).Updates(map[string]interface{}{
		"access_token":    accessToken,
		"refresh_token":   refreshToken,
		"token_expire_at": expireAt,
		"token_status":    status,
	}).Error
}

func (r *TenantRepository) ListNeedRefresh(ctx context.Context) ([]model.Tenant, error) {
	var tenants []model.Tenant
	err := r.db.WithContext(ctx).
		Where("status = 1 AND token_status IN ('active','expiring') AND token_expire_at < DATE_ADD(NOW(), INTERVAL 2 HOUR)").
		Find(&tenants).Error
	return tenants, err
}
