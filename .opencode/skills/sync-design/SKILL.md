---
name: sync-design
description: Use ONLY when the user asks to sync design docs, update design.md, or says /sync-design. Analyzes git diff to detect code changes and updates docs/design.md to keep it in sync with actual code. Covers file-to-section mapping for Go backend (model, handler, router, service, middleware) and Vue3 frontend (router, views, layouts, stores, api, types).
---

# Sync Design — 技术方案文档同步

当用户请求同步技术方案文档时（如 `/sync-design` 或 "同步设计文档"、"更新design.md"），执行以下流程。

## 同步流程

### Step 1: 确定变更范围
执行 `git diff --name-only` 获取未提交的变更文件，或 `git diff HEAD~1 --name-only` 获取最近一次提交的变更。

### Step 2: 代码 → 文档章节映射

| 代码文件路径 | 文档章节 | 更新内容 |
|---|---|---|
| `backend/cmd/*/main.go` | 2.1 服务入口 | 服务初始化流程 |
| `backend/internal/config/*.go` | 2.2 配置管理 | 配置项说明 |
| `backend/internal/model/*.go` | 2.3 数据模型 | 新增/修改的 struct，更新 Mermaid ER 图 |
| `backend/internal/handler/*.go` | 2.4 API 接口 | 接口清单、请求/响应结构 |
| `backend/internal/service/*.go` | 2.5 业务逻辑层 | 服务方法说明 |
| `backend/internal/middleware/*.go` | 2.6 中间件 | 中间件清单和功能 |
| `backend/internal/router/*.go` | 2.7 路由注册 | 路由表 |
| `frontend/src/router/index.ts` | 3.1 路由设计 | 同步路由表 |
| `frontend/src/views/**/*.vue` | 3.2 页面清单 | 页面路径、功能描述 |
| `frontend/src/layouts/**/*.vue` | 3.3 布局组件 | 布局结构 |
| `frontend/src/components/**/*.vue` | 3.4 公共组件 | 组件列表、Props/Emits |
| `frontend/src/stores/*.ts` | 3.5 状态管理 | Store 清单、State/Actions |
| `frontend/src/api/*.ts` | 3.6 API 调用层 | 请求封装、拦截器 |
| `frontend/src/types/*.d.ts` | 3.7 类型定义 | 新增的 interface/type |
| `vite.config.ts` | 4.1 构建配置 | 代理、端口、别名 |
| `deploy/**` | 5 部署方案 | Dockerfile、部署流程 |

### Step 3: 更新文档内容

对每个受影响的章节：
1. **读取实际代码**：读取变更后的源文件，获取真实的类型、函数、组件结构
2. **对比现有文档**：检查当前文档内容是否与实际代码一致
3. **更新对应章节**：用代码中的实际定义替换文档中过时的内容
4. **更新 Mermaid 图**：如果涉及架构关系变化，更新对应的 Mermaid 图表
5. **标注时间戳**：在更新的章节末尾添加 `<!-- synced: YYYY-MM-DD -->` 注释

### Step 4: 输出同步报告

向用户报告：
- 哪些文件发生了变更
- 哪些文档章节被更新
- 新增/删除/修改了哪些内容

## 文档格式约定

- 使用 `<!-- synced: YYYY-MM-DD -->` 标记最后同步时间
- 接口定义使用 TypeScript 语法（前端）或 Go 语法（后端）
- Mermaid 图表使用 ````mermaid` 代码块
- 保持章节编号结构不变，不随意增删一级章节

## 注意事项

- 只更新变更相关的章节，不要重写整个文档
- 保持文档格式一致，不改变缩进、表格风格
- 如果代码被删除，文档中标注 `[已移除]` 而不是直接删除内容
- 如果检测到新增模块，新增对应章节
