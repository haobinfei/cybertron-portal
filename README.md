# Cybertron Portal

基于 Go + Vue3 的运维管理平台，提供服务监控、资源管理、自动化运维等功能。

## 技术栈

| 层级     | 技术选型                            |
| -------- | ----------------------------------- |
| 后端     | Go 1.22+, Gin, GORM, Viper, Zap     |
| 前端     | Vue 3, TypeScript, Vite, Pinia, Element Plus |
| 数据库   | MySQL 8.0 / PostgreSQL              |
| 缓存     | Redis                               |
| 部署     | Docker, Nginx                       |

## 目录结构

```
cybertron_portal/
├── README.md                   # 项目说明文档
├── Makefile                    # 项目统一构建脚本
├── docker-compose.yml          # 本地容器编排
│
├── backend/                    # Go 后端服务
│   ├── cmd/
│   │   └── server/
│   │       └── main.go         # 程序入口
│   ├── internal/               # 内部模块（不可外部引用）
│   │   ├── config/             # 配置加载与解析
│   │   ├── handler/            # HTTP 请求处理器（Controller 层）
│   │   ├── middleware/         # 中间件（认证、日志、CORS 等）
│   │   ├── model/              # 数据模型定义
│   │   ├── repository/         # 数据访问层（DAO）
│   │   ├── router/             # 路由注册
│   │   └── service/            # 业务逻辑层
│   ├── pkg/                    # 可复用的公共包
│   │   └── utils/              # 工具函数
│   ├── config/
│   │   └── config.yaml         # 默认配置文件
│   ├── go.mod
│   └── go.sum
│
├── frontend/                   # Vue3 前端项目
│   ├── src/
│   │   ├── api/                # API 请求封装
│   │   ├── assets/             # 静态资源（图片、样式）
│   │   ├── components/         # 公共组件
│   │   ├── composables/        # 组合式函数（Hooks）
│   │   ├── layouts/            # 布局组件
│   │   ├── router/             # 路由配置
│   │   ├── stores/             # Pinia 状态管理
│   │   ├── types/              # TypeScript 类型定义
│   │   ├── utils/              # 工具函数
│   │   ├── views/              # 页面视图
│   │   ├── App.vue             # 根组件
│   │   └── main.ts             # 前端入口
│   ├── public/                 # 静态文件（不经过构建）
│   ├── index.html              # HTML 模板
│   ├── package.json
│   ├── tsconfig.json
│   ├── vite.config.ts
│   └── .env                    # 环境变量
│
├── deploy/                     # 部署相关配置
│   ├── Dockerfile.backend      # 后端 Docker 镜像
│   ├── Dockerfile.frontend     # 前端 Docker 镜像
│   └── nginx/
│       └── nginx.conf          # Nginx 反向代理配置
│
└── docs/                       # 项目文档
    └── api.md                  # API 接口文档
```

## 功能模块（规划）

- **仪表盘** — 系统概览、关键指标展示
- **服务管理** — 服务注册、发现、启停控制
- **资源监控** — CPU、内存、磁盘、网络实时监控
- **告警中心** — 告警规则配置、告警通知
- **任务管理** — 定时任务、脚本下发、批量执行
- **权限管理** — 用户、角色、权限控制
- **操作审计** — 操作日志记录与查询
- **配置管理** — 配置中心、配置下发

## 快速开始

### 环境要求

- Go 1.22+
- Node.js 18+
- MySQL 8.0
- Redis 7.0+

### 本地开发

**1. 克隆项目**

```bash
git clone git@github.com:haobinfei/cybertron-portal.git
cd cybertron-portal
```

**2. 启动后端**

```bash
cd backend
cp config/config.yaml config/config.local.yaml  # 修改为本地配置
go mod tidy
go run cmd/server/main.go
```

**3. 启动前端**

```bash
cd frontend
npm install
npm run dev
```

**4. 使用 Docker Compose 一键启动（推荐）**

```bash
docker-compose up -d
```

## 开发规范

### 分支管理

- `master` — 生产分支
- `develop` — 开发分支
- `feature/*` — 功能分支
- `fix/*` — 修复分支

### 提交规范

```
feat: 新功能
fix: 修复 Bug
docs: 文档更新
style: 代码格式调整
refactor: 重构
test: 测试
chore: 构建/工具变动
```

## 部署

后端编译为二进制文件后通过 Docker 部署，前端构建为静态文件后通过 Nginx 提供服务。

```bash
# 后端构建
cd backend && go build -o server cmd/server/main.go

# 前端构建
cd frontend && npm run build
```
