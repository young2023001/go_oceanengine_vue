# 前端集成实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 补全前端 12 个新页面 + 改造 9 个现有页面 + 路由/导航，使 SaaS 平台端到端可跑通。

**Architecture:** 新增 3 个 API 模块（group/template/scope），12 个 Vue 页面，改造 6 个现有页面数据源，更新路由和侧边栏。所有页面使用 Vue 3 Composition API + TypeScript + TailwindCSS，API 调用使用已有的 request.ts 封装。

**Tech Stack:** Vue 3 + TypeScript + Vite + TailwindCSS + vue-chartjs + axios

---

## 文件结构

### 新建文件（12 个）

| 路径 | 职责 |
|------|------|
| `frontend/src/api/group.ts` | 分组 API |
| `frontend/src/api/template.ts` | 模板 API |
| `frontend/src/api/scope.ts` | 投手权限 API |
| `frontend/src/views/local/AccountList.vue` | 账户列表 + JSON 导入 |
| `frontend/src/views/local/GroupManage.vue` | 分组 CRUD + 成员管理 |
| `frontend/src/views/local/StoreList.vue` | 门店列表 |
| `frontend/src/views/local/TemplateProject.vue` | 项目模板 CRUD |
| `frontend/src/views/local/TemplatePromotion.vue` | 单元模板 CRUD |
| `frontend/src/views/local/AnalyticsMultiDim.vue` | 多维分析 |
| `frontend/src/views/local/AnalyticsExport.vue` | 报表导出 |
| `frontend/src/views/system/TenantManage.vue` | 租户管理 + OAuth |
| `frontend/src/views/system/ScopeManage.vue` | 投手权限分配 |

### 改造文件（9 个）

| 路径 | 改造内容 |
|------|---------|
| `frontend/src/views/local/LocalDashboard.vue` | 替换为 analytics API 数据看板 |
| `frontend/src/views/local/ProjectList.vue` | 去掉手动 token，改调 projectApi |
| `frontend/src/views/local/ProjectCreate.vue` | 改调新 API |
| `frontend/src/views/local/ProjectDetail.vue` | 改调 projectApi.getProject |
| `frontend/src/views/local/PromotionList.vue` | 改调 projectApi.getPromotions |
| `frontend/src/views/local/PromotionCreate.vue` | 改调新 API |
| `frontend/src/views/local/ReportProject.vue` | 改调 analyticsApi.getDetail |
| `frontend/src/router/routes.ts` | 注册所有新路由 |
| `frontend/src/components/layout/AppSidebar.vue` | 更新菜单结构 |

---

## Task 1: 新增 API 模块（group.ts / template.ts / scope.ts） ✅ (已完成 commit ba01b38)

> **补充：** 已创建的 account.ts 中 Account 接口需要补充 balance 和 status 字段：
> ```typescript
> export interface Account {
>   id: number
>   local_account_id: number
>   name: string
>   balance?: number  // 补充
>   status?: string   // 补充
>   created_at: string
> }
> ```

**Files:**
- Create: `frontend/src/api/group.ts`
- Create: `frontend/src/api/template.ts`
- Create: `frontend/src/api/scope.ts`

- [ ] **Step 1: 创建 group.ts**

```typescript
// frontend/src/api/group.ts
import request from './request'

export interface Group {
  id: number
  name: string
  type: 'franchisee' | 'region' | 'custom'
  parent_id: number
  sort: number
  created_at: string
}

export interface CreateGroupParams {
  name: string
  type: 'franchisee' | 'region' | 'custom'
  parent_id?: number
  sort?: number
}

export interface UpdateGroupParams {
  name?: string
  sort?: number
}

export interface MembersParams {
  account_ids: number[]
}

export const groupApi = {
  create(data: CreateGroupParams) {
    return request.post<Group>('/groups', data)
  },

  getList(params?: { type?: string }) {
    return request.get<Group[]>('/groups', params)
  },

  update(id: number, data: UpdateGroupParams) {
    return request.put<void>(`/groups/${id}`, data)
  },

  remove(id: number) {
    return request.delete<void>(`/groups/${id}`)
  },

  addMembers(id: number, data: MembersParams) {
    return request.post<void>(`/groups/${id}/members`, data)
  },

  removeMembers(id: number, data: MembersParams) {
    // 注意：后端 DELETE /groups/:id/members 需要 JSON body
    // 但前端 request.delete 不支持 body，需要后端新增 POST 端点
    // 临时方案：直接用 axios instance 发 DELETE + body
    return request.post<void>(`/groups/${id}/members/delete`, data)
  }
}
```

- [ ] **Step 2: 创建 template.ts**

```typescript
// frontend/src/api/template.ts
import request from './request'

export interface ProjectTemplate {
  id: number
  name: string
  config: Record<string, unknown>
  use_count: number
  created_by: number
  created_at: string
}

export interface PromotionTemplate {
  id: number
  name: string
  config: Record<string, unknown>
  use_count: number
  created_by: number
  created_at: string
}

export interface CreateTemplateParams {
  name: string
  config: Record<string, unknown>
}

export const templateApi = {
  createProject(data: CreateTemplateParams) {
    return request.post<ProjectTemplate>('/templates/projects', data)
  },

  listProjects() {
    return request.get<ProjectTemplate[]>('/templates/projects')
  },

  getProject(id: number) {
    return request.get<ProjectTemplate>(`/templates/projects/${id}`)
  },

  updateProject(id: number, data: CreateTemplateParams) {
    return request.put<void>(`/templates/projects/${id}`, data)
  },

  deleteProject(id: number) {
    return request.delete<void>(`/templates/projects/${id}`)
  },

  createPromotion(data: CreateTemplateParams) {
    return request.post<PromotionTemplate>('/templates/promotions', data)
  },

  listPromotions() {
    return request.get<PromotionTemplate[]>('/templates/promotions')
  },

  deletePromotion(id: number) {
    return request.delete<void>(`/templates/promotions/${id}`)
  }
}
```

- [ ] **Step 3: 创建 scope.ts**

```typescript
// frontend/src/api/scope.ts
import request from './request'

export interface UserScope {
  account_ids: number[]
}

export const scopeApi = {
  getScope(userId: number) {
    return request.get<UserScope>(`/system/users/${userId}/scope`)
  },

  setScope(userId: number, data: { account_ids: number[] }) {
    return request.post<void>(`/system/users/${userId}/scope`, data)
  }
}
```

- [ ] **Step 4: 提交**

```bash
git add frontend/src/api/group.ts frontend/src/api/template.ts frontend/src/api/scope.ts
git commit -m "feat: add group, template, scope API modules"
```

---

## Task 2: 更新路由注册 ✅ (已完成 commit 0fd22aa)

**Files:**
- Modify: `frontend/src/router/routes.ts`

- [ ] **Step 1: 读取现有 routes.ts，在 protectedRoutes 的 children 中替换/添加本地推相关路由**

替换现有的 local 相关路由为：

```typescript
// 工作台
{
  path: '',
  name: 'Dashboard',
  component: () => import('@/views/local/LocalDashboard.vue'),
  meta: { title: '工作台' }
},
// 本地推 - 投放管理
{
  path: 'local/project',
  name: 'LocalProjectList',
  component: () => import('@/views/local/ProjectList.vue'),
  meta: { title: '项目列表' }
},
{
  path: 'local/project/create',
  name: 'LocalProjectCreate',
  component: () => import('@/views/local/ProjectCreate.vue'),
  meta: { title: '创建项目' }
},
{
  path: 'local/project/:id',
  name: 'LocalProjectDetail',
  component: () => import('@/views/local/ProjectDetail.vue'),
  meta: { title: '项目详情' }
},
{
  path: 'local/promotion',
  name: 'LocalPromotionList',
  component: () => import('@/views/local/PromotionList.vue'),
  meta: { title: '单元列表' }
},
{
  path: 'local/promotion/create',
  name: 'LocalPromotionCreate',
  component: () => import('@/views/local/PromotionCreate.vue'),
  meta: { title: '创建单元' }
},
{
  path: 'local/promotion/:id',
  name: 'LocalPromotionDetail',
  component: () => import('@/views/local/PromotionDetail.vue'),
  meta: { title: '单元详情' }
},
// 本地推 - 批量操作
{
  path: 'local/batch/create',
  name: 'BatchCreate',
  component: () => import('@/views/batch/BatchCreate.vue'),
  meta: { title: '批量创建' }
},
{
  path: 'local/batch/tasks',
  name: 'BatchTaskList',
  component: () => import('@/views/batch/BatchTaskList.vue'),
  meta: { title: '任务记录' }
},
// 本地推 - 模板管理
{
  path: 'local/template/project',
  name: 'TemplateProject',
  component: () => import('@/views/local/TemplateProject.vue'),
  meta: { title: '项目模板' }
},
{
  path: 'local/template/promotion',
  name: 'TemplatePromotion',
  component: () => import('@/views/local/TemplatePromotion.vue'),
  meta: { title: '单元模板' }
},
// 本地推 - 账户管理
{
  path: 'local/account',
  name: 'AccountList',
  component: () => import('@/views/local/AccountList.vue'),
  meta: { title: '账户列表' }
},
{
  path: 'local/group',
  name: 'GroupManage',
  component: () => import('@/views/local/GroupManage.vue'),
  meta: { title: '分组管理' }
},
{
  path: 'local/store',
  name: 'StoreList',
  component: () => import('@/views/local/StoreList.vue'),
  meta: { title: '门店列表' }
},
// 本地推 - 数据分析
{
  path: 'local/analytics',
  name: 'AnalyticsDashboard',
  component: () => import('@/views/analytics/AnalyticsDashboard.vue'),
  meta: { title: '投放看板' }
},
{
  path: 'local/analytics/multi',
  name: 'AnalyticsMultiDim',
  component: () => import('@/views/local/AnalyticsMultiDim.vue'),
  meta: { title: '多维分析' }
},
{
  path: 'local/analytics/export',
  name: 'AnalyticsExport',
  component: () => import('@/views/local/AnalyticsExport.vue'),
  meta: { title: '报表导出' }
},
// 本地推 - 线索
{
  path: 'local/clue',
  name: 'LocalClue',
  component: () => import('@/views/local/ClueList.vue'),
  meta: { title: '线索管理' }
},
{
  path: 'local/clue/:id',
  name: 'LocalClueDetail',
  component: () => import('@/views/local/ClueDetail.vue'),
  meta: { title: '线索详情' }
},
// 系统设置
{
  path: 'system/tenant',
  name: 'TenantManage',
  component: () => import('@/views/system/TenantManage.vue'),
  meta: { title: '租户管理', requiresAdmin: true }
},
{
  path: 'system/scope',
  name: 'ScopeManage',
  component: () => import('@/views/system/ScopeManage.vue'),
  meta: { title: '投手权限', requiresAdmin: true }
},
```

保留现有的 system/user、system/role、system/log 等路由不动。
保留工具与服务模块的所有路由不动。

- [ ] **Step 2: 验证无 TypeScript 错误**

Run: `cd frontend && npx vue-tsc --noEmit 2>&1 | head -20`
如果有类型错误，修复后继续。

- [ ] **Step 3: 提交**

```bash
git add frontend/src/router/routes.ts
git commit -m "feat: register all new SaaS routes"
```

---

## Task 3: 更新侧边栏导航 ✅ (已完成 commit 65ca347)

**Files:**
- Modify: `frontend/src/components/layout/AppSidebar.vue`

- [ ] **Step 1: 替换 productLines 数组**

读取 AppSidebar.vue，将 productLines 数组替换为 spec 第 5 节定义的结构。关键变化：
- 概览模块：工作台 path 改为 '/'
- 本地推模块：重新组织为 投放管理/批量操作/模板管理/账户管理/数据分析/线索管理
- 删除企业号模块
- 新增系统设置模块（含租户管理和投手权限）
- 工具与服务保留不动

同时添加权限控制：从 localStorage 读取 JWT 解析 role_key，系统设置模块仅 admin 可见。

```typescript
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

// 从 localStorage 获取用户角色
const getUserRole = (): string => {
  const token = localStorage.getItem('access_token')
  if (!token) return ''
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    return payload.role_key || ''
  } catch {
    return ''
  }
}

const isAdmin = computed(() => {
  const role = getUserRole()
  return role === 'admin' || role === 'super_admin'
})
```

在 template 中，system 模块用 v-if 控制显示：
```html
<div v-for="product in productLines" :key="product.id" class="mb-1"
     v-show="product.id !== 'system' || isAdmin">
```

- [ ] **Step 2: 验证页面渲染正常**

Run: `cd frontend && npm run dev`
在浏览器中确认侧边栏正确显示新菜单结构。

- [ ] **Step 3: 提交**

```bash
git add frontend/src/components/layout/AppSidebar.vue
git commit -m "feat: update sidebar navigation for SaaS modules"
```

---

## Task 4: 新建 TenantManage.vue

**Files:**
- Create: `frontend/src/views/system/TenantManage.vue`

- [ ] **Step 1: 创建租户管理页面**

```vue
<!-- frontend/src/views/system/TenantManage.vue -->
<template>
  <div class="p-6">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold">租户管理</h1>
      <button @click="showCreateDialog = true" class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700">
        新建租户
      </button>
    </div>

    <div class="bg-white rounded-lg shadow">
      <table class="min-w-full">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">租户名称</th>
            <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">AppID</th>
            <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">Token 状态</th>
            <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">操作</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200">
          <tr v-for="tenant in tenants" :key="tenant.id">
            <td class="px-4 py-3 font-medium">{{ tenant.name }}</td>
            <td class="px-4 py-3 text-sm text-gray-500">{{ tenant.oauth_app_id }}</td>
            <td class="px-4 py-3">
              <span :class="statusClass(tenant.token_status)" class="px-2 py-1 text-xs rounded-full">
                {{ statusLabel(tenant.token_status) }}
              </span>
            </td>
            <td class="px-4 py-3 space-x-2">
              <button v-if="tenant.token_status === 'expired' || tenant.token_status === 'failed'"
                @click="authorize(tenant)" class="text-blue-600 text-sm hover:underline">
                {{ tenant.token_status === 'expired' ? '发起授权' : '重新授权' }}
              </button>
              <span v-else class="text-green-600 text-sm">已授权</span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- 新建弹窗 -->
    <div v-if="showCreateDialog" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg p-6 w-96">
        <h2 class="text-lg font-bold mb-4">新建租户</h2>
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium mb-1">品牌名称</label>
            <input v-model="form.name" class="w-full border rounded-lg px-3 py-2" placeholder="输入品牌名称" />
          </div>
          <div>
            <label class="block text-sm font-medium mb-1">OAuth AppID</label>
            <input v-model="form.oauth_app_id" class="w-full border rounded-lg px-3 py-2" placeholder="巨量引擎应用 AppID" />
          </div>
          <div>
            <label class="block text-sm font-medium mb-1">OAuth Secret</label>
            <input v-model="form.oauth_secret" type="password" class="w-full border rounded-lg px-3 py-2" placeholder="应用 Secret" />
          </div>
        </div>
        <div class="flex justify-end gap-2 mt-6">
          <button @click="showCreateDialog = false" class="px-4 py-2 border rounded-lg">取消</button>
          <button @click="createTenant" class="px-4 py-2 bg-blue-600 text-white rounded-lg">确认</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { tenantApi, type Tenant } from '@/api/tenant'

const tenants = ref<Tenant[]>([])
const showCreateDialog = ref(false)
const form = ref({ name: '', oauth_app_id: '', oauth_secret: '' })

const statusClass = (status: string) => ({
  'bg-green-100 text-green-800': status === 'active',
  'bg-yellow-100 text-yellow-800': status === 'expiring',
  'bg-red-100 text-red-800': status === 'expired' || status === 'failed',
}[status] || 'bg-gray-100 text-gray-800')

const statusLabel = (status: string) => ({
  active: '正常', expiring: '即将过期', expired: '已过期', failed: '授权失败'
}[status] || status)

async function loadTenants() {
  const res = await tenantApi.getList()
  tenants.value = Array.isArray(res) ? res : []
}

async function createTenant() {
  await tenantApi.create(form.value)
  showCreateDialog.value = false
  form.value = { name: '', oauth_app_id: '', oauth_secret: '' }
  loadTenants()
}

async function authorize(tenant: Tenant) {
  const res = await tenantApi.getOAuthURL(tenant.id)
  const url = (res as any).auth_url || (res as any).url
  if (url) {
    window.open(url, '_blank')
  }
}

onMounted(loadTenants)
</script>
```

- [ ] **Step 2: 提交**

```bash
git add frontend/src/views/system/TenantManage.vue
git commit -m "feat: add tenant management page with OAuth authorization"
```

---

## Task 5: 新建 AccountList.vue

**Files:**
- Create: `frontend/src/views/local/AccountList.vue`

- [ ] **Step 1: 创建账户列表页面（含 JSON 导入）**

```vue
<!-- frontend/src/views/local/AccountList.vue -->
<template>
  <div class="p-6">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold">账户列表</h1>
      <button @click="showImportDialog = true" class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700">
        批量导入
      </button>
    </div>

    <div class="bg-white rounded-lg shadow">
      <table class="min-w-full">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">ID</th>
            <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">巨量账户ID</th>
            <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">名称</th>
            <th class="px-4 py-3 text-right text-sm font-medium text-gray-500">余额</th>
            <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">状态</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200">
          <tr v-for="acc in accounts" :key="acc.id">
            <td class="px-4 py-3 text-sm">{{ acc.id }}</td>
            <td class="px-4 py-3 text-sm">{{ acc.local_account_id }}</td>
            <td class="px-4 py-3 font-medium">{{ acc.name }}</td>
            <td class="px-4 py-3 text-right text-sm">{{ acc.balance?.toFixed(2) || '0.00' }}</td>
            <td class="px-4 py-3">
              <span class="px-2 py-1 text-xs rounded-full bg-green-100 text-green-800">{{ acc.status || 'active' }}</span>
            </td>
          </tr>
          <tr v-if="accounts.length === 0">
            <td colspan="5" class="px-4 py-8 text-center text-gray-400">暂无账户，请点击批量导入</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- 导入弹窗 -->
    <div v-if="showImportDialog" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg p-6 w-[500px]">
        <h2 class="text-lg font-bold mb-4">批量导入账户</h2>
        <p class="text-sm text-gray-500 mb-2">JSON 格式：[{"local_account_id": 100001, "name": "门店A"}, ...]</p>
        <textarea v-model="importJson" rows="8" class="w-full border rounded-lg px-3 py-2 font-mono text-sm"
          placeholder='[{"local_account_id": 100001, "name": "门店A"}]'></textarea>
        <p v-if="importError" class="text-red-500 text-sm mt-2">{{ importError }}</p>
        <div class="flex justify-end gap-2 mt-4">
          <button @click="showImportDialog = false" class="px-4 py-2 border rounded-lg">取消</button>
          <button @click="handleImport" :disabled="importing" class="px-4 py-2 bg-blue-600 text-white rounded-lg">
            {{ importing ? '导入中...' : '确认导入' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { accountApi, type Account } from '@/api/account'
// 注意：Account 接口需要补充 balance 和 status 字段（在 Task 1 执行时已创建的 account.ts 中添加）

const accounts = ref<Account[]>([])
const showImportDialog = ref(false)
const importJson = ref('')
const importError = ref('')
const importing = ref(false)

async function loadAccounts() {
  const res = await accountApi.getList()
  accounts.value = res.list || res as any || []
}

async function handleImport() {
  importError.value = ''
  try {
    const items = JSON.parse(importJson.value)
    if (!Array.isArray(items) || items.length === 0) {
      importError.value = 'JSON 必须是非空数组'
      return
    }
    importing.value = true
    await accountApi.importAccounts({ items })
    showImportDialog.value = false
    importJson.value = ''
    loadAccounts()
  } catch (e: any) {
    importError.value = e.message || '导入失败'
  } finally {
    importing.value = false
  }
}

onMounted(loadAccounts)
</script>
```

- [ ] **Step 2: 提交**

```bash
git add frontend/src/views/local/AccountList.vue
git commit -m "feat: add account list page with JSON import"
```

---

> **后端适配说明：** group.ts 的 removeMembers 使用 POST /groups/:id/members/delete，
> 需要在后端 router.go 中新增此路由（复用现有 RemoveMembers handler）。
> 或者修改前端 request.ts 的 delete 方法支持 body。实施时二选一。

## Task 6: 新建 GroupManage.vue

**Files:**
- Create: `frontend/src/views/local/GroupManage.vue`

- [ ] **Step 1: 创建分组管理页面**

左右布局：左侧分组列表（支持类型筛选），右侧选中分组的成员账户列表。

核心逻辑：
- `groupApi.getList({ type })` 获取分组列表
- `groupApi.create({ name, type })` 新建分组
- `groupApi.update(id, { name })` 编辑分组
- `groupApi.remove(id)` 删除分组
- `groupApi.addMembers(id, { account_ids })` 添加成员
- `groupApi.removeMembers(id, { account_ids })` 移除成员
- `accountApi.getList()` 获取可选账户列表

参考现有页面 `ProjectList.vue` 的表格和弹窗模式。

- [ ] **Step 2: 提交**

```bash
git add frontend/src/views/local/GroupManage.vue
git commit -m "feat: add group management page"
```

---

## Task 7: 新建 StoreList.vue

**Files:**
- Create: `frontend/src/views/local/StoreList.vue`

- [ ] **Step 1: 创建门店列表页面**

只读表格：门店名称 | 所属账户 | 城市 | 区域 | 地址 | POI ID。
数据来源：暂时从 `accountApi.getList()` 获取（门店信息附在账户上），或显示空状态提示"门店数据将在账户导入时自动关联"。

- [ ] **Step 2: 提交**

```bash
git add frontend/src/views/local/StoreList.vue
git commit -m "feat: add store list page"
```

---

## Task 8: 新建 TemplateProject.vue + TemplatePromotion.vue

**Files:**
- Create: `frontend/src/views/local/TemplateProject.vue`
- Create: `frontend/src/views/local/TemplatePromotion.vue`

- [ ] **Step 1: 创建项目模板页面**

表格：模板名称 | 使用次数 | 创建时间 | 操作（编辑/删除/使用）。
新建/编辑弹窗：名称输入 + JSON 配置 textarea（marketing_goal, bid_type, bid, budget 等字段）。
"使用"按钮跳转到 `/local/batch/create?template_id=xxx`。

核心 API：
- `templateApi.listProjects()` 列表
- `templateApi.createProject({ name, config })` 新建
- `templateApi.updateProject(id, { name, config })` 编辑
- `templateApi.deleteProject(id)` 删除

- [ ] **Step 2: 创建单元模板页面**

与项目模板结构相同，调用 `templateApi.listPromotions()` / `createPromotion` / `deletePromotion`。

- [ ] **Step 3: 提交**

```bash
git add frontend/src/views/local/TemplateProject.vue frontend/src/views/local/TemplatePromotion.vue
git commit -m "feat: add project and promotion template pages"
```

---

## Task 9: 新建 ScopeManage.vue

**Files:**
- Create: `frontend/src/views/system/ScopeManage.vue`

- [ ] **Step 1: 创建投手权限页面**

左右布局：
- 左侧：用户列表（从现有 system API 获取，筛选 role=operator）
- 右侧：选中用户的账户权限范围（checkbox 列表）

核心逻辑：
- 获取用户列表：`request.get('/system/users', { role: 'operator' })`
- 获取用户权限：`scopeApi.getScope(userId)` → 返回 account_ids
- 保存权限：`scopeApi.setScope(userId, { account_ids: [...] })`
- 获取所有账户：`accountApi.getList()` 用于 checkbox 选项

- [ ] **Step 2: 提交**

```bash
git add frontend/src/views/system/ScopeManage.vue
git commit -m "feat: add scope management page for operator permissions"
```

---

## Task 10: 新建 AnalyticsMultiDim.vue

**Files:**
- Create: `frontend/src/views/local/AnalyticsMultiDim.vue`

- [ ] **Step 1: 创建多维分析页面**

顶部：日期范围选择器 + 维度切换下拉（按门店/按加盟商/按区域/按投手）。
表格：维度名称 | 消耗 | 展示 | 点击 | 转化 | 转化成本。

核心 API：
- `analyticsApi.getRank({ start_date, end_date, order_by: 'cost', limit: 50 })`
- 维度切换实现：后端 getRank 按 account_id 聚合，前端通过不同的展示方式区分维度：
  - 按门店：直接显示 getRank 结果（account = 门店）
  - 按加盟商/区域：调用 groupApi.getList 获取分组，前端按分组聚合显示
  - 按投手：需要后端新增 /analytics/rank-by-user 端点（标注为后续迭代）
- 当前版本先实现按门店维度，其他维度显示即将上线

- [ ] **Step 2: 提交**

```bash
git add frontend/src/views/local/AnalyticsMultiDim.vue
git commit -m "feat: add multi-dimension analytics page"
```

---

## Task 11: 新建 AnalyticsExport.vue

**Files:**
- Create: `frontend/src/views/local/AnalyticsExport.vue`

- [ ] **Step 1: 创建报表导出页面**

上半部分：导出表单（标题 + 日期范围 + 提交按钮）。
下半部分：导出任务列表（标题 | 状态 | 行数 | 创建时间 | 操作）。
状态为 completed 时显示"下载"链接。

核心 API：
- `analyticsApi.createExport({ title, start_date, end_date })` 发起导出
- `analyticsApi.getExports()` 获取任务列表
- 下载：`window.open('/api/v1/analytics/exports/' + id + '/download')`

- [ ] **Step 2: 提交**

```bash
git add frontend/src/views/local/AnalyticsExport.vue
git commit -m "feat: add analytics export page"
```

---

## Task 12: 改造 LocalDashboard.vue

**Files:**
- Modify: `frontend/src/views/local/LocalDashboard.vue`

- [ ] **Step 1: 替换为 analytics 数据看板**

读取现有 LocalDashboard.vue，将内容替换为：
- 汇总卡片（今日消耗/展示/点击/转化）→ `analyticsApi.getOverview()`
- 趋势图（最近 7 天）→ `analyticsApi.getTrend({ start_date, end_date })`
- 排行榜（TOP 10）→ `analyticsApi.getRank({ start_date, end_date, limit: 10 })`
- 环比数据 → `analyticsApi.getCompare({ date, compare_type: 'yesterday' })`

可复用已有的 `AnalyticsDashboard.vue` 的逻辑，或直接将其内容移入。

- [ ] **Step 2: 提交**

```bash
git add frontend/src/views/local/LocalDashboard.vue
git commit -m "feat: replace local dashboard with analytics data"
```

---

## Task 13: 改造 ProjectList.vue + ProjectCreate.vue

**Files:**
- Modify: `frontend/src/views/local/ProjectList.vue`
- Modify: `frontend/src/views/local/ProjectCreate.vue`

- [ ] **Step 1: 改造 ProjectList.vue**

读取现有文件，做以下修改：
1. 删除顶部的 token/advertiser_id 输入框
2. 将 `localApi.getProjectList(token, advertiserId, ...)` 替换为 `projectApi.getList({ status, page, page_size })`
3. 将状态操作 `localApi.updateProjectStatus(token, ...)` 替换为 `projectApi.updateStatus(id, status)`
4. 数据自动从新 API 获取，后端按 DataScope 过滤

- [ ] **Step 2: 改造 ProjectCreate.vue**

将 SDK 透传调用替换为新 API：
- 通过 `batchApi.createTask({ task_type: 'create_project', account_ids: [selectedAccountId], config: formData })`（后端无单独的 project create API，统一走批量任务）

- [ ] **Step 3: 提交**

```bash
git add frontend/src/views/local/ProjectList.vue frontend/src/views/local/ProjectCreate.vue
git commit -m "feat: refactor project pages to use new SaaS API"
```

---

## Task 14: 改造 PromotionList.vue + PromotionCreate.vue

**Files:**
- Modify: `frontend/src/views/local/PromotionList.vue`
- Modify: `frontend/src/views/local/PromotionCreate.vue`

- [ ] **Step 1: 改造 PromotionList.vue**

同 ProjectList 改造模式：
1. 删除手动 token 输入
2. 替换为 `projectApi.getPromotions({ project_id, page, page_size })`

- [ ] **Step 2: 改造 PromotionCreate.vue**

替换 SDK 调用为新 API。

- [ ] **Step 3: 提交**

```bash
git add frontend/src/views/local/PromotionList.vue frontend/src/views/local/PromotionCreate.vue
git commit -m "feat: refactor promotion pages to use new SaaS API"
```

---

## Task 15: 改造 BatchCreate.vue（接入模板选择）

**Files:**
- Modify: `frontend/src/views/batch/BatchCreate.vue`

- [ ] **Step 1: 在 Step 1 增加模板选择**

读取现有 BatchCreate.vue，在操作类型选择之后增加：
- "从模板创建"选项：显示模板列表（`templateApi.listProjects()`）
- 选择模板后自动填充 config 字段
- 支持 URL 参数 `?template_id=xxx` 预选模板

- [ ] **Step 2: 提交**

```bash
git add frontend/src/views/batch/BatchCreate.vue
git commit -m "feat: integrate template selection into batch create wizard"
```

---

## 完成标准

全部 15 个 Task 完成后，用户可以：
1. 管理员登录 → 系统设置 → 创建租户 → OAuth 授权
2. 管理员 → 账户管理 → 批量导入账户 → 创建分组 → 分配投手权限
3. 投手登录 → 工作台看到数据看板
4. 投手 → 项目列表看到自己权限范围内的项目
5. 投手 → 模板管理 → 创建模板 → 批量操作 → 选模板创建项目
6. 投手 → 任务记录 → 查看进度/取消/重试
7. 投手 → 数据分析 → 多维分析 → 导出报表

---

## 依赖关系

```
Task 1 (API 模块) ──→ Task 4-11 (新页面都依赖 API)
Task 2 (路由) ──→ 所有页面可访问
Task 3 (侧边栏) ──→ 所有页面可导航
Task 4-11 (新页面) ──→ 互相独立，可并行
Task 12-15 (改造页面) ──→ 依赖 Task 1 (API 模块)
```
