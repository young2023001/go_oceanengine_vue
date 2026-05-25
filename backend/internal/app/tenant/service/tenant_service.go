package service

import (
	"context"
	"time"

	"oceanengine-backend/internal/app/tenant/model"
	"oceanengine-backend/internal/app/tenant/repository"
)

type TenantService struct {
	repo *repository.TenantRepository
}

func NewTenantService(repo *repository.TenantRepository) *TenantService {
	return &TenantService{repo: repo}
}

type CreateTenantRequest struct {
	Name        string `json:"name" binding:"required"`
	OAuthAppID  string `json:"oauth_app_id" binding:"required"`
	OAuthSecret string `json:"oauth_secret" binding:"required"`
}

func (s *TenantService) Create(ctx context.Context, req CreateTenantRequest) (*model.Tenant, error) {
	tenant := &model.Tenant{
		Name:        req.Name,
		OAuthAppID:  req.OAuthAppID,
		OAuthSecret: req.OAuthSecret,
		TokenStatus: model.TokenStatusExpired,
		Status:      1,
	}
	if err := s.repo.Create(ctx, tenant); err != nil {
		return nil, err
	}
	return tenant, nil
}

func (s *TenantService) GetByID(ctx context.Context, id uint64) (*model.Tenant, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TenantService) List(ctx context.Context) ([]model.Tenant, error) {
	return s.repo.List(ctx)
}

func (s *TenantService) SaveOAuthToken(ctx context.Context, tenantID uint64, accessToken, refreshToken string, expiresIn int64) error {
	expireAt := time.Now().Add(time.Duration(expiresIn) * time.Second)
	return s.repo.UpdateToken(ctx, tenantID, accessToken, refreshToken, &expireAt, model.TokenStatusActive)
}
