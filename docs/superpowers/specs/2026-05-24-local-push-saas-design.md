# 本地推投流管理 SaaS 平台 — 系统设计文档

## 1. 项目概述

### 1.1 业务背景

- **平台方（你）：** 技术公司，提供 SaaS 平台
- **客户：** 品牌方/流量运营公司，帮线下门店进行本地推广告投流
- **账户结构：** 品牌在巨量引擎有一个组织，组织下创建多个本地推账户（local_account_id），每个账户关联一个门店（业务规则隔离，非巨量强制）
- **规模：** 单个品牌客户可能有数百个账户/门店
- **分组需求：** 按加盟商/区域等维度分组分析（平台自定义，巨量无此概念）
- **权限需求：** 投手按账户范围隔离，能看到的就能操作

### 1.2 核心目标

1. **投放效率** — 批量创建项目/单元、批量调预算/出价/启停，解决一人管 100+ 账户的效率问题
2. **数据分析** — 按门店/加盟商/区域/投手多维聚合分析，支持环比，指导优化决策
3. **权限管理** — 多投手协作，按账户范围隔离

### 1.3 系统边界

- 本期做：本地推投放管理 + 数据分析 + 账户/权限体系
- 后续扩展：抖音来客对接、线索管理、智能分析引擎

### 1.4 架构选型

多租户 SaaS，单体架构（Go/Gin），未来按需拆分。

---

## 2. 巨量本地推账户层级

```
品牌组织 (Organization)
  │  一次 OAuth 授权 → 一个 access_token
  │
  ├── 账户 A (local_account_id) = 门店 A
  │     ├── 项目 1 (project_id)
  │     │     ├── 单元 1 (promotion_id)
  │     │     └── 单元 2
  │     └── 项目 2
  │
  ├── 账户 B (local_account_id) = 门店 B
  │     └── ...
  │
  └── 账户 C ... (数百个)
```

**关键约束：**
- 所有 SDK API 以 `local_account_id` 为操作入口
- 一个品牌一次 OAuth 授权，Token 可操作该组织下所有账户
- 产品设计上需兼容"一个账户对应多个门店"的情况（虽然客户当前 1:1）

---

## 3. 整体架构

```
前端 (Vue 3 + Vite + TailwindCSS)
    │ HTTP API
主服务 (Go/Gin 单体)
    │
    ├── 投放管理（项目/单元 CRUD + 批量操作）
    ├── 数据分析（报表同步 + 多维聚合）
    ├── 账户体系（租户/账户/分组/权限）
    ├── 基础设施（限流器/Token管理/定时任务）
    │
    ├── MySQL（业务数据）
    ├── Redis（限流/缓存）
    └── 巨量引擎 API（通过官方 SDK 调用）
```

**为什么单体：**
- 当前只有本地推一个业务域
- 一个品牌的 API 调用量在 10 QPS 限制内
- 定时任务作为主服务内 goroutine 运行
- 部署简单，未来加来客等业务时再按需拆分

---

## 4. 数据模型

### 4.1 租户与账户

```sql
CREATE TABLE tenant (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(200) NOT NULL COMMENT '品牌名称',
    oauth_app_id VARCHAR(100) NOT NULL COMMENT '巨量应用 AppID',
    oauth_secret VARCHAR(200) NOT NULL COMMENT '巨量应用 Secret（加密存储）',
    access_token VARCHAR(500),
    refresh_token VARCHAR(500),
    token_expire_at DATETIME,
    token_status ENUM('active','expiring','expired','failed') DEFAULT 'active',
    organization_id BIGINT UNSIGNED DEFAULT 0 COMMENT '巨量组织ID',
    status TINYINT DEFAULT 1,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    deleted_at DATETIME DEFAULT NULL,
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE local_account (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_id BIGINT UNSIGNED NOT NULL,
    local_account_id BIGINT UNSIGNED NOT NULL COMMENT '巨量本地推账户ID',
    name VARCHAR(200) NOT NULL,
    status VARCHAR(50) DEFAULT 'active',
    balance DECIMAL(12,2) DEFAULT 0 COMMENT '账户余额（元）',
    last_sync_at DATETIME,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    deleted_at DATETIME DEFAULT NULL,
    UNIQUE INDEX uk_tenant_account (tenant_id, local_account_id),
    INDEX idx_tenant (tenant_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE store (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_id BIGINT UNSIGNED NOT NULL,
    account_id BIGINT UNSIGNED NOT NULL COMMENT '关联 local_account.id',
    poi_id BIGINT UNSIGNED DEFAULT 0 COMMENT '巨量 POI ID',
    name VARCHAR(200) NOT NULL,
    city VARCHAR(50),
    district VARCHAR(50),
    address VARCHAR(500),
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    INDEX idx_tenant (tenant_id),
    INDEX idx_account (account_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### 4.2 分组

```sql
CREATE TABLE account_group (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_id BIGINT UNSIGNED NOT NULL,
    name VARCHAR(100) NOT NULL,
    type ENUM('franchisee','region','custom') NOT NULL DEFAULT 'custom',
    parent_id BIGINT UNSIGNED DEFAULT 0,
    sort INT DEFAULT 0,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    deleted_at DATETIME DEFAULT NULL,
    INDEX idx_tenant_type (tenant_id, type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE account_group_member (
    account_id BIGINT UNSIGNED NOT NULL,
    group_id BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (account_id, group_id),
    INDEX idx_group (group_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### 4.3 权限

```sql
CREATE TABLE user_account_scope (
    user_id BIGINT UNSIGNED NOT NULL,
    account_id BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (user_id, account_id),
    INDEX idx_account (account_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### 4.4 投放数据（从巨量同步到本地）

```sql
CREATE TABLE local_project (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_id BIGINT UNSIGNED NOT NULL,
    account_id BIGINT UNSIGNED NOT NULL,
    project_id BIGINT UNSIGNED NOT NULL COMMENT '巨量项目ID',
    name VARCHAR(200),
    marketing_goal VARCHAR(50),
    delivery_scene VARCHAR(50),
    ad_type VARCHAR(50),
    external_action VARCHAR(50),
    bid_type VARCHAR(20),
    bid INT DEFAULT 0 COMMENT '出价（分）',
    budget_mode VARCHAR(30),
    budget INT DEFAULT 0 COMMENT '预算（分）',
    status_first VARCHAR(50),
    status_second VARCHAR(50),
    start_time VARCHAR(20),
    end_time VARCHAR(20),
    poi_id BIGINT UNSIGNED DEFAULT 0,
    poi_name VARCHAR(200),
    product_id BIGINT UNSIGNED DEFAULT 0,
    synced_at DATETIME,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    UNIQUE INDEX uk_tenant_project (tenant_id, project_id),
    INDEX idx_account (account_id),
    INDEX idx_status (status_first)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE local_promotion (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_id BIGINT UNSIGNED NOT NULL,
    account_id BIGINT UNSIGNED NOT NULL,
    project_id BIGINT UNSIGNED NOT NULL,
    promotion_id BIGINT UNSIGNED NOT NULL COMMENT '巨量单元ID',
    name VARCHAR(200),
    status_first VARCHAR(50),
    status_second VARCHAR(50),
    aweme_id VARCHAR(100),
    synced_at DATETIME,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    UNIQUE INDEX uk_tenant_promotion (tenant_id, promotion_id),
    INDEX idx_project (project_id),
    INDEX idx_account (account_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### 4.5 报表

```sql
CREATE TABLE report_daily (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_id BIGINT UNSIGNED NOT NULL,
    account_id BIGINT UNSIGNED NOT NULL,
    project_id BIGINT UNSIGNED DEFAULT 0,
    promotion_id BIGINT UNSIGNED DEFAULT 0,
    stat_date DATE NOT NULL,
    cost DECIMAL(12,2) DEFAULT 0,
    show_cnt BIGINT DEFAULT 0,
    click_cnt BIGINT DEFAULT 0,
    ctr DECIMAL(8,4) DEFAULT 0,
    cpc DECIMAL(8,2) DEFAULT 0,
    cpm DECIMAL(8,2) DEFAULT 0,
    convert_cnt BIGINT DEFAULT 0,
    conversion_rate DECIMAL(8,4) DEFAULT 0,
    conversion_cost DECIMAL(10,2) DEFAULT 0,
    form_cnt BIGINT DEFAULT 0,
    phone_confirm_cnt BIGINT DEFAULT 0,
    phone_connect_cnt BIGINT DEFAULT 0,
    message_action_cnt BIGINT DEFAULT 0,
    clue_pay_order_cnt BIGINT DEFAULT 0,
    clue_message_count BIGINT DEFAULT 0,
    dy_like BIGINT DEFAULT 0,
    dy_comment BIGINT DEFAULT 0,
    dy_share BIGINT DEFAULT 0,
    dy_follow BIGINT DEFAULT 0,
    total_play BIGINT DEFAULT 0,
    play_over_rate DECIMAL(8,4) DEFAULT 0,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    UNIQUE INDEX uk_report (tenant_id, account_id, project_id, promotion_id, stat_date),
    INDEX idx_date (stat_date),
    INDEX idx_account_date (account_id, stat_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### 4.6 批量任务

```sql
CREATE TABLE batch_task (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_id BIGINT UNSIGNED NOT NULL,
    task_type VARCHAR(50) NOT NULL,
    title VARCHAR(200) NOT NULL,
    config JSON NOT NULL,
    status ENUM('pending','running','completed','partial_success','failed','cancelled') DEFAULT 'pending',
    total_count INT DEFAULT 0,
    success_count INT DEFAULT 0,
    failed_count INT DEFAULT 0,
    created_by BIGINT UNSIGNED,
    started_at DATETIME,
    completed_at DATETIME,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    INDEX idx_tenant_status (tenant_id, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE batch_task_item (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    task_id BIGINT UNSIGNED NOT NULL,
    account_id BIGINT UNSIGNED NOT NULL,
    local_account_id BIGINT UNSIGNED NOT NULL,
    config_override JSON,
    status ENUM('pending','running','success','failed','skipped','cancelled') DEFAULT 'pending',
    result_id BIGINT UNSIGNED DEFAULT 0,
    error_msg VARCHAR(500),
    executed_at DATETIME,
    INDEX idx_task (task_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE project_template (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_id BIGINT UNSIGNED NOT NULL,
    name VARCHAR(200) NOT NULL,
    config JSON NOT NULL,
    use_count INT DEFAULT 0,
    created_by BIGINT UNSIGNED,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    deleted_at DATETIME DEFAULT NULL,
    INDEX idx_tenant (tenant_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

---

## 5. 核心业务流程

### 5.1 批量创建项目

```
投手选择模板（或手动配置公共参数）
  → 选择目标账户（按分组勾选或手动选）
  → 为每个账户设置差异化字段（预算/门店POI）
  → 预览影响范围
  → 确认提交 → 创建 batch_task + N 条 item
  → Worker 逐条执行（限流 10 QPS）
  → 调用 SDK: project.Create(ctx, sdkClient, token, req)
  → 成功记录 project_id，失败记录 error_msg
  → 前端轮询进度
  → 完成后同步到 local_project 表
```

### 5.2 数据同步

| 任务 | 频率 | SDK 调用 |
|------|------|---------|
| 日报表同步 | 每天凌晨 2:00 | report.ProjectGet + PromotionGet |
| 实时消耗 | 每 30 分钟 | report.ProjectGet（今日） |
| 项目/单元状态 | 每小时 | project.List + promotion.List |
| Token 续期 | 每 5 分钟 | OAuth refresh |

### 5.3 数据分析聚合

| 维度 | SQL 实现 |
|------|---------|
| 按门店 | GROUP BY account_id |
| 按加盟商 | JOIN account_group_member → GROUP BY group_id |
| 按区域 | JOIN store → GROUP BY city/district |
| 按投手 | JOIN user_account_scope → GROUP BY user_id |
| 环比 | 同维度对比昨日/上周同期 |

### 5.4 权限隔离

- 中间件从 JWT 提取 tenant_id + user_id
- 管理员：只过滤 tenant_id
- 投手：过滤 tenant_id + account_id IN (user_account_scope)
- Repository 层统一注入 WHERE 条件

---

## 6. API 设计

### 账户管理
- POST /api/v1/accounts/sync
- GET /api/v1/accounts
- GET /api/v1/accounts/:id
- GET /api/v1/accounts/:id/balance

### 分组管理
- CRUD /api/v1/groups
- POST /api/v1/groups/:id/members
- DELETE /api/v1/groups/:id/members

### 投放管理
- CRUD /api/v1/projects
- PUT /api/v1/projects/:id/status
- CRUD /api/v1/promotions
- PUT /api/v1/promotions/:id/status

### 批量操作
- POST /api/v1/batch/projects
- POST /api/v1/batch/promotions
- POST /api/v1/batch/budget
- POST /api/v1/batch/bid
- POST /api/v1/batch/status
- GET /api/v1/batch/tasks
- GET /api/v1/batch/tasks/:id
- POST /api/v1/batch/tasks/:id/cancel
- POST /api/v1/batch/tasks/:id/retry

### 模板管理
- CRUD /api/v1/templates

### 数据分析
- GET /api/v1/analytics/overview
- GET /api/v1/analytics/trend
- GET /api/v1/analytics/rank
- GET /api/v1/analytics/compare
- GET /api/v1/analytics/detail

### 系统管理
- POST /api/v1/users/:id/scope
- GET /api/v1/users/:id/scope
- GET /api/v1/sync/status
- POST /api/v1/sync/trigger

---

## 7. 技术栈

| 层 | 技术 |
|---|---|
| 前端 | Vue 3 + TypeScript + Vite + TailwindCSS + vue-chartjs |
| 后端 | Go 1.26 + Gin + GORM |
| SDK | github.com/bububa/oceanengine (项目内 sdk/ 目录) |
| 数据库 | MySQL 9 |
| 缓存/限流 | Redis 8 (Lua 令牌桶) |
| 迁移 | golang-migrate |

---

## 8. 目录结构

```
backend/
├── cmd/server/main.go
├── config/settings.yml
├── internal/
│   ├── router/router.go
│   ├── middleware/
│   │   ├── jwt.go
│   │   ├── data_scope.go          (新)
│   │   └── operation_log.go
│   ├── app/
│   │   ├── admin/                  (复用+改造)
│   │   ├── account/                (新)
│   │   ├── group/                  (新)
│   │   ├── store/                  (新)
│   │   ├── project/                (新)
│   │   ├── promotion/              (新)
│   │   ├── batch/                  (新)
│   │   ├── template/               (新)
│   │   └── analytics/              (新)
│   ├── batch/worker.go
│   └── scheduler/scheduler.go
├── pkg/
│   ├── ocean/sdk.go                (新)
│   ├── ratelimiter/limiter.go      (新)
│   ├── database/
│   ├── response/
│   └── errcode/
└── migrations/
```

---

## 9. 与现有代码的融合

| 复用 | 改造 | 不需要 |
|------|------|--------|
| Gin 框架 + 路由结构 | pkg/oceanengine/local.go → 替换为 SDK | internal/app/qianchuan/ |
| JWT + 用户管理 | internal/app/local/ → 重写为 project/promotion | internal/app/star/ |
| RBAC 权限 | internal/app/report/ → 重写为 analytics | internal/app/dpa/ |
| 操作日志中间件 | cmd/task/ → 改为主服务内 goroutine | internal/app/enterprise/ |
| 配置管理 (Viper) | 数据库模型加 tenant_id | internal/app/clue/ |
| 数据库连接 (GORM) | | internal/app/campaign/ |
| 响应格式/错误码 | | internal/app/ad/ |
| 前端 Vue3 骨架 | | |
| 前端通用组件 | | |

---

## 10. 扩展性预留

| 未来能力 | 扩展方式 |
|---------|---------|
| 抖音来客 | tenant 表加来客 Token；新增 internal/app/laike/ 模块 |
| 线索管理 | 新增 clue 表 + internal/app/clue/；复用账户/权限体系 |
| 智能分析 | 新增 internal/app/rule/；如果计算量大，拆为独立进程 |
| 多品牌大规模 | 限流器按 tenant 隔离；报表表按月分区 |

---

## 11. 实施分段建议

```
Phase 1 (Week 1-3):  账户体系
  租户管理 + OAuth + 账户同步 + 分组 + 门店 + 投手权限 + 限流器

Phase 2 (Week 4-6):  投放管理
  项目/单元 CRUD + 批量操作引擎 + 模板 + 数据同步定时任务

Phase 3 (Week 7-9):  数据分析
  报表同步 + 汇总看板 + 多维聚合 + 环比 + 排行 + 导出

Phase 4 (Week 10-11): 打磨
  前端体验优化 + 性能调优 + 缓冲
```

---

## 12. 审查修正（5 轮审查，17 个问题）

### 修正 1：store 与 account 关系明确

store 表支持多条记录指向同一个 account_id（一对多）。权限隔离以 account_id 为单位，投手能看到某个账户就能看到该账户下所有门店。这是设计意图，不做门店级权限细分。

### 修正 2：OAuth 回调流程

```
1. 管理员在平台发起授权 → 跳转巨量 OAuth 页面
2. 用户授权后回调 /api/v1/advertisers/oauth/callback
3. 后端用 auth_code 换取 access_token + refresh_token
4. 写入 tenant 表（token_expire_at = now + expires_in）
5. 自动触发一次账户同步（拉取该组织下所有 local_account_id）
```

### 修正 3：限流器按租户隔离

每个租户独立限流，互不影响：

```
Redis Key: ratelimit:{tenant_id}:ocean
每个租户 10 QPS（巨量引擎的限制是按应用维度，每个租户对应一个 AppID）
```

如果多个租户共用同一个 AppID（不推荐），则需要共享限流池。当前设计假设每个租户有独立的 AppID。

### 修正 4：report_daily 存储策略

报表按粒度分层存储：

```
project_id = 0, promotion_id = 0 → 账户级日汇总
project_id > 0, promotion_id = 0 → 项目级日汇总
project_id > 0, promotion_id > 0 → 单元级日明细
```

唯一索引 `uk_report (tenant_id, account_id, project_id, promotion_id, stat_date)` 保证每个粒度每天只有一条记录。同步时分别拉取三个粒度的数据。

### 修正 5：数据同步吞吐量验证

假设 3 个租户各 200 个账户 = 600 个账户：

| 同步任务 | 调用次数 | 耗时（10 QPS/租户） |
|---------|---------|-------------------|
| 项目状态（每小时） | 600 次 | 每租户 200/10 = 20s |
| 实时消耗（每 30 分钟） | 600 次 | 每租户 20s |
| 日报表（凌晨 2 点） | 600×3 粒度 = 1800 次 | 每租户 60s |

结论：每个租户的同步任务在 1-2 分钟内完成，远小于轮询间隔。多租户并行执行（各自有独立限流池），互不阻塞。

### 修正 6：batch_task_item 增加 tenant_id

```sql
ALTER TABLE batch_task_item ADD COLUMN tenant_id BIGINT UNSIGNED NOT NULL AFTER id;
ALTER TABLE batch_task_item ADD INDEX idx_tenant (tenant_id);
```

### 修正 7：现有 user 表改造

现有 `sys_user` 表增加 `tenant_id` 字段：

```sql
ALTER TABLE sys_user ADD COLUMN tenant_id BIGINT UNSIGNED NOT NULL DEFAULT 0 AFTER id;
ALTER TABLE sys_user ADD INDEX idx_tenant (tenant_id);
```

- tenant_id = 0 表示平台超级管理员（可管理所有租户）
- tenant_id > 0 表示租户内用户（只能看到自己租户的数据）

### 修正 8：账户同步来源

巨量引擎 SDK 中没有直接的"获取组织下所有 local_account_id"接口。账户同步的实现方式：

**方案 A（推荐）：手动录入 + 批量导入**
- 管理员通过 Excel 导入账户列表（local_account_id + 名称 + 门店信息）
- 或在巨量后台导出账户列表，上传到平台

**方案 B：通过代理商 API**
- 如果品牌方通过代理商开户，可用 `/agent/advertiser/select/` 获取子账户列表
- 需要确认品牌方的开户方式

**方案 C：通过项目列表反推**
- 调用 project.List 时不传 filtering，返回的项目中包含 local_account_id
- 遍历所有项目可以收集到活跃的账户列表

Phase 1 实现时先用方案 A（手动导入），后续验证方案 B/C 是否可行。

`POST /api/v1/accounts/sync` 改为 `POST /api/v1/accounts/import`（Excel 批量导入）。

### 修正 9：增加单元模板

```sql
CREATE TABLE promotion_template (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_id BIGINT UNSIGNED NOT NULL,
    name VARCHAR(200) NOT NULL,
    config JSON NOT NULL COMMENT '单元配置（素材、抖音号等）',
    use_count INT DEFAULT 0,
    created_by BIGINT UNSIGNED,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    deleted_at DATETIME DEFAULT NULL,
    INDEX idx_tenant (tenant_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

模板管理 API 扩展：
- CRUD /api/v1/templates/projects — 项目模板
- CRUD /api/v1/templates/promotions — 单元模板

### 修正 10：Token 续期边界处理

```
Token 生命周期管理：
  - access_token 有效期 24h → 每 5 分钟检查，提前 2h 续期
  - refresh_token 有效期 30 天 → 提前 7 天告警"即将需要重新授权"
  - refresh_token 过期 → token_status 标记为 expired，前端显示"需要重新授权"
  - 管理员点击"重新授权" → 跳转 OAuth 流程
```

### 修正 11：批量任务取消和重试

```
取消规则：
  - pending 状态的 item → 标记 cancelled
  - running 状态的 item → 等待执行完成，不中断
  - 已 success 的 item → 不回滚（巨量 API 无事务）

重试规则：
  - POST /api/v1/batch/tasks/:id/retry → 重试该任务下所有 failed 状态的 item
  - 重试时 item 状态从 failed 改回 pending，重新进入执行队列
```

### 修正 12：报表导出策略

```
导出限制：
  - 单次导出最多 10000 行
  - 超过时提示"请缩小时间范围或筛选条件"
  - 导出为异步任务：请求导出 → 后台生成文件 → 通知下载
  - 文件格式：Excel (.xlsx)
  - 文件保留 24h 后自动清理
```

### 修正 13：account_id 字段命名规范

统一命名规范，避免混淆：

| 字段名 | 含义 | 类型 |
|--------|------|------|
| `local_account.id` | 本平台内部自增 ID | 内部引用 |
| `local_account.local_account_id` | 巨量引擎的账户 ID | 调 SDK 时使用 |
| 其他表的 `account_id` | 指向 `local_account.id` | 外键关联 |

所有业务表（local_project, report_daily 等）的 `account_id` 存的是 `local_account.id`（内部 ID），调用 SDK 时通过 JOIN 或缓存获取 `local_account_id`（巨量 ID）。

### 修正 14：同步与用户操作的一致性

```
策略：同步时不覆盖"新鲜"数据

规则：
  - local_project 表增加 `local_updated_at` 字段（平台侧最后操作时间）
  - 同步时检查：如果 local_updated_at > synced_at - 5min，跳过该记录
  - 即：如果用户刚操作过（5 分钟内），不用同步数据覆盖
  - 下一轮同步时（1 小时后），巨量侧数据已经更新，正常覆盖
```

### 修正 15：批量任务进度推送

采用**短轮询**（前端每 2 秒请求一次 `GET /api/v1/batch/tasks/:id`），原因：
- 批量任务通常 10-60 秒完成（100 个 item / 10 QPS = 10s）
- 短轮询实现简单，不需要 SSE/WebSocket 基础设施
- 任务完成后停止轮询

### 修正 16：限流器 Redis Key 设计

```
Key: ratelimit:{tenant_id}:ocean:{second_bucket}
TTL: 2s
Value: 当前秒内已消耗的令牌数

Lua 脚本逻辑：
  KEYS[1] = ratelimit:{tenant_id}:ocean:{current_second}
  if GET(key) < 10 then INCR(key); return 0 (获得令牌)
  else return 1 (等待)
```

### 修正 17：SDK import path

项目内 `sdk/` 目录使用 `go.mod` 的 replace 指令：

```go
// backend/go.mod
replace github.com/bububa/oceanengine => ../sdk
```

这样代码中 import 路径保持 `github.com/bububa/oceanengine/...`，但实际引用本地 sdk/ 目录。
