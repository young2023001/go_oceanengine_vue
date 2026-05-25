# Phase 3: 数据分析实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 实现报表数据同步、汇总看板、多维聚合分析（按门店/加盟商/区域/投手）、趋势环比、排行榜、报表导出。

**Architecture:** 新增 analytics 模块 + scheduler/sync_reports.go 定时同步。报表按三个粒度存储（账户级/项目级/单元级），聚合查询通过 SQL JOIN + GROUP BY 实现。导出为异步任务生成 Excel。

**Tech Stack:** Go 1.21 + Gin + GORM + MySQL + Redis + excelize (Excel生成)

---

## 文件结构

### 新建文件

| 路径 | 职责 |
|------|------|
| backend/internal/app/analytics/model/report.go | ReportDaily + ExportTask 模型 |
| backend/internal/app/analytics/repository/report_repo.go | 报表查询（聚合/趋势/排行） |
| backend/internal/app/analytics/repository/export_repo.go | 导出任务 CRUD |
| backend/internal/app/analytics/service/analytics_service.go | 分析业务逻辑 |
| backend/internal/app/analytics/service/export_service.go | 导出任务执行 |
| backend/internal/app/analytics/api/analytics_handler.go | 分析 HTTP handler |
| backend/internal/scheduler/sync_reports.go | 报表同步定时任务 |

### 修改文件

| 路径 | 改动 |
|------|------|
| backend/internal/router/router.go | 注册 analytics 路由 |
| backend/cmd/server/main.go | 启动报表同步 goroutine |
| backend/go.mod | 添加 excelize 依赖 |

---

## Task 0: AutoMigrate 注册 Phase 3 新表

**Files:**
- Modify: backend/cmd/server/main.go

- [ ] **Step 1: 添加 import 和 AutoMigrate**

```go
import (
	analyticsModel "oceanengine-backend/internal/app/analytics/model"
)

// 在现有 AutoMigrate 中追加
db.AutoMigrate(
	&analyticsModel.ReportDaily{},
	&analyticsModel.ExportTask{},
)
```

- [ ] **Step 2: 编译验证并提交**

```bash
git add cmd/server/main.go
git commit -m "feat: register phase3 tables in AutoMigrate"
```

---

## Task 1: 报表数据模型

**Files:**
- Create: backend/internal/app/analytics/model/report.go

- [ ] **Step 1: 创建 report model**

```go
// backend/internal/app/analytics/model/report.go
package model

import (
	"time"
)

type ReportDaily struct {
	ID              uint64    `gorm:"primaryKey" json:"id"`
	TenantID        uint64    `gorm:"not null;uniqueIndex:uk_report" json:"tenant_id"`
	AccountID       uint64    `gorm:"not null;uniqueIndex:uk_report;index:idx_account_date" json:"account_id"`
	ProjectID       uint64    `gorm:"default:0;uniqueIndex:uk_report" json:"project_id"`
	PromotionID     uint64    `gorm:"default:0;uniqueIndex:uk_report" json:"promotion_id"`
	StatDate        time.Time `gorm:"type:date;not null;uniqueIndex:uk_report;index:idx_account_date" json:"stat_date"`
	Cost            float64   `gorm:"type:decimal(12,2);default:0" json:"cost"`
	ShowCnt         int64     `gorm:"default:0" json:"show_cnt"`
	ClickCnt        int64     `gorm:"default:0" json:"click_cnt"`
	Ctr             float64   `gorm:"type:decimal(8,4);default:0" json:"ctr"`
	Cpc             float64   `gorm:"type:decimal(8,2);default:0" json:"cpc"`
	Cpm             float64   `gorm:"type:decimal(8,2);default:0" json:"cpm"`
	ConvertCnt      int64     `gorm:"default:0" json:"convert_cnt"`
	ConversionRate  float64   `gorm:"type:decimal(8,4);default:0" json:"conversion_rate"`
	ConversionCost  float64   `gorm:"type:decimal(10,2);default:0" json:"conversion_cost"`
	FormCnt         int64     `gorm:"default:0" json:"form_cnt"`
	PhoneConfirmCnt int64     `gorm:"default:0" json:"phone_confirm_cnt"`
	PhoneConnectCnt int64     `gorm:"default:0" json:"phone_connect_cnt"`
	MessageActionCnt int64    `gorm:"default:0" json:"message_action_cnt"`
	DyLike          int64     `gorm:"default:0" json:"dy_like"`
	DyComment       int64     `gorm:"default:0" json:"dy_comment"`
	DyShare         int64     `gorm:"default:0" json:"dy_share"`
	DyFollow        int64     `gorm:"default:0" json:"dy_follow"`
	TotalPlay       int64     `gorm:"default:0" json:"total_play"`
	PlayOverRate    float64   `gorm:"type:decimal(8,4);default:0" json:"play_over_rate"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (ReportDaily) TableName() string {
	return "report_daily"
}

type ExportTask struct {
	ID        uint64     `gorm:"primaryKey" json:"id"`
	TenantID  uint64     `gorm:"not null;index" json:"tenant_id"`
	Title     string     `gorm:"size:200" json:"title"`
	Status    string     `gorm:"size:20;default:pending" json:"status"`
	FilePath  string     `gorm:"size:500" json:"file_path"`
	RowCount  int        `gorm:"default:0" json:"row_count"`
	ExpireAt  *time.Time `json:"expire_at"`
	CreatedBy uint64     `gorm:"default:0" json:"created_by"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (ExportTask) TableName() string {
	return "export_task"
}
```

- [ ] **Step 2: 编译验证**

Run: cd backend && go build ./internal/app/analytics/...
Expected: 编译成功

- [ ] **Step 3: 提交**

```bash
git add internal/app/analytics/model/
git commit -m "feat: add report_daily and export_task models"
```

---

## Task 2: 报表 Repository（聚合查询）

**Files:**
- Create: backend/internal/app/analytics/repository/report_repo.go

- [ ] **Step 1: 创建 report repository**

```go
// backend/internal/app/analytics/repository/report_repo.go
package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
	"oceanengine-backend/internal/app/analytics/model"
)

type ReportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

type OverviewResult struct {
	TotalCost       float64 `json:"total_cost"`
	TotalShow       int64   `json:"total_show"`
	TotalClick      int64   `json:"total_click"`
	TotalConvert    int64   `json:"total_convert"`
	AvgCtr          float64 `json:"avg_ctr"`
	AvgCpc          float64 `json:"avg_cpc"`
	AvgConvCost     float64 `json:"avg_conversion_cost"`
}

func (r *ReportRepository) GetOverview(ctx context.Context, tenantID uint64, accountIDs []uint64, date time.Time) (*OverviewResult, error) {
	query := r.db.WithContext(ctx).Model(&model.ReportDaily{}).
		Select("COALESCE(SUM(cost),0) as total_cost, COALESCE(SUM(show_cnt),0) as total_show, COALESCE(SUM(click_cnt),0) as total_click, COALESCE(SUM(convert_cnt),0) as total_convert").
		Where("tenant_id = ? AND stat_date = ? AND project_id = 0 AND promotion_id = 0", tenantID, date)
	if len(accountIDs) > 0 {
		query = query.Where("account_id IN ?", accountIDs)
	}
	var result OverviewResult
	err := query.Scan(&result).Error
	return &result, err
}

type TrendPoint struct {
	StatDate   string  `json:"stat_date"`
	Cost       float64 `json:"cost"`
	ShowCnt    int64   `json:"show_cnt"`
	ClickCnt   int64   `json:"click_cnt"`
	ConvertCnt int64   `json:"convert_cnt"`
}

func (r *ReportRepository) GetTrend(ctx context.Context, tenantID uint64, accountIDs []uint64, startDate, endDate time.Time) ([]TrendPoint, error) {
	query := r.db.WithContext(ctx).Model(&model.ReportDaily{}).
		Select("stat_date, SUM(cost) as cost, SUM(show_cnt) as show_cnt, SUM(click_cnt) as click_cnt, SUM(convert_cnt) as convert_cnt").
		Where("tenant_id = ? AND stat_date BETWEEN ? AND ? AND project_id = 0 AND promotion_id = 0", tenantID, startDate, endDate)
	if len(accountIDs) > 0 {
		query = query.Where("account_id IN ?", accountIDs)
	}
	query = query.Group("stat_date").Order("stat_date ASC")
	var points []TrendPoint
	err := query.Scan(&points).Error
	return points, err
}

type RankItem struct {
	AccountID  uint64  `json:"account_id"`
	Name       string  `json:"name"`
	Cost       float64 `json:"cost"`
	ConvertCnt int64   `json:"convert_cnt"`
}

func (r *ReportRepository) GetRank(ctx context.Context, tenantID uint64, accountIDs []uint64, startDate, endDate time.Time, orderBy string, limit int) ([]RankItem, error) {
	query := r.db.WithContext(ctx).Table("report_daily r").
		Select("r.account_id, a.name, SUM(r.cost) as cost, SUM(r.convert_cnt) as convert_cnt").
		Joins("JOIN local_account a ON a.id = r.account_id").
		Where("r.tenant_id = ? AND r.stat_date BETWEEN ? AND ? AND r.project_id = 0 AND r.promotion_id = 0", tenantID, startDate, endDate).
		Group("r.account_id, a.name")
	if len(accountIDs) > 0 {
		query = query.Where("r.account_id IN ?", accountIDs)
	}
	if orderBy == "convert" {
		query = query.Order("convert_cnt DESC")
	} else {
		query = query.Order("cost DESC")
	}
	var items []RankItem
	err := query.Limit(limit).Scan(&items).Error
	return items, err
}

type GroupAggResult struct {
	GroupID    uint64  `json:"group_id"`
	GroupName  string  `json:"group_name"`
	Cost       float64 `json:"cost"`
	ShowCnt    int64   `json:"show_cnt"`
	ClickCnt   int64   `json:"click_cnt"`
	ConvertCnt int64   `json:"convert_cnt"`
}

func (r *ReportRepository) GetByGroup(ctx context.Context, tenantID uint64, accountIDs []uint64, startDate, endDate time.Time) ([]GroupAggResult, error) {
	query := r.db.WithContext(ctx).Table("report_daily r").
		Select("g.id as group_id, g.name as group_name, SUM(r.cost) as cost, SUM(r.show_cnt) as show_cnt, SUM(r.click_cnt) as click_cnt, SUM(r.convert_cnt) as convert_cnt").
		Joins("JOIN account_group_member m ON m.account_id = r.account_id").
		Joins("JOIN account_group g ON g.id = m.group_id").
		Where("r.tenant_id = ? AND r.stat_date BETWEEN ? AND ? AND r.project_id = 0 AND r.promotion_id = 0", tenantID, startDate, endDate).
		Group("g.id, g.name")
	if len(accountIDs) > 0 {
		query = query.Where("r.account_id IN ?", accountIDs)
	}
	var results []GroupAggResult
	err := query.Scan(&results).Error
	return results, err
}

func (r *ReportRepository) GetDetail(ctx context.Context, tenantID uint64, accountIDs []uint64, startDate, endDate time.Time, page, pageSize int) ([]model.ReportDaily, int64, error) {
	query := r.db.WithContext(ctx).Where("tenant_id = ? AND stat_date BETWEEN ? AND ?", tenantID, startDate, endDate)
	if len(accountIDs) > 0 {
		query = query.Where("account_id IN ?", accountIDs)
	}
	var total int64
	query.Model(&model.ReportDaily{}).Count(&total)
	var reports []model.ReportDaily
	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Order("stat_date DESC, account_id ASC").Find(&reports).Error
	return reports, total, err
}

func (r *ReportRepository) Upsert(ctx context.Context, report *model.ReportDaily) error {
	return r.db.WithContext(ctx).Save(report).Error
}
```

- [ ] **Step 2: 编译验证**

Run: cd backend && go build ./internal/app/analytics/...
Expected: 编译成功

- [ ] **Step 3: 提交**

```bash
git add internal/app/analytics/repository/report_repo.go
git commit -m "feat: add report repository with aggregation queries"
```

---

## Task 3: Analytics Service

**Files:**
- Create: backend/internal/app/analytics/service/analytics_service.go

- [ ] **Step 1: 创建 analytics service**

```go
// backend/internal/app/analytics/service/analytics_service.go
package service

import (
	"context"
	"time"

	"oceanengine-backend/internal/app/analytics/repository"
)

type AnalyticsService struct {
	repo *repository.ReportRepository
}

func NewAnalyticsService(repo *repository.ReportRepository) *AnalyticsService {
	return &AnalyticsService{repo: repo}
}

type OverviewRequest struct {
	Date string `form:"date"`
}

func (s *AnalyticsService) GetOverview(ctx context.Context, tenantID uint64, scopeIDs []uint64, req OverviewRequest) (*repository.OverviewResult, error) {
	date := time.Now()
	if req.Date != "" {
		parsed, err := time.Parse("2006-01-02", req.Date)
		if err == nil {
			date = parsed
		}
	}
	return s.repo.GetOverview(ctx, tenantID, scopeIDs, date)
}

type TrendRequest struct {
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
}

func (s *AnalyticsService) GetTrend(ctx context.Context, tenantID uint64, scopeIDs []uint64, req TrendRequest) ([]repository.TrendPoint, error) {
	start, _ := time.Parse("2006-01-02", req.StartDate)
	end, _ := time.Parse("2006-01-02", req.EndDate)
	return s.repo.GetTrend(ctx, tenantID, scopeIDs, start, end)
}

type RankRequest struct {
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
	OrderBy   string `form:"order_by,default=cost"`
	Limit     int    `form:"limit,default=20"`
}

func (s *AnalyticsService) GetRank(ctx context.Context, tenantID uint64, scopeIDs []uint64, req RankRequest) ([]repository.RankItem, error) {
	start, _ := time.Parse("2006-01-02", req.StartDate)
	end, _ := time.Parse("2006-01-02", req.EndDate)
	if req.Limit == 0 {
		req.Limit = 20
	}
	return s.repo.GetRank(ctx, tenantID, scopeIDs, start, end, req.OrderBy, req.Limit)
}

type CompareRequest struct {
	Date      string `form:"date" binding:"required"`
	CompareType string `form:"compare_type,default=yesterday"`
}

func (s *AnalyticsService) GetCompare(ctx context.Context, tenantID uint64, scopeIDs []uint64, req CompareRequest) (map[string]*repository.OverviewResult, error) {
	current, _ := time.Parse("2006-01-02", req.Date)
	var compareDate time.Time
	if req.CompareType == "last_week" {
		compareDate = current.AddDate(0, 0, -7)
	} else {
		compareDate = current.AddDate(0, 0, -1)
	}
	currentData, err := s.repo.GetOverview(ctx, tenantID, scopeIDs, current)
	if err != nil {
		return nil, err
	}
	compareData, err := s.repo.GetOverview(ctx, tenantID, scopeIDs, compareDate)
	if err != nil {
		return nil, err
	}
	return map[string]*repository.OverviewResult{
		"current": currentData,
		"compare": compareData,
	}, nil
}

type DetailRequest struct {
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
	Page      int    `form:"page,default=1"`
	PageSize  int    `form:"page_size,default=20"`
}

func (s *AnalyticsService) GetDetail(ctx context.Context, tenantID uint64, scopeIDs []uint64, req DetailRequest) (interface{}, int64, error) {
	start, _ := time.Parse("2006-01-02", req.StartDate)
	end, _ := time.Parse("2006-01-02", req.EndDate)
	if req.Page == 0 { req.Page = 1 }
	if req.PageSize == 0 { req.PageSize = 20 }
	return s.repo.GetDetail(ctx, tenantID, scopeIDs, start, end, req.Page, req.PageSize)
}
```

- [ ] **Step 2: 编译验证并提交**

```bash
git add internal/app/analytics/service/analytics_service.go
git commit -m "feat: add analytics service with overview/trend/rank/compare"
```

---

## Task 4: Analytics Handler + 路由

**Files:**
- Create: backend/internal/app/analytics/api/analytics_handler.go
- Modify: backend/internal/router/router.go

- [ ] **Step 1: 创建 analytics handler**

```go
// backend/internal/app/analytics/api/analytics_handler.go
package api

import (
	"github.com/gin-gonic/gin"
	"oceanengine-backend/internal/app/analytics/service"
	"oceanengine-backend/pkg/ginutil"
	"oceanengine-backend/pkg/response"
)

type AnalyticsHandler struct {
	svc *service.AnalyticsService
}

func NewAnalyticsHandler(svc *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{svc: svc}
}

func (h *AnalyticsHandler) GetOverview(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	var scopeIDs []uint64
	if v, exists := c.Get("scope_account_ids"); exists {
		scopeIDs, _ = v.([]uint64)
	}
	var req service.OverviewRequest
	c.ShouldBindQuery(&req)
	result, err := h.svc.GetOverview(c.Request.Context(), tenantID, scopeIDs, req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *AnalyticsHandler) GetTrend(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	var scopeIDs []uint64
	if v, exists := c.Get("scope_account_ids"); exists {
		scopeIDs, _ = v.([]uint64)
	}
	var req service.TrendRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	result, err := h.svc.GetTrend(c.Request.Context(), tenantID, scopeIDs, req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *AnalyticsHandler) GetRank(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	var scopeIDs []uint64
	if v, exists := c.Get("scope_account_ids"); exists {
		scopeIDs, _ = v.([]uint64)
	}
	var req service.RankRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	result, err := h.svc.GetRank(c.Request.Context(), tenantID, scopeIDs, req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *AnalyticsHandler) GetCompare(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	var scopeIDs []uint64
	if v, exists := c.Get("scope_account_ids"); exists {
		scopeIDs, _ = v.([]uint64)
	}
	var req service.CompareRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	result, err := h.svc.GetCompare(c.Request.Context(), tenantID, scopeIDs, req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *AnalyticsHandler) GetDetail(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	var scopeIDs []uint64
	if v, exists := c.Get("scope_account_ids"); exists {
		scopeIDs, _ = v.([]uint64)
	}
	var req service.DetailRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	result, total, err := h.svc.GetDetail(c.Request.Context(), tenantID, scopeIDs, req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OKWithList(c, result, total, req.Page, req.PageSize)
}
```

- [ ] **Step 2: 注册路由**

```go
func (r *Router) registerAnalyticsRoutes(rg *gin.RouterGroup) {
	repo := analyticsRepo.NewReportRepository(r.db)
	svc := analyticsService.NewAnalyticsService(repo)
	handler := analyticsApi.NewAnalyticsHandler(svc)

	analytics := rg.Group("/analytics")
	{
		analytics.GET("/overview", handler.GetOverview)
		analytics.GET("/trend", handler.GetTrend)
		analytics.GET("/rank", handler.GetRank)
		analytics.GET("/compare", handler.GetCompare)
		analytics.GET("/detail", handler.GetDetail)
	}
}
```

- [ ] **Step 3: 编译验证并提交**

```bash
git add internal/app/analytics/api/ internal/router/router.go
git commit -m "feat: add analytics API endpoints"
```

---

## Task 5: 报表导出服务

**Files:**
- Create: backend/internal/app/analytics/repository/export_repo.go
- Create: backend/internal/app/analytics/service/export_service.go

- [ ] **Step 1: 创建 export repository**

```go
// backend/internal/app/analytics/repository/export_repo.go
package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
	"oceanengine-backend/internal/app/analytics/model"
)

type ExportRepository struct {
	db *gorm.DB
}

func NewExportRepository(db *gorm.DB) *ExportRepository {
	return &ExportRepository{db: db}
}

func (r *ExportRepository) Create(ctx context.Context, task *model.ExportTask) error {
	return r.db.WithContext(ctx).Create(task).Error
}

func (r *ExportRepository) GetByID(ctx context.Context, tenantID, id uint64) (*model.ExportTask, error) {
	var task model.ExportTask
	err := r.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).First(&task).Error
	return &task, err
}

func (r *ExportRepository) List(ctx context.Context, tenantID uint64) ([]model.ExportTask, error) {
	var tasks []model.ExportTask
	err := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).Order("id DESC").Limit(20).Find(&tasks).Error
	return tasks, err
}

func (r *ExportRepository) UpdateStatus(ctx context.Context, id uint64, status, filePath string, rowCount int) error {
	return r.db.WithContext(ctx).Model(&model.ExportTask{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":    status,
		"file_path": filePath,
		"row_count": rowCount,
	}).Error
}

func (r *ExportRepository) CleanExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("expire_at < ?", time.Now()).Delete(&model.ExportTask{}).Error
}
```

- [ ] **Step 2: 创建 export service**

```go
// backend/internal/app/analytics/service/export_service.go
package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/xuri/excelize/v2"
	"oceanengine-backend/internal/app/analytics/model"
	"oceanengine-backend/internal/app/analytics/repository"
)

type ExportService struct {
	exportRepo *repository.ExportRepository
	reportRepo *repository.ReportRepository
}

func NewExportService(exportRepo *repository.ExportRepository, reportRepo *repository.ReportRepository) *ExportService {
	return &ExportService{exportRepo: exportRepo, reportRepo: reportRepo}
}

type CreateExportRequest struct {
	Title     string `json:"title" binding:"required"`
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
}

func (s *ExportService) CreateExport(ctx context.Context, tenantID, userID uint64, scopeIDs []uint64, req CreateExportRequest) (*model.ExportTask, error) {
	expireAt := time.Now().Add(24 * time.Hour)
	task := &model.ExportTask{
		TenantID:  tenantID,
		Title:     req.Title,
		Status:    "pending",
		ExpireAt:  &expireAt,
		CreatedBy: userID,
	}
	if err := s.exportRepo.Create(ctx, task); err != nil {
		return nil, err
	}
	go s.executeExport(task.ID, tenantID, scopeIDs, req)
	return task, nil
}

func (s *ExportService) executeExport(taskID, tenantID uint64, scopeIDs []uint64, req CreateExportRequest) {
	ctx := context.Background()
	start, _ := time.Parse("2006-01-02", req.StartDate)
	end, _ := time.Parse("2006-01-02", req.EndDate)

	reports, total, err := s.reportRepo.GetDetail(ctx, tenantID, scopeIDs, start, end, 1, 10000)
	if err != nil {
		s.exportRepo.UpdateStatus(ctx, taskID, "failed", "", 0)
		return
	}
	if total > 10000 {
		s.exportRepo.UpdateStatus(ctx, taskID, "failed", "", 0)
		return
	}

	filePath := fmt.Sprintf("/tmp/exports/%d_%d.xlsx", tenantID, taskID)

	// 生成 Excel
	f := excelize.NewFile()
	sheet := "Sheet1"
	headers := []string{"日期", "账户ID", "消耗", "展示", "点击", "转化", "转化成本"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}
	reportList := reports.([]model.ReportDaily)
	for row, r := range reportList {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row+2), r.StatDate.Format("2006-01-02"))
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row+2), r.AccountID)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row+2), r.Cost)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row+2), r.ShowCnt)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row+2), r.ClickCnt)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row+2), r.ConvertCnt)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row+2), r.ConversionCost)
	}
	os.MkdirAll("/tmp/exports", 0755)
	if err := f.SaveAs(filePath); err != nil {
		s.exportRepo.UpdateStatus(ctx, taskID, "failed", "", 0)
		return
	}
	s.exportRepo.UpdateStatus(ctx, taskID, "completed", filePath, int(total))
}

func (s *ExportService) ListExports(ctx context.Context, tenantID uint64) ([]model.ExportTask, error) {
	return s.exportRepo.List(ctx, tenantID)
}

func (s *ExportService) GetExport(ctx context.Context, tenantID, id uint64) (*model.ExportTask, error) {
	return s.exportRepo.GetByID(ctx, tenantID, id)
}
```

- [ ] **Step 3: 在 handler 中添加导出端点**

```go
// 添加到 analytics_handler.go
func (h *AnalyticsHandler) CreateExport(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	userID := ginutil.GetUserID(c)
	var scopeIDs []uint64
	if v, exists := c.Get("scope_account_ids"); exists {
		scopeIDs, _ = v.([]uint64)
	}
	var req service.CreateExportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	task, err := h.exportSvc.CreateExport(c.Request.Context(), tenantID, userID, scopeIDs, req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, task)
}

func (h *AnalyticsHandler) ListExports(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	tasks, err := h.exportSvc.ListExports(c.Request.Context(), tenantID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, tasks)
}

func (h *AnalyticsHandler) DownloadExport(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	task, err := h.exportSvc.GetExport(c.Request.Context(), tenantID, id)
	if err != nil || task.Status != "completed" {
		response.BadRequest(c, "export not ready")
		return
	}
	c.File(task.FilePath)
}
```

路由追加：

```go
analytics.POST("/export", handler.CreateExport)
analytics.GET("/exports", handler.ListExports)
analytics.GET("/exports/:id/download", handler.DownloadExport)
```

- [ ] **Step 4: 编译验证并提交**

```bash
git add internal/app/analytics/
git commit -m "feat: add report export with async Excel generation"
```

---

## Task 6: 报表同步定时任务

**Files:**
- Create: backend/internal/scheduler/sync_reports.go
- Modify: backend/cmd/server/main.go

- [ ] **Step 1: 创建报表同步任务**

```go
// backend/internal/scheduler/sync_reports.go
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
	// 调用 SDK 获取报表数据（此处为伪代码，实际调用 SDK report API）
	// resp, err := sdk.GetProjectReport(ctx, tenant.AccessToken, acc.LocalAccountID, date)
	// 将返回数据写入 report_daily 表

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
	// UPSERT: ON DUPLICATE KEY UPDATE
	s.db.Where("tenant_id = ? AND account_id = ? AND project_id = 0 AND promotion_id = 0 AND stat_date = ?",
		tenant.ID, acc.ID, date).Assign(report).FirstOrCreate(&report)
}
```

- [ ] **Step 2: 在 main.go 启动同步 goroutine**

```go
reportSyncer := scheduler.NewReportSyncer(db, rateLimiter, logger)
go reportSyncer.StartDaily(ctx)
go reportSyncer.StartRealtime(ctx)
```

- [ ] **Step 3: 编译验证并提交**

```bash
git add internal/scheduler/sync_reports.go cmd/server/main.go
git commit -m "feat: add daily and realtime report sync scheduler"
```

---

## 完成标准

Phase 3 完成后：
1. 每天凌晨 2 点自动同步昨日报表（三个粒度）
2. 每 30 分钟同步今日实时消耗
3. 汇总看板显示今日/指定日期的核心指标
4. 趋势图按天展示，支持昨日/上周环比
5. 排行榜按消耗/转化 TOP N
6. 按加盟商分组聚合分析
7. 明细数据分页查询
8. 异步导出 Excel，24h 后自动清理
