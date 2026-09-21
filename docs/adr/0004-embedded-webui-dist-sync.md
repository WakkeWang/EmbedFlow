# ADR-0004: 嵌入式前端走双目录同步，部署必须走完整三步

日期: 2026-09-21
状态: 已采纳

## 背景

服务端二进制经 `go:embed` 内嵌前端产物（单文件分发，Distribution Plan）。嵌入点是 `server/internal/webui/dist`（git 入库的副本），而 Vite 构建写到 `web/dist`。两个目录内容相同但**没有自动同步**——`web/vite.config.ts` 未配置 outDir 指向嵌入点，README 的快速开始也只写了 `cd web && npm run build`。

生产 192.168.0.80 上实际踩中（2026-09-21）：按 README 流程构建部署后，验证环节发现线上入口 `index-*.js` 文件名与本地新构建不一致——二进制里嵌的是**上一次**的 dist，当天全部前端改动静默丢失。`go:embed` 对旧文件照常打包，无任何报错。

## 决策

构建部署的完整流程（顺序关键）：

1. `cd web && npm run build` → 产物在 `web/dist`
2. **手动同步嵌入点**：删除 `server/internal/webui/dist`，将 `web/dist` 拷贝过去
3. `CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath -o embedflow-server-linux-arm64 ./server/cmd/embedflow-server`（纯 Go 依赖，双架构同理）
4. 上传目标机 → 原子替换（`/tmp` 中转 `mv`）→ `systemctl restart`

**部署后必须验证**：`curl http://<server>:8420/` 返回的 `index-*.js` 文件名须与本地 `web/dist/index.html` 一致。只看 HTTP 200 或服务 active 抓不住这个错。

## 否决的替代方案

- **vite outDir 直接指向嵌入点**：省掉拷贝，但 `web/dist` 与嵌入点二选一后，另一个目录的存在只会误导；且仓库惯例（git 入库的嵌入点 + 构建产物不入库的 web/dist）会被打破。保持双目录 + 显式同步。
- **Makefile / 脚本封装**：迟早要做，但流程正确性优先于自动化；先在文档与记忆中固定，工具化作为后续改进。

## 后果

- 约束：任何部署（哪怕只改了 Go 代码）都建议跑一遍验证比对；改了前端则第 2 步绝不可省。
- 仓库里 `server/internal/webui/dist` 是**唯一入库的前端产物**，更新它需单独一个 `chore: rebuild embedded webui dist` 提交（既有惯例）。
