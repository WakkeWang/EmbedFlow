# ADR-0001: 批次与构建记录 ID 由内核发放，落库时显式写入

日期: 2026-09-21
状态: 已采纳

## 背景

批次（batches）与构建记录（build_records）的 ID 原本有两个独立来源：

1. **批次内核**（`server/internal/batch`，纯内存状态机）用 `nextBatch`/`nextRecord` 计数器发号；
2. **SQLite** 靠 AUTOINCREMENT 在 INSERT 时发号。

两者仅靠启动时的「播种仪式」维持一致：`sweepStartup` 从持久化数据读 `MAX(id)` 回填内核计数器（`TestBuildE2E_KernelCountersSurviveRestart` 针对过一次相邻缺陷）。

这条脆弱的等式在生产 192.168.0.80 上两次断裂：

- **第一次（e64893c 修的形态）**：历史全为终态时重启，内核计数器从 0 重新发号，新批次落库撞历史主键，批次成僵尸「running」。
- **第二次（本 ADR 的直接动因，2026-09-21）**：操作员删光批次历史后重启。SQLite 的 `sqlite_sequence` 水位不回退（AUTOINCREMENT 语义），下一个自增 ID 仍是 14；而内核从 `MAX(id)=0` 播种后从 1 发号。于是内核批次 #1 落库成 #14，`kernel.PlanRecords(14)` 查不到该批次 → 一条构建记录都不落库 → 批次永远「running」，且取消无效（`kernel.Cancel(14)` 同样查不到内核批次）。内核 id=1 的记录事件又因「no rows」失败，泄漏一个并行槽位。

## 决策

**统一 ID 空间：内核发号即数据库 ID。** `CreateBatch`/`CreateBuildRecord` 支持 `ID > 0` 时显式写入该值（AUTOINCREMENT 列接受显式值，并自动抬高 `sqlite_sequence`）；`buildOrchestrator.CreateBatch` 先 `kernel.Enqueue` 拿到内核发放的批次 ID 与记录 ID，再带着这些 ID 落库，最后应用内核事件。

配套改动：

- 落库失败时 `kernel.Cancel` 回收内核批次并**丢弃其事件**（行不存在，事件不能应用），避免泄漏并行槽位。
- 播种仪式保留为双保险（`MAX(id)` 播种 + 显式 ID 落库），但一致性不再依赖它。
- 临时目录以构建记录 ID 命名（`tmp/build-<记录ID>`），ID 不复用后不会与「失败保留 7 天」的旧目录撞名（历史上曾发生 `destination already exists` 撞名失败）。

## 否决的替代方案

- **播种改读 `sqlite_sequence`**：能修第二次事故，但「一个实体两套 ID + 三个补丁仪式」的结构不变，第四次断裂只是时间问题。拒绝。
- **数据库 ID 回填内核**（落库后再把自增 ID 告诉内核）：需要在落库和内核之间做二次映射表，事件载荷携带的内核 ID 仍需翻译，复杂度更高且事件顺序更难保证。拒绝。

## 后果

- 正向：内核 ID = 数据库 ID 成为结构不变量，任何「删历史 + 重启」组合都不再产生分叉；编号对用户单调递增，删除历史不影响新编号。
- 约束：任何新代码插入批次/记录行时**必须携带内核发放的 ID**，禁止走自增路径（自增路径仅为测试/种子保留）。测试 `TestBuildE2E_BatchIDsSurviveHistoryDelete` 复刻「跑一批 → 删光 → 重启 → 再触发」守护此不变量。
- 迁移：无需（旧行不动，新行从下一次创建起携带内核 ID）。
