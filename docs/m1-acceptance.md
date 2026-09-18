# M1 真机验收手册（issue #12）

> M1 验收金线：把真实 VBS 刷机脚本经编辑 UI 录入成 expect 规则，真机端到端跑通。
> 本手册把验收拆成带勾选框的步骤，按顺序执行，全部勾完即可关闭 #12。
> 执行前提：#2 的延迟 spike 已出 GO 结论（见 `docs/spike-latency.md`）。

## 0. 环境准备

- [ ] 编译服务器上跑服务端：`./embedflow-server -addr :8420`（不加 `-demo`）
- [ ] 测试机 Windows 上跑客户端：`serial-client -server http://<server>:8420 -user U -pass P -device N -port COM3 -baud 115200`
- [ ] Web 界面设备列表：该设备显示「在线」（串口已共享）
- [ ] 修改默认 admin 密码，另建一个成员账号

## 1. 录入规则（必须走编辑 UI，不许手写 JSON 灌库）

把 VBS 脚本（P0133 系列（内部代称），`reboot → uboot → ramdisk → recovery.sh → 二次登录`）翻译成步骤序列。对照表：

| VBS 原文语义 | expect 步骤 | 要点 |
|---|---|---|
| `crt.Screen.WaitForString "Loading Environment from EEPROM... OK"` | await=`Loading Environment from EEPROM... OK`，超时 10s | reboot 后第一拍 |
| 发 uboot 密码 | send=`<密码>\r`，**勾 secret**（issue #14：日志掩码） | 密码步骤一律勾 secret |
| `WaitForString "<board>#"` | await=`<board>#`，超时 30s | uboot 阶段输出慢，超时放宽 |
| 发 `usb reset; fatload usb 0 $loadaddr <ver>/<folder>/board.itb && run ramargs && bootm` + `\r` | send 原文（`\r` 结尾） | 路径按实际版本目录改 |
| `WaitForString "Please press Enter to activate this console."` → 发回车 | await + send=`\r` | |
| `WaitForString "<board> login:"` → 发用户名 | await + send=`<user>\r` | |
| 发密码 | send=`<密码>\r`，**勾 secret** | |
| 等 shell 提示符 → 发 `sh /run/media/sda1/.../recovery.sh` | await 提示符 + send | |
| `WaitForString "install over"` | await=`install over`，超时 **600s**（安装最慢的一步，勿用默认 15s） | |
| 发 `reboot`，等 `localhost login:`，二次登录 | 同上模式 | |

编辑器操作检查项：

- [ ] 步骤卡片增删/上移/下移/复制正常
- [ ] 发送内容用按钮插入 `\r`（不手打转义）
- [ ] 保存后刷新页面，规则仍在（持久化）
- [ ] secret 步骤在 UI 与普通步骤无差别（掩码只发生在日志层）

## 2. 端到端刷机

- [ ] 打开设备会话，点 Run，选中刚录入的规则
- [ ] 进度条逐步骤推进（issue #8），步骤名可读
- [ ] **uboot 阶段无超时误报**（验收核心：时序余量，依赖 spike GO 结论）
- [ ] 序列走到「序列完成」徽章（DS-3A 终态）
- [ ] 运行期间人工键入被拒绝（CEO-17A 输入闸）

## 3. 日志与人机混合步骤

- [ ] 会话历史里找到本次任务会话，日志可下载，含完整 RX/TX
- [ ] **日志里两处密码步骤显示 `TX | ***`，无明文密码**（issue #14）
- [ ] **确认目标设备密码不回显**：下载日志搜密码明文（RX 行）。部分 uboot/bootloader 会回显输入字符——若回显，TX 掩码形同虚设，需要在 issue #14 补 RX 侧掩码策略
- [ ] **手动会话不要手输密码**：Web 终端的每一次按键都明文进日志（M1 无手动输入掩码）；密码一律走 expect 规则的 secret 步骤。此限制记入 issue #14，M2/M3 给终端加密码输入模式
- [ ] 会话中插一张人为确认卡（如「LED 灯亮？」），PASS + 备注，卡片进记录
- [ ] 人为确认时序：刷机中途插卡 → 运行不受影响 → 卡片保持 pending 直到人工点 PASS/FAIL

## 4. 异常路径（各验一次）

- [ ] 刷机中拔网线 10 秒再插回：任务会话经 30s 容忍窗存活（CEO-16A），序列继续
- [ ] 刷机中拔网线 >30 秒：会话标失败终止，日志尾部标注，设备回空闲
- [ ] 中止路径：点中止 → 二次确认（半刷警告）→ 会话 FAILED
- [ ] 忙拒绝：第二浏览器开会话 → 显示占用者信息 → 只读跟随可见同一终端

## 5. 验收记录归档

- [ ] 全部勾选后：截图（终端运行态 + 进度条 + 完成徽章 + 日志掩码行）存 `docs/assets/m1-acceptance/`
- [ ] 在 issue #12 贴验收结论（日期、硬件型号、规则步骤数、总耗时）
- [ ] 关闭 #12

## 已知边界（验收时不阻塞）

- 客户端目前是 CLI（Wails GUI 在 M2 补），测试机用命令行启动即可
- 会话日志密码掩码为整行掩码（`TX | ***`），不保留长度信息——防止长度侧信道
- 空闲超时只作用于手动会话；任务会话按 expect 步骤超时走（CEO-1A）
