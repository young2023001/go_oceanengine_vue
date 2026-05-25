# Phase 4: 打磨与前端实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 实现前端全部页面（Vue 3）、后端性能优化（缓存/索引）、容错机制（重试/断点续传/幂等）。

**Architecture:** 前端使用 Vue 3 Composition API + TypeScript + TailwindCSS + vue-chartjs。后端增加 Redis 缓存层、SDK 调用重试、批量任务断点续传。

**Tech Stack:** Vue 3 + TypeScript + Vite + TailwindCSS + vue-chartjs + axios | Go + Redis cache

---

## 文件结构

### 前端新建文件

| 路径 | 职责 |
|------|------|
| src/views/tenant/TenantList.vue | 租户管理页 |
| src/views/account/AccountList.vue | 账户列表 + 导入 |
| src/views/account/GroupManage.vue | 分组管理 |
| src/views/project/ProjectList.vue | 项目列表 |
| src/views/project/PromotionList.vue | 单元列表 |
| src/views/batch/BatchCreate.vue | 批量操作向导 |
| src/views/batch/BatchTaskList.vue | 任务列表 + 进度 |
| src/views/analytics/Dashboard.vue | 数据看板 |
| src/views/project/ProjectList.vue | 项目列表页 |
| src/views/account/AccountList.vue | 账户列表 + 导入页 |
| src/views/scope/ScopeManage.vue | 投手权限分配 |
| src/api/tenant.ts | 租户 API |
| src/api/account.ts | 账户 API |
| src/api/project.ts | 项目 API |
| src/api/batch.ts | 批量操作 API |
| src/api/analytics.ts | 分析 API |
| src/api/scope.ts | 权限 API |

### 后端新建/修改文件

| 路径 | 职责 |
|------|------|
| backend/pkg/cache/report_cache.go | 报表查询 Redis 缓存 |
| backend/pkg/retry/retry.go | SDK 调用指数退避重试 |
| backend/internal/app/batch/service/recovery.go | 批量任务断点续传 |
| backend/cmd/server/main.go | 启动时恢复 running 任务 |

---

## Task 1: SDK 调用重试机制

**Files:**
- Create: backend/pkg/retry/retry.go

- [ ] **Step 1: 创建重试工具**

```go
// backend/pkg/retry/retry.go
package retry

import (
	"context"
	"math"
	"time"
)

type Config struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
}

func DefaultConfig() Config {
	return Config{
		MaxAttempts: 3,
		BaseDelay:   500 * time.Millisecond,
		MaxDelay:    5 * time.Second,
	}
}

func Do(ctx context.Context, cfg Config, fn func() error) error {
	var lastErr error
	for attempt := 0; attempt < cfg.MaxAttempts; attempt++ {
		lastErr = fn()
		if lastErr == nil {
			return nil
		}
		if attempt == cfg.MaxAttempts-1 {
			break
		}
		delay := time.Duration(float64(cfg.BaseDelay) * math.Pow(2, float64(attempt)))
		if delay > cfg.MaxDelay {
			delay = cfg.MaxDelay
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
	}
	return lastErr
}
```

- [ ] **Step 2: 编译验证并提交**

```bash
git add pkg/retry/retry.go
git commit -m "feat: add exponential backoff retry utility"
```

---

## Task 2: 报表查询缓存

**Files:**
- Create: backend/pkg/cache/report_cache.go

- [ ] **Step 1: 创建报表缓存层**

```go
// backend/pkg/cache/report_cache.go
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type ReportCache struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewReportCache(rdb *redis.Client) *ReportCache {
	return &ReportCache{rdb: rdb, ttl: 5 * time.Minute}
}

func (c *ReportCache) Get(ctx context.Context, key string, dest interface{}) bool {
	val, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		return false
	}
	json.Unmarshal([]byte(val), dest)
	return true
}

func (c *ReportCache) Set(ctx context.Context, key string, value interface{}) {
	data, _ := json.Marshal(value)
	c.rdb.Set(ctx, key, string(data), c.ttl)
}

func (c *ReportCache) Invalidate(ctx context.Context, tenantID uint64) {
	pattern := fmt.Sprintf("report:%d:*", tenantID)
	keys, _ := c.rdb.Keys(ctx, pattern).Result()
	if len(keys) > 0 {
		c.rdb.Del(ctx, keys...)
	}
}

func OverviewKey(tenantID uint64, date string) string {
	return fmt.Sprintf("report:%d:overview:%s", tenantID, date)
}

func TrendKey(tenantID uint64, start, end string) string {
	return fmt.Sprintf("report:%d:trend:%s_%s", tenantID, start, end)
}
```

- [ ] **Step 2: 在 analytics service 中集成缓存**

修改 `backend/internal/app/analytics/service/analytics_service.go` 的 GetOverview 方法：

```go
func (s *AnalyticsService) GetOverview(ctx context.Context, tenantID uint64, scopeIDs []uint64, req OverviewRequest) (*repository.OverviewResult, error) {
	date := time.Now()
	if req.Date != "" {
		parsed, err := time.Parse("2006-01-02", req.Date)
		if err == nil {
			date = parsed
		}
	}

	// 先查缓存
	cacheKey := cache.OverviewKey(tenantID, date.Format("2006-01-02"))
	var cached repository.OverviewResult
	if s.cache.Get(ctx, cacheKey, &cached) {
		return &cached, nil
	}

	// 缓存 miss，查 DB
	result, err := s.repo.GetOverview(ctx, tenantID, scopeIDs, date)
	if err != nil {
		return nil, err
	}

	// 写入缓存
	s.cache.Set(ctx, cacheKey, result)
	return result, nil
}
```

同时修改 AnalyticsService 结构体，注入 cache。在 import 中添加：

```go
import (
	"oceanengine-backend/pkg/cache"
)
```

修改结构体：

```go
type AnalyticsService struct {
	repo  *repository.ReportRepository
	cache *cache.ReportCache
}

func NewAnalyticsService(repo *repository.ReportRepository, c *cache.ReportCache) *AnalyticsService {
	return &AnalyticsService{repo: repo, cache: c}
}
```

- [ ] **Step 3: 编译验证并提交**

```bash
git add pkg/cache/report_cache.go
git commit -m "feat: add Redis report cache with 5min TTL"
```

---

## Task 3: 批量任务断点续传

**Files:**
- Create: backend/internal/app/batch/service/recovery.go
- Modify: backend/cmd/server/main.go

- [ ] **Step 1: 创建 recovery 逻辑**

```go
// backend/internal/app/batch/service/recovery.go
package service

import (
	"context"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"oceanengine-backend/internal/app/batch/model"
)

type TaskRecovery struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewTaskRecovery(db *gorm.DB, logger *zap.Logger) *TaskRecovery {
	return &TaskRecovery{db: db, logger: logger}
}

func (r *TaskRecovery) RecoverOnStartup(ctx context.Context) {
	var tasks []model.BatchTask
	r.db.Where("status = ?", model.TaskStatusRunning).Find(&tasks)

	r.logger.Info("recovering batch tasks", zap.Int("count", len(tasks)))

	for _, task := range tasks {
		// 将 running 状态的 item 重置为 pending
		r.db.Model(&model.BatchTaskItem{}).
			Where("task_id = ? AND status = ?", task.ID, model.ItemStatusRunning).
			Update("status", model.ItemStatusPending)
	}
}
```

- [ ] **Step 2: 在 main.go 启动时调用**

```go
recovery := batchService.NewTaskRecovery(db, logger)
recovery.RecoverOnStartup(ctx)
```

- [ ] **Step 3: 编译验证并提交**

```bash
git add internal/app/batch/service/recovery.go cmd/server/main.go
git commit -m "feat: add batch task recovery on server restart"
```

---

## Task 4: 健康检查增强

**Files:**
- Modify: backend/internal/router/router.go

- [ ] **Step 1: 增强 /health 端点**

```go
func (r *Router) healthCheck(c *gin.Context) {
	// 检查 MySQL
	sqlDB, err := r.db.DB()
	if err != nil || sqlDB.Ping() != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "mysql": "down"})
		return
	}
	// 检查 Redis（如果有 redis client 注入）
	c.JSON(http.StatusOK, gin.H{"status": "ok", "mysql": "up", "redis": "up"})
}
```

- [ ] **Step 2: 提交**

```bash
git add internal/router/router.go
git commit -m "feat: enhance health check with MySQL/Redis connectivity"
```

---

## Task 5: 前端 API 层

**Files:**
- Create: src/api/tenant.ts
- Create: src/api/account.ts
- Create: src/api/project.ts
- Create: src/api/batch.ts
- Create: src/api/analytics.ts

- [ ] **Step 1: 创建 API 模块**

```typescript
// src/api/tenant.ts
import request from '@/utils/request'

export function getTenants() {
  return request.get('/tenants')
}
export function createTenant(data: { name: string; oauth_app_id: string; oauth_secret: string }) {
  return request.post('/tenants', data)
}
export function getTenantOAuthURL(id: number) {
  return request.get(`/tenants/${id}/oauth/url`)
}
```

```typescript
// src/api/account.ts
import request from '@/utils/request'

export function getAccounts(params?: { page?: number; page_size?: number }) {
  return request.get('/accounts', { params })
}
export function importAccounts(data: { items: Array<{ local_account_id: number; name: string }> }) {
  return request.post('/accounts/import', data)
}
```

```typescript
// src/api/project.ts
import request from '@/utils/request'

export function getProjects(params?: { status?: string; page?: number; page_size?: number }) {
  return request.get('/projects', { params })
}
export function getProject(id: number) {
  return request.get(`/projects/${id}`)
}
export function updateProjectStatus(id: number, status: string) {
  return request.put(`/projects/${id}/status`, { status })
}
export function getPromotions(params?: { project_id?: number; page?: number; page_size?: number }) {
  return request.get('/promotions', { params })
}
```

```typescript
// src/api/batch.ts
import request from '@/utils/request'

export function createBatchTask(data: { task_type: string; title: string; config: object; account_ids: number[] }) {
  return request.post('/batch/projects', data)
}
export function getBatchTasks(params?: { page?: number; page_size?: number }) {
  return request.get('/batch/tasks', { params })
}
export function getBatchTask(id: number) {
  return request.get(`/batch/tasks/${id}`)
}
export function cancelBatchTask(id: number) {
  return request.post(`/batch/tasks/${id}/cancel`)
}
export function retryBatchTask(id: number) {
  return request.post(`/batch/tasks/${id}/retry`)
}
```

```typescript
// src/api/analytics.ts
import request from '@/utils/request'

export function getOverview(params?: { date?: string }) {
  return request.get('/analytics/overview', { params })
}
export function getTrend(params: { start_date: string; end_date: string }) {
  return request.get('/analytics/trend', { params })
}
export function getRank(params: { start_date: string; end_date: string; order_by?: string; limit?: number }) {
  return request.get('/analytics/rank', { params })
}
export function getCompare(params: { date: string; compare_type?: string }) {
  return request.get('/analytics/compare', { params })
}
export function getDetail(params: { start_date: string; end_date: string; page?: number; page_size?: number }) {
  return request.get('/analytics/detail', { params })
}
export function createExport(data: { title: string; start_date: string; end_date: string }) {
  return request.post('/analytics/export', data)
}
export function getExports() {
  return request.get('/analytics/exports')
}
```

- [ ] **Step 2: 提交**

```bash
git add src/api/
git commit -m "feat: add frontend API modules for all endpoints"
```

---

## Task 6: 数据看板页面

**Files:**
- Create: src/views/analytics/Dashboard.vue

- [ ] **Step 1: 创建看板页面**

```vue
<!-- src/views/analytics/Dashboard.vue -->
<template>
  <div class="p-6 space-y-6">
    <h1 class="text-2xl font-bold">数据看板</h1>

    <!-- 汇总卡片 -->
    <div class="grid grid-cols-4 gap-4">
      <div v-for="card in summaryCards" :key="card.label" class="bg-white rounded-lg shadow p-4">
        <p class="text-sm text-gray-500">{{ card.label }}</p>
        <p class="text-2xl font-bold mt-1">{{ card.value }}</p>
        <p v-if="card.change !== undefined" class="text-xs mt-1" :class="card.change >= 0 ? 'text-green-500' : 'text-red-500'">
          {{ card.change >= 0 ? '+' : '' }}{{ card.change }}%
        </p>
      </div>
    </div>

    <!-- 趋势图 -->
    <div class="bg-white rounded-lg shadow p-4">
      <h2 class="text-lg font-semibold mb-4">消耗趋势</h2>
      <Line :data="chartData" :options="chartOptions" />
    </div>

    <!-- 排行榜 -->
    <div class="bg-white rounded-lg shadow p-4">
      <h2 class="text-lg font-semibold mb-4">门店排行</h2>
      <table class="w-full text-sm">
        <thead>
          <tr class="border-b">
            <th class="text-left py-2">排名</th>
            <th class="text-left py-2">门店</th>
            <th class="text-right py-2">消耗</th>
            <th class="text-right py-2">转化</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(item, idx) in rankList" :key="item.account_id" class="border-b">
            <td class="py-2">{{ idx + 1 }}</td>
            <td class="py-2">{{ item.name }}</td>
            <td class="text-right py-2">{{ item.cost.toFixed(2) }}</td>
            <td class="text-right py-2">{{ item.convert_cnt }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { Line } from 'vue-chartjs'
import { Chart as ChartJS, CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend } from 'chart.js'
import { getOverview, getTrend, getRank, getCompare } from '@/api/analytics'

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend)

const overview = ref<any>({})
const compare = ref<any>({})
const trendData = ref<any[]>([])
const rankList = ref<any[]>([])

const summaryCards = computed(() => [
  { label: '今日消耗', value: formatMoney(overview.value.total_cost), change: calcChange('total_cost') },
  { label: '展示量', value: formatNum(overview.value.total_show), change: calcChange('total_show') },
  { label: '点击量', value: formatNum(overview.value.total_click), change: calcChange('total_click') },
  { label: '转化量', value: formatNum(overview.value.total_convert), change: calcChange('total_convert') },
])

const chartData = computed(() => ({
  labels: trendData.value.map(p => p.stat_date),
  datasets: [{
    label: '消耗',
    data: trendData.value.map(p => p.cost),
    borderColor: '#3b82f6',
    tension: 0.3,
  }]
}))

const chartOptions = { responsive: true, plugins: { legend: { display: false } } }

function formatMoney(v: number) { return v ? '¥' + v.toFixed(2) : '¥0' }
function formatNum(v: number) { return v ? v.toLocaleString() : '0' }
function calcChange(field: string) {
  const cur = overview.value[field] || 0
  const prev = compare.value?.compare?.[field] || 0
  if (prev === 0) return 0
  return Math.round((cur - prev) / prev * 100)
}

onMounted(async () => {
  const today = new Date().toISOString().slice(0, 10)
  const weekAgo = new Date(Date.now() - 7 * 86400000).toISOString().slice(0, 10)

  const [ovRes, cmpRes, trendRes, rankRes] = await Promise.all([
    getOverview({ date: today }),
    getCompare({ date: today, compare_type: 'yesterday' }),
    getTrend({ start_date: weekAgo, end_date: today }),
    getRank({ start_date: weekAgo, end_date: today, order_by: 'cost', limit: 10 }),
  ])
  overview.value = ovRes.data.data
  compare.value = cmpRes.data.data
  trendData.value = trendRes.data.data
  rankList.value = rankRes.data.data
})
</script>
```

- [ ] **Step 2: 提交**

```bash
git add src/views/analytics/Dashboard.vue
git commit -m "feat: add analytics dashboard with cards, trend chart, rank table"
```

---

## Task 7: 批量操作向导页面

**Files:**
- Create: src/views/batch/BatchCreate.vue

- [ ] **Step 1: 创建批量操作向导**

```vue
<!-- src/views/batch/BatchCreate.vue -->
<template>
  <div class="p-6 max-w-4xl mx-auto">
    <h1 class="text-2xl font-bold mb-6">批量操作</h1>

    <!-- 步骤指示器 -->
    <div class="flex mb-8">
      <div v-for="(s, i) in steps" :key="i" class="flex-1 text-center"
           :class="i <= currentStep ? 'text-blue-600 font-semibold' : 'text-gray-400'">
        {{ s }}
      </div>
    </div>

    <!-- Step 1: 选择操作类型 -->
    <div v-if="currentStep === 0" class="space-y-4">
      <div v-for="opt in taskTypes" :key="opt.value"
           class="border rounded-lg p-4 cursor-pointer hover:border-blue-500"
           :class="form.task_type === opt.value ? 'border-blue-500 bg-blue-50' : ''"
           @click="form.task_type = opt.value">
        <p class="font-medium">{{ opt.label }}</p>
        <p class="text-sm text-gray-500">{{ opt.desc }}</p>
      </div>
    </div>

    <!-- Step 2: 选择账户 -->
    <div v-if="currentStep === 1" class="space-y-4">
      <div v-for="acc in accounts" :key="acc.id" class="flex items-center gap-2">
        <input type="checkbox" :value="acc.id" v-model="form.account_ids" />
        <span>{{ acc.name }} ({{ acc.local_account_id }})</span>
      </div>
    </div>

    <!-- Step 3: 配置参数 -->
    <div v-if="currentStep === 2" class="space-y-4">
      <div v-if="form.task_type === 'update_budget'">
        <label class="block text-sm font-medium">预算（元）</label>
        <input type="number" v-model="form.config.budget" class="mt-1 border rounded px-3 py-2 w-full" />
      </div>
      <div v-if="form.task_type === 'update_bid'">
        <label class="block text-sm font-medium">出价（分）</label>
        <input type="number" v-model="form.config.bid" class="mt-1 border rounded px-3 py-2 w-full" />
      </div>
      <div v-if="form.task_type === 'update_status'">
        <label class="block text-sm font-medium">状态</label>
        <select v-model="form.config.opt_status" class="mt-1 border rounded px-3 py-2 w-full">
          <option value="enable">启用</option>
          <option value="disable">暂停</option>
        </select>
      </div>
    </div>

    <!-- Step 4: 确认 -->
    <div v-if="currentStep === 3" class="space-y-2">
      <p>操作类型：{{ form.task_type }}</p>
      <p>影响账户：{{ form.account_ids.length }} 个</p>
      <p>配置：{{ JSON.stringify(form.config) }}</p>
    </div>

    <!-- 导航按钮 -->
    <div class="flex justify-between mt-8">
      <button v-if="currentStep > 0" @click="currentStep--" class="px-4 py-2 border rounded">上一步</button>
      <div v-else></div>
      <button v-if="currentStep < 3" @click="currentStep++" class="px-4 py-2 bg-blue-600 text-white rounded"
              :disabled="!canNext">下一步</button>
      <button v-else @click="submit" class="px-4 py-2 bg-green-600 text-white rounded"
              :disabled="submitting">{{ submitting ? '提交中...' : '确认提交' }}</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { getAccounts } from '@/api/account'
import { createBatchTask } from '@/api/batch'

const router = useRouter()
const steps = ['选择操作', '选择账户', '配置参数', '确认提交']
const currentStep = ref(0)
const submitting = ref(false)
const accounts = ref<any[]>([])

const taskTypes = [
  { value: 'create_project', label: '批量创建项目', desc: '为多个账户创建相同配置的项目' },
  { value: 'update_budget', label: '批量调预算', desc: '批量修改项目日预算' },
  { value: 'update_bid', label: '批量调出价', desc: '批量修改项目出价' },
  { value: 'update_status', label: '批量启停', desc: '批量启用或暂停项目' },
]

const form = reactive({
  task_type: '',
  account_ids: [] as number[],
  config: {} as Record<string, any>,
})

const canNext = computed(() => {
  if (currentStep.value === 0) return form.task_type !== ''
  if (currentStep.value === 1) return form.account_ids.length > 0
  return true
})

async function submit() {
  submitting.value = true
  try {
    await createBatchTask({
      task_type: form.task_type,
      title: taskTypes.find(t => t.value === form.task_type)?.label || '',
      config: form.config,
      account_ids: form.account_ids,
    })
    router.push('/batch/tasks')
  } finally {
    submitting.value = false
  }
}

onMounted(async () => {
  const res = await getAccounts({ page: 1, page_size: 500 })
  accounts.value = res.data.data || []
})
</script>
```

- [ ] **Step 2: 提交**

```bash
git add src/views/batch/BatchCreate.vue
git commit -m "feat: add batch operation wizard with 4-step flow"
```

---

## Task 8: 批量任务进度页面

**Files:**
- Create: src/views/batch/BatchTaskList.vue

- [ ] **Step 1: 创建任务列表页（含轮询进度）**

```vue
<!-- src/views/batch/BatchTaskList.vue -->
<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold mb-4">批量任务</h1>
    <table class="w-full text-sm">
      <thead>
        <tr class="border-b bg-gray-50">
          <th class="text-left py-2 px-3">标题</th>
          <th class="text-left py-2 px-3">类型</th>
          <th class="text-left py-2 px-3">状态</th>
          <th class="text-right py-2 px-3">进度</th>
          <th class="text-right py-2 px-3">操作</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="task in tasks" :key="task.id" class="border-b">
          <td class="py-2 px-3">{{ task.title }}</td>
          <td class="py-2 px-3">{{ task.task_type }}</td>
          <td class="py-2 px-3">
            <span class="px-2 py-1 rounded text-xs" :class="statusClass(task.status)">{{ task.status }}</span>
          </td>
          <td class="text-right py-2 px-3">{{ task.success_count + task.failed_count }}/{{ task.total_count }}</td>
          <td class="text-right py-2 px-3 space-x-2">
            <button v-if="task.status === 'running'" @click="cancel(task.id)" class="text-red-500 text-xs">取消</button>
            <button v-if="task.status === 'partial_success' || task.status === 'failed'" @click="retry(task.id)" class="text-blue-500 text-xs">重试</button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { getBatchTasks, cancelBatchTask, retryBatchTask } from '@/api/batch'

const tasks = ref<any[]>([])
let pollTimer: ReturnType<typeof setInterval> | null = null

function statusClass(status: string) {
  const map: Record<string, string> = {
    running: 'bg-blue-100 text-blue-700',
    completed: 'bg-green-100 text-green-700',
    failed: 'bg-red-100 text-red-700',
    partial_success: 'bg-yellow-100 text-yellow-700',
    cancelled: 'bg-gray-100 text-gray-700',
  }
  return map[status] || 'bg-gray-100'
}

async function loadTasks() {
  const res = await getBatchTasks({ page: 1, page_size: 50 })
  tasks.value = res.data.data || []
}

async function cancel(id: number) {
  await cancelBatchTask(id)
  loadTasks()
}

async function retry(id: number) {
  await retryBatchTask(id)
  loadTasks()
}

onMounted(() => {
  loadTasks()
  pollTimer = setInterval(loadTasks, 2000)
})

onUnmounted(() => {
  if (pollTimer) clearInterval(pollTimer)
})
</script>
```

- [ ] **Step 2: 提交**

```bash
git add src/views/batch/BatchTaskList.vue
git commit -m "feat: add batch task list with 2s polling and cancel/retry"
```

---

## Task 9: Worker 集成 retry

**Files:**
- Modify: backend/internal/app/batch/service/worker.go

- [ ] **Step 1: 在 Worker 的 SDK 调用中使用 retry**

修改 worker.go 中所有 SDK 调用，包裹 retry.Do：

```go
import "oceanengine-backend/pkg/retry"

func (w *Worker) execCreateProject(ctx context.Context, task model.BatchTask, item model.BatchTaskItem) (uint64, string) {
	accessToken, _ := w.tokenStore.GetToken(ctx, task.TenantID)
	var req ocean.CreateProjectRequest
	json.Unmarshal([]byte(task.Config), &req)
	if item.ConfigOverride != "" {
		json.Unmarshal([]byte(item.ConfigOverride), &req)
	}
	req.LocalAccountID = item.LocalAccountID

	var resp *ocean.CreateProjectResponse
	err := retry.Do(ctx, retry.DefaultConfig(), func() error {
		var e error
		resp, e = w.sdk.CreateProject(ctx, accessToken, req)
		return e
	})
	if err != nil {
		return 0, err.Error()
	}
	return resp.Data.ProjectID, ""
}
```

对 execUpdateBudget/execUpdateBid/execUpdateStatus 做同样改造。

- [ ] **Step 2: 编译验证并提交**

```bash
git add internal/app/batch/service/worker.go
git commit -m "feat: integrate retry into batch worker SDK calls"
```

---

## 完成标准

Phase 4 完成后：
1. SDK 调用失败自动重试（指数退避，最多 3 次）
2. 报表查询有 Redis 缓存（5min TTL）
3. 服务重启后自动恢复 running 状态的批量任务
4. /health 检查 MySQL + Redis 连通性
5. 前端 API 层完整覆盖所有后端接口
6. 数据看板页面：汇总卡片 + 趋势图 + 排行榜
7. 批量操作向导：4 步流程（选类型→选账户→配参数→确认）
8. 批量任务列表：2 秒轮询进度 + 取消/重试操作
