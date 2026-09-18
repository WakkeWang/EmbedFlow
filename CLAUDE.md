# EmbedFlow

中小嵌入式团队的内部构建/测试/发布系统，开源。需求文档见 `docs/需求文档.md`，术语表见 `CONTEXT.md`，设计文档见 `docs/designs/office-hours-design.md`。

## Agent skills

### Issue tracker

Issues live in this repo's GitHub Issues, via the `gh` CLI. See `docs/agents/issue-tracker.md`.

### Triage labels

The five canonical triage roles use their default label strings (`needs-triage`, `needs-info`, `ready-for-agent`, `ready-for-human`, `wontfix`). See `docs/agents/triage-labels.md`.

### Domain docs

Single-context layout: root `CONTEXT.md` plus `docs/adr/`. See `docs/agents/domain.md`.

## 输出注意

- 向用户提问（AskUserQuestion）时，**问题文本和选项只用中文、英文、数字、半角括号 () 和半角标点**。禁用全角括号（）、「」『』、箭头 →、带圈数字 ①②③、省略号 …… 等一切特殊 Unicode 符号——用户终端对这些字符渲染会出现乱码。列表用 1234 或 a)b)c) 编号。

## 技能路由

当用户请求与可用技能匹配时，通过 Skill 工具调用。拿不准时也调用。

关键路由规则：
- 产品想法/脑暴 → /office-hours
- 战略/范围 → /plan-ceo-review
- 架构 → /plan-eng-review
- 设计系统/方案评审 → /design-consultation 或 /plan-design-review
- 完整评审流水线 → /autoplan
- Bug/报错 → /investigate
- QA/测试站点行为 → /qa 或 /qa-only
- 代码评审/diff 检查 → /review
- 视觉打磨 → /design-review
- 发布/部署/PR → /ship 或 /land-and-deploy
- 保存进度 → /context-save
- 恢复上下文 → /context-restore
- 撰写可入 backlog 的规格/issue → /spec
