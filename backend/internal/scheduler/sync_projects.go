package scheduler

import (
	"context"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	accountModel "oceanengine-backend/internal/app/account/model"
	projectModel "oceanengine-backend/internal/app/project/model"
	tenantModel "oceanengine-backend/internal/app/tenant/model"
	"oceanengine-backend/pkg/ocean"
	"oceanengine-backend/pkg/ratelimiter"
)

type ProjectSyncer struct {
	db      *gorm.DB
	sdk     *ocean.Client
	limiter *ratelimiter.TenantRateLimiter
	logger  *zap.Logger
}

func NewProjectSyncer(db *gorm.DB, sdk *ocean.Client, limiter *ratelimiter.TenantRateLimiter, logger *zap.Logger) *ProjectSyncer {
	return &ProjectSyncer{db: db, sdk: sdk, limiter: limiter, logger: logger}
}

func (s *ProjectSyncer) Start(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.syncAll(ctx)
		}
	}
}

func (s *ProjectSyncer) syncAll(ctx context.Context) {
	var tenants []tenantModel.Tenant
	s.db.Where("status = 1 AND token_status = ?", tenantModel.TokenStatusActive).Find(&tenants)

	for _, tenant := range tenants {
		s.syncTenant(ctx, tenant)
	}
}

func (s *ProjectSyncer) syncTenant(ctx context.Context, tenant tenantModel.Tenant) {
	var accounts []accountModel.LocalAccount
	s.db.Where("tenant_id = ?", tenant.ID).Find(&accounts)

	for _, acc := range accounts {
		if err := s.limiter.Wait(ctx, tenant.ID); err != nil {
			return
		}
		s.syncAccount(ctx, tenant, acc)
	}
}

func (s *ProjectSyncer) syncAccount(ctx context.Context, tenant tenantModel.Tenant, acc accountModel.LocalAccount) {
	resp, err := s.sdk.ListProjects(ctx, tenant.AccessToken, ocean.ListProjectsRequest{
		LocalAccountID: acc.LocalAccountID,
		Page:           1,
		PageSize:       100,
	})
	if err != nil {
		s.logger.Error("sync projects failed", zap.Uint64("account_id", acc.ID), zap.Error(err))
		return
	}

	now := time.Now()
	fiveMinAgo := now.Add(-5 * time.Minute)

	for _, p := range resp.Data.List {
		var existing projectModel.LocalProject
		result := s.db.Where("tenant_id = ? AND project_id = ?", tenant.ID, p.ProjectID).First(&existing)

		if result.Error == nil && existing.LocalUpdatedAt != nil && existing.LocalUpdatedAt.After(fiveMinAgo) {
			continue
		}

		project := projectModel.LocalProject{
			TenantID:    tenant.ID,
			AccountID:   acc.ID,
			ProjectID:   p.ProjectID,
			Name:        p.Name,
			StatusFirst: p.StatusFirst,
			SyncedAt:    &now,
		}
		if result.Error == nil {
			project.ID = existing.ID
		}
		s.db.Save(&project)
	}
}
