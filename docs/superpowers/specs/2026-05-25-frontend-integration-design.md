# 本地推 SaaS 前端集成设计文档

## 1. 目标

将后端 Phase 1-4 实现的 SaaS 能力完整暴露到前端，使用户能跑通完整业务流程：
管理员创建租户 → OAuth 授权 → 导入账户 → 分组 → 分配权限 → 投手操作投放 → 查看数据

## 2. 改造范围

### 2.1 改造现有页面（替换数据源）

| 页面 | 当前逻辑 | 改造为 |
|------|---------|--------|
| `local/ProjectList.vue` | 直接调 SDK（需手动填 token） | 调 `/api/v1/projects`，自动认证 |
| `local/ProjectCreate.vue` | 直接调 SDK | 调 `/api/v1/batch/projects`（单账户模式） |
| `local/PromotionList.vue` | 直接调 SDK | 调 `/api/v1/promotions` |
| `local/PromotionCreate.vue` | 直接调 SDK | 调新 API |
| `local/ReportProject.vue` | 直接调 SDK 报表 | 调 `/api/v1/analytics/detail` |
| `local/LocalDashboard.vue` | 简单概览 | 替换为投放看板（调 analytics API） |

### 2.2 新增页面

| 页面 | 功能 |
|------|------|
| `local/AccountList.vue` | 账户列表 + JSON 批量导入 |
| `local/GroupManage.vue` | 分组 CRUD + 成员管理 |
| `local/StoreList.vue` | 门店列表（只读） |
| `local/TemplateProject.vue` | 项目模板 CRUD |
| `local/TemplatePromotion.vue` | 单元模板 CRUD |
| `local/AnalyticsMultiDim.vue` | 多维分析（按门店/加盟商/区域/投手） |
| `local/AnalyticsExport.vue` | 报表导出（任务列表+下载） |
| `system/TenantManage.vue` | 租户管理 + OAuth 授权 |
| `system/ScopeManage.vue` | 投手权限分配 |

### 2.3 路由 + 导航

- `router/routes.ts` — 注册所有新路由，移除旧的 SDK 透传路由
- `AppSidebar.vue` — 更新为确认的菜单结构

## 3. 菜单结构

```
概览
└── 工作台（投放数据看板）

本地推
├── 投放管理
│   ├── 项目列表
│   ├── 创建项目
│   ├── 单元列表
│   └── 创建单元
├── 批量操作
│   ├── 批量创建
│   └── 任务记录
├── 模板管理
│   ├── 项目模板
│   └── 单元模板
├── 账户管理
│   ├── 账户列表
│   ├── 分组管理
│   └── 门店列表
├── 数据分析
│   ├── 投放看板
│   ├── 多维分析
│   └── 报表导出
└── 线索管理

工具与服务（保留不动）

系统设置（管理员）
├── 租户管理
├── 投手权限
├── 用户管理
├── 角色管理
└── 操作日志
```

## 4. 用户操作流程

### 管理员首次配置
1. 登录 → 系统设置 → 租户管理 → 点击"新建租户"
2. 填写品牌名称 + OAuth AppID + Secret → 保存
3. 点击"发起授权" → 跳转巨量引擎 OAuth 页面 → 用户授权 → 回调成功
4. 本地推 → 账户管理 → 点击"导入账户" → 粘贴/上传 JSON → 确认导入
5. 分组管理 → 创建分组（如"华东区"）→ 添加账户到分组
6. 系统设置 → 投手权限 → 选择投手 → 勾选可操作的账户 → 保存

### 投手日常操作
1. 登录 → 工作台（看到自己权限范围内的数据看板）
2. 投放管理 → 项目列表（只看到自己负责的账户下的项目）
3. 模板管理 → 创建项目模板（配置公共参数）
4. 批量操作 → 批量创建 → 选模板 → 选账户 → 差异化配置 → 提交
5. 任务记录 → 查看进度 → 失败项重试
6. 数据分析 → 投放看板 / 多维分析 / 导出

### OAuth 授权提示
- 租户管理页面显示 Token 状态（active/expiring/expired/failed）
- expired 状态时显示红色警告 + "重新授权"按钮
- 授权成功后自动刷新状态

## 5. 技术方案

- 所有新页面使用 Vue 3 Composition API + TypeScript + TailwindCSS
- API 调用使用已创建的 `src/api/` 模块（tenant/account/project/batch/analytics）
- 新增 `src/api/group.ts` 和 `src/api/template.ts`
- 认证通过现有 JWT 机制，token 存 localStorage
- 权限控制：侧边栏根据 role_key 显示/隐藏管理员菜单

## 6. 不改动的部分

- 工具与服务模块（代理商/财务/广告资产/素材/基础工具）
- 登录页（已有 JWT 认证）
- 系统管理中的用户/角色/菜单/日志（复用现有）
- 线索管理（保留现有 SDK 透传，后续迭代）
