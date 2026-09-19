# EmbedFlow

**Build, flash, test, release -- for small embedded teams, with humans as first-class citizens of the flow.**

EmbedFlow is a self-hosted build/test/release platform for small embedded product teams. It keeps the serial port on the tester's Windows PC, tunnels it through a client to the server, and drives configurable expect flash rules and human-in-the-loop steps from the browser.

> 中文说明见下方 [中文](#中文) 部分。

## Why EmbedFlow

Existing HIL/CI frameworks (LAVA, LabGrid, pytest-embedded, Renode) assume the serial port is on the server and everything can be automated. Small teams live in a different reality:

- The serial cable is plugged into a **tester's Windows PC**, not a server rack
- Every version ships with **human steps**: watch the LED, unplug one power line, confirm the beep
- Flash scripts live as **SecureCRT VBS files on one person's machine** -- not shared, not reviewable, not recorded

EmbedFlow makes **the human a first-class citizen of the test flow**: human observations and interventions are formal steps with prompts, web confirmations, and a place in the report. That is the fundamental fork from every existing framework.

### What you get

- **Build module**: turn your build scripts into project-scoped build items -- git or server-local sources, build commands with timeout, artifact globs with sha256/md5 checksums, version extraction. Trigger batches (up to N in parallel) over prerequisite groups (OR within a group, AND across groups); every build gets a full phase-segmented log, a commit snapshot, and archived artifacts with recorded checksums. Dirty sources refuse to build.
- **Browser terminal**: open a session against a shared serial device from any browser (xterm.js), with device mutex, busy-rejection with occupier info, idle timeouts, and read-only follow for a second viewer
- **Configurable expect flashing**: turn your VBS flash scripts into reusable, reviewable step sequences (await string + optional send + timeout + failure action + delay + control-char encoding), editable in web forms, executed with live progress and a guarded abort
- **Human confirmation cards**: record "LED on?", "cable unplugged" as first-class steps -- PASS/FAIL with notes, never auto-dismissed, stored in the session record
- **Session logs that survive**: byte-accurate, timestamped RX/TX lines, buffered + time-flushed, download with HTTP Range resume; write failures mark the session "log incomplete" instead of silently losing data
- **Virtual demo device**: run the whole loop with zero hardware -- expect flashing AND a real build (cloned git repo, archived artifact)

## Quick start (zero hardware)

Requires Go 1.27+ and Node.js (only for building the web assets; prebuilt releases embed them).

```bash
# 1. Build the server (frontend assets included via embed)
cd web && npm install && npm run build && cd ..
go build -o embedflow-server ./server/cmd/embedflow-server

# 2. Run with the built-in virtual device
./embedflow-server -demo -addr :8420

# 3. Open http://localhost:8420, log in (default admin/admin -- change it)
```

With `-demo` the server registers a `demo-virtual` board, seeds a `demo-flash` expect rule, and seeds a `demo-build` build item (a self-contained git repo under the data dir). Open the device, press Run, and watch a full uboot-style flash flow execute in the terminal; open the Build section and trigger a real build with an archived, checksummed artifact -- no device, no serial cable.

## Connecting a real device

1. Run the server on your build machine (`./embedflow-server`, no `-demo`)
2. On the tester's Windows PC, run the serial client (`serial-client -server http://<server>:8420 -user U -pass P -device N -port COM3`) -- it registers as the device's serial transport, pure passthrough, no business logic
3. From any browser, open the device and start a session

## Architecture

```
┌────────────┐  serial   ┌───────────────┐  WebSocket  ┌──────────────────┐
│ Device      ├──────────►  serial-client  ├────────────►│  embedflow-server │
│ (target hw) │  (COM)    │  (Windows PC) │  dual-frame │  (build machine)  │
└────────────┘           └───────────────┘  control:JSON └────────┬─────────┘
                                                     data:binary  │  SQLite + logs
                                                            ┌─────┴─────┐
                                                            │  Browser   │
                                                            │  (xterm.js)│
                                                            └───────────┘
```

- **Domain kernel** (`server/internal/session`): pure state machine, no IO -- device mutex, session lifecycle, disconnect branches, dual-threshold liveness
- **Executors** (`transport`, `sessionlog`, `expect`, `serialport`): wire, persistence, engine, port abstraction
- **Shared protocol** (`pkg/protocol`): one wire definition for server and client (WebSocket dual-frame: JSON control frames, binary data frames)

## Development

```bash
go work sync
go test ./...              # all modules
cd web && npm run dev      # Vite dev server with API proxy
```

Layout:

```
pkg/protocol/       shared wire protocol
server/             Go server (kernel, executors, REST+WS API, demo device)
serial-client/      Windows serial tunnel client (core + CLI; Wails GUI planned)
web/                Vue3 + Naive UI frontend
docs/               requirements, design, ADRs
```

## Roadmap

M1 (session layer) and M2 (build module) are done: sessions, logs, expect runs, human confirmations, build items, batches, records and artifacts, plus the virtual demo device covering both. Next: deploy module with SSH and flash rules per project (M3), test module with reports (M4), release module (M5). See `docs/designs/office-hours-design.md` for the full design and decision log.

## License

MIT

---

# 中文

**构建、刷机、测试、发布——为中小嵌入式团队打造，把人作为流程的一等公民。**

EmbedFlow 是面向中小嵌入式产品团队的自托管构建/测试/发布平台。串口留在测试人员的 Windows 电脑上，经客户端纯透传到服务器，在浏览器里驱动可配置的 expect 刷机规则与人机混合步骤。

## 为什么是 EmbedFlow

现有 HIL/CI 框架（LAVA、LabGrid、pytest-embedded、Renode）都假设「串口在服务器上、一切可自动化」。中小团队的现实是：

- 串口插在**测试人员的 Windows 电脑**上，不在机架服务器上
- 每个版本都带**人为步骤**：看 LED、拔一根电源线、听蜂鸣器
- 刷机脚本是**某个人电脑上的 SecureCRT VBS 文件**——不共享、不可审查、无记录

EmbedFlow 把**人作为测试流程的一等公民**：人为观察与人为干预是正式的流程步骤，有提示词、有 Web 确认、进报告。这是与所有现有框架的根本分叉。

### 核心能力

- **构建模块**：把构建脚本变成工程内的构建项目——git 或服务器本地源、带超时的构建命令、产物通配声明与 sha256/md5 校验和、版本提取。按前置分组（组内 OR、组间 AND）触发批次（最多 N 个并行）；每次构建有完整的分段日志、commit 快照与归档产物及校验和。dirty 源拒建。
- **浏览器终端**：任何浏览器对共享串口设备开会话（xterm.js），设备互斥、忙拒绝显示占用者、空闲超时、第二查看者只读跟随
- **可配置 expect 刷机**：把 VBS 刷机脚本变成可复用、可审查的步骤序列（等待字符串 + 可选发送 + 超时 + 失败动作 + 延时 + 控制字符编码），网页表单编辑，实时进度、二次确认中止
- **人为确认卡**：「LED 亮了吗」「拔掉电源线」作为正式步骤记录——PASS/FAIL 加备注，永不超时消失，进会话记录
- **可靠的会话日志**：逐字节、带时间戳的 RX/TX 行，缓冲 + 定时双刷新，HTTP Range 断点续传下载；写失败标记「日志不完整」而不是静默丢数据
- **虚拟演示设备**：零硬件跑通全流程——expect 刷机与真实构建（git 克隆、产物归档）都能跑

## 快速开始（零硬件）

需要 Go 1.27+ 与 Node.js（仅构建前端；预构建版本已内嵌）。

```bash
# 1. 构建服务器（前端产物经 embed 打进二进制）
cd web && npm install && npm run build && cd ..
go build -o embedflow-server ./server/cmd/embedflow-server

# 2. 以内置虚拟设备启动
./embedflow-server -demo -addr :8420

# 3. 打开 http://localhost:8420，登录（默认 admin/admin——请立即改密码）
```

`-demo` 模式会注册一块 `demo-virtual` 虚拟板卡、预置 `demo-flash` expect 规则，并预置 `demo-build` 构建项目（源是数据目录下自带的 git 仓库）。打开设备、点运行，就能在终端里看到完整的 uboot 风格刷机流程；打开构建板块触发一次真实构建，拿到带校验和的归档产物——不需要真机，不需要串口线。

## 连接真实设备

1. 在编译服务器上运行 `./embedflow-server`（不加 `-demo`）
2. 在测试人员的 Windows 电脑上运行串口客户端（`serial-client -server http://<服务器>:8420 -user U -pass P -device N -port COM3`）——它作为设备的串口传输层，纯透传，无业务逻辑
3. 任意浏览器打开设备、开会话

## 架构

```
┌────────────┐  串口    ┌───────────────┐  WebSocket  ┌──────────────────┐
│  目标设备    ├─────────►│  串口客户端     ├────────────►│   服务器          │
│            │  (COM)   │  (Windows PC) │ 控制帧:JSON  │  (编译服务器)      │
└────────────┘          └───────────────┘ 数据帧:binary └────────┬─────────┘
                                                     │  SQLite + 日志
                                                            ┌─────┴─────┐
                                                            │  浏览器     │
                                                            │ (xterm.js) │
                                                            └───────────┘
```

- **领域内核**（`server/internal/session`）：纯状态机、无 IO——设备互斥、会话生命周期、断线分支、双阈值存活检测
- **执行器层**（`transport` / `sessionlog` / `expect` / `serialport`）：网络、持久化、引擎、串口抽象
- **共享协议**（`pkg/protocol`）：服务端与客户端共用一套线上协议（WebSocket 双帧：控制帧 JSON、数据帧 binary）

## 开发

```bash
go work sync
go test ./...              # 全部模块
cd web && npm run dev      # Vite dev server（带 API 代理）
```

## 路线图

M1（会话层）与 M2（构建模块）已完成：会话、日志、expect 执行、人为确认、构建项目、批次、构建记录与产物，虚拟演示设备同时覆盖刷机与构建。后续：SSH 部署模块与工程级刷机规则（M3）、测试模块与报告（M4）、发布模块（M5）。完整设计与决策记录见 `docs/designs/office-hours-design.md`。

## 许可证

MIT
