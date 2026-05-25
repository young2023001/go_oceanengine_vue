package service

import (
	"context"
	"time"

	"oceanengine-backend/internal/app/analytics/repository"
)

const dateLayout = "2006-01-02"

type OverviewRequest struct {
	Date string `json:"date" form:"date"`
}

type TrendRequest struct {
	StartDate string `json:"start_date" form:"start_date" binding:"required"`
	EndDate   string `json:"end_date" form:"end_date" binding:"required"`
}

type RankRequest struct {
	StartDate string `json:"start_date" form:"start_date" binding:"required"`
	EndDate   string `json:"end_date" form:"end_date" binding:"required"`
	OrderBy   string `json:"order_by" form:"order_by"`
	Limit     int    `json:"limit" form:"limit"`
}

type CompareRequest struct {
	Date        string `json:"date" form:"date"`
	CompareType string `json:"compare_type" form:"compare_type"`
}

type DetailRequest struct {
	StartDate string `json:"start_date" form:"start_date" binding:"required"`
	EndDate   string `json:"end_date" form:"end_date" binding:"required"`
	Page      int    `json:"page" form:"page"`
	PageSize  int    `json:"page_size" form:"page_size"`
}

type AnalyticsService struct {
	repo *repository.ReportRepository
}

func NewAnalyticsService(repo *repository.ReportRepository) *AnalyticsService {
	return &AnalyticsService{repo: repo}
}

func (s *AnalyticsService) GetOverview(ctx context.Context, tenantID uint64, scopeIDs []uint64, req OverviewRequest) (*repository.OverviewResult, error) {
	date := time.Now()
	if req.Date != "" {
		parsed, err := time.Parse(dateLayout, req.Date)
		if err != nil {
			return nil, err
		}
		date = parsed
	}
	return s.repo.GetOverview(ctx, tenantID, scopeIDs, date)
}

func (s *AnalyticsService) GetTrend(ctx context.Context, tenantID uint64, scopeIDs []uint64, req TrendRequest) ([]repository.TrendPoint, error) {
	startDate, err := time.Parse(dateLayout, req.StartDate)
	if err != nil {
		return nil, err
	}
	endDate, err := time.Parse(dateLayout, req.EndDate)
	if err != nil {
		return nil, err
	}
	return s.repo.GetTrend(ctx, tenantID, scopeIDs, startDate, endDate)
}

func (s *AnalyticsService) GetRank(ctx context.Context, tenantID uint64, scopeIDs []uint64, req RankRequest) ([]repository.RankItem, error) {
	startDate, err := time.Parse(dateLayout, req.StartDate)
	if err != nil {
		return nil, err
	}
	endDate, err := time.Parse(dateLayout, req.EndDate)
	if err != nil {
		return nil, err
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	return s.repo.GetRank(ctx, tenantID, scopeIDs, startDate, endDate, req.OrderBy, limit)
}

func (s *AnalyticsService) GetCompare(ctx context.Context, tenantID uint64, scopeIDs []uint64, req CompareRequest) (map[string]*repository.OverviewResult, error) {
	current := time.Now()
	if req.Date != "" {
		parsed, err := time.Parse(dateLayout, req.Date)
		if err != nil {
			return nil, err
		}
		current = parsed
	}

	var previous time.Time
	switch req.CompareType {
	case "week":
		previous = current.AddDate(0, 0, -7)
	default:
		previous = current.AddDate(0, 0, -1)
	}

	currentResult, err := s.repo.GetOverview(ctx, tenantID, scopeIDs, current)
	if err != nil {
		return nil, err
	}
	previousResult, err := s.repo.GetOverview(ctx, tenantID, scopeIDs, previous)
	if err != nil {
		return nil, err
	}

	result := map[string]*repository.OverviewResult{
		"current":  currentResult,
		"previous": previousResult,
	}
	return result, nil
}

func (s *AnalyticsService) GetDetail(ctx context.Context, tenantID uint64, scopeIDs []uint64, req DetailRequest) ([]interface{}, int64, error) {
	startDate, err := time.Parse(dateLayout, req.StartDate)
	if err != nil {
		return nil, 0, err
	}
	endDate, err := time.Parse(dateLayout, req.EndDate)
	if err != nil {
		return nil, 0, err
	}
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	records, total, err := s.repo.GetDetail(ctx, tenantID, scopeIDs, startDate, endDate, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	items := make([]interface{}, len(records))
	for i, r := range records {
		items[i] = r
	}
	return items, total, nil
}
