# Web 网页变更监听与推送系统 — 需求规格说明书

## 1. 概述

定时监听网页内容变更，通过 jQuery 脚本提取数据，变更时邮件推送通知。

| 项目 | 内容 |
|------|------|
| 技术栈 | Go + Gin + chromedp / Vue 3 + Element Plus |
| 数据库 | SQLite 3 |
| 浏览器引擎 | Chromium (headless, chromedp) |
| 部署 | Docker 单容器 |

## 2. 功能模块

### FM-01 监听任务管理

- 任务列表：搜索、分页、新增、编辑、删除、手动执行、查看历史
- 任务字段：名称、目标URL、执行频率(5分钟~每天7档)、状态、邮件通知开关、备注

### FM-02 jQuery脚本规则

每条任务包含多个步骤，按序执行。

| 字段 | 说明 |
|------|------|
| 步骤排序 | 自动编号，支持拖拽调整 |
| jQuery脚本 | 用户编写的jQuery代码 |
| 执行模式 | 必须成功(失败报错) / 失败跳过 |

**链式执行**：上一步结果自动存入 `window.__prev_result`(文本) 和 `window.__prev_html`(HTML)，下一步可直接引用。

**元素计数**：测试执行时自动检测脚本中 `$('...')` 选择器匹配的元素数量，匹配0个时明确提示。

### FM-03 扫描历史

- 历史列表：扫描时间、变更状态、内容预览、邮件状态
- 详情对比：本次内容 + Diff对比

### FM-04 系统配置

- SMTP邮件配置：PLAIN / TLS 两种加密方式
- 测试发送

### FM-05 定时调度

- cron表达式映射7档频率
- 并发限制 + 防重入

## 3. 数据模型

**monitor_task**: id, name, target_url, freq_code, status, email_notify, last_scan_time, last_scan_content, remark

**monitor_record**: id, task_id, target_url, scan_time, scan_result, is_changed, error_msg, email_sent

**scan_rule**: id, task_id, step_order, rule_content(jQuery脚本), rule_mode

**sys_config**: id, config_key, config_value(SMTP配置)

## 4. API

| 方法 | 路径 | 说明 |
|------|------|------|
| GET/POST | /api/tasks | 列表/新增 |
| GET/PUT/DELETE | /api/tasks/:id | 详情/编辑/删除 |
| POST | /api/tasks/:id/execute | 手动扫描 |
| GET/POST | /api/tasks/:id/rules | 规则列表/新增 |
| PUT/DELETE | /api/rules/:id | 编辑/删除规则 |
| POST | /api/rules/test | 测试jQuery脚本 |
| GET | /api/tasks/:id/records | 扫描历史 |
| GET | /api/records/:id | 历史详情 |
| GET/PUT | /api/config/smtp | SMTP配置 |
| POST | /api/config/smtp/test | 测试邮件 |
| POST | /api/cache | 缓存页面HTML |
