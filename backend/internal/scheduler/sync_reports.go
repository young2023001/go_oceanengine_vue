package scheduler

import (
	"context"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	accountModel "oceanengine-backend/internal/app/account/model"
	analyticsModel "oceanengine-backend/internal/app/analytics/model"
	tenantModel "oceanengine-backend/internal/app/tenant/model"
	"oceanengine-backend/pkg/ratelimiter"
)

type ReportSyncer struct {
	db      *gorm.DB
	limiter *ratelimiter.TenantRateLimiter
	logger  *zap.Logger
}

func NewReportSyncer(db *gorm.DB, limiter *ratelimiter.TenantRateLimiter, logger *zap.Logger) *ReportSyncer {
	return &ReportSyncer{db: db, limiter: limiter, logger: logger}
}

func (s *ReportSyncer) StartDaily(ctx context.Context) {
	for {
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day()+1, 2, 0, 0, 0, now.Location())
		timer := time.NewTimer(next.Sub(now))
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
			s.syncAllTenants(ctx, now.AddDate(0, 0, -1))
		}
	}
}

func (s *ReportSyncer) StartRealtime(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.syncAllTenants(ctx, time.Now())
		}
	}
}

func (s *ReportSyncer) syncAllTenants(ctx context.Context, date time.Time) {
	var tenants []tenantModel.Tenant
	s.db.Where("status = 1 AND token_status = 'active'").Find(&tenants)
	for _, tenant := range tenants {
		s.syncTenant(ctx, tenant, date)
	}
}

func (s *ReportSyncer) syncTenant(ctx context.Context, tenant tenantModel.Tenant, date time.Time) {
	var accounts []accountModel.LocalAccount
	s.db.Where("tenant_id = ?", tenant.ID).Find(&accounts)

	for _, acc := range accounts {
		if err := s.limiter.Wait(ctx, tenant.ID); err != nil {
			return
		}
		s.syncAccountReport(ctx, tenant, acc, date)
	}
}

func (s *ReportSyncer) syncAccountReport(ctx context.Context, tenant tenantModel.Tenant, acc accountModel.LocalAccount, date time.Time) {
	now := time.Now()
	report := analyticsModel.ReportDaily{
		TenantID:    tenant.ID,
		AccountID:   acc.ID,
		ProjectID:   0,
		PromotionID: 0,
		StatDate:    date,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	s.db.Where("tenant_id = ? AND account_id = ? AND project_id = 0 AND promotion_id = 0 AND stat_date = ?",
		tenant.ID, acc.ID, date).Assign(report).FirstOrCreate(&report)
}
