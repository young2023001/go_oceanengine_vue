# Phase 2: 投放管理实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 实现本地推项目/单元 CRUD、批量操作引擎（创建/调预算/调出价/启停）、项目模板/单元模板、数据同步定时任务。

**Architecture:** 新增 project/promotion/batch/template 模块。批量操作通过 batch_task + batch_task_item 表驱动，Worker goroutine 逐条执行并通过 Phase 1 的 TenantRateLimiter 限流。数据同步作为主服务内 goroutine 定时运行。

**Tech Stack:** Go 1.21 + Gin + GORM + MySQL + Redis + 巨量引擎 SDK

---

## 文件结构

### 新建文件

| 路径 | 职责 |
|------|------|
| `backend/internal/app/project/model/project.go` | LocalProject + LocalPromotion 模型 |
| `backend/internal/app/project/repository/project_repo.go` | 项目/单元 CRUD |
| `backend/internal/app/project/service/project_service.go` | 项目/单元业务逻辑 |
| `backend/internal/app/project/api/project_handler.go` | 项目/单元 HTTP handler |
| `backend/internal/app/batch/model/batch.go` | BatchTask + BatchTaskItem 模型 |
| `backend/internal/app/batch/repository/batch_repo.go` | 批量任务 CRUD |
| `backend/internal/app/batch/service/batch_service.go` | 批量任务创建/查询 |
| `backend/internal/app/batch/service/worker.go` | Worker 执行引擎 |
| `backend/internal/app/batch/api/batch_handler.go` | 批量操作 HTTP handler |
| `backend/internal/app/template/model/template.go` | ProjectTemplate + PromotionTemplate 模型 |
| `backend/internal/app/template/repository/template_repo.go` | 模板 CRUD |
| `backend/internal/app/template/service/template_service.go` | 模板业务逻辑 |
| `backend/internal/app/template/api/template_handler.go` | 模板 HTTP handler |
| `backend/internal/scheduler/sync_projects.go` | 项目/单元状态同步定时任务 |
| `backend/pkg/ocean/sdk.go` | SDK 封装层（统一调用入口） |

### 修改文件

| 路径 | 改动 |
|------|------|
| `backend/internal/router/router.go` | 注册 project/batch/template 路由 |
| `backend/cmd/server/main.go` | 启动 Worker + 同步定时任务 goroutine |

---

## Task 1: 项目/单元数据模型

**Files:**
- Create: `backend/internal/app/project/model/project.go`

- [ ] **Step 1: 创建 project model 文件**

```go
// backend/internal/app/project/model/project.go
package model

import (
	"time"

	"gorm.io/gorm"
)

type LocalProject struct {
	ID             uint64         `gorm:"primaryKey" json:"id"`
	TenantID       uint64         `gorm:"not null;uniqueIndex:uk_tenant_project" json:"tenant_id"`
	AccountID      uint64         `gorm:"not null;index" json:"account_id"`
	ProjectID      uint64         `gorm:"not null;uniqueIndex:uk_tenant_project" json:"project_id"`
	Name           string         `gorm:"size:200" json:"name"`
	MarketingGoal  string         `gorm:"size:50" json:"marketing_goal"`
	DeliveryScene  string         `gorm:"size:50" json:"delivery_scene"`
	AdType         string         `gorm:"size:50" json:"ad_type"`
	ExternalAction string         `gorm:"size:50" json:"external_action"`
	BidType        string         `gorm:"size:20" json:"bid_type"`
	Bid            int            `gorm:"default:0" json:"bid"`
	BudgetMode     string         `gorm:"size:30" json:"budget_mode"`
	Budget         int            `gorm:"default:0" json:"budget"`
	StatusFirst    string         `gorm:"size:50;index" json:"status_first"`
	StatusSecond   string         `gorm:"size:50" json:"status_second"`
	StartTime      string         `gorm:"size:20" json:"start_time"`
	EndTime        string         `gorm:"size:20" json:"end_time"`
	PoiID          uint64         `gorm:"default:0" json:"poi_id"`
	PoiName        string         `gorm:"size:200" json:"poi_name"`
	ProductID      uint64         `gorm:"default:0" json:"product_id"`
	SyncedAt       *time.Time     `json:"synced_at"`
	LocalUpdatedAt *time.Time     `json:"local_updated_at"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (LocalProject) TableName() string {
	return "local_project"
}

type LocalPromotion struct {
	ID           uint64         `gorm:"primaryKey" json:"id"`
	TenantID     uint64         `gorm:"not null;uniqueIndex:uk_tenant_promotion" json:"tenant_id"`
	AccountID    uint64         `gorm:"not null;index" json:"account_id"`
	ProjectID    uint64         `gorm:"not null;index" json:"project_id"`
	PromotionID  uint64         `gorm:"not null;uniqueIndex:uk_tenant_promotion" json:"promotion_id"`
	Name         string         `gorm:"size:200" json:"name"`
	StatusFirst  string         `gorm:"size:50" json:"status_first"`
	StatusSecond string         `gorm:"size:50" json:"status_second"`
	AwemeID      string         `gorm:"size:100" json:"aweme_id"`
	SyncedAt     *time.Time     `json:"synced_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (LocalPromotion) TableName() string {
	return "local_promotion"
}
```

- [ ] **Step 2: 编译验证**

Run: `cd /Users/leoyang/Desktop/LeoValue/qikuai/go_oceanengine_vue/backend && go build ./internal/app/project/...`
Expected: 编译成功

- [ ] **Step 3: 提交**

```bash
git add internal/app/project/model/project.go
git commit -m "feat: add local_project and local_promotion models"
```

---

## Task 2: 批量任务数据模型

**Files:**
- Create: `backend/internal/app/batch/model/batch.go`

- [ ] **Step 1: 创建 batch model 文件**

```go
// backend/internal/app/batch/model/batch.go
package model

import (
	"time"

	"gorm.io/gorm"
)

type TaskStatus string

const (
	TaskStatusPending        TaskStatus = "pending"
	TaskStatusRunning        TaskStatus = "running"
	TaskStatusCompleted      TaskStatus = "completed"
	TaskStatusPartialSuccess TaskStatus = "partial_success"
	TaskStatusFailed         TaskStatus = "failed"
	TaskStatusCancelled      TaskStatus = "cancelled"
)

type ItemStatus string

const (
	ItemStatusPending   ItemStatus = "pending"
	ItemStatusRunning   ItemStatus = "running"
	ItemStatusSuccess   ItemStatus = "success"
	ItemStatusFailed    ItemStatus = "failed"
	ItemStatusSkipped   ItemStatus = "skipped"
	ItemStatusCancelled ItemStatus = "cancelled"
)

type BatchTask struct {
	ID           uint64         `gorm:"primaryKey" json:"id"`
	TenantID     uint64         `gorm:"not null;index:idx_tenant_status" json:"tenant_id"`
	TaskType     string         `gorm:"size:50;not null" json:"task_type"`
	Title        string         `gorm:"size:200;not null" json:"title"`
	Config       string         `gorm:"type:json;not null" json:"config"`
	Status       TaskStatus     `gorm:"size:30;default:pending;index:idx_tenant_status" json:"status"`
	TotalCount   int            `gorm:"default:0" json:"total_count"`
	SuccessCount int            `gorm:"default:0" json:"success_count"`
	FailedCount  int            `gorm:"default:0" json:"failed_count"`
	CreatedBy    uint64         `gorm:"default:0" json:"created_by"`
	StartedAt    *time.Time     `json:"started_at"`
	CompletedAt  *time.Time     `json:"completed_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (BatchTask) TableName() string {
	return "batch_task"
}

type BatchTaskItem struct {
	ID             uint64     `gorm:"primaryKey" json:"id"`
	TenantID       uint64     `gorm:"not null;index" json:"tenant_id"`
	TaskID         uint64     `gorm:"not null;index" json:"task_id"`
	AccountID      uint64     `gorm:"not null" json:"account_id"`
	LocalAccountID uint64     `gorm:"not null" json:"local_account_id"`
	ConfigOverride string     `gorm:"type:json" json:"config_override"`
	Status         ItemStatus `gorm:"size:20;default:pending;index" json:"status"`
	ResultID       uint64     `gorm:"default:0" json:"result_id"`
	ErrorMsg       string     `gorm:"size:500" json:"error_msg"`
	ExecutedAt     *time.Time `json:"executed_at"`
}

func (BatchTaskItem) TableName() string {
	return "batch_task_item"
}
```

- [ ] **Step 2: 编译验证**

Run: `cd /Users/leoyang/Desktop/LeoValue/qikuai/go_oceanengine_vue/backend && go build ./internal/app/batch/...`
Expected: 编译成功

- [ ] **Step 3: 提交**

```bash
git add internal/app/batch/model/batch.go
git commit -m "feat: add batch_task and batch_task_item models"
```

---

## Task 3: 模板数据模型

**Files:**
- Create: `backend/internal/app/template/model/template.go`

- [ ] **Step 1: 创建 template model 文件**

```go
// backend/internal/app/template/model/template.go
package model

import (
	"time"

	"gorm.io/gorm"
)

type ProjectTemplate struct {
	ID        uint64         `gorm:"primaryKey" json:"id"`
	TenantID  uint64         `gorm:"not null;index" json:"tenant_id"`
	Name      string         `gorm:"size:200;not null" json:"name"`
	Config    string         `gorm:"type:json;not null" json:"config"`
	UseCount  int            `gorm:"default:0" json:"use_count"`
	CreatedBy uint64         `gorm:"default:0" json:"created_by"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ProjectTemplate) TableName() string {
	return "project_template"
}

type PromotionTemplate struct {
	ID        uint64         `gorm:"primaryKey" json:"id"`
	TenantID  uint64         `gorm:"not null;index" json:"tenant_id"`
	Name      string         `gorm:"size:200;not null" json:"name"`
	Config    string         `gorm:"type:json;not null" json:"config"`
	UseCount  int            `gorm:"default:0" json:"use_count"`
	CreatedBy uint64         `gorm:"default:0" json:"created_by"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (PromotionTemplate) TableName() string {
	return "promotion_template"
}
```

- [ ] **Step 2: 编译验证**

Run: `cd /Users/leoyang/Desktop/LeoValue/qikuai/go_oceanengine_vue/backend && go build ./internal/app/template/...`
Expected: 编译成功

- [ ] **Step 3: 提交**

```bash
git add internal/app/template/model/template.go
git commit -m "feat: add project_template and promotion_template models"
```

---

## Task 4: AutoMigrate 注册新表

**Files:**
- Modify: `backend/cmd/server/main.go`

- [ ] **Step 1: 添加 import 和 AutoMigrate**

```go
import (
	projectModel "oceanengine-backend/internal/app/project/model"
	batchModel "oceanengine-backend/internal/app/batch/model"
	templateModel "oceanengine-backend/internal/app/template/model"
)

// 在现有 AutoMigrate 调用中追加
db.AutoMigrate(
	&projectModel.LocalProject{},
	&projectModel.LocalPromotion{},
	&batchModel.BatchTask{},
	&batchModel.BatchTaskItem{},
	&templateModel.ProjectTemplate{},
	&templateModel.PromotionTemplate{},
)
```

- [ ] **Step 2: 编译验证**

Run: `cd /Users/leoyang/Desktop/LeoValue/qikuai/go_oceanengine_vue/backend && go build ./cmd/server/...`
Expected: 编译成功

- [ ] **Step 3: 提交**

```bash
git add cmd/server/main.go
git commit -m "feat: register phase2 tables in AutoMigrate"
```

---

## Task 5: SDK 封装层

**Files:**
- Create: `backend/pkg/ocean/sdk.go`

- [ ] **Step 1: 创建 SDK 统一封装**

```go
// backend/pkg/ocean/sdk.go
package ocean

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{httpClient: &http.Client{}}
}

type CreateProjectRequest struct {
	LocalAccountID uint64 `json:"advertiser_id"`
	Name           string `json:"name"`
	MarketingGoal  string `json:"marketing_goal"`
	BidType        string `json:"bid_type"`
	Bid            int    `json:"bid"`
	BudgetMode     string `json:"budget_mode"`
	Budget         int    `json:"budget"`
	DeliveryScene  string `json:"delivery_scene"`
	ExternalAction string `json:"external_action"`
}

type CreateProjectResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		ProjectID uint64 `json:"project_id"`
	} `json:"data"`
}

func (c *Client) CreateProject(ctx context.Context, accessToken string, req CreateProjectRequest) (*CreateProjectResponse, error) {
	return c.doPost(ctx, accessToken, "https://open.oceanengine.com/open_api/v1.0/local/project/create/", req)
}

type UpdateBudgetRequest struct {
	LocalAccountID uint64 `json:"advertiser_id"`
	ProjectID      uint64 `json:"project_id"`
	Budget         int    `json:"budget"`
}

func (c *Client) UpdateProjectBudget(ctx context.Context, accessToken string, req UpdateBudgetRequest) error {
	_, err := c.doPost(ctx, accessToken, "https://open.oceanengine.com/open_api/v1.0/local/project/update/", req)
	return err
}

type UpdateBidRequest struct {
	LocalAccountID uint64 `json:"advertiser_id"`
	ProjectID      uint64 `json:"project_id"`
	Bid            int    `json:"bid"`
}

func (c *Client) UpdateProjectBid(ctx context.Context, accessToken string, req UpdateBidRequest) error {
	_, err := c.doPost(ctx, accessToken, "https://open.oceanengine.com/open_api/v1.0/local/project/update/", req)
	return err
}

type UpdateStatusRequest struct {
	LocalAccountID uint64 `json:"advertiser_id"`
	ProjectID      uint64 `json:"project_id"`
	OptStatus      string `json:"opt_status"`
}

func (c *Client) UpdateProjectStatus(ctx context.Context, accessToken string, req UpdateStatusRequest) error {
	_, err := c.doPost(ctx, accessToken, "https://open.oceanengine.com/open_api/v1.0/local/project/update/status/", req)
	return err
}

type ListProjectsRequest struct {
	LocalAccountID uint64 `json:"advertiser_id"`
	Page           int    `json:"page"`
	PageSize       int    `json:"page_size"`
}

type ProjectInfo struct {
	ProjectID    uint64 `json:"project_id"`
	Name         string `json:"name"`
	StatusFirst  string `json:"first_status"`
	StatusSecond string `json:"second_status"`
}

type ListProjectsResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		List     []ProjectInfo `json:"list"`
		PageInfo struct {
			Total int `json:"total"`
		} `json:"page_info"`
	} `json:"data"`
}

func (c *Client) ListProjects(ctx context.Context, accessToken string, req ListProjectsRequest) (*ListProjectsResponse, error) {
	url := fmt.Sprintf("https://open.oceanengine.com/open_api/v1.0/local/project/list/?advertiser_id=%d&page=%d&page_size=%d",
		req.LocalAccountID, req.Page, req.PageSize)
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Access-Token", accessToken)
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result ListProjectsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) doPost(ctx context.Context, accessToken, url string, body interface{}) (*CreateProjectResponse, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Access-Token", accessToken)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result CreateProjectResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if result.Code != 0 {
		return nil, fmt.Errorf("ocean API error: %s", result.Message)
	}
	return &result, nil
}
```

- [ ] **Step 2: 编译验证**

Run: `cd /Users/leoyang/Desktop/LeoValue/qikuai/go_oceanengine_vue/backend && go build ./pkg/ocean/...`
Expected: 编译成功

- [ ] **Step 3: 提交**

```bash
git add pkg/ocean/sdk.go
git commit -m "feat: add ocean SDK wrapper for local push API"
```

---

## Task 6: Project Repository + Service + Handler

**Files:**
- Create: `backend/internal/app/project/repository/project_repo.go`
- Create: `backend/internal/app/project/service/project_service.go`
- Create: `backend/internal/app/project/api/project_handler.go`

- [ ] **Step 1: 创建 project repository**

```go
// backend/internal/app/project/repository/project_repo.go
package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
	"oceanengine-backend/internal/app/project/model"
)

type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) ListProjects(ctx context.Context, tenantID uint64, accountIDs []uint64, status string, page, pageSize int) ([]model.LocalProject, int64, error) {
	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if len(accountIDs) > 0 {
		query = query.Where("account_id IN ?", accountIDs)
	}
	if status != "" {
		query = query.Where("status_first = ?", status)
	}
	var total int64
	query.Model(&model.LocalProject{}).Count(&total)
	var projects []model.LocalProject
	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Order("id DESC").Find(&projects).Error
	return projects, total, err
}

func (r *ProjectRepository) GetByID(ctx context.Context, tenantID, id uint64) (*model.LocalProject, error) {
	var p model.LocalProject
	err := r.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProjectRepository) Upsert(ctx context.Context, project *model.LocalProject) error {
	return r.db.WithContext(ctx).Save(project).Error
}

func (r *ProjectRepository) UpdateStatus(ctx context.Context, tenantID, id uint64, status string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&model.LocalProject{}).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Updates(map[string]interface{}{"status_first": status, "local_updated_at": now}).Error
}

func (r *ProjectRepository) ListPromotions(ctx context.Context, tenantID uint64, projectID uint64, accountIDs []uint64, page, pageSize int) ([]model.LocalPromotion, int64, error) {
	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if projectID > 0 {
		query = query.Where("project_id = ?", projectID)
	}
	if len(accountIDs) > 0 {
		query = query.Where("account_id IN ?", accountIDs)
	}
	var total int64
	query.Model(&model.LocalPromotion{}).Count(&total)
	var promotions []model.LocalPromotion
	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Order("id DESC").Find(&promotions).Error
	return promotions, total, err
}

func (r *ProjectRepository) UpsertPromotion(ctx context.Context, promo *model.LocalPromotion) error {
	return r.db.WithContext(ctx).Save(promo).Error
}
```

- [ ] **Step 2: 创建 project service**

```go
// backend/internal/app/project/service/project_service.go
package service

import (
	"context"

	"oceanengine-backend/internal/app/project/model"
	"oceanengine-backend/internal/app/project/repository"
)

type ProjectService struct {
	repo *repository.ProjectRepository
}

func NewProjectService(repo *repository.ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

type ListProjectsRequest struct {
	Status   string `form:"status"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=20"`
}

func (s *ProjectService) ListProjects(ctx context.Context, tenantID uint64, scopeIDs []uint64, req ListProjectsRequest) ([]model.LocalProject, int64, error) {
	return s.repo.ListProjects(ctx, tenantID, scopeIDs, req.Status, req.Page, req.PageSize)
}

func (s *ProjectService) GetProject(ctx context.Context, tenantID, id uint64) (*model.LocalProject, error) {
	return s.repo.GetByID(ctx, tenantID, id)
}

func (s *ProjectService) UpdateStatus(ctx context.Context, tenantID, id uint64, status string) error {
	return s.repo.UpdateStatus(ctx, tenantID, id, status)
}

type ListPromotionsRequest struct {
	ProjectID uint64 `form:"project_id"`
	Page      int    `form:"page,default=1"`
	PageSize  int    `form:"page_size,default=20"`
}

func (s *ProjectService) ListPromotions(ctx context.Context, tenantID uint64, scopeIDs []uint64, req ListPromotionsRequest) ([]model.LocalPromotion, int64, error) {
	return s.repo.ListPromotions(ctx, tenantID, req.ProjectID, scopeIDs, req.Page, req.PageSize)
}
```

- [ ] **Step 3: 创建 project handler**

```go
// backend/internal/app/project/api/project_handler.go
package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"oceanengine-backend/internal/app/project/service"
	"oceanengine-backend/pkg/ginutil"
	"oceanengine-backend/pkg/response"
)

type ProjectHandler struct {
	svc *service.ProjectService
}

func NewProjectHandler(svc *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{svc: svc}
}

func (h *ProjectHandler) ListProjects(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	var scopeIDs []uint64
	if v, exists := c.Get("scope_account_ids"); exists {
		scopeIDs, _ = v.([]uint64)
	}
	var req service.ListProjectsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	projects, total, err := h.svc.ListProjects(c.Request.Context(), tenantID, scopeIDs, req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OKWithList(c, projects, total, req.Page, req.PageSize)
}

func (h *ProjectHandler) GetProject(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	project, err := h.svc.GetProject(c.Request.Context(), tenantID, id)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, project)
}

type UpdateStatusReq struct {
	Status string `json:"status" binding:"required"`
}

func (h *ProjectHandler) UpdateProjectStatus(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	var req UpdateStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.UpdateStatus(c.Request.Context(), tenantID, id, req.Status); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *ProjectHandler) ListPromotions(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	var scopeIDs []uint64
	if v, exists := c.Get("scope_account_ids"); exists {
		scopeIDs, _ = v.([]uint64)
	}
	var req service.ListPromotionsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	promotions, total, err := h.svc.ListPromotions(c.Request.Context(), tenantID, scopeIDs, req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OKWithList(c, promotions, total, req.Page, req.PageSize)
}
```

- [ ] **Step 4: 注册路由**

```go
func (r *Router) registerProjectRoutes(rg *gin.RouterGroup) {
	repo := projectRepo.NewProjectRepository(r.db)
	svc := projectService.NewProjectService(repo)
	handler := projectApi.NewProjectHandler(svc)

	projects := rg.Group("/projects")
	{
		projects.GET("", handler.ListProjects)
		projects.GET("/:id", handler.GetProject)
		projects.PUT("/:id/status", handler.UpdateProjectStatus)
	}

	promotions := rg.Group("/promotions")
	{
		promotions.GET("", handler.ListPromotions)
	}
}
```

- [ ] **Step 5: 编译验证**

Run: `cd /Users/leoyang/Desktop/LeoValue/qikuai/go_oceanengine_vue/backend && go build ./...`
Expected: 编译成功

- [ ] **Step 6: 提交**

```bash
git add internal/app/project/ internal/router/router.go
git commit -m "feat: add project and promotion list/status API"
```

---

## Task 7: Batch Worker 引擎

**Files:**
- Create: `backend/internal/app/batch/repository/batch_repo.go`
- Create: `backend/internal/app/batch/service/batch_service.go`
- Create: `backend/internal/app/batch/service/worker.go`

- [ ] **Step 1: 创建 batch repository**

```go
// backend/internal/app/batch/repository/batch_repo.go
package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
	"oceanengine-backend/internal/app/batch/model"
)

type BatchRepository struct {
	db *gorm.DB
}

func NewBatchRepository(db *gorm.DB) *BatchRepository {
	return &BatchRepository{db: db}
}

func (r *BatchRepository) CreateTask(ctx context.Context, task *model.BatchTask, items []model.BatchTaskItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(task).Error; err != nil {
			return err
		}
		for i := range items {
			items[i].TaskID = task.ID
			items[i].TenantID = task.TenantID
		}
		return tx.CreateInBatches(items, 100).Error
	})
}

func (r *BatchRepository) GetTask(ctx context.Context, tenantID, taskID uint64) (*model.BatchTask, error) {
	var task model.BatchTask
	err := r.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", tenantID, taskID).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *BatchRepository) ListTasks(ctx context.Context, tenantID uint64, page, pageSize int) ([]model.BatchTask, int64, error) {
	var total int64
	r.db.WithContext(ctx).Model(&model.BatchTask{}).Where("tenant_id = ?", tenantID).Count(&total)
	var tasks []model.BatchTask
	err := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).
		Order("id DESC").Offset((page-1)*pageSize).Limit(pageSize).Find(&tasks).Error
	return tasks, total, err
}

func (r *BatchRepository) ListRunningTasks(ctx context.Context) ([]model.BatchTask, error) {
	var tasks []model.BatchTask
	err := r.db.WithContext(ctx).Where("status = ?", model.TaskStatusRunning).Find(&tasks).Error
	return tasks, err
}

func (r *BatchRepository) GetTaskByID(ctx context.Context, taskID uint64) (*model.BatchTask, error) {
	var task model.BatchTask
	err := r.db.WithContext(ctx).First(&task, taskID).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *BatchRepository) GetPendingItems(ctx context.Context, taskID uint64, limit int) ([]model.BatchTaskItem, error) {
	var items []model.BatchTaskItem
	err := r.db.WithContext(ctx).Where("task_id = ? AND status = ?", taskID, model.ItemStatusPending).
		Limit(limit).Find(&items).Error
	return items, err
}

func (r *BatchRepository) UpdateItemStatus(ctx context.Context, itemID uint64, status model.ItemStatus, resultID uint64, errMsg string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&model.BatchTaskItem{}).Where("id = ?", itemID).
		Updates(map[string]interface{}{
			"status":      status,
			"result_id":   resultID,
			"error_msg":   errMsg,
			"executed_at": now,
		}).Error
}

func (r *BatchRepository) UpdateTaskStatus(ctx context.Context, taskID uint64, status model.TaskStatus) error {
	updates := map[string]interface{}{"status": status}
	if status == model.TaskStatusRunning {
		now := time.Now()
		updates["started_at"] = now
	}
	if status == model.TaskStatusCompleted || status == model.TaskStatusPartialSuccess || status == model.TaskStatusFailed {
		now := time.Now()
		updates["completed_at"] = now
	}
	return r.db.WithContext(ctx).Model(&model.BatchTask{}).Where("id = ?", taskID).Updates(updates).Error
}

func (r *BatchRepository) IncrSuccess(ctx context.Context, taskID uint64) error {
	return r.db.WithContext(ctx).Model(&model.BatchTask{}).Where("id = ?", taskID).
		UpdateColumn("success_count", gorm.Expr("success_count + 1")).Error
}

func (r *BatchRepository) IncrFailed(ctx context.Context, taskID uint64) error {
	return r.db.WithContext(ctx).Model(&model.BatchTask{}).Where("id = ?", taskID).
		UpdateColumn("failed_count", gorm.Expr("failed_count + 1")).Error
}

func (r *BatchRepository) CancelPendingItems(ctx context.Context, taskID uint64) (int64, error) {
	result := r.db.WithContext(ctx).Model(&model.BatchTaskItem{}).
		Where("task_id = ? AND status = ?", taskID, model.ItemStatusPending).
		Update("status", model.ItemStatusCancelled)
	return result.RowsAffected, result.Error
}

func (r *BatchRepository) RetryFailedItems(ctx context.Context, taskID uint64) (int64, error) {
	result := r.db.WithContext(ctx).Model(&model.BatchTaskItem{}).
		Where("task_id = ? AND status = ?", taskID, model.ItemStatusFailed).
		Updates(map[string]interface{}{"status": model.ItemStatusPending, "error_msg": ""})
	return result.RowsAffected, result.Error
}
```

- [ ] **Step 2: 创建 Worker 引擎**

```go
// backend/internal/app/batch/service/worker.go
package service

import (
	"context"
	"encoding/json"
	"time"

	"go.uber.org/zap"
	"oceanengine-backend/internal/app/batch/model"
	"oceanengine-backend/internal/app/batch/repository"
	"oceanengine-backend/pkg/ocean"
	"oceanengine-backend/pkg/ratelimiter"
)

type TokenStore interface {
	GetToken(ctx context.Context, tenantID uint64) (string, error)
}

type Worker struct {
	repo       *repository.BatchRepository
	sdk        *ocean.Client
	limiter    *ratelimiter.TenantRateLimiter
	tokenStore TokenStore
	logger     *zap.Logger
}

func NewWorker(repo *repository.BatchRepository, sdk *ocean.Client, limiter *ratelimiter.TenantRateLimiter, tokenStore TokenStore, logger *zap.Logger) *Worker {
	return &Worker{repo: repo, sdk: sdk, limiter: limiter, tokenStore: tokenStore, logger: logger}
}

func (w *Worker) Start(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.processNext(ctx)
		}
	}
}

func (w *Worker) processNext(ctx context.Context) {
	tasks, err := w.repo.ListRunningTasks(ctx)
	if err != nil {
		w.logger.Error("list running tasks", zap.Error(err))
		return
	}
	for _, task := range tasks {
		w.processTask(ctx, task)
	}
}

func (w *Worker) processTask(ctx context.Context, task model.BatchTask) {
	items, err := w.repo.GetPendingItems(ctx, task.ID, 10)
	if err != nil {
		w.logger.Error("get pending items", zap.Error(err))
		return
	}
	if len(items) == 0 {
		w.finalizeTask(ctx, task.ID)
		return
	}
	for _, item := range items {
		if err := w.limiter.Wait(ctx, task.TenantID); err != nil {
			return
		}
		w.executeItem(ctx, task, item)
	}
}

func (w *Worker) executeItem(ctx context.Context, task model.BatchTask, item model.BatchTaskItem) {
	w.repo.UpdateItemStatus(ctx, item.ID, model.ItemStatusRunning, 0, "")

	var errMsg string
	var resultID uint64

	switch task.TaskType {
	case "create_project":
		resultID, errMsg = w.execCreateProject(ctx, task, item)
	case "update_budget":
		errMsg = w.execUpdateBudget(ctx, task, item)
	case "update_bid":
		errMsg = w.execUpdateBid(ctx, task, item)
	case "update_status":
		errMsg = w.execUpdateStatus(ctx, task, item)
	}

	if errMsg != "" {
		w.repo.UpdateItemStatus(ctx, item.ID, model.ItemStatusFailed, 0, errMsg)
		w.repo.IncrFailed(ctx, task.ID)
	} else {
		w.repo.UpdateItemStatus(ctx, item.ID, model.ItemStatusSuccess, resultID, "")
		w.repo.IncrSuccess(ctx, task.ID)
	}
}

func (w *Worker) execCreateProject(ctx context.Context, task model.BatchTask, item model.BatchTaskItem) (uint64, string) {
	var req ocean.CreateProjectRequest
	json.Unmarshal([]byte(task.Config), &req)
	if item.ConfigOverride != "" {
		json.Unmarshal([]byte(item.ConfigOverride), &req)
	}
	req.LocalAccountID = item.LocalAccountID
	accessToken, _ := w.tokenStore.GetToken(ctx, item.TenantID)
	resp, err := w.sdk.CreateProject(ctx, accessToken, req)
	if err != nil {
		return 0, err.Error()
	}
	return resp.Data.ProjectID, ""
}

func (w *Worker) execUpdateBudget(ctx context.Context, task model.BatchTask, item model.BatchTaskItem) string {
	accessToken, _ := w.tokenStore.GetToken(ctx, task.TenantID)
	var req ocean.UpdateBudgetRequest
	json.Unmarshal([]byte(task.Config), &req)
	req.LocalAccountID = item.LocalAccountID
	if err := w.sdk.UpdateProjectBudget(ctx, accessToken, req); err != nil {
		return err.Error()
	}
	return ""
}

func (w *Worker) execUpdateBid(ctx context.Context, task model.BatchTask, item model.BatchTaskItem) string {
	accessToken, _ := w.tokenStore.GetToken(ctx, task.TenantID)
	var req ocean.UpdateBidRequest
	json.Unmarshal([]byte(task.Config), &req)
	req.LocalAccountID = item.LocalAccountID
	if err := w.sdk.UpdateProjectBid(ctx, accessToken, req); err != nil {
		return err.Error()
	}
	return ""
}

func (w *Worker) execUpdateStatus(ctx context.Context, task model.BatchTask, item model.BatchTaskItem) string {
	accessToken, _ := w.tokenStore.GetToken(ctx, task.TenantID)
	var req ocean.UpdateStatusRequest
	json.Unmarshal([]byte(task.Config), &req)
	req.LocalAccountID = item.LocalAccountID
	if err := w.sdk.UpdateProjectStatus(ctx, accessToken, req); err != nil {
		return err.Error()
	}
	return ""
}

func (w *Worker) finalizeTask(ctx context.Context, taskID uint64) {
	task, _ := w.repo.GetTaskByID(ctx, taskID)
	if task == nil {
		return
	}
	if task.FailedCount == 0 {
		w.repo.UpdateTaskStatus(ctx, taskID, model.TaskStatusCompleted)
	} else if task.SuccessCount == 0 {
		w.repo.UpdateTaskStatus(ctx, taskID, model.TaskStatusFailed)
	} else {
		w.repo.UpdateTaskStatus(ctx, taskID, model.TaskStatusPartialSuccess)
	}
}
```

- [ ] **Step 3: 编译验证**

Run: `cd /Users/leoyang/Desktop/LeoValue/qikuai/go_oceanengine_vue/backend && go build ./internal/app/batch/...`
Expected: 编译成功

- [ ] **Step 4: 提交**

```bash
git add internal/app/batch/repository/ internal/app/batch/service/
git commit -m "feat: add batch task repository and worker engine"
```

---

## Task 8: Batch Handler + 路由

**Files:**
- Create: `backend/internal/app/batch/service/batch_service.go`
- Create: `backend/internal/app/batch/api/batch_handler.go`
- Modify: `backend/internal/router/router.go`

- [ ] **Step 1: 创建 batch service**

```go
// backend/internal/app/batch/service/batch_service.go
package service

import (
	"context"
	"encoding/json"

	"oceanengine-backend/internal/app/batch/model"
	"oceanengine-backend/internal/app/batch/repository"
)

type BatchService struct {
	repo *repository.BatchRepository
}

func NewBatchService(repo *repository.BatchRepository) *BatchService {
	return &BatchService{repo: repo}
}

type CreateBatchRequest struct {
	TaskType   string                   `json:"task_type" binding:"required"`
	Title      string                   `json:"title" binding:"required"`
	Config     map[string]interface{}   `json:"config" binding:"required"`
	AccountIDs []uint64                 `json:"account_ids" binding:"required,min=1"`
	Overrides  map[uint64]interface{}   `json:"overrides"`
}

func (s *BatchService) CreateTask(ctx context.Context, tenantID, userID uint64, localAccountMap map[uint64]uint64, req CreateBatchRequest) (*model.BatchTask, error) {
	configJSON, _ := json.Marshal(req.Config)
	task := &model.BatchTask{
		TenantID:   tenantID,
		TaskType:   req.TaskType,
		Title:      req.Title,
		Config:     string(configJSON),
		Status:     model.TaskStatusRunning,
		TotalCount: len(req.AccountIDs),
		CreatedBy:  userID,
	}
	items := make([]model.BatchTaskItem, 0, len(req.AccountIDs))
	for _, accID := range req.AccountIDs {
		item := model.BatchTaskItem{
			AccountID:      accID,
			LocalAccountID: localAccountMap[accID],
			Status:         model.ItemStatusPending,
		}
		if override, ok := req.Overrides[accID]; ok {
			overrideJSON, _ := json.Marshal(override)
			item.ConfigOverride = string(overrideJSON)
		}
		items = append(items, item)
	}
	if err := s.repo.CreateTask(ctx, task, items); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *BatchService) GetTask(ctx context.Context, tenantID, taskID uint64) (*model.BatchTask, error) {
	return s.repo.GetTask(ctx, tenantID, taskID)
}

func (s *BatchService) ListTasks(ctx context.Context, tenantID uint64, page, pageSize int) ([]model.BatchTask, int64, error) {
	return s.repo.ListTasks(ctx, tenantID, page, pageSize)
}

func (s *BatchService) CancelTask(ctx context.Context, tenantID, taskID uint64) error {
	task, err := s.repo.GetTask(ctx, tenantID, taskID)
	if err != nil {
		return err
	}
	if task.Status != model.TaskStatusRunning && task.Status != model.TaskStatusPending {
		return nil
	}
	s.repo.CancelPendingItems(ctx, taskID)
	return s.repo.UpdateTaskStatus(ctx, taskID, model.TaskStatusCancelled)
}

func (s *BatchService) RetryTask(ctx context.Context, tenantID, taskID uint64) error {
	task, err := s.repo.GetTask(ctx, tenantID, taskID)
	if err != nil {
		return err
	}
	s.repo.RetryFailedItems(ctx, taskID)
	task.FailedCount = 0
	return s.repo.UpdateTaskStatus(ctx, taskID, model.TaskStatusRunning)
}
```

- [ ] **Step 2: 创建 batch handler**

```go
// backend/internal/app/batch/api/batch_handler.go
package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"oceanengine-backend/internal/app/batch/service"
	"oceanengine-backend/pkg/ginutil"
	"oceanengine-backend/pkg/response"
)

type BatchHandler struct {
	svc        *service.BatchService
	accountDB  *gorm.DB
}

func NewBatchHandler(svc *service.BatchService, db *gorm.DB) *BatchHandler {
	return &BatchHandler{svc: svc, accountDB: db}
}

func (h *BatchHandler) CreateTask(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	userID := ginutil.GetUserID(c)
	var req service.CreateBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 查询 account_id -> local_account_id 映射
	type accountMapping struct {
		ID             uint64
		LocalAccountID uint64
	}
	var mappings []accountMapping
	h.accountDB.Table("local_account").
		Select("id, local_account_id").
		Where("tenant_id = ? AND id IN ?", tenantID, req.AccountIDs).
		Scan(&mappings)

	localAccountMap := make(map[uint64]uint64, len(mappings))
	for _, m := range mappings {
		localAccountMap[m.ID] = m.LocalAccountID
	}

	task, err := h.svc.CreateTask(c.Request.Context(), tenantID, userID, localAccountMap, req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, task)
}

func (h *BatchHandler) GetTask(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	task, err := h.svc.GetTask(c.Request.Context(), tenantID, id)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, task)
}

func (h *BatchHandler) ListTasks(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	tasks, total, err := h.svc.ListTasks(c.Request.Context(), tenantID, page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OKWithList(c, tasks, total, page, pageSize)
}

func (h *BatchHandler) CancelTask(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.CancelTask(c.Request.Context(), tenantID, id); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *BatchHandler) RetryTask(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.RetryTask(c.Request.Context(), tenantID, id); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, nil)
}
```

- [ ] **Step 3: 注册路由**

```go
func (r *Router) registerBatchRoutes(rg *gin.RouterGroup) {
	repo := batchRepo.NewBatchRepository(r.db)
	svc := batchService.NewBatchService(repo)
	handler := batchApi.NewBatchHandler(svc, r.db)

	batch := rg.Group("/batch")
	{
		batch.POST("/projects", handler.CreateTask)
		batch.POST("/promotions", handler.CreateTask)
		batch.POST("/budget", handler.CreateTask)
		batch.POST("/bid", handler.CreateTask)
		batch.POST("/status", handler.CreateTask)
		batch.GET("/tasks", handler.ListTasks)
		batch.GET("/tasks/:id", handler.GetTask)
		batch.POST("/tasks/:id/cancel", handler.CancelTask)
		batch.POST("/tasks/:id/retry", handler.RetryTask)
	}
}
```

- [ ] **Step 4: 编译验证**

Run: `cd /Users/leoyang/Desktop/LeoValue/qikuai/go_oceanengine_vue/backend && go build ./...`
Expected: 编译成功

- [ ] **Step 5: 提交**

```bash
git add internal/app/batch/ internal/router/router.go
git commit -m "feat: add batch task API with cancel and retry"
```

---

## Task 9: Template CRUD

**Files:**
- Create: `backend/internal/app/template/repository/template_repo.go`
- Create: `backend/internal/app/template/service/template_service.go`
- Create: `backend/internal/app/template/api/template_handler.go`

- [ ] **Step 1: 创建 template repository**

```go
// backend/internal/app/template/repository/template_repo.go
package repository

import (
	"context"

	"gorm.io/gorm"
	"oceanengine-backend/internal/app/template/model"
)

type TemplateRepository struct {
	db *gorm.DB
}

func NewTemplateRepository(db *gorm.DB) *TemplateRepository {
	return &TemplateRepository{db: db}
}

func (r *TemplateRepository) CreateProject(ctx context.Context, t *model.ProjectTemplate) error {
	return r.db.WithContext(ctx).Create(t).Error
}

func (r *TemplateRepository) ListProjects(ctx context.Context, tenantID uint64) ([]model.ProjectTemplate, error) {
	var templates []model.ProjectTemplate
	err := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).Order("use_count DESC").Find(&templates).Error
	return templates, err
}

func (r *TemplateRepository) GetProject(ctx context.Context, tenantID, id uint64) (*model.ProjectTemplate, error) {
	var t model.ProjectTemplate
	err := r.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).First(&t).Error
	return &t, err
}

func (r *TemplateRepository) UpdateProject(ctx context.Context, t *model.ProjectTemplate) error {
	return r.db.WithContext(ctx).Save(t).Error
}

func (r *TemplateRepository) DeleteProject(ctx context.Context, tenantID, id uint64) error {
	return r.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).Delete(&model.ProjectTemplate{}).Error
}

func (r *TemplateRepository) IncrUseCount(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Model(&model.ProjectTemplate{}).Where("id = ?", id).
		UpdateColumn("use_count", gorm.Expr("use_count + 1")).Error
}

func (r *TemplateRepository) CreatePromotion(ctx context.Context, t *model.PromotionTemplate) error {
	return r.db.WithContext(ctx).Create(t).Error
}

func (r *TemplateRepository) ListPromotions(ctx context.Context, tenantID uint64) ([]model.PromotionTemplate, error) {
	var templates []model.PromotionTemplate
	err := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).Order("use_count DESC").Find(&templates).Error
	return templates, err
}

func (r *TemplateRepository) DeletePromotion(ctx context.Context, tenantID, id uint64) error {
	return r.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).Delete(&model.PromotionTemplate{}).Error
}
```

- [ ] **Step 2: 创建 template service**

```go
// backend/internal/app/template/service/template_service.go
package service

import (
	"context"
	"encoding/json"

	"oceanengine-backend/internal/app/template/model"
	"oceanengine-backend/internal/app/template/repository"
)

type TemplateService struct {
	repo *repository.TemplateRepository
}

func NewTemplateService(repo *repository.TemplateRepository) *TemplateService {
	return &TemplateService{repo: repo}
}

type CreateTemplateRequest struct {
	Name   string                 `json:"name" binding:"required"`
	Config map[string]interface{} `json:"config" binding:"required"`
}

func (s *TemplateService) CreateProject(ctx context.Context, tenantID, userID uint64, req CreateTemplateRequest) (*model.ProjectTemplate, error) {
	configJSON, _ := json.Marshal(req.Config)
	t := &model.ProjectTemplate{
		TenantID:  tenantID,
		Name:      req.Name,
		Config:    string(configJSON),
		CreatedBy: userID,
	}
	if err := s.repo.CreateProject(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *TemplateService) ListProjects(ctx context.Context, tenantID uint64) ([]model.ProjectTemplate, error) {
	return s.repo.ListProjects(ctx, tenantID)
}

func (s *TemplateService) GetProject(ctx context.Context, tenantID, id uint64) (*model.ProjectTemplate, error) {
	return s.repo.GetProject(ctx, tenantID, id)
}

func (s *TemplateService) DeleteProject(ctx context.Context, tenantID, id uint64) error {
	return s.repo.DeleteProject(ctx, tenantID, id)
}

func (s *TemplateService) CreatePromotion(ctx context.Context, tenantID, userID uint64, req CreateTemplateRequest) (*model.PromotionTemplate, error) {
	configJSON, _ := json.Marshal(req.Config)
	t := &model.PromotionTemplate{
		TenantID:  tenantID,
		Name:      req.Name,
		Config:    string(configJSON),
		CreatedBy: userID,
	}
	if err := s.repo.CreatePromotion(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *TemplateService) ListPromotions(ctx context.Context, tenantID uint64) ([]model.PromotionTemplate, error) {
	return s.repo.ListPromotions(ctx, tenantID)
}

func (s *TemplateService) DeletePromotion(ctx context.Context, tenantID, id uint64) error {
	return s.repo.DeletePromotion(ctx, tenantID, id)
}
```

- [ ] **Step 3: 创建 template handler**

```go
// backend/internal/app/template/api/template_handler.go
package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"oceanengine-backend/internal/app/template/service"
	"oceanengine-backend/pkg/ginutil"
	"oceanengine-backend/pkg/response"
)

type TemplateHandler struct {
	svc *service.TemplateService
}

func NewTemplateHandler(svc *service.TemplateService) *TemplateHandler {
	return &TemplateHandler{svc: svc}
}

func (h *TemplateHandler) CreateProject(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	userID := ginutil.GetUserID(c)
	var req service.CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	t, err := h.svc.CreateProject(c.Request.Context(), tenantID, userID, req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, t)
}

func (h *TemplateHandler) ListProjects(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	list, err := h.svc.ListProjects(c.Request.Context(), tenantID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, list)
}

func (h *TemplateHandler) GetProject(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	t, err := h.svc.GetProject(c.Request.Context(), tenantID, id)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, t)
}

func (h *TemplateHandler) UpdateProject(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req service.CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	t, err := h.svc.GetProject(c.Request.Context(), tenantID, id)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, t)
}

func (h *TemplateHandler) DeleteProject(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.DeleteProject(c.Request.Context(), tenantID, id); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *TemplateHandler) CreatePromotion(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	userID := ginutil.GetUserID(c)
	var req service.CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	t, err := h.svc.CreatePromotion(c.Request.Context(), tenantID, userID, req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, t)
}

func (h *TemplateHandler) ListPromotions(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	list, err := h.svc.ListPromotions(c.Request.Context(), tenantID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, list)
}

func (h *TemplateHandler) DeletePromotion(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.DeletePromotion(c.Request.Context(), tenantID, id); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, nil)
}
```

- [ ] **Step 4: 注册路由**

handler 注册路由：

```go
func (r *Router) registerTemplateRoutes(rg *gin.RouterGroup) {
	repo := templateRepo.NewTemplateRepository(r.db)
	svc := templateService.NewTemplateService(repo)
	handler := templateApi.NewTemplateHandler(svc)

	tpl := rg.Group("/templates")
	{
		projects := tpl.Group("/projects")
		{
			projects.POST("", handler.CreateProject)
			projects.GET("", handler.ListProjects)
			projects.GET("/:id", handler.GetProject)
			projects.PUT("/:id", handler.UpdateProject)
			projects.DELETE("/:id", handler.DeleteProject)
		}
		promotions := tpl.Group("/promotions")
		{
			promotions.POST("", handler.CreatePromotion)
			promotions.GET("", handler.ListPromotions)
			promotions.DELETE("/:id", handler.DeletePromotion)
		}
	}
}
```

- [ ] **Step 3: 编译验证并提交**

```bash
git add internal/app/template/ internal/router/router.go
git commit -m "feat: add project and promotion template CRUD"
```

---

## Task 10: 数据同步定时任务

**Files:**
- Create: `backend/internal/scheduler/sync_projects.go`
- Modify: `backend/cmd/server/main.go`

- [ ] **Step 1: 创建同步任务**

```go
// backend/internal/scheduler/sync_projects.go
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
	s.db.Where("status = 1 AND token_status = 'active'").Find(&tenants)

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
		s.logger.Error("sync projects", zap.Uint64("account_id", acc.ID), zap.Error(err))
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
```

- [ ] **Step 2: 在 main.go 启动同步 goroutine**

```go
syncer := scheduler.NewProjectSyncer(db, oceanClient, rateLimiter, logger)
go syncer.Start(ctx)
```

- [ ] **Step 3: 编译验证并提交**

```bash
git add internal/scheduler/ cmd/server/main.go
git commit -m "feat: add hourly project sync scheduler"
```

---

## 完成标准

Phase 2 完成后：
1. 本地项目/单元列表可查询（从同步数据）
2. 批量创建项目、调预算、调出价、启停 — 通过 batch_task 异步执行
3. 批量任务支持取消和重试
4. 前端可轮询任务进度
5. 项目/单元模板 CRUD 可用
6. 每小时自动同步项目状态，5 分钟内有用户操作的记录不被覆盖
