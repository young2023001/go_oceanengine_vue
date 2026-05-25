# 本地推 SaaS 前端集成 — 端到端设计方案

## 1. 目标

补全前端，使整个 SaaS 平台可端到端跑通完整业务流程。后端 API 已全部就绪，本方案聚焦前端改造。

---

## 2. 端到端业务流程

### 流程 A：管理员首次配置

```
步骤1: 登录（现有功能）
  前端: LoginView.vue → POST /api/v1/auth/login → 获取 JWT
  
步骤2: 创建租户
  前端: system/TenantManage.vue → tenantApi.create()
  后端: POST /api/v1/tenants
  
步骤3: OAuth 授权
  前端: TenantManage.vue 点击"发起授权" → tenantApi.getOAuthURL(id)
  后端: GET /api/v1/tenants/:id/oauth/url → 返回巨量 OAuth URL
  前端: window.open(url) → 用户在巨量页面授权
  后端: GET /api/v1/tenants/oauth/callback → 自动保存 Token
  前端: 轮询租户状态，token_status 变为 active 表示成功

步骤4: 导入账户
  前端: local/AccountList.vue → accountApi.importAccounts()
  后端: POST /api/v1/accounts/import
  
步骤5: 创建分组
  前端: local/GroupManage.vue → groupApi.create()
  后端: POST /api/v1/groups
  
步骤6: 添加账户到分组
  前端: GroupManage.vue → groupApi.addMembers(groupId, accountIds)
  后端: POST /api/v1/groups/:id/members
  
步骤7: 分配投手权限
  前端: system/ScopeManage.vue → scopeApi.setScope(userId, accountIds)
  后端: POST /api/v1/system/users/:id/scope
```

### 流程 B：投手日常投放

```
步骤1: 登录 → 进入工作台
  前端: local/LocalDashboard.vue（改造）→ analyticsApi.getOverview()
  后端: GET /api/v1/analytics/overview（自动按 scope 过滤）

步骤2: 查看项目列表
  前端: local/ProjectList.vue（改造）→ projectApi.getList()
  后端: GET /api/v1/projects（DataScope 中间件自动过滤）

步骤3: 创建项目模板
  前端: local/TemplateProject.vue → templateApi.createProject()
  后端: POST /api/v1/templates/projects

步骤4: 批量创建项目
  前端: batch/BatchCreate.vue → batchApi.createTask()
  后端: POST /api/v1/batch/projects

步骤5: 查看任务进度
  前端: batch/BatchTaskList.vue → batchApi.getList()（5s 轮询）
  后端: GET /api/v1/batch/tasks

步骤6: 查看数据分析
  前端: analytics/AnalyticsDashboard.vue → analyticsApi.getOverview/getTrend/getRank
  后端: GET /api/v1/analytics/overview + trend + rank

步骤7: 导出报表
  前端: local/AnalyticsExport.vue → analyticsApi.createExport()
  后端: POST /api/v1/analytics/export → 异步生成 → 下载
```

---

## 3. 前后端端点对应表

### 3.1 已有前端 API（需确认页面连通）

| 后端路由 | 前端 API 方法 | 前端页面 | 当前状态 |
|---------|-------------|---------|---------|
| POST /api/v1/tenants | tenantApi.create | TenantManage.vue | 缺页面 |
| GET /api/v1/tenants | tenantApi.getList | TenantManage.vue | 缺页面 |
| GET /api/v1/tenants/:id/oauth/url | tenantApi.getOAuthURL | TenantManage.vue | 缺页面 |
| POST /api/v1/accounts/import | accountApi.importAccounts | AccountList.vue | 缺页面 |
| GET /api/v1/accounts | accountApi.getList | AccountList.vue | 缺页面 |
| GET /api/v1/projects | projectApi.getList | ProjectList.vue | 需改造 |
| GET /api/v1/projects/:id | projectApi.getDetail | ProjectDetail.vue | 需改造 |
| PUT /api/v1/projects/:id/status | projectApi.updateStatus | ProjectList.vue | 需改造 |
| GET /api/v1/promotions | projectApi.getPromotions | PromotionList.vue | 需改造 |
| POST /api/v1/batch/projects | batchApi.createTask | BatchCreate.vue | 缺路由 |
| GET /api/v1/batch/tasks | batchApi.getList | BatchTaskList.vue | 缺路由 |
| GET /api/v1/batch/tasks/:id | batchApi.getDetail | BatchTaskList.vue | 缺路由 |
| POST /api/v1/batch/tasks/:id/cancel | batchApi.cancel | BatchTaskList.vue | 缺路由 |
| POST /api/v1/batch/tasks/:id/retry | batchApi.retry | BatchTaskList.vue | 缺路由 |
| GET /api/v1/analytics/overview | analyticsApi.getOverview | AnalyticsDashboard.vue | 缺路由 |
| GET /api/v1/analytics/trend | analyticsApi.getTrend | AnalyticsDashboard.vue | 缺路由 |
| GET /api/v1/analytics/rank | analyticsApi.getRank | AnalyticsDashboard.vue | 缺路由 |
| GET /api/v1/analytics/compare | analyticsApi.getCompare | AnalyticsDashboard.vue | 缺路由 |
| GET /api/v1/analytics/detail | analyticsApi.getDetail | AnalyticsExport.vue | 缺页面 |
| POST /api/v1/analytics/export | analyticsApi.createExport | AnalyticsExport.vue | 缺页面 |
| GET /api/v1/analytics/exports | analyticsApi.getExports | AnalyticsExport.vue | 缺页面 |

### 3.2 缺失的前端 API 模块

| 文件 | 对应后端路由 | 方法 |
|------|------------|------|
| src/api/group.ts | POST /api/v1/groups | create(data) |
| | GET /api/v1/groups | getList(params) |
| | PUT /api/v1/groups/:id | update(id, data) |
| | DELETE /api/v1/groups/:id | remove(id) |
| | POST /api/v1/groups/:id/members | addMembers(id, accountIds) |
| | DELETE /api/v1/groups/:id/members | removeMembers(id, accountIds) |
| src/api/template.ts | POST /api/v1/templates/projects | createProject(data) |
| | GET /api/v1/templates/projects | listProjects() |
| | GET /api/v1/templates/projects/:id | getProject(id) |
| | PUT /api/v1/templates/projects/:id | updateProject(id, data) |
| | DELETE /api/v1/templates/projects/:id | deleteProject(id) |
| | POST /api/v1/templates/promotions | createPromotion(data) |
| | GET /api/v1/templates/promotions | listPromotions() |
| | DELETE /api/v1/templates/promotions/:id | deletePromotion(id) |
| src/api/scope.ts | POST /api/v1/system/users/:id/scope | setScope(userId, accountIds) |
| | GET /api/v1/system/users/:id/scope | getScope(userId) |

---

## 4. 完整文件改动清单

### 4.1 新建文件（15 个）

| 文件路径 | 功能 |
|---------|------|
| src/api/group.ts | 分组 API |
| src/api/template.ts | 模板 API |
| src/api/scope.ts | 投手权限 API |
| src/views/local/AccountList.vue | 账户列表 + JSON 导入 |
| src/views/local/GroupManage.vue | 分组 CRUD + 成员管理 |
| src/views/local/StoreList.vue | 门店列表（只读） |
| src/views/local/TemplateProject.vue | 项目模板 CRUD |
| src/views/local/TemplatePromotion.vue | 单元模板 CRUD |
| src/views/local/AnalyticsMultiDim.vue | 多维分析（按门店/加盟商/区域/投手） |
| src/views/local/AnalyticsExport.vue | 报表导出（发起+列表+下载） |
| src/views/system/TenantManage.vue | 租户管理 + OAuth 授权 + Token 状态 |
| src/views/system/ScopeManage.vue | 投手权限分配（选用户→选账户） |

### 4.2 改造文件（8 个）

| 文件路径 | 改造内容 |
|---------|---------|
| src/views/local/LocalDashboard.vue | 替换为调 analyticsApi.getOverview/getTrend/getRank |
| src/views/local/ProjectList.vue | 去掉手动 token 输入，改调 projectApi.getList |
| src/views/local/ProjectCreate.vue | 改调新 API（单账户创建项目） |
| src/views/local/PromotionList.vue | 改调 projectApi.getPromotions |
| src/views/local/PromotionCreate.vue | 改调新 API |
| src/views/local/ReportProject.vue | 改调 analyticsApi.getDetail |
| src/router/routes.ts | 注册所有新路由（见 4.3） |
| src/components/layout/AppSidebar.vue | 更新菜单结构（见第 5 节） |

### 4.3 路由注册清单

```typescript
// 本地推 - 投放管理
{ path: 'local', component: LocalDashboard, meta: { title: '工作台' } }
{ path: 'local/project', component: ProjectList, meta: { title: '项目列表' } }
{ path: 'local/project/create', component: ProjectCreate, meta: { title: '创建项目' } }
{ path: 'local/project/:id', component: ProjectDetail, meta: { title: '项目详情' } }
{ path: 'local/promotion', component: PromotionList, meta: { title: '单元列表' } }
{ path: 'local/promotion/create', component: PromotionCreate, meta: { title: '创建单元' } }

// 本地推 - 批量操作
{ path: 'local/batch/create', component: BatchCreate, meta: { title: '批量创建' } }
{ path: 'local/batch/tasks', component: BatchTaskList, meta: { title: '任务记录' } }

// 本地推 - 模板管理
{ path: 'local/template/project', component: TemplateProject, meta: { title: '项目模板' } }
{ path: 'local/template/promotion', component: TemplatePromotion, meta: { title: '单元模板' } }

// 本地推 - 账户管理
{ path: 'local/account', component: AccountList, meta: { title: '账户列表' } }
{ path: 'local/group', component: GroupManage, meta: { title: '分组管理' } }
{ path: 'local/store', component: StoreList, meta: { title: '门店列表' } }

// 本地推 - 数据分析
{ path: 'local/analytics', component: AnalyticsDashboard, meta: { title: '投放看板' } }
{ path: 'local/analytics/multi', component: AnalyticsMultiDim, meta: { title: '多维分析' } }
{ path: 'local/analytics/export', component: AnalyticsExport, meta: { title: '报表导出' } }

// 本地推 - 线索
{ path: 'local/clue', component: ClueList, meta: { title: '线索管理' } }

// 系统设置
{ path: 'system/tenant', component: TenantManage, meta: { title: '租户管理', admin: true } }
{ path: 'system/scope', component: ScopeManage, meta: { title: '投手权限', admin: true } }
```

---

## 5. 侧边栏菜单结构

```typescript
const productLines = [
  {
    id: 'overview',
    name: '概览',
    sections: [{ title: '', items: [
      { name: '工作台', path: '/local' }
    ]}]
  },
  {
    id: 'local',
    name: '本地推',
    sections: [
      { title: '投放管理', items: [
        { name: '项目列表', path: '/local/project' },
        { name: '创建项目', path: '/local/project/create' },
        { name: '单元列表', path: '/local/promotion' },
        { name: '创建单元', path: '/local/promotion/create' },
      ]},
      { title: '批量操作', items: [
        { name: '批量创建', path: '/local/batch/create' },
        { name: '任务记录', path: '/local/batch/tasks' },
      ]},
      { title: '模板管理', items: [
        { name: '项目模板', path: '/local/template/project' },
        { name: '单元模板', path: '/local/template/promotion' },
      ]},
      { title: '账户管理', items: [
        { name: '账户列表', path: '/local/account' },
        { name: '分组管理', path: '/local/group' },
        { name: '门店列表', path: '/local/store' },
      ]},
      { title: '数据分析', items: [
        { name: '投放看板', path: '/local/analytics' },
        { name: '多维分析', path: '/local/analytics/multi' },
        { name: '报表导出', path: '/local/analytics/export' },
      ]},
      { title: '线索管理', items: [
        { name: '线索列表', path: '/local/clue' },
      ]},
    ]
  },
  {
    id: 'tools',
    name: '工具与服务',
    // 保留现有结构不动
  },
  {
    id: 'system',
    name: '系统设置',
    sections: [
      { title: '平台管理', items: [
        { name: '租户管理', path: '/system/tenant' },
        { name: '投手权限', path: '/system/scope' },
      ]},
      { title: '系统管理', items: [
        { name: '用户管理', path: '/system/user' },
        { name: '角色管理', path: '/system/role' },
        { name: '操作日志', path: '/system/log' },
      ]},
    ]
  }
]
```

---

## 6. 各页面核心交互设计

### TenantManage.vue（租户管理）
- 表格：租户名称 | AppID | Token状态(active/expired/failed) | 操作
- 操作：新建 | 发起授权 | 重新授权
- Token 状态用颜色标签：active=绿, expiring=黄, expired=红, failed=红
- "发起授权"点击后 window.open 跳转巨量 OAuth 页面

### AccountList.vue（账户列表）
- 表格：账户ID | 巨量账户ID | 名称 | 余额 | 状态 | 导入时间
- 操作：批量导入（弹窗，粘贴 JSON 或上传文件）
- JSON 格式：`[{"local_account_id": 100001, "name": "门店A"}, ...]`

### GroupManage.vue（分组管理）
- 左侧：分组列表（树形，支持加盟商/区域/自定义类型筛选）
- 右侧：选中分组的成员列表（账户）
- 操作：新建分组 | 编辑 | 删除 | 添加成员 | 移除成员

### TemplateProject.vue（项目模板）
- 表格：模板名称 | 使用次数 | 创建人 | 创建时间 | 操作
- 操作：新建 | 编辑 | 删除 | 使用（跳转批量创建并预填模板）
- 新建/编辑弹窗：名称 + JSON 配置编辑器（marketing_goal, bid_type, bid, budget 等）

### ScopeManage.vue（投手权限）
- 左侧：投手列表（系统用户，role=operator）
- 右侧：选中投手的账户权限范围（checkbox 列表）
- 操作：勾选/取消账户 → 保存

### ProjectList.vue（改造）
- 去掉顶部的 token/advertiser_id 输入框
- 数据自动从 `/api/v1/projects` 获取（后端按 scope 过滤）
- 保留：筛选（状态）、分页、启停操作

### BatchCreate.vue（已有，需接入模板）
- Step 1 增加"从模板创建"选项，选择模板后自动填充配置
- 其余流程不变

### LocalDashboard.vue（改造）
- 替换为调 analyticsApi 的数据看板
- 汇总卡片 + 趋势图 + 排行榜（复用 AnalyticsDashboard.vue 的逻辑）

---

## 7. 权限控制

| 菜单/页面 | 可见条件 |
|----------|---------|
| 系统设置/租户管理 | role_key = admin |
| 系统设置/投手权限 | role_key = admin |
| 本地推/所有页面 | 所有已登录用户（数据按 scope 过滤） |
| 工具与服务 | 所有已登录用户 |

前端实现：侧边栏根据 JWT 中的 role_key 决定是否渲染"系统设置"模块。

---

## 8. 不改动的部分

- 后端 API（全部已实现，无需修改）
- 登录/注册流程
- 工具与服务模块
- 线索管理（保留现有 SDK 透传）
- JWT 认证机制
