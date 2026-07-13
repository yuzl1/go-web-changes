# Web 网页变更监听与推送系统 — 开发任务与测试任务

> 基于 [REQUIREMENTS.md](./REQUIREMENTS.md) 和 [DETAILED_DESIGN.md](./DETAILED_DESIGN.md) 拆分。
> 勾选框：开发完成后在 `[x]` 中标记 `x`。

---

## 一、后端项目搭建

### 1.1 项目骨架

- [x] **1.1.1** Go Module 初始化 — 创建 `server/` 目录，`go mod init monitor-server`
- [x] **1.1.2** 目录结构搭建 — 按照详设 3.1 节创建 `internal/handler/`、`internal/service/`、`internal/dao/`、`internal/model/`、`internal/engine/`、`internal/scheduler/`、`internal/notifier/`、`internal/middleware/`、`pkg/crypto/`、`pkg/response/`、`config/`、`data/`、`logs/`
- [x] **1.1.3** 配置文件 — 实现 `config/config.go` + `config/config.yaml`，使用 viper 加载
- [x] **1.1.4** Gin 路由骨架 — `main.go` 启动 Gin 服务，注册 CORS 中间件、日志中间件
- [ ] **1.1.5** 统一响应封装 — `pkg/response/response.go`（code/message/data 结构 + 分页封装）

### 1.2 数据库

- [ ] **1.2.1** GORM 初始化 — 连接 SQLite，启用 `PRAGMA foreign_keys = ON`
- [ ] **1.2.2** AutoMigrate — 使用 GORM AutoMigrate 自动建表（monitor_task、monitor_record、scan_rule、sys_config）
- [ ] **1.2.3** 初始数据 — sys_config 表写入 SMTP 配置默认行

---

## 二、监听任务 CRUD（后端）

### 2.1 Model 层

- [ ] **2.1.1** `internal/model/task.go` — MonitorTask 结构体定义（GORM 标签）
- [ ] **2.1.2** `internal/model/record.go` — MonitorRecord 结构体定义
- [ ] **2.1.3** `internal/model/rule.go` — ScanRule 结构体定义
- [ ] **2.1.4** `internal/model/config.go` — SysConfig 结构体定义

### 2.2 DAO 层

- [ ] **2.2.1** `internal/dao/task.go` — 分页查询（关键字模糊搜索+状态筛选）、按ID查询、创建、更新、删除
- [ ] **2.2.2** `internal/dao/rule.go` — 按 task_id 查规则列表、批量创建、更新、按ID删除、按 task_id 批量删除
- [ ] **2.2.3** `internal/dao/record.go` — 分页查询（按 task_id）、按ID查询（含上一条记录）、创建
- [ ] **2.2.4** `internal/dao/config.go` — 按 key 查询、按 key 批量查询、更新/插入

### 2.3 Service 层

- [ ] **2.3.1** `internal/service/task.go` — 列表查询、详情查询（含规则）、创建（含规则，事务）、更新（删除旧规则+插入新规则，事务）、删除（级联）
- [ ] **2.3.2** `internal/service/rule.go` — 列表查询、新增、编辑、删除、重新编号
- [ ] **2.3.3** `internal/service/record.go` — 历史列表、历史详情（含 prev_scan_result）
- [ ] **2.3.4** `internal/service/config.go` — 获取 SMTP 配置、保存 SMTP 配置（密码加密）

### 2.4 Handler 层

- [ ] **2.4.1** `internal/handler/task.go` — 注册路由：GET /api/tasks、POST /api/tasks、GET /api/tasks/:id、PUT /api/tasks/:id、DELETE /api/tasks/:id
- [ ] **2.4.2** `internal/handler/rule.go` — 注册路由：GET /api/tasks/:id/rules、POST /api/tasks/:id/rules、PUT /api/rules/:id、DELETE /api/rules/:id
- [ ] **2.4.3** `internal/handler/record.go` — 注册路由：GET /api/tasks/:id/records、GET /api/records/:id
- [ ] **2.4.4** `internal/handler/config.go` — 注册路由：GET /api/config/smtp、PUT /api/config/smtp

### 2.5 参数校验

- [ ] **2.5.1** URL 格式校验（必须以 http:// 或 https:// 开头）
- [ ] **2.5.2** 任务创建时校验 rules 至少包含一条结果型规则（类型1或4）
- [ ] **2.5.3** 规则内容非空校验、step_order 合法性校验

---

## 三、扫描引擎

### 3.1 chromedp 集成

- [ ] **3.1.1** `internal/engine/engine.go` — chromedp 上下文创建（超时 30s）、浏览器 Tab 管理、资源释放
- [ ] **3.1.2** Chrome 进程管理 — 启动时检查 Chrome 路径，配置 headless 参数（`--headless`、`--no-sandbox`、`--disable-gpu`）
- [ ] **3.1.3** 并发控制 — 信号量 channel 限制最大 3 并发扫描

### 3.2 规则执行器

- [ ] **3.2.1** `internal/engine/rules.go` — 规则执行入口函数 `ExecuteRules(ctx, targetURL, rules)`
- [ ] **3.2.2** 类型1 执行 — chromedp.Text() 实现 CSS 选择器文本提取
- [ ] **3.2.3** 类型2 执行 — chromedp.Evaluate() 实现 JS 代码执行
- [ ] **3.2.4** 类型3 执行 — chromedp.WaitVisible() 实现等待元素出现（最长等待 10s）
- [ ] **3.2.5** 类型4 执行 — chromedp.AttributeValue() 实现属性值提取
- [ ] **3.2.6** 结果汇总 — 取最后一条结果型规则的输出作为 final_result

### 3.3 完整扫描流程

- [ ] **3.3.1** `ExecuteScan(task, rules)` — 页面导航 → 规则执行 → 结果提取 → 变更比对
- [ ] **3.3.2** 变更检测 — trim 后字符串对比
- [ ] **3.3.3** 扫描记录写入 — is_changed=1/0/-1 三种状态处理
- [ ] **3.3.4** 异常处理 — 超时/选择器未匹配/JS错误/chromedp崩溃，记录 error_msg

---

## 四、定时调度引擎

- [ ] **4.1** `internal/scheduler/scheduler.go` — robfig/cron 初始化 + freq_code 到 cron 表达式的映射
- [ ] **4.2** 启动时加载所有 status=1 的任务并注册定时器
- [ ] **4.3** 任务执行入口 `scanTask(task)` — 获取规则 → 调用扫描引擎 → 处理结果 → 发送通知
- [ ] **4.4** 优雅关闭 — `cron.Stop()` + 等待正在执行的任务完成
- [ ] **4.5** 任务运行时状态保护 — 同一任务同一时间只允许一个扫描实例（sync.Map 防重入）

---

## 五、邮件通知（SMTP PLAIN / TLS）

### 5.1 SMTP 发送

- [ ] **5.1.1** `internal/notifier/email.go` — Notifier 接口实现
- [ ] **5.1.2** PLAIN 模式 — gomail.Dialer 不启用 TLS，明文连接
- [ ] **5.1.3** TLS 模式 — 端口 587 STARTTLS / 端口 465 隐式 TLS 自动判断
- [ ] **5.1.4** 邮件模板 — HTML 格式邮件（任务名称、目标URL、扫描时间、变更内容截取前500字符）
- [ ] **5.1.5** 发送失败处理 — 记录日志，不重试，不阻断扫描流程

### 5.2 密码加密

- [ ] **5.2.1** `pkg/crypto/aes.go` — AES-256-GCM 加密/解密函数
- [ ] **5.2.2** 密钥管理 — 优先读取配置文件 secret_key，无配置时使用默认密钥（开发环境）

### 5.3 测试发送

- [ ] **5.3.1** `POST /api/config/smtp/test` — 使用当前数据库 SMTP 配置发送测试邮件
- [ ] **5.3.2** 成功返回成功提示，失败返回具体错误原因

---

## 六、规则在线测试（后端）

- [ ] **6.1** `POST /api/rules/test` Handler — 接收 target_url + rules 数组
- [ ] **6.2** 测试扫描 — 创建独立 chromedp context（超时 20s），执行规则
- [ ] **6.3** 逐步骤返回 — 每步返回 status(success/error)、elapsed_ms、output、error
- [ ] **6.4** 最终结果汇总 — final_result = 最后一条成功的结果型规则输出
- [ ] **6.5** 测试不写库 — 不写入 monitor_record，不更新 monitor_task，不发送邮件

---

## 七、手动执行扫描（后端）

- [ ] **7.1** `POST /api/tasks/:id/execute` Handler — 异步执行扫描
- [ ] **7.2** 返回 scan_result、is_changed、scan_time、record_id
- [ ] **7.3** 变更时触发邮件通知

---

## 八、前端项目搭建

### 8.1 项目初始化

- [ ] **8.1.1** Vite + Vue 3 + TypeScript 项目创建 — `npm create vite@latest`
- [ ] **8.1.2** 安装依赖 — Element Plus、Vue Router、Pinia、Axios、diff、@element-plus/icons-vue、sass
- [ ] **8.1.3** `utils/request.ts` — Axios 封装（baseURL、拦截器、错误处理、统一响应解包）
- [ ] **8.1.4** `types/index.ts` — TypeScript 类型定义（MonitorTask、ScanRule、MonitorRecord、SmtpConfig、RuleTestResult 等）

### 8.2 路由与布局

- [ ] **8.2.1** `router/index.ts` — 路由配置（/tasks、/tasks/new、/tasks/:id/edit、/tasks/:id/records、/config）
- [ ] **8.2.2** `App.vue` — Element Plus Container 布局（左侧导航菜单 + 右侧内容区）
- [ ] **8.2.3** 导航菜单 — 监听列表、系统配置两个菜单项

### 8.3 API 接口层

- [ ] **8.3.1** `api/task.ts` — 任务 CRUD + 手动执行接口封装
- [ ] **8.3.2** `api/rule.ts` — 规则 CRUD + 在线测试接口封装
- [ ] **8.3.3** `api/record.ts` — 扫描历史列表 + 详情接口封装
- [ ] **8.3.4** `api/config.ts` — SMTP 配置获取 + 保存 + 测试发送接口封装

---

## 九、前端页面开发

### 9.1 监听列表页 `TaskList.vue`

- [ ] **9.1.1** 搜索栏 — 关键字输入框 + 搜索按钮 + 新增按钮
- [ ] **9.1.2** 表格展示 — 序号、监听名称、目标URL、执行频率（中文显示）、上次扫描时间、状态开关、操作列
- [ ] **9.1.3** 分页 — Element Plus Pagination，每页10条
- [ ] **9.1.4** 操作按钮 — 编辑（跳转编辑页）、删除（二次确认弹窗 + 级联删除提示）、立即执行（调用 execute 接口 + 刷新列表）、历史（跳转历史页）
- [ ] **9.1.5** 手动执行反馈 — Loading + 执行结果 Toast（有无变更）

### 9.2 新增/编辑任务页 `TaskForm.vue`

- [ ] **9.2.1** 基本信息表单 — 监听名称、目标URL（URL格式校验）、执行频率（下拉选择，显示中文）、状态开关、备注
- [ ] **9.2.2** 新增/编辑模式切换 — 根据路由参数 `:id` 判断；编辑模式下回填数据 + 预加载规则列表
- [ ] **9.2.3** 规则列表区域 `RuleTable.vue` — 表格（步骤排序、规则类型中文名、规则内容截断、操作列）
- [ ] **9.2.4** 规则操作 — 上移/下移调整步骤、编辑、删除
- [ ] **9.2.5** 提交 — 调用创建/更新接口，成功后跳转回列表页
- [ ] **9.2.6** 取消 — 返回列表页，未保存提示确认

### 9.3 规则编辑弹窗 `RuleForm.vue`

- [ ] **9.3.1** 规则类型下拉 — CSS选择器提取 / JavaScript执行 / 等待元素出现 / 属性值提取
- [ ] **9.3.2** 规则内容输入 — 文本域，根据规则类型动态切换 placeholder
- [ ] **9.3.3** 属性名称输入 — 仅规则类型=4时显示
- [ ] **9.3.4** 测试执行按钮 — 调用 `/api/rules/test` 接口

### 9.4 规则测试结果面板 `RuleTestPanel.vue`

- [ ] **9.4.1** Loading 状态 — 测试执行中显示加载动画
- [ ] **9.4.2** 步骤结果列表 — 每步显示：步骤序号、规则类型中文名、成功✓/失败✗图标、耗时(ms)
- [ ] **9.4.3** 失败步骤 — 红色高亮，显示 error 信息
- [ ] **9.4.4** 最终结果展示 — 代码块显示提取到的文本内容（截取前 1000 字符）
- [ ] **9.4.5** 重新测试 — 支持修改规则后再次点击测试执行
- [ ] **9.4.6** 目标URL未填写时 — 测试按钮点击后提示"请先在基本信息中填写目标URL"

### 9.5 扫描历史页 `RecordList.vue`

- [ ] **9.5.1** 面包屑导航 + 任务摘要卡片（任务名称、目标URL）
- [ ] **9.5.2** 历史表格 — 扫描时间、目标URL、是否有变更（绿色✅/灰色➖/红色✗）、操作（查看详情）
- [ ] **9.5.3** 分页 + 按扫描时间倒序

### 9.6 历史详情对比 `RecordDetail.vue`（抽屉/弹窗）

- [ ] **9.6.1** 基本信息 — 扫描时间、目标URL、是否有变更
- [ ] **9.6.2** `DiffViewer.vue` — 左右双栏对比（上次扫描内容 vs 本次扫描内容）
- [ ] **9.6.3** 差异高亮 — 使用 `diff` 库做文本差异计算，新增绿色标记、删除红色标记
- [ ] **9.6.4** 无上次内容时 — 左栏显示"（首次扫描，无历史内容）"

### 9.7 系统配置页 `SysConfig.vue`

- [ ] **9.7.1** SMTP 表单 — 服务器、端口（数字）、加密方式（PLAIN/TLS下拉）、发件人邮箱、发件人名称、邮箱账号、邮箱密码、收件人邮箱
- [ ] **9.7.2** 端口自动提示 — 选择 PLAIN 时提示"常用端口: 25"，选择 TLS 时提示"常用端口: 587/465"
- [ ] **9.7.3** 邮箱格式校验 — 发件人邮箱、收件人邮箱格式校验（收件人支持逗号分隔多个）
- [ ] **9.7.4** 密码 placeholder — 编辑时显示"留空则不修改密码"
- [ ] **9.7.5** 保存按钮 — 调用 PUT 接口
- [ ] **9.7.6** 测试发送按钮 — 调用 POST /api/config/smtp/test，成功 Toast 提示，失败弹窗显示错误

---

## 十、Docker 部署

- [ ] **10.1** Go embed 配置 — `//go:embed web/dist/*` 嵌入前端静态资源，Gin NoRoute fallback 到 index.html
- [ ] **10.2** `server/Dockerfile` — 多阶段构建（node 构建前端 + golang 构建后端 + debian 运行镜像含 Chromium）
- [ ] **10.3** `docker-compose.yml` — 端口映射 + Volume 挂载（data/、logs/）+ 资源限制
- [ ] **10.4** Chromium 环境变量 — `CHROME_BIN=/usr/bin/chromium`、`CHROME_PATH=/usr/bin/chromium`
- [ ] **10.5** `.dockerignore` — 排除 node_modules、.git、logs、data
- [ ] **10.6** 构建脚本 — `build.sh`（一键构建前端+后端+Docker镜像）

---

## 十一、测试任务

### 11.1 后端接口测试

- [ ] **11.1.1** 任务 CRUD 测试 — 创建任务 → 查询列表/详情 → 编辑 → 删除 → 验证级联删除规则和记录
- [ ] **11.1.2** 规则 CRUD 测试 — 新增规则 → 查询列表 → 编辑 → 删除 → 验证 step_order 唯一性
- [ ] **11.1.3** 参数校验测试 — 空名称、非法URL、空规则列表、无结果型规则等边界情况
- [ ] **11.1.4** 分页测试 — 超过10条数据时分页正确性

### 11.2 扫描引擎测试

- [ ] **11.2.1** 静态页面 CSS 提取测试 — 准备一个简单的 HTML 页面，验证类型1选择器提取正确
- [ ] **11.2.2** JS 执行测试 — 验证类型2 `document.querySelector('.btn').click()` 正确执行
- [ ] **11.2.3** 等待元素测试 — SPA 页面（含 setTimeout 动态插入元素），验证类型3等待成功和超时两种场景
- [ ] **11.2.4** 属性提取测试 — 验证类型4提取 href/src/data-* 属性正确
- [ ] **11.2.5** 多步骤组合测试 — 等待元素 → JS点击 → CSS提取 完整链路
- [ ] **11.2.6** 选择器未匹配测试 — 使用不存在的选择器，验证 error_msg 返回
- [ ] **11.2.7** 页面超时测试 — 访问一个慢速页面（>15s），验证超时处理
- [ ] **11.2.8** URL 不可达测试 — 访问不存在的域名，验证异常处理

### 11.3 规则在线测试

- [ ] **11.3.1** 正常测试 — 提交有效 rules + target_url，验证返回每步状态和最终结果
- [ ] **11.3.2** 失败步骤测试 — 某条规则失败，验证返回 error 信息且不影响步骤列表完整性
- [ ] **11.3.3** 不写库验证 — 测试执行后检查 monitor_record 表无新增记录、monitor_task 无更新

### 11.4 定时调度测试

- [ ] **11.4.1** Cron 注册测试 — 验证不同 freq_code 注册的 cron 表达式正确
- [ ] **11.4.2** 定时触发测试 — 设置每5分钟任务，等待触发，验证扫描记录自动生成
- [ ] **11.4.3** 禁用任务测试 — status=0 的任务不会被调度
- [ ] **11.4.4** 防重入测试 — 同一任务上一次扫描未完成时不会重复触发

### 11.5 邮件通知测试

- [ ] **11.5.1** PLAIN 模式测试 — 使用本地 SMTP 服务（如 MailHog/Mailpit）测试明文发送
- [ ] **11.5.2** TLS 模式测试 — 使用 QQ邮箱/163邮箱 测试 STARTTLS 发送
- [ ] **11.5.3** 测试发送功能 — 系统配置页点击测试发送，验证收到邮件
- [ ] **11.5.4** 变更通知测试 — 执行扫描发现变更后，验证自动发送通知邮件
- [ ] **11.5.5** 邮件发送失败测试 — 断开网络/错误密码，验证不影响扫描记录存储
- [ ] **11.5.6** 密码加密存储测试 — 验证数据库中 smtp_password 为密文，解密后可用

### 11.6 前端页面测试

- [ ] **11.6.1** 监听列表页 — 搜索、分页、新增/编辑/删除/手动执行/历史 按钮功能正常
- [ ] **11.6.2** 新增任务 — 填写表单 + 添加多条规则 + 提交 → 列表刷新
- [ ] **11.6.3** 编辑任务 — 回填数据正确、修改规则后保存生效
- [ ] **11.6.4** 删除任务 — 二次确认弹窗 + 级联删除验证
- [ ] **11.6.5** 规则测试面板 — 测试执行 Loading → 步骤结果展示 → 失败错误提示 → 最终结果展示
- [ ] **11.6.6** 目标URL未填写时测试 — 弹窗提示"请先填写目标URL"
- [ ] **11.6.7** 扫描历史页 — 列表正确、变更状态图标正确
- [ ] **11.6.8** 历史详情对比 — Diff 左右双栏、差异高亮（红/绿标记）
- [ ] **11.6.9** 系统配置页 — SMTP 表单填写、端口自动提示、测试发送按钮、密码留空不修改
- [ ] **11.6.10** 响应式布局 — 窗口缩小时布局不错乱

### 11.7 端到端测试

- [ ] **11.7.1** 完整链路 — 创建任务 → 配置规则 → 测试规则 → 手动执行 → 查看历史 → 对比详情
- [ ] **11.7.2** 定时执行链路 — 创建任务（每5分钟）→ 等待触发 → 查看历史是否有新增记录
- [ ] **11.7.3** 变更通知链路 — 修改目标网页内容 → 手动执行 → 检查是否收到邮件 → 查看历史对比
- [ ] **11.7.4** 并发扫描 — 同时手动执行 5 个任务，验证最多 3 个并发，其余排队

### 11.8 Docker 部署测试

- [ ] **11.8.1** 镜像构建 — `docker build` 成功，镜像大小在合理范围
- [ ] **11.8.2** 容器启动 — 容器正常启动，8080 端口可访问
- [ ] **11.8.3** 前端访问 — 浏览器打开 `http://localhost:8080` 正常显示页面
- [ ] **11.8.4** API 访问 — `curl http://localhost:8080/api/tasks` 返回正常
- [ ] **11.8.5** Chrome 可用 — 手动执行扫描成功（确认容器内 Chromium 正常工作）
- [ ] **11.8.6** 数据持久化 — 重启容器后，任务和数据不丢失
- [ ] **11.8.7** 日志输出 — `docker logs` 查看日志正常输出

---

## 十二、规则增强（开发中新增）

### 12.1 执行模式 (rule_mode)

- [x] **12.1.1** `ScanRule` 模型增加 `rule_mode` 字段（1=必须存在, 2=存在则执行）
- [x] **12.1.2** 规则执行器支持跳过模式 — 元素不存在时 skip 而非 error
- [x] **12.1.3** 测试结果增加 `count`（匹配元素数）、`skipped`、`skip_reason` 字段
- [x] **12.1.4** 规则编辑弹窗增加执行模式下拉（必须存在/存在则执行）
- [ ] **12.1.5** 测试：验证 mode=1 元素不存在时报错，mode=2 元素不存在时跳过

### 12.2 页面元素探查

- [x] **12.2.1** 后端 `POST /api/rules/inspect` 接口 — 探查页面 CSS 选择器
- [x] **12.2.2** `InspectPage()` 引擎函数 — JS 收集 id/class/标签信息
- [x] **12.2.3** 前端 `inspectPage()` API 封装
- [x] **12.2.4** 规则编辑弹窗增加「查看页面元素」按钮
- [x] **12.2.5** 页面元素探查弹窗 — 分类展示 ID/Class/标签，点击填入规则
- [ ] **12.2.6** 测试：验证探查结果正确展示，点击选择器正确填入

### 12.3 测试结果增强

- [x] **12.3.1** 步骤结果展示匹配元素数量（count）
- [x] **12.3.2** 步骤结果区分为 success / error / skipped 三种状态
- [x] **12.3.3** 跳过步骤显示黄色警告 + 跳过原因
- [x] **12.3.4** 成功步骤展示提取内容预览
- [x] **12.3.5** 最终汇总结果绿色左边框高亮
- [ ] **12.3.6** 测试：验证三种状态的视觉区分正确

---

## 十三、jQuery 注入与提取（开发中新增）

### 13.1 jQuery 自动注入

- [x] **13.1.1** `MonitorTask` 模型增加 `inject_jquery` 字段（1=注入, 0=不注入）
- [x] **13.1.2** 引擎 `injectJQuery()` 函数 — CDN 加载 jQuery 3.7.1 + noConflict
- [x] **13.1.3** `navigatePage()` 统一入口 — 导航后按需注入 jQuery
- [x] **13.1.4** ExecuteScan / ExecuteTest / InspectPage 统一走 navigatePage
- [x] **13.1.5** 前端任务表单增加「注入jQuery」开关，默认开启
- [ ] **13.1.6** 测试：验证启用注入后页面 $ 可用，不注入时 $ 为 undefined

### 13.2 规则类型5: jQuery选择器提取

- [x] **13.2.1** `RuleTypeMap` 增加类型5: jQuery选择器提取
- [x] **13.2.2** `executeJQueryExtract()` 执行函数 — 自动检测元素存在、自动补 .text()
- [x] **13.2.3** IsResultType 包含类型5（结果型规则）
- [x] **13.2.4** 校验逻辑更新（ruleType 范围 1-5）
- [x] **13.2.5** 前端规则弹窗增加类型5选项和 placeholder
- [ ] **13.2.6** 测试：验证 `$('selector').text()` / `.html()` / `.attr()` 语法正确执行
- [ ] **13.2.7** 测试：验证纯选择器自动加 .text() 逻辑

---

## 任务统计

| 模块             | 开发任务数 | 测试任务数 |
| ---------------- | ---------- | ---------- |
| 后端项目搭建     | 8          | 0          |
| 任务 CRUD        | 13         | 0          |
| 扫描引擎         | 12         | 0          |
| 定时调度         | 5          | 0          |
| 邮件通知         | 7          | 0          |
| 规则在线测试     | 5          | 0          |
| 手动执行         | 3          | 0          |
| 前端项目搭建     | 9          | 0          |
| 前端页面         | 22         | 0          |
| Docker 部署      | 6          | 0          |
| 规则增强(新增)   | 13         | 0          |
| jQuery注入(新增) | 6          | 0          |
| jQuery提取(新增) | 5          | 0          |
| 后端接口测试     | 0          | 4          |
| 扫描引擎测试     | 0          | 8          |
| 规则测试测试     | 0          | 3          |
| 定时调度测试     | 0          | 4          |
| 邮件通知测试     | 0          | 6          |
| 前端页面测试     | 0          | 10         |
| 端到端测试       | 0          | 4          |
| Docker 部署测试  | 0          | 7          |
| 规则增强测试     | 0          | 3          |
| jQuery测试       | 0          | 3          |
| **合计**         | **114**    | **52**     |
| **总计**         |            | **166**    |
