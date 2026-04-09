# CLAUDE.md

本文件为 Claude Code (claude.ai/code) 在此仓库中工作提供指引。

完整的全栈架构、后端命令、环境变量和 API 端点请参考上级 [../CLAUDE.md](../CLAUDE.md)。
回复的内容实用中文
## 语言规范

`frontend/` 目录下的代码、注释和文档全部使用英文，包括提交信息、变量名、JSDoc 以及返回给前端的字符串。完整语言规范见上级 CLAUDE.md。

## 常用命令

```bash
npm install          # 安装依赖
npm run build        # 生产构建 → dist/
npm run preview      # 本地预览生产构建
```

无 linter、TypeScript、格式化工具或测试框架——纯 JavaScript 项目。

## 技术栈

- Vue 3（Composition API，`<script setup>`）+ Vite 6
- ethers.js v6 用于钱包交互
- 纯 CSS 暗色主题（`src/styles/main.css`），无 Tailwind/SCSS
- 无状态管理库（Vuex/Pinia）——参见下方单例模式

## 目录结构

```
src/
├── App.vue                  # 根组件：组合 WalletConnect + AccessCheck
├── composables/
│   ├── useWallet.js         # 钱包连接、签名、JWT 存储
│   └── useApi.js            # 后端 API 调用封装（challenge、verify、check-access）
├── components/
│   ├── WalletConnect.vue    # MetaMask 连接按钮 + 状态显示
│   ├── AccessCheck.vue      # ERC-20 合约地址输入 + 查询触发
│   └── ResultDisplay.vue    # 查询结果展示（有权限/无权限/错误）
└── styles/
    └── main.css             # 全局样式，暗色主题
```

## 关键模式

**模块级 ref 单例状态**：`useWallet.js` 在导出函数外部声明 `currentAddress`、`jwtToken`、`isConnecting` 为 `ref()`。所有调用 `useWallet()` 的组件共享同一份响应式状态，无需状态管理库。

**ethers.js 动态导入**：`connect()` 内部通过 `await import('ethers')` 按需加载，减小初始包体积。

**API 基地址**：`useApi.js` 使用 `window.location.origin`。开发环境下 Vite proxy（`vite.config.js`）将 `/auth`、`/check-access`、`/health` 转发到 `localhost:8080`；生产环境由 Go 后端直接提供 `dist/` 静态文件，同源访问无需额外配置。

**默认合约地址**：`AccessCheck.vue` 默认填入 USDT 合约地址（`0xdAC17F958D2ee523a2206206994597C13D831ec7`），方便快速测试。

## 开发流程

- 先启动 Go 后端（`:8080`），再运行 `npm run dev` 启动前端热更新
- 开发环境下 Vite proxy 处理跨域，无需配置 `ALLOWED_ORIGINS`
- 钱包功能需要安装 MetaMask 浏览器扩展
