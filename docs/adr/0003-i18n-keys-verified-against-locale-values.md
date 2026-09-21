# ADR-0003: 界面文案一律走 i18n key，借用他板块 key 前核对语义

日期: 2026-09-21
状态: 已采纳

## 背景

构建记录页两处借用 `history.download`（会话历史页的「下载完整日志」）渲染产物下载与日志下载链接；另一处引用了**根本不存在**的 key `build.logNoun`。生产上用户看到「下载完整日志build.logNoun」「产物下载链接写着下载完整日志」两处文案错误。

vue-i18n 的回退行为是把缺失 key 的**key 名原样渲染**给用户——不是报错、不是空白。这类缺陷静态检查（类型检查、API 测试）全部抓不到：key 存在性、数据链路、接口返回全对，错的只是渲染出来的字符串。

## 决策

1. 所有界面文案一律走 i18n key（`web/src/locales/zh.ts` / `en.ts`），禁止在模板里硬编码可见文案。
2. key 与消费页面一一对应；跨板块借用前先核对该 key 的实际文案值是否语义一致（`history.download` = 「下载完整日志」只适用于会话历史页；构建记录页的产物下载用 `build.downloadArtifacts`，日志下载用 `history.downloadLog`）。
3. 新增文案时 zh/en 两文件必须同补；删除功能时顺带确认 key 无消费方再删。

## 验证手段（防复发）

- 改动文案后检索渲染产物：构建后在 `web/dist/assets/index-*.js` 里 grep 新文案的中文值（不是 key 名），确认进了 bundle；
- 全量扫描脚本思路：抽取所有模板 `t('x.y')` 引用，与 locale 文件定义比对，报告缺失 key（2026-09-21 手工执行过一次，扫出 `deploy.sshDevice`/`paramName`/`paramDefault` 三个缺失并已补）。

## 后果

- 正向：用户可见文案不再出现裸 key 名；双语完整性可机器校验。
- 约束：code review 时凡涉及 `t('...')` 的新增/变更，须同时确认 zh/en 两处定义与语义。
