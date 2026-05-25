# Phase 1: 账户体系实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 实现多租户账户体系，包括租户管理、OAuth Token 管理（含回调授权流程）、本地推账户导入、分组管理、门店管理、投手权限隔离、API 限流器。

**Architecture:** 在现有 Go/Gin 单体架构上新增 tenant/account/group/scope 模块，复用现有 JWT 认证和 RBAC 体系。sys_user 表增加 tenant_id 实现多租户隔离。权限通过 user_account_scope 表实现投手级账户隔离，data_scope 中间件统一注入 WHERE 条件。

**Tech Stack:** Go 1.21 + Gin + GORM + MySQL + Redis (令牌桶限流)

---

## 文件结构

### 新建文件

| 路径 | 职责 |
|------|------|
| `backend/internal/app/tenant/model/tenant.go` | Tenant 模型定义 |
| `backend/internal/app/tenant/repository/tenant_repo.go` | Tenant CRUD 数据访问 |
| `backend/internal/app/tenant/service/tenant_service.go` | 租户业务逻辑 + OAuth Token 管理 |
| `backend/internal/app/tenant/service/oauth_client.go` | OAuth 客户端（调用巨量 API 换 Token） |
| `backend/internal/app/tenant/service/token_refresher.go` | Token 自动续期 goroutine |
| `backend/internal/app/tenant/api/tenant_handler.go` | 租户 HTTP handler（含 OAuth 回调） |
| `backend/internal/app/account/model/account.go` | LocalAccount + Store 模型 |
| `backend/internal/app/account/repository/account_repo.go` | 账户 CRUD |
| `backend/internal/app/account/service/account_service.go` | 账户导入逻辑 |
| `backend/internal/app/account/api/account_handler.go` | 账户 HTTP handler |
| `backend/internal/app/group/model/group.go` | AccountGroup + AccountGroupMember 模型 |
| `backend/internal/app/group/repository/group_repo.go` | 分组 CRUD |
| `backend/internal/app/group/service/group_service.go` | 分组业务逻辑 |
| `backend/internal/app/group/api/group_handler.go` | 分组 HTTP handler |
| `backend/internal/app/scope/model/scope.go` | UserAccountScope 模型 |
| `backend/internal/app/scope/repository/scope_repo.go` | 权限范围 CRUD |
| `backend/internal/app/scope/service/scope_service.go` | 权限分配逻辑 |
| `backend/internal/app/scope/api/scope_handler.go` | 权限 HTTP handler |
| `backend/internal/middleware/data_scope.go` | 数据权限中间件 |
| `backend/internal/middleware/data_scope_test.go` | 数据权限中间件测试 |
| `backend/pkg/ratelimiter/limiter.go` | Redis 令牌桶限流器 |
| `backend/pkg/ratelimiter/limiter_test.go` | 限流器单元测试 |
| `backend/test/integration/tenant_test.go` | 租户集成测试 |
| `backend/test/integration/account_test.go` | 账户集成测试 |
| `backend/test/integration/group_test.go` | 分组集成测试 |

### 修改文件

| 路径 | 改动 |
|------|------|
| `backend/internal/app/admin/model/user.go` | User 结构体增加 TenantID 字段 |
| `backend/internal/router/router.go` | 注册新模块路由 + import aliases |
| `backend/internal/middleware/jwt.go` | JWT claims 增加 tenant_id，context 注入 |
| `backend/pkg/auth/jwt.go` | Claims 结构体增加 TenantID |
| `backend/internal/app/admin/service/auth.go` | Login 方法传入 TenantID |
| `backend/cmd/server/main.go` | AutoMigrate + 启动 Token 续期 goroutine |

---

## 通用约定

**获取 context 中的 tenant_id/user_id：** Gin 的 `c.Get()` 返回 `interface{}`，需要类型断言。计划中所有 handler 统一使用以下辅助函数：

```go
// backend/pkg/ginutil/context.go
package ginutil

import "github.com/gin-gonic/gin"

func GetTenantID(c *gin.Context) uint64 {
	v, _ := c.Get("tenant_id")
	id, _ := v.(uint64)
	return id
}

func GetUserID(c *gin.Context) uint64 {
	v, _ := c.Get("user_id")
	id, _ := v.(uint64)
	return id
}

func GetRoleKey(c *gin.Context) string {
	v, _ := c.Get("role_key")
	s, _ := v.(string)
	return s
}
```

所有 handler 中 `tenantID := ginutil.GetTenantID(c)` 替代 `c.GetUint64("tenant_id")`。

---

## Task 1: 通用辅助 + 数据库模型（Tenant / LocalAccount / Store）

**Files:**
- Create: `backend/pkg/ginutil/context.go`
- Create: `backend/internal/app/tenant/model/tenant.go`
- Create: `backend/internal/app/account/model/account.go`

- [ ] **Step 1: 创建 ginutil 辅助包**

```go
// backend/pkg/ginutil/context.go
package ginutil

import "github.com/gin-gonic/gin"

func GetTenantID(c *gin.Context) uint64 {
	v, _ := c.Get("tenant_id")
	id, _ := v.(uint64)
	return id
}

func GetUserID(c *gin.Context) uint64 {
	v, _ := c.Get("user_id")
	id, _ := v.(uint64)
	return id
}

func GetRoleKey(c *gin.Context) string {
	v, _ := c.Get("role_key")
	s, _ := v.(string)
	return s
}
```

- [ ] **Step 2: 创建 tenant model**

```go
// backend/internal/app/tenant/model/tenant.go
package model

import (
	"time"

	"gorm.io/gorm"
)

type TokenStatus string

const (
	TokenStatusActive   TokenStatus = "active"
	TokenStatusExpiring TokenStatus = "expiring"
	TokenStatusExpired  TokenStatus = "expired"
	TokenStatusFailed   TokenStatus = "failed"
)

type Tenant struct {
	ID             uint64         `gorm:"primaryKey" json:"id"`
	Name           string         `gorm:"size:200;not null" json:"name"`
	OAuthAppID     string         `gorm:"column:oauth_app_id;size:100;not null" json:"oauth_app_id"`
	OAuthSecret    string         `gorm:"column:oauth_secret;size:200;not null" json:"-"`
	AccessToken    string         `gorm:"size:500" json:"-"`
	RefreshToken   string         `gorm:"size:500" json:"-"`
	TokenExpireAt  *time.Time     `json:"-"`
	TokenStatus    TokenStatus    `gorm:"size:20;default:active" json:"token_status"`
	OrganizationID uint64         `gorm:"default:0" json:"organization_id"`
	Status         int8           `gorm:"default:1;index" json:"status"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Tenant) TableName() string {
	return "tenant"
}
```

- [ ] **Step 3: 创建 account model（含唯一索引）**

```go
// backend/internal/app/account/model/account.go
package model

import (
	"time"

	"gorm.io/gorm"
)

type LocalAccount struct {
	ID             uint64         `gorm:"primaryKey" json:"id"`
	TenantID       uint64         `gorm:"not null;uniqueIndex:uk_tenant_account" json:"tenant_id"`
	LocalAccountID uint64         `gorm:"not null;uniqueIndex:uk_tenant_account" json:"local_account_id"`
	Name           string         `gorm:"size:200;not null" json:"name"`
	Status         string         `gorm:"size:50;default:active" json:"status"`
	Balance        float64        `gorm:"type:decimal(12,2);default:0" json:"balance"`
	LastSyncAt     *time.Time     `json:"last_sync_at"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (LocalAccount) TableName() string {
	return "local_account"
}

type Store struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	TenantID  uint64    `gorm:"not null;index" json:"tenant_id"`
	AccountID uint64    `gorm:"not null;index" json:"account_id"`
	PoiID     uint64    `gorm:"default:0" json:"poi_id"`
	Name      string    `gorm:"size:200;not null" json:"name"`
	City      string    `gorm:"size:50" json:"city"`
	District  string    `gorm:"size:50" json:"district"`
	Address   string    `gorm:"size:500" json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Store) TableName() string {
	return "store"
}
```

- [ ] **Step 4: 编译验证**

Run: `cd /Users/leoyang/Desktop/LeoValue/qikuai/go_oceanengine_vue/backend && go build ./pkg/ginutil/... ./internal/app/tenant/... ./internal/app/account/...`
Expected: 编译成功

- [ ] **Step 5: 提交**

```bash
git add pkg/ginutil/context.go internal/app/tenant/model/tenant.go internal/app/account/model/account.go
git commit -m "feat: add ginutil helpers, tenant and account model definitions"
```

---

## Task 2: 数据库模型（Group / Scope）

**Files:**
- Create: `backend/internal/app/group/model/group.go`
- Create: `backend/internal/app/scope/model/scope.go`

- [ ] **Step 1: 创建 group model（含复合索引）**

```go
// backend/internal/app/group/model/group.go
package model

import (
	"time"

	"gorm.io/gorm"
)

type GroupType string

const (
	GroupTypeFranchisee GroupType = "franchisee"
	GroupTypeRegion     GroupType = "region"
	GroupTypeCustom     GroupType = "custom"
)

type AccountGroup struct {
	ID        uint64         `gorm:"primaryKey" json:"id"`
	TenantID  uint64         `gorm:"not null;index:idx_tenant_type" json:"tenant_id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Type      GroupType      `gorm:"size:20;not null;default:custom;index:idx_tenant_type" json:"type"`
	ParentID  uint64         `gorm:"default:0" json:"parent_id"`
	Sort      int            `gorm:"default:0" json:"sort"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (AccountGroup) TableName() string {
	return "account_group"
}

type AccountGroupMember struct {
	AccountID uint64 `gorm:"primaryKey" json:"account_id"`
	GroupID   uint64 `gorm:"primaryKey;index" json:"group_id"`
}

func (AccountGroupMember) TableName() string {
	return "account_group_member"
}
```

- [ ] **Step 2: 创建 scope model**

```go
// backend/internal/app/scope/model/scope.go
package model

type UserAccountScope struct {
	UserID    uint64 `gorm:"primaryKey" json:"user_id"`
	AccountID uint64 `gorm:"primaryKey;index" json:"account_id"`
}

func (UserAccountScope) TableName() string {
	return "user_account_scope"
}
```

- [ ] **Step 3: 编译验证**

Run: `cd /Users/leoyang/Desktop/LeoValue/qikuai/go_oceanengine_vue/backend && go build ./internal/app/group/... ./internal/app/scope/...`
Expected: 编译成功

- [ ] **Step 4: 提交**

```bash
git add internal/app/group/model/group.go internal/app/scope/model/scope.go
git commit -m "feat: add group and user_account_scope model definitions"
```

---

## Task 3: sys_user 增加 tenant_id + JWT Claims 改造

**Files:**
- Modify: `backend/internal/app/admin/model/user.go`
- Modify: `backend/pkg/auth/jwt.go`
- Modify: `backend/internal/middleware/jwt.go`
- Modify: `backend/internal/app/admin/service/auth.go`

- [ ] **Step 1: 在 User 结构体中增加 TenantID 字段**

在 `backend/internal/app/admin/model/user.go` 的 User 结构体 `ID` 字段之后添加：

```go
TenantID uint64 `gorm:"default:0;index" json:"tenant_id"`
```

- [ ] **Step 2: 修改 pkg/auth/jwt.go 中的 Claims 结构体**

查看现有 Claims 结构体，增加 TenantID 字段：

```go
type Claims struct {
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
	RoleKey  string `json:"role_key"`
	TenantID uint64 `json:"tenant_id"`
	jwt.RegisteredClaims
}
```

同时修改 `GenerateToken` 方法，从 User 对象读取 TenantID 写入 claims。

- [ ] **Step 3: 修改 JWT 中间件提取 tenant_id 到 context**

在 `backend/internal/middleware/jwt.go` 的认证成功后添加：

```go
c.Set("tenant_id", claims.TenantID)
c.Set("user_id", claims.UserID)
c.Set("role_key", claims.RoleKey)
```

- [ ] **Step 4: 修改 AuthService.Login 方法**

在 `backend/internal/app/admin/service/auth.go` 中，登录时将 user.TenantID 传入 token 生成：

```go
token, err := s.jwtManager.GenerateToken(user.ID, user.Username, role.Key, user.TenantID)
```

- [ ] **Step 5: 编译验证**

Run: `cd /Users/leoyang/Desktop/LeoValue/qikuai/go_oceanengine_vue/backend && go build ./...`
Expected: 编译成功

- [ ] **Step 6: 提交**

```bash
git add internal/app/admin/model/user.go pkg/auth/jwt.go internal/middleware/jwt.go internal/app/admin/service/auth.go
git commit -m "feat: add tenant_id to user model and JWT claims"
```

---

## Task 4: GORM AutoMigrate 注册新表

**Files:**
- Modify: `backend/cmd/server/main.go`

- [ ] **Step 1: 在 main.go 顶部添加 import**

```go
import (
	// ... 现有 imports ...
	accountModel "oceanengine-backend/internal/app/account/model"
	groupModel "oceanengine-backend/internal/app/group/model"
	scopeModel "oceanengine-backend/internal/app/scope/model"
	tenantModel "oceanengine-backend/internal/app/tenant/model"
)
```

- [ ] **Step 2: 在数据库初始化之后添加 AutoMigrate**

```go
if err := db.AutoMigrate(
	&tenantModel.Tenant{},
	&accountModel.LocalAccount{},
	&accountModel.Store{},
	&groupModel.AccountGroup{},
	&groupModel.AccountGroupMember{},
	&scopeModel.UserAccountScope{},
); err != nil {
	logger.Fatal("auto migrate failed", zap.Error(err))
}
```

- [ ] **Step 3: 编译验证**

Run: `cd /Users/leoyang/Desktop/LeoValue/qikuai/go_oceanengine_vue/backend && go build ./cmd/server/...`
Expected: 编译成功

- [ ] **Step 4: 提交**

```bash
git add cmd/server/main.go
git commit -m "feat: register new tables in AutoMigrate"
```

---

## Task 5: Redis 令牌桶限流器

**Files:**
- Create: `backend/pkg/ratelimiter/limiter.go`
- Create: `backend/pkg/ratelimiter/limiter_test.go`

- [ ] **Step 1: 编写限流器测试**

```go
// backend/pkg/ratelimiter/limiter_test.go
package ratelimiter

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func setupTestRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   15,
	})
}

func TestTenantRateLimiter_Allow(t *testing.T) {
	rdb := setupTestRedis()
	ctx := context.Background()
	rdb.FlushDB(ctx)
	defer rdb.FlushDB(ctx)

	limiter := NewTenantRateLimiter(rdb, 10)

	for i := 0; i < 10; i++ {
		allowed, err := limiter.Allow(ctx, 1)
		assert.NoError(t, err)
		assert.True(t, allowed, "request %d should be allowed", i+1)
	}

	allowed, err := limiter.Allow(ctx, 1)
	assert.NoError(t, err)
	assert.False(t, allowed, "request 11 should be denied")

	time.Sleep(1100 * time.Millisecond)
	allowed, err = limiter.Allow(ctx, 1)
	assert.NoError(t, err)
	assert.True(t, allowed, "request after 1s should be allowed")
}

func TestTenantRateLimiter_IsolatesTenants(t *testing.T) {
	rdb := setupTestRedis()
	ctx := context.Background()
	rdb.FlushDB(ctx)
	defer rdb.FlushDB(ctx)

	limiter := NewTenantRateLimiter(rdb, 2)

	limiter.Allow(ctx, 1)
	limiter.Allow(ctx, 1)
	allowed, _ := limiter.Allow(ctx, 1)
	assert.False(t, allowed, "tenant 1 should be denied")

	allowed, _ = limiter.Allow(ctx, 2)
	assert.True(t, allowed, "tenant 2 should not be affected")
}
```

- [ ] **Step 2: 运行测试确认失败（RED）**

Run: `cd /Users/leoyang/Desktop/LeoValue/qikuai/go_oceanengine_vue/backend && go test ./pkg/ratelimiter/ -v -run TestTenantRateLimiter`
Expected: 编译失败 — `NewTenantRateLimiter` 未定义

- [ ] **Step 3: 实现限流器**

```go
// backend/pkg/ratelimiter/limiter.go
package ratelimiter

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var luaScript = redis.NewScript(`
local key = KEYS[1]
local limit = tonumber(ARGV[1])
local current = tonumber(redis.call('GET', key) or '0')
if current < limit then
    redis.call('INCR', key)
    redis.call('EXPIRE', key, 2)
    return 0
end
return 1
`)

type TenantRateLimiter struct {
	rdb   *redis.Client
	limit int
}

func NewTenantRateLimiter(rdb *redis.Client, qpsPerTenant int) *TenantRateLimiter {
	return &TenantRateLimiter{
		rdb:   rdb,
		limit: qpsPerTenant,
	}
}

func (l *TenantRateLimiter) Allow(ctx context.Context, tenantID uint64) (bool, error) {
	now := time.Now().Unix()
	key := fmt.Sprintf("ratelimit:%d:ocean:%d", tenantID, now)

	result, err := luaScript.Run(ctx, l.rdb, []string{key}, l.limit).Int()
	if err != nil {
		return false, err
	}
	return result == 0, nil
}

func (l *TenantRateLimiter) Wait(ctx context.Context, tenantID uint64) error {
	for {
		allowed, err := l.Allow(ctx, tenantID)
		if err != nil {
			return err
		}
		if allowed {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(100 * time.Millisecond):
		}
	}
}
```

- [ ] **Step 4: 运行测试确认通过（GREEN）**

Run: `cd /Users/leoyang/Desktop/LeoValue/qikuai/go_oceanengine_vue/backend && go test ./pkg/ratelimiter/ -v -run TestTenantRateLimiter`
Expected: PASS

- [ ] **Step 5: 提交**

```bash
git add pkg/ratelimiter/
git commit -m "feat: add Redis-based tenant rate limiter with Lua script"
```

---

## Task 6: Tenant Repository + Service

**Files:**
- Create: `backend/internal/app/tenant/repository/tenant_repo.go`
- Create: `backend/internal/app/tenant/service/tenant_service.go`

- [ ] **Step 1: 创建 tenant repository**

```go
// backend/internal/app/tenant/repository/tenant_repo.go
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
```

- [ ] **Step 2: 创建 tenant service**

```go
// backend/internal/app/tenant/service/tenant_service.go
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
```

- [ ] **Step 3: 编译验证**

Run: `cd /Users/leoyang/Desktop/LeoValue/qikuai/go_oceanengine_vue/backend && go build ./internal/app/tenant/...`
Expected: 编译成功

- [ ] **Step 4: 提交**

```bash
git add internal/app/tenant/repository/tenant_repo.go internal/app/tenant/service/tenant_service.go
git commit -m "feat: add tenant repository and service layer"
```

---

## Task 7: OAuth 客户端 + Tenant Handler（含 OAuth 回调）

**Files:**
- Create: `backend/internal/app/tenant/service/oauth_client.go`
- Create: `backend/internal/app/tenant/api/tenant_handler.go`
- Modify: `backend/internal/router/router.go`

- [ ] **Step 1: 创建 OAuth 客户端接口和实现**

```go
// backend/internal/app/tenant/service/oauth_client.go
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type OAuthClient interface {
	GetAuthURL(appID, redirectURI, state string) string
	ExchangeToken(ctx context.Context, appID, secret, authCode string) (accessToken, refreshToken string, expiresIn int64, err error)
	RefreshToken(ctx context.Context, appID, secret, refreshToken string) (newAccessToken, newRefreshToken string, expiresIn int64, err error)
}

type oceanOAuthClient struct {
	httpClient *http.Client
}

func NewOAuthClient() OAuthClient {
	return &oceanOAuthClient{httpClient: &http.Client{}}
}

func (c *oceanOAuthClient) GetAuthURL(appID, redirectURI, state string) string {
	return fmt.Sprintf(
		"https://open.oceanengine.com/audit/oauth.html?app_id=%s&redirect_uri=%s&state=%s&scope=",
		appID, url.QueryEscape(redirectURI), state,
	)
}

type oauthTokenResponse struct {
	Data struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int64  `json:"expires_in"`
	} `json:"data"`
	Message string `json:"message"`
}

func (c *oceanOAuthClient) ExchangeToken(ctx context.Context, appID, secret, authCode string) (string, string, int64, error) {
	reqURL := "https://open.oceanengine.com/open_api/oauth2/access_token/"
	body := fmt.Sprintf(`{"app_id":"%s","secret":"%s","grant_type":"auth_code","auth_code":"%s"}`, appID, secret, authCode)

	req, err := http.NewRequestWithContext(ctx, "POST", reqURL, strings.NewReader(body))
	if err != nil {
		return "", "", 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", "", 0, err
	}
	defer resp.Body.Close()

	var result oauthTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", 0, err
	}
	if result.Data.AccessToken == "" {
		return "", "", 0, fmt.Errorf("oauth exchange failed: %s", result.Message)
	}
	return result.Data.AccessToken, result.Data.RefreshToken, result.Data.ExpiresIn, nil
}

func (c *oceanOAuthClient) RefreshToken(ctx context.Context, appID, secret, refreshToken string) (string, string, int64, error) {
	reqURL := "https://open.oceanengine.com/open_api/oauth2/refresh_token/"
	body := fmt.Sprintf(`{"app_id":"%s","secret":"%s","grant_type":"refresh_token","refresh_token":"%s"}`, appID, secret, refreshToken)

	req, err := http.NewRequestWithContext(ctx, "POST", reqURL, strings.NewReader(body))
	if err != nil {
		return "", "", 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", "", 0, err
	}
	defer resp.Body.Close()

	var result oauthTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", 0, err
	}
	if result.Data.AccessToken == "" {
		return "", "", 0, fmt.Errorf("oauth refresh failed: %s", result.Message)
	}
	return result.Data.AccessToken, result.Data.RefreshToken, result.Data.ExpiresIn, nil
}
```

- [ ] **Step 2: 创建 tenant handler（含 OAuth 授权 URL + 回调）**

```go
// backend/internal/app/tenant/api/tenant_handler.go
package api

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"oceanengine-backend/internal/app/tenant/service"
	"oceanengine-backend/pkg/response"
)

type TenantHandler struct {
	svc   *service.TenantService
	oauth service.OAuthClient
}

func NewTenantHandler(svc *service.TenantService, oauth service.OAuthClient) *TenantHandler {
	return &TenantHandler{svc: svc, oauth: oauth}
}

func (h *TenantHandler) Create(c *gin.Context) {
	var req service.CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	tenant, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, tenant)
}

func (h *TenantHandler) List(c *gin.Context) {
	tenants, err := h.svc.List(c.Request.Context())
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, tenants)
}

func (h *TenantHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	tenant, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, tenant)
}

func (h *TenantHandler) GetOAuthURL(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	tenant, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	redirectURI := fmt.Sprintf("%s/api/v1/tenants/oauth/callback", c.Request.Host)
	state := fmt.Sprintf("%d", tenant.ID)
	authURL := h.oauth.GetAuthURL(tenant.OAuthAppID, redirectURI, state)
	response.OK(c, gin.H{"auth_url": authURL})
}

func (h *TenantHandler) OAuthCallback(c *gin.Context) {
	authCode := c.Query("auth_code")
	state := c.Query("state")
	if authCode == "" || state == "" {
		response.BadRequest(c, "missing auth_code or state")
		return
	}

	tenantID, err := strconv.ParseUint(state, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid state")
		return
	}

	tenant, err := h.svc.GetByID(c.Request.Context(), tenantID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	accessToken, refreshToken, expiresIn, err := h.oauth.ExchangeToken(
		c.Request.Context(), tenant.OAuthAppID, tenant.OAuthSecret, authCode,
	)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	if err := h.svc.SaveOAuthToken(c.Request.Context(), tenantID, accessToken, refreshToken, expiresIn); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "授权成功"})
}
```

注意：`GetOAuthURL` 中不需要从 context 获取 tenant_id（通过 URL param 获取），所以 import 中不需要 `ginutil`。

- [ ] **Step 3: 在 router.go 中注册租户路由**

在 router.go 顶部添加 import：

```go
import (
	// ... 现有 imports ...
	tenantApi "oceanengine-backend/internal/app/tenant/api"
	tenantRepo "oceanengine-backend/internal/app/tenant/repository"
	tenantService "oceanengine-backend/internal/app/tenant/service"
)
```

在 `registerProtectedRoutes` 中添加：

```go
r.registerTenantRoutes(rg)
```

新增方法：

```go
func (r *Router) registerTenantRoutes(rg *gin.RouterGroup) {
	repo := tenantRepo.NewTenantRepository(r.db)
	svc := tenantService.NewTenantService(repo)
	oauth := tenantService.NewOAuthClient()
	handler := tenantApi.NewTenantHandler(svc, oauth)

	tenants := rg.Group("/tenants")
	{
		tenants.POST("", handler.Create)
		tenants.GET("", handler.List)
		tenants.GET("/:id", handler.GetByID)
		tenants.GET("/:id/oauth/url", handler.GetOAuthURL)
	}
}
```

在 `registerPublicRoutes` 中添加 OAuth 回调（公开路由）：

```go
// 租户 OAuth 回调（公开）
tenantOAuthRepo := tenantRepo.NewTenantRepository(r.db)
tenantOAuthSvc := tenantService.NewTenantService(tenantOAuthRepo)
tenantOAuthClient := tenantService.NewOAuthClient()
tenantOAuthHandler := tenantApi.NewTenantHandler(tenantOAuthSvc, tenantOAuthClient)
rg.GET("/tenants/oauth/callback", tenantOAuthHandler.OAuthCallback)
```

- [ ] **Step 4: 编译验证**

Run: `cd /Users/leoyang/Desktop/LeoValue/qikuai/go_oceanengine_vue/backend && go build ./...`
Expected: 编译成功

- [ ] **Step 5: 提交**

```bash
git add internal/app/tenant/service/oauth_client.go internal/app/tenant/api/tenant_handler.go internal/router/router.go
git commit -m "feat: add tenant handler with OAuth authorization flow"
```

---

## Task 8: Account Repository + Service + Handler

**Files:**
- Create: `backend/internal/app/account/repository/account_repo.go`
- Create: `backend/internal/app/account/service/account_service.go`
- Create: `backend/internal/app/account/api/account_handler.go`
- Modify: `backend/internal/router/router.go`

- [ ] **Step 1: 创建 account repository**

```go
// backend/internal/app/account/repository/account_repo.go
package repository

import (
	"context"

	"gorm.io/gorm"
	"oceanengine-backend/internal/app/account/model"
)

type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) BatchCreate(ctx context.Context, accounts []model.LocalAccount) error {
	return r.db.WithContext(ctx).CreateInBatches(accounts, 100).Error
}

func (r *AccountRepository) ListByTenant(ctx context.Context, tenantID uint64, accountIDs []uint64) ([]model.LocalAccount, error) {
	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if len(accountIDs) > 0 {
		query = query.Where("id IN ?", accountIDs)
	}
	var accounts []model.LocalAccount
	err := query.Find(&accounts).Error
	return accounts, err
}

func (r *AccountRepository) GetByID(ctx context.Context, tenantID, id uint64) (*model.LocalAccount, error) {
	var account model.LocalAccount
	err := r.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *AccountRepository) CreateStore(ctx context.Context, store *model.Store) error {
	return r.db.WithContext(ctx).Create(store).Error
}

func (r *AccountRepository) ListStores(ctx context.Context, tenantID uint64, accountID uint64) ([]model.Store, error) {
	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if accountID > 0 {
		query = query.Where("account_id = ?", accountID)
	}
	var stores []model.Store
	err := query.Find(&stores).Error
	return stores, err
}
```

- [ ] **Step 2: 创建 account service**

```go
// backend/internal/app/account/service/account_service.go
package service

import (
	"context"
	"fmt"

	"oceanengine-backend/internal/app/account/model"
	"oceanengine-backend/internal/app/account/repository"
)

type AccountService struct {
	repo *repository.AccountRepository
}

func NewAccountService(repo *repository.AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}

type ImportAccountItem struct {
	LocalAccountID uint64 `json:"local_account_id" binding:"required"`
	Name           string `json:"name" binding:"required"`
	StoreName      string `json:"store_name"`
	City           string `json:"city"`
	District       string `json:"district"`
	Address        string `json:"address"`
}

type ImportAccountsRequest struct {
	Items []ImportAccountItem `json:"items" binding:"required,min=1"`
}

func (s *AccountService) Import(ctx context.Context, tenantID uint64, req ImportAccountsRequest) (int, error) {
	accounts := make([]model.LocalAccount, 0, len(req.Items))
	for _, item := range req.Items {
		accounts = append(accounts, model.LocalAccount{
			TenantID:       tenantID,
			LocalAccountID: item.LocalAccountID,
			Name:           item.Name,
			Status:         "active",
		})
	}
	if err := s.repo.BatchCreate(ctx, accounts); err != nil {
		return 0, fmt.Errorf("batch create accounts: %w", err)
	}
	return len(accounts), nil
}

func (s *AccountService) List(ctx context.Context, tenantID uint64, scopeAccountIDs []uint64) ([]model.LocalAccount, error) {
	return s.repo.ListByTenant(ctx, tenantID, scopeAccountIDs)
}

func (s *AccountService) GetByID(ctx context.Context, tenantID, id uint64) (*model.LocalAccount, error) {
	return s.repo.GetByID(ctx, tenantID, id)
}
```

- [ ] **Step 3: 创建 account handler**

```go
// backend/internal/app/account/api/account_handler.go
package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"oceanengine-backend/internal/app/account/service"
	"oceanengine-backend/pkg/ginutil"
	"oceanengine-backend/pkg/response"
)

type AccountHandler struct {
	svc *service.AccountService
}

func NewAccountHandler(svc *service.AccountService) *AccountHandler {
	return &AccountHandler{svc: svc}
}

func (h *AccountHandler) Import(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	var req service.ImportAccountsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	count, err := h.svc.Import(c.Request.Context(), tenantID, req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, gin.H{"imported": count})
}

func (h *AccountHandler) List(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	var scopeIDs []uint64
	if v, exists := c.Get("scope_account_ids"); exists {
		scopeIDs, _ = v.([]uint64)
	}
	accounts, err := h.svc.List(c.Request.Context(), tenantID, scopeIDs)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, accounts)
}

func (h *AccountHandler) GetByID(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	account, err := h.svc.GetByID(c.Request.Context(), tenantID, id)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, account)
}
```

- [ ] **Step 4: 在 router.go 注册账户路由**

添加 import：

```go
accountApi "oceanengine-backend/internal/app/account/api"
accountRepo "oceanengine-backend/internal/app/account/repository"
accountService "oceanengine-backend/internal/app/account/service"
```

新增方法：

```go
func (r *Router) registerAccountRoutes(rg *gin.RouterGroup) {
	repo := accountRepo.NewAccountRepository(r.db)
	svc := accountService.NewAccountService(repo)
	handler := accountApi.NewAccountHandler(svc)

	accounts := rg.Group("/accounts")
	{
		accounts.POST("/import", handler.Import)
		accounts.GET("", handler.List)
		accounts.GET("/:id", handler.GetByID)
	}
}
```

- [ ] **Step 5: 编译验证**

Run: `cd /Users/leoyang/Desktop/LeoValue/qikuai/go_oceanengine_vue/backend && go build ./...`
Expected: 编译成功

- [ ] **Step 6: 提交**

```bash
git add internal/app/account/ internal/router/router.go
git commit -m "feat: add account import, list, get API"
```

---

## Task 9: Group Repository + Service + Handler

**Files:**
- Create: `backend/internal/app/group/repository/group_repo.go`
- Create: `backend/internal/app/group/service/group_service.go`
- Create: `backend/internal/app/group/api/group_handler.go`
- Modify: `backend/internal/router/router.go`

- [ ] **Step 1: 创建 group repository**

```go
// backend/internal/app/group/repository/group_repo.go
package repository

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"oceanengine-backend/internal/app/group/model"
)

type GroupRepository struct {
	db *gorm.DB
}

func NewGroupRepository(db *gorm.DB) *GroupRepository {
	return &GroupRepository{db: db}
}

func (r *GroupRepository) Create(ctx context.Context, group *model.AccountGroup) error {
	return r.db.WithContext(ctx).Create(group).Error
}

func (r *GroupRepository) GetByID(ctx context.Context, tenantID, id uint64) (*model.AccountGroup, error) {
	var group model.AccountGroup
	err := r.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).First(&group).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *GroupRepository) List(ctx context.Context, tenantID uint64, groupType string) ([]model.AccountGroup, error) {
	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if groupType != "" {
		query = query.Where("type = ?", groupType)
	}
	var groups []model.AccountGroup
	err := query.Order("sort ASC, id ASC").Find(&groups).Error
	return groups, err
}

func (r *GroupRepository) Update(ctx context.Context, group *model.AccountGroup) error {
	return r.db.WithContext(ctx).Save(group).Error
}

func (r *GroupRepository) Delete(ctx context.Context, tenantID, id uint64) error {
	return r.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).Delete(&model.AccountGroup{}).Error
}

func (r *GroupRepository) AddMembers(ctx context.Context, groupID uint64, accountIDs []uint64) error {
	members := make([]model.AccountGroupMember, 0, len(accountIDs))
	for _, aid := range accountIDs {
		members = append(members, model.AccountGroupMember{AccountID: aid, GroupID: groupID})
	}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(members, 100).Error
}

func (r *GroupRepository) RemoveMembers(ctx context.Context, groupID uint64, accountIDs []uint64) error {
	return r.db.WithContext(ctx).Where("group_id = ? AND account_id IN ?", groupID, accountIDs).Delete(&model.AccountGroupMember{}).Error
}

func (r *GroupRepository) ListMembers(ctx context.Context, groupID uint64) ([]uint64, error) {
	var accountIDs []uint64
	err := r.db.WithContext(ctx).Model(&model.AccountGroupMember{}).Where("group_id = ?", groupID).Pluck("account_id", &accountIDs).Error
	return accountIDs, err
}
```

- [ ] **Step 2: 创建 group service**

```go
// backend/internal/app/group/service/group_service.go
package service

import (
	"context"

	"oceanengine-backend/internal/app/group/model"
	"oceanengine-backend/internal/app/group/repository"
)

type GroupService struct {
	repo *repository.GroupRepository
}

func NewGroupService(repo *repository.GroupRepository) *GroupService {
	return &GroupService{repo: repo}
}

type CreateGroupRequest struct {
	Name     string          `json:"name" binding:"required"`
	Type     model.GroupType `json:"type" binding:"required,oneof=franchisee region custom"`
	ParentID uint64          `json:"parent_id"`
	Sort     int             `json:"sort"`
}

type UpdateGroupRequest struct {
	Name string `json:"name"`
	Sort int    `json:"sort"`
}

type MembersRequest struct {
	AccountIDs []uint64 `json:"account_ids" binding:"required,min=1"`
}

func (s *GroupService) Create(ctx context.Context, tenantID uint64, req CreateGroupRequest) (*model.AccountGroup, error) {
	group := &model.AccountGroup{
		TenantID: tenantID,
		Name:     req.Name,
		Type:     req.Type,
		ParentID: req.ParentID,
		Sort:     req.Sort,
	}
	if err := s.repo.Create(ctx, group); err != nil {
		return nil, err
	}
	return group, nil
}

func (s *GroupService) List(ctx context.Context, tenantID uint64, groupType string) ([]model.AccountGroup, error) {
	return s.repo.List(ctx, tenantID, groupType)
}

func (s *GroupService) Update(ctx context.Context, tenantID, id uint64, req UpdateGroupRequest) error {
	group, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if req.Name != "" {
		group.Name = req.Name
	}
	group.Sort = req.Sort
	return s.repo.Update(ctx, group)
}

func (s *GroupService) Delete(ctx context.Context, tenantID, id uint64) error {
	return s.repo.Delete(ctx, tenantID, id)
}

func (s *GroupService) AddMembers(ctx context.Context, tenantID, groupID uint64, req MembersRequest) error {
	if _, err := s.repo.GetByID(ctx, tenantID, groupID); err != nil {
		return err
	}
	return s.repo.AddMembers(ctx, groupID, req.AccountIDs)
}

func (s *GroupService) RemoveMembers(ctx context.Context, tenantID, groupID uint64, req MembersRequest) error {
	return s.repo.RemoveMembers(ctx, groupID, req.AccountIDs)
}
```

- [ ] **Step 3: 创建 group handler**

```go
// backend/internal/app/group/api/group_handler.go
package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"oceanengine-backend/internal/app/group/service"
	"oceanengine-backend/pkg/ginutil"
	"oceanengine-backend/pkg/response"
)

type GroupHandler struct {
	svc *service.GroupService
}

func NewGroupHandler(svc *service.GroupService) *GroupHandler {
	return &GroupHandler{svc: svc}
}

func (h *GroupHandler) Create(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	var req service.CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	group, err := h.svc.Create(c.Request.Context(), tenantID, req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, group)
}

func (h *GroupHandler) List(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	groupType := c.Query("type")
	groups, err := h.svc.List(c.Request.Context(), tenantID, groupType)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, groups)
}

func (h *GroupHandler) Update(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	var req service.UpdateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.Update(c.Request.Context(), tenantID, id, req); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *GroupHandler) Delete(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), tenantID, id); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *GroupHandler) AddMembers(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	var req service.MembersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.AddMembers(c.Request.Context(), tenantID, id, req); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *GroupHandler) RemoveMembers(c *gin.Context) {
	tenantID := ginutil.GetTenantID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	var req service.MembersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.RemoveMembers(c.Request.Context(), tenantID, id, req); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, nil)
}
```

- [ ] **Step 4: 在 router.go 注册分组路由**

添加 import：

```go
groupApi "oceanengine-backend/internal/app/group/api"
groupRepo "oceanengine-backend/internal/app/group/repository"
groupService "oceanengine-backend/internal/app/group/service"
```

新增方法：

```go
func (r *Router) registerGroupRoutes(rg *gin.RouterGroup) {
	repo := groupRepo.NewGroupRepository(r.db)
	svc := groupService.NewGroupService(repo)
	handler := groupApi.NewGroupHandler(svc)

	groups := rg.Group("/groups")
	{
		groups.POST("", handler.Create)
		groups.GET("", handler.List)
		groups.PUT("/:id", handler.Update)
		groups.DELETE("/:id", handler.Delete)
		groups.POST("/:id/members", handler.AddMembers)
		groups.DELETE("/:id/members", handler.RemoveMembers)
	}
}
```

- [ ] **Step 5: 编译验证**

Run: `cd /Users/leoyang/Desktop/LeoValue/qikuai/go_oceanengine_vue/backend && go build ./...`
Expected: 编译成功

- [ ] **Step 6: 提交**

```bash
git add internal/app/group/ internal/router/router.go
git commit -m "feat: add group CRUD and member management API"
```

---

## Task 10: Scope Handler + 路由注册

**Files:**
- Create: `backend/internal/app/scope/repository/scope_repo.go`
- Create: `backend/internal/app/scope/service/scope_service.go`
- Create: `backend/internal/app/scope/api/scope_handler.go`
- Modify: `backend/internal/router/router.go`

- [ ] **Step 1: 创建 scope repository**

```go
// backend/internal/app/scope/repository/scope_repo.go
package repository

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"oceanengine-backend/internal/app/scope/model"
)

type ScopeRepository struct {
	db *gorm.DB
}

func NewScopeRepository(db *gorm.DB) *ScopeRepository {
	return &ScopeRepository{db: db}
}

func (r *ScopeRepository) SetUserScope(ctx context.Context, userID uint64, accountIDs []uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&model.UserAccountScope{}).Error; err != nil {
			return err
		}
		if len(accountIDs) == 0 {
			return nil
		}
		scopes := make([]model.UserAccountScope, 0, len(accountIDs))
		for _, aid := range accountIDs {
			scopes = append(scopes, model.UserAccountScope{UserID: userID, AccountID: aid})
		}
		return tx.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(scopes, 100).Error
	})
}

func (r *ScopeRepository) GetUserScope(ctx context.Context, userID uint64) ([]uint64, error) {
	var accountIDs []uint64
	err := r.db.WithContext(ctx).Model(&model.UserAccountScope{}).Where("user_id = ?", userID).Pluck("account_id", &accountIDs).Error
	return accountIDs, err
}
```

- [ ] **Step 2: 创建 scope service**

```go
// backend/internal/app/scope/service/scope_service.go
package service

import (
	"context"

	"oceanengine-backend/internal/app/scope/repository"
)

type ScopeService struct {
	repo *repository.ScopeRepository
}

func NewScopeService(repo *repository.ScopeRepository) *ScopeService {
	return &ScopeService{repo: repo}
}

type SetScopeRequest struct {
	AccountIDs []uint64 `json:"account_ids" binding:"required"`
}

func (s *ScopeService) SetUserScope(ctx context.Context, userID uint64, req SetScopeRequest) error {
	return s.repo.SetUserScope(ctx, userID, req.AccountIDs)
}

func (s *ScopeService) GetUserScope(ctx context.Context, userID uint64) ([]uint64, error) {
	return s.repo.GetUserScope(ctx, userID)
}
```

- [ ] **Step 3: 创建 scope handler**

```go
// backend/internal/app/scope/api/scope_handler.go
package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"oceanengine-backend/internal/app/scope/service"
	"oceanengine-backend/pkg/response"
)

type ScopeHandler struct {
	svc *service.ScopeService
}

func NewScopeHandler(svc *service.ScopeService) *ScopeHandler {
	return &ScopeHandler{svc: svc}
}

func (h *ScopeHandler) SetScope(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}
	var req service.SetScopeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.SetUserScope(c.Request.Context(), userID, req); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *ScopeHandler) GetScope(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}
	accountIDs, err := h.svc.GetUserScope(c.Request.Context(), userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, gin.H{"account_ids": accountIDs})
}
```

- [ ] **Step 4: 在 router.go 的 registerSystemRoutes 中注册 scope 端点**

添加 import：

```go
scopeApi "oceanengine-backend/internal/app/scope/api"
scopeRepo "oceanengine-backend/internal/app/scope/repository"
scopeService "oceanengine-backend/internal/app/scope/service"
```

在 `registerSystemRoutes` 方法中，users group 内添加：

```go
// 投手权限范围管理
scopeRepository := scopeRepo.NewScopeRepository(r.db)
scopeSvc := scopeService.NewScopeService(scopeRepository)
scopeHandler := scopeApi.NewScopeHandler(scopeSvc)
users.POST("/:id/scope", scopeHandler.SetScope)
users.GET("/:id/scope", scopeHandler.GetScope)
```

- [ ] **Step 5: 编译验证**

Run: `cd /Users/leoyang/Desktop/LeoValue/qikuai/go_oceanengine_vue/backend && go build ./...`
Expected: 编译成功

- [ ] **Step 6: 提交**

```bash
git add internal/app/scope/ internal/router/router.go
git commit -m "feat: add user account scope management API"
```

---

## Task 11: data_scope 中间件

**Files:**
- Create: `backend/internal/middleware/data_scope.go`
- Create: `backend/internal/middleware/data_scope_test.go`
- Modify: `backend/internal/router/router.go`

- [ ] **Step 1: 编写 data_scope 中间件测试（含 admin bypass 和普通用户场景）**

```go
// backend/internal/middleware/data_scope_test.go
package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestDataScope_AdminBypassesScope(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint64(1))
		c.Set("tenant_id", uint64(1))
		c.Set("role_key", "admin")
		c.Next()
	})
	r.Use(DataScope(nil))
	r.GET("/test", func(c *gin.Context) {
		_, exists := c.Get("scope_account_ids")
		assert.False(t, exists, "admin should not have scope_account_ids set")
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDataScope_NonAdminWithoutScope_Returns403(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint64(99))
		c.Set("tenant_id", uint64(1))
		c.Set("role_key", "operator")
		c.Next()
	})
	// 传 nil db 模拟无 scope 记录（Pluck 返回空）
	r.Use(DataScope(nil))
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}
```

注意：第二个测试中 db=nil 会 panic，需要在实现中处理。改为使用 ScopeRepository 接口注入。

- [ ] **Step 2: 运行测试确认失败（RED）**

Run: `cd /Users/leoyang/Desktop/LeoValue/qikuai/go_oceanengine_vue/backend && go test ./internal/middleware/ -v -run TestDataScope`
Expected: FAIL — DataScope 未定义

- [ ] **Step 3: 实现 data_scope 中间件**

```go
// backend/internal/middleware/data_scope.go
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	scopeModel "oceanengine-backend/internal/app/scope/model"
	"oceanengine-backend/pkg/ginutil"
)

func DataScope(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleKey := ginutil.GetRoleKey(c)
		if roleKey == "admin" || roleKey == "super_admin" {
			c.Next()
			return
		}

		userID := ginutil.GetUserID(c)
		if userID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if db == nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "no account scope assigned"})
			return
		}

		var accountIDs []uint64
		db.Model(&scopeModel.UserAccountScope{}).
			Where("user_id = ?", userID).
			Pluck("account_id", &accountIDs)

		if len(accountIDs) == 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "no account scope assigned"})
			return
		}

		c.Set("scope_account_ids", accountIDs)
		c.Next()
	}
}
```

- [ ] **Step 4: 运行测试确认通过（GREEN）**

Run: `cd /Users/leoyang/Desktop/LeoValue/qikuai/go_oceanengine_vue/backend && go test ./internal/middleware/ -v -run TestDataScope`
Expected: PASS

- [ ] **Step 5: 在路由中应用中间件**

在 `registerProtectedRoutes` 中，将需要数据隔离的路由放到 scoped 组：

```go
// 需要数据权限隔离的路由
scoped := rg.Group("")
scoped.Use(middleware.DataScope(r.db))
r.registerAccountRoutes(scoped)
r.registerGroupRoutes(scoped)
```

- [ ] **Step 6: 提交**

```bash
git add internal/middleware/data_scope.go internal/middleware/data_scope_test.go internal/router/router.go
git commit -m "feat: add data_scope middleware for account-level isolation"
```

---

## Task 12: Token 自动续期

**Files:**
- Create: `backend/internal/app/tenant/service/token_refresher.go`
- Modify: `backend/cmd/server/main.go`

- [ ] **Step 1: 创建 token refresher**

```go
// backend/internal/app/tenant/service/token_refresher.go
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
```

- [ ] **Step 2: 在 main.go 中启动 token refresher goroutine**

在服务启动部分（路由注册之后、ListenAndServe 之前）添加：

```go
import (
	tenantRepository "oceanengine-backend/internal/app/tenant/repository"
	tenantService "oceanengine-backend/internal/app/tenant/service"
)

// 创建带 cancel 的 context
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// Token 自动续期
tRepo := tenantRepository.NewTenantRepository(db)
oauthClient := tenantService.NewOAuthClient()
tokenRefresher := tenantService.NewTokenRefresher(tRepo, oauthClient, logger)
go tokenRefresher.Start(ctx)
```

- [ ] **Step 3: 编译验证**

Run: `cd /Users/leoyang/Desktop/LeoValue/qikuai/go_oceanengine_vue/backend && go build ./...`
Expected: 编译成功

- [ ] **Step 4: 提交**

```bash
git add internal/app/tenant/service/token_refresher.go cmd/server/main.go
git commit -m "feat: add automatic OAuth token refresh goroutine"
```

---

## Task 13: 集成测试

**Files:**
- Create: `backend/test/integration/tenant_test.go`
- Create: `backend/test/integration/account_test.go`
- Create: `backend/test/integration/group_test.go`
- Modify: `backend/test/integration/testutil.go`

- [ ] **Step 1: 扩展 testutil.go 添加辅助函数**

在现有 `backend/test/integration/testutil.go` 中添加（参考现有测试的 setup 模式）：

```go
package integration

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"oceanengine-backend/pkg/auth"
)

var testRouter *gin.Engine
var testJWTManager *auth.JWTManager

func init() {
	// 使用测试密钥初始化 JWT manager
	testJWTManager = auth.NewJWTManager("test-secret-key-for-integration", 24*time.Hour)
}

func doRequest(t *testing.T, method, path, body, token string) *httptest.ResponseRecorder {
	t.Helper()
	var req *http.Request
	if body != "" {
		req, _ = http.NewRequest(method, path, bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, _ = http.NewRequest(method, path, nil)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	return w
}

func assertStatus(t *testing.T, w *httptest.ResponseRecorder, expected int) {
	t.Helper()
	assert.Equal(t, expected, w.Code, "response body: %s", w.Body.String())
}

func adminToken() string {
	// tenant_id=0 表示平台管理员, role_key=admin
	token, _ := testJWTManager.GenerateToken(1, "admin", "admin", 0)
	return token
}

func userToken() string {
	// tenant_id=1, role_key=operator
	token, _ := testJWTManager.GenerateToken(2, "operator1", "operator", 1)
	return token
}
```

注意：`testRouter` 的初始化需要在 `TestMain` 中完成，使用 SQLite 内存数据库 + 测试用 JWT secret。具体实现参考现有 `testutil.go` 中的 setup 模式，确保 `testJWTManager` 与路由中使用的是同一个实例。

- [ ] **Step 2: 编写租户 CRUD 集成测试**

```go
// backend/test/integration/tenant_test.go
package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTenantCreate(t *testing.T) {
	body := `{"name":"测试品牌","oauth_app_id":"test_app_001","oauth_secret":"test_secret_001"}`
	w := doRequest(t, "POST", "/api/v1/tenants", body, adminToken())
	assertStatus(t, w, http.StatusOK)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "测试品牌", data["name"])
	assert.Equal(t, "expired", data["token_status"])
}

func TestTenantList(t *testing.T) {
	w := doRequest(t, "GET", "/api/v1/tenants", "", adminToken())
	assertStatus(t, w, http.StatusOK)
}

func TestTenantGetByID(t *testing.T) {
	w := doRequest(t, "GET", "/api/v1/tenants/1", "", adminToken())
	assertStatus(t, w, http.StatusOK)
}
```

- [ ] **Step 3: 编写账户导入集成测试**

```go
// backend/test/integration/account_test.go
package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccountImport(t *testing.T) {
	body := `{"items":[{"local_account_id":100001,"name":"门店A"},{"local_account_id":100002,"name":"门店B"}]}`
	w := doRequest(t, "POST", "/api/v1/accounts/import", body, userToken())
	assertStatus(t, w, http.StatusOK)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, float64(2), data["imported"])
}

func TestAccountList(t *testing.T) {
	w := doRequest(t, "GET", "/api/v1/accounts", "", userToken())
	assertStatus(t, w, http.StatusOK)
}

func TestAccountGetByID(t *testing.T) {
	w := doRequest(t, "GET", "/api/v1/accounts/1", "", userToken())
	assertStatus(t, w, http.StatusOK)
}
```

- [ ] **Step 4: 编写分组管理集成测试**

```go
// backend/test/integration/group_test.go
package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroupCreate(t *testing.T) {
	body := `{"name":"华东区","type":"region"}`
	w := doRequest(t, "POST", "/api/v1/groups", body, userToken())
	assertStatus(t, w, http.StatusOK)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "华东区", data["name"])
	assert.Equal(t, "region", data["type"])
}

func TestGroupAddMembers(t *testing.T) {
	body := `{"account_ids":[1,2]}`
	w := doRequest(t, "POST", "/api/v1/groups/1/members", body, userToken())
	assertStatus(t, w, http.StatusOK)
}

func TestGroupList(t *testing.T) {
	w := doRequest(t, "GET", "/api/v1/groups?type=region", "", userToken())
	assertStatus(t, w, http.StatusOK)
}

func TestGroupRemoveMembers(t *testing.T) {
	body := `{"account_ids":[2]}`
	w := doRequest(t, "DELETE", "/api/v1/groups/1/members", body, userToken())
	assertStatus(t, w, http.StatusOK)
}

func TestGroupDelete(t *testing.T) {
	w := doRequest(t, "DELETE", "/api/v1/groups/1", "", userToken())
	assertStatus(t, w, http.StatusOK)
}
```

- [ ] **Step 5: 运行集成测试**

Run: `cd /Users/leoyang/Desktop/LeoValue/qikuai/go_oceanengine_vue/backend && go test ./test/integration/ -v -run "TestTenant|TestAccount|TestGroup"`
Expected: PASS

- [ ] **Step 6: 提交**

```bash
git add test/integration/tenant_test.go test/integration/account_test.go test/integration/group_test.go test/integration/testutil.go
git commit -m "test: add integration tests for tenant, account, group modules"
```

---

## 完成标准

Phase 1 完成后，系统应具备以下能力：

1. 平台管理员可以创建租户（品牌），配置 OAuth AppID/Secret
2. 管理员可以发起 OAuth 授权，回调自动保存 Token，Token 每 5 分钟自动续期
3. 投手可以通过 JSON 批量导入本地推账户
4. 账户可以按加盟商/区域/自定义维度分组，支持批量添加/移除成员
5. 管理员可以为投手分配账户权限范围
6. 投手只能看到和操作自己权限范围内的账户（data_scope 中间件）
7. API 调用按租户隔离限流（10 QPS/租户，Redis Lua 令牌桶）
8. 所有 API 通过集成测试验证

---

## 依赖关系

```
Task 1 (models + ginutil) ─┐
Task 2 (group/scope model) ┼──→ Task 4 (AutoMigrate)
Task 3 (user tenant_id + JWT) ─┘       │
                                        ▼
Task 5 (rate limiter) ──────── 独立，无依赖
                                        │
Task 4 ──→ Task 6 (tenant repo/svc) ──→ Task 7 (tenant handler + OAuth)
Task 4 ──→ Task 8 (account repo/svc/handler)
Task 4 ──→ Task 9 (group repo/svc/handler)
Task 4 ──→ Task 10 (scope handler)
Task 10 ─→ Task 11 (data_scope middleware)
Task 7 ──→ Task 12 (token refresher)
All above → Task 13 (integration tests)
```
