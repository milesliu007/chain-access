# chain-access

Web3 身份与访问管理



cd frontend && npm run build
cd frontend && npm run build


开发流程：

    1. 编辑 frontend/src/admin/ 下的 Vue 文件
    2. 构建：
    cd frontend && npm run build
    2. 输出自动到 frontend/dist/，Go 后端直接读取这个目录
    3. 重启 Go 服务：
    go run .

    开发时用热更新（不需要手动复制）：
    # 终端1：启动 Go 后端
    go run .

    # 终端2：启动 Vite 开发服务器（带热更新）
    cd frontend && npm run dev

    开发服务器访问 http://localhost:5173/index-admin.html，API 请求通过 Vite proxy 转发到 :8080。

    生产环境只需 npm run build，Go 直接 serve frontend/dist/，无需额外复制步骤。