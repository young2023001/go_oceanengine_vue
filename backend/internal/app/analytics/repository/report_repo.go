package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"oceanengine-backend/internal/app/analytics/model"
)

type ReportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

type OverviewResult struct {
	TotalCost    float64 `json:"total_cost"`
	TotalShow    int64   `json:"total_show"`
	TotalClick   int64   `json:"total_click"`
	TotalConvert int64   `json:"total_convert"`
}

// GetOverview 获取汇总数据（账户级，project_id=0, promotion_id=0）
func (r *ReportRepository) GetOverview(ctx context.Context, tenantID uint64, accountIDs []uint64, date time.Time) (*OverviewResult, error) {
	var result OverviewResult
	err := r.db.WithContext(ctx).
		Model(&model.ReportDaily{}).
		Select("COALESCE(SUM(cost),0) AS total_cost, COALESCE(SUM(show_cnt),0) AS total_show, COALESCE(SUM(click_cnt),0) AS total_click, COALESCE(SUM(convert_cnt),0) AS total_convert").
		Where("tenant_id = ? AND account_id IN ? AND stat_date = ? AND project_id = 0 AND promotion_id = 0", tenantID, accountIDs, date).
		Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

type TrendPoint struct {
	StatDate   string  `json:"stat_date"`
	Cost       float64 `json:"cost"`
	ShowCnt    int64   `json:"show_cnt"`
	ClickCnt   int64   `json:"click_cnt"`
	ConvertCnt int64   `json:"convert_cnt"`
}

// GetTrend 获取趋势数据（按天聚合）
func (r *ReportRepository) GetTrend(ctx context.Context, tenantID uint64, accountIDs []uint64, startDate, endDate time.Time) ([]TrendPoint, error) {
	var points []TrendPoint
	err := r.db.WithContext(ctx).
		Model(&model.ReportDaily{}).
		Select("DATE_FORMAT(stat_date,'%Y-%m-%d') AS stat_date, SUM(cost) AS cost, SUM(show_cnt) AS show_cnt, SUM(click_cnt) AS click_cnt, SUM(convert_cnt) AS convert_cnt").
		Where("tenant_id = ? AND account_id IN ? AND stat_date BETWEEN ? AND ? AND project_id = 0 AND promotion_id = 0", tenantID, accountIDs, startDate, endDate).
		Group("stat_date").
		Order("stat_date ASC").
		Scan(&points).Error
	if err != nil {
		return nil, err
	}
	return points, nil
}

type RankItem struct {
	AccountID  uint64  `json:"account_id"`
	Name       string  `json:"name"`
	Cost       float64 `json:"cost"`
	ConvertCnt int64   `json:"convert_cnt"`
}

// GetRank 获取排行榜（JOIN local_account 获取名称）
func (r *ReportRepository) GetRank(ctx context.Context, tenantID uint64, accountIDs []uint64, startDate, endDate time.Time, orderBy string, limit int) ([]RankItem, error) {
	allowedOrders := map[string]string{
		"cost":        "cost DESC",
		"convert_cnt": "convert_cnt DESC",
	}
	order, ok := allowedOrders[orderBy]
	if !ok {
		order = "cost DESC"
	}

	var items []RankItem
	err := r.db.WithContext(ctx).
		Table("report_daily r").
		Select("r.account_id, a.name, SUM(r.cost) AS cost, SUM(r.convert_cnt) AS convert_cnt").
		Joins("JOIN local_account a ON a.id = r.account_id").
		Where("r.tenant_id = ? AND r.account_id IN ? AND r.stat_date BETWEEN ? AND ? AND r.project_id = 0 AND r.promotion_id = 0", tenantID, accountIDs, startDate, endDate).
		Group("r.account_id, a.name").
		Order(order).
		Limit(limit).
		Scan(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

// GetDetail 获取明细数据（分页）
func (r *ReportRepository) GetDetail(ctx context.Context, tenantID uint64, accountIDs []uint64, startDate, endDate time.Time, page, pageSize int) ([]model.ReportDaily, int64, error) {
	var total int64
	query := r.db.WithContext(ctx).
		Model(&model.ReportDaily{}).
		Where("tenant_id = ? AND account_id IN ? AND stat_date BETWEEN ? AND ?", tenantID, accountIDs, startDate, endDate)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var records []model.ReportDaily
	offset := (page - 1) * pageSize
	err := query.
		Order("stat_date DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&records).Error
	if err != nil {
		return nil, 0, err
	}
	return records, total, nil
}

// Upsert 插入或更新报表记录
func (r *ReportRepository) Upsert(ctx context.Context, report *model.ReportDaily) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "tenant_id"},
				{Name: "account_id"},
				{Name: "project_id"},
				{Name: "promotion_id"},
				{Name: "stat_date"},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"cost", "show_cnt", "click_cnt", "ctr", "cpc", "cpm",
				"convert_cnt", "conversion_rate", "conversion_cost",
				"form_cnt", "phone_confirm_cnt", "phone_connect_cnt",
				"message_action_cnt", "dy_like", "dy_comment", "dy_share",
				"dy_follow", "total_play", "play_over_rate", "updated_at",
			}),
		}).
		Create(report).Error
}
