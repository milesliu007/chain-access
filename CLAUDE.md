# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 语言规范
- 所有文档（README、注释、文档文件等）使用中文
- 控制台输出使用中文
- ai 回复的内容实用中文
- 代码中的变量名、函数名等标识符使用英文
- 输出的日志使用英文
- frontend/ 目录下，全部使用英文
- 代码中的返回值，对外输出的，必须英文
- 所有给前端可以看到的，必须使用英文

## 常用命令

### 后端（Go）
```bash
go build -o chain-access .     # 编译
go run .                       # 本地运行（需要 .env 中的环境变量）
go test ./...                  # 运行所有测试
go test ./api/service/...      # 运行单个包的测试
go test -run TestFuncName ./api/service/...  # 运行单个测试
```

### 前端（Vue 3 + Vite）
```bash
cd frontend
npm install                    # 安装依赖
npm run dev                    # 开发服务器（代理 /auth、/check-access、/health 到 :8080）
npm run build                  # 生产构建（输出到 dist/）
npm run preview                # 预览生产构建
```

无 linter、TypeScript 或格式化工具配置——纯 JavaScript 项目。目前无测试文件。

### Docker
```bash
docker build -t chain-access .
docker run -p 8080:8080 --env-file .env chain-access
```

## 环境变量

必须设置（缺失时服务拒绝启动）：
- `JWT_SECRET` — JWT 签名密钥
- `INFURA_API_KEY` — Infura RPC 接口密钥

可选：
- `PORT` — 服务端口（默认 8080）
- `ALLOWED_ORIGINS` — CORS 允许的来源（默认 `http://localhost:8080`）

参考 `.env.example` 配置。

## 架构概览

Go 模块名：`chain-access`

Web3 代币门控应用：用户通过以太坊钱包签名登录，后端验证签名后签发 JWT，再通过链上查询判断用户是否持有指定 ERC-20 代币。

### 请求流程

```
前端 (Vue 3) → Gin Router → Middleware (JWT) → Controller → Service → 链上/内存
```

### 后端分层（`api/` 目录）

| 层 | 目录 | 职责 |
|---|---|---|
| Config | `api/config/` | 从环境变量加载配置，启动时校验必填项 |
| Router | `api/router/` | 路由注册，区分公开路由和受保护路由 |
| Middleware | `api/middleware/` | JWT 认证中间件，将钱包地址注入 context |
| Controller | `api/controller/` | HTTP 请求/响应处理，参数绑定与校验 |
| Service | `api/service/` | 核心业务逻辑（均为接口定义，便于测试） |
| Repository | `api/repository/` | 数据持久化（当前为内存实现） |
| Model | `api/model/` | 请求/响应 DTO |

### 核心服务

- **AuthService** — 挑战-签名认证流程：生成 UUID 挑战（5 分钟 TTL）→ 验证 EIP-191 签名 → 签发 24 小时 JWT
- **EthereumService** — 通过 Infura RPC 调用 ERC-20 `balanceOf()`，判断地址是否持有代币

### API 端点

| 方法 | 路径 | 认证 | 说明 |
|---|---|---|---|
| POST | `/auth/challenge` | 无 | 获取签名挑战 |
| POST | `/auth/verify` | 无 | 提交签名，获取 JWT |
| GET | `/check-access` | JWT | 检查代币持有状态 |

### 前端（`frontend/` 目录）

Vue 3 + Vite 单页应用，通过 MetaMask 连接钱包。核心逻辑在 composables 中：
- `useWallet.js` — 钱包连接与签名
- `useApi.js` — 后端 API 调用封装

关键模式：
- `useWallet.js` 中的 `ref` 定义在函数外部（模块级），形成全局单例共享状态，无需 Vuex/Pinia
- ethers.js 通过 `await import('ethers')` 动态导入，减小初始包体积
- `useApi.js` 使用 `window.location.origin` 作为 API 基地址，开发环境通过 Vite proxy（`vite.config.js`）转发到后端
- 样式为纯 CSS（`src/styles/main.css`），暗色主题，无 Tailwind/SCSS

生产环境下由 Go 后端以静态文件方式提供服务（`frontend/dist/`）。

### 关键设计决策

- 无数据库依赖：挑战存储使用带 TTL 的内存 map（`sync.Map` + 过期清理）
- 服务层均通过接口定义，controller 依赖接口而非实现
- Dockerfile 采用多阶段构建：Node 构建前端 → Go 构建后端 → 最终镜像仅含二进制和静态文件
- `EthereumService` 手动编码 ERC-20 ABI 选择器（`0x70a08231`），无 ABI codegen 依赖
- `/check-access` 强制校验查询的 `address` 参数必须与 JWT 中的 `address` 一致（大小写不敏感），防止越权
- `EthereumService` 中 Infura RPC 通过本地代理 `127.0.0.1:7897` 访问（开发环境网络代理）