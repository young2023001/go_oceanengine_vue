package service

import (
	"context"
	"time"

	"go.uber.org/zap"
	"oceanengine-backend/internal/app/tenant/model"
	"oceanengine-backend/internal/app/tenant/repository"
)

type TokenRefresher struct {
	repo   *repository.TenantRepository
	oauth  OAuthClient
	logger *zap.Logger
}

func NewTokenRefresher(repo *repository.TenantRepository, oauth OAuthClient, logger *zap.Logger) *TokenRefresher {
	return &TokenRefresher{repo: repo, oauth: oauth, logger: logger}
}

func (r *TokenRefresher) Start(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.refreshAll(ctx)
		}
	}
}

func (r *TokenRefresher) refreshAll(ctx context.Context) {
	tenants, err := r.repo.ListNeedRefresh(ctx)
	if err != nil {
		r.logger.Error("list tenants for refresh", zap.Error(err))
		return
	}

	for i := range tenants {
		r.refreshOne(ctx, tenants[i])
	}
}

func (r *TokenRefresher) refreshOne(ctx context.Context, t model.Tenant) {
	accessToken, refreshToken, expiresIn, err := r.oauth.RefreshToken(ctx, t.OAuthAppID, t.OAuthSecret, t.RefreshToken)
	if err != nil {
		r.logger.Error("refresh token failed", zap.Uint64("tenant_id", t.ID), zap.Error(err))
		r.repo.UpdateToken(ctx, t.ID, t.AccessToken, t.RefreshToken, t.TokenExpireAt, model.TokenStatusFailed)
		return
	}

	expireAt := time.Now().Add(time.Duration(expiresIn) * time.Second)
	if err := r.repo.UpdateToken(ctx, t.ID, accessToken, refreshToken, &expireAt, model.TokenStatusActive); err != nil {
		r.logger.Error("save refreshed token", zap.Uint64("tenant_id", t.ID), zap.Error(err))
	}
}
