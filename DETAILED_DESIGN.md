# Web 网页变更监听与推送系统 — 详细设计

## 1. 架构

```
Vue 3 前端 → Gin API → Service → DAO/GORM → SQLite
                        ↓
                   chromedp → Chromium (headless)
                        ↓
                   jQuery 脚本执行
                        ↓
                   SMTP 邮件通知
```

## 2. 技术选型

| 层 | 技术 | 说明 |
|----|------|------|
| 后端 | Go 1.21+ / Gin | HTTP框架 |
| ORM | GORM + SQLite | 数据持久化 |
| 浏览器 | chromedp + Chromium | 无头浏览器，执行jQuery |
| 调度 | robfig/cron v3 | 定时任务 |
| 邮件 | gomail v2 | SMTP PLAIN/TLS |
| 前端 | Vue 3 + Element Plus + TypeScript | SPA |
| 构建 | Vite | 前端打包 |

## 3. 目录结构

```
server/
├── main.go              # 入口,路由,embed前端
├── config/              # viper配置
├── internal/
│   ├── model/           # MonitorTask,ScanRule,MonitorRecord,SysConfig
│   ├── dao/             # 数据访问
│   ├── service/         # 业务逻辑+校验
│   ├── handler/         # HTTP控制器(task/rule/record/scan/config)
│   ├── engine/          # chromedp + jQuery执行引擎
│   ├── scheduler/       # cron定时调度
│   ├── notifier/        # SMTP邮件
│   └── middleware/      # CORS/Logger
├── pkg/crypto/          # AES密码加密
├── pkg/response/        # 统一响应
└── web/dist/            # 前端构建产物(embed)
web/
└── src/
    ├── views/           # TaskList,TaskForm,RecordList,SysConfig
    ├── components/      # RuleFormDialog,DiffViewer
    ├── api/             # Axios封装
    ├── types/           # TypeScript类型
    └── router/          # Vue Router
```

## 4. 核心引擎

### jQuery执行流程

```
1. chromedp 启动 headless Chrome
2. 导航到目标URL
3. 检测并注入 jQuery 3.7.1 (CDN, 轮询等待加载)
4. 初始化 window.__prev_result = ''; window.__prev_html = ''
5. 按步骤执行:
   for each rule:
     chromedp.Evaluate(rule.RuleContent)  →  result
     存入 window.__prev_result = result
     存入 window.__prev_html = result
     自动检测 $(...) 选择器匹配数量
6. 返回最终结果
```

### 元素计数

正则提取脚本中第一个 `$('...')` 或 `$("...")` 选择器，执行 `.length` 获取匹配数量。

### 链式引用

```javascript
// 步骤1: 提取文章HTML
$('.article-body').html()

// 步骤2: 从HTML中提取链接 (window.__prev_html = 步骤1的HTML)
$(window.__prev_html).find('a').attr('href')

// 步骤3: 文本替换 (window.__prev_result = 步骤2的链接)
window.__prev_result.replace('http://', 'https://')
```

## 5. 数据库

```sql
CREATE TABLE monitor_task (id INTEGER PK, name VARCHAR(100), target_url VARCHAR(2048),
  freq_code INTEGER DEFAULT 4, status INTEGER DEFAULT 1, email_notify INTEGER DEFAULT 1,
  last_scan_time DATETIME, last_scan_content TEXT, remark VARCHAR(500),
  created_at DATETIME, updated_at DATETIME);

CREATE TABLE scan_rule (id INTEGER PK, task_id INTEGER FK, step_order INTEGER,
  rule_content TEXT, rule_mode INTEGER DEFAULT 1, created_at DATETIME);

CREATE TABLE monitor_record (id INTEGER PK, task_id INTEGER FK, target_url VARCHAR(2048),
  scan_time DATETIME, scan_result TEXT, is_changed INTEGER DEFAULT 0,
  error_msg VARCHAR(1000), email_sent INTEGER DEFAULT 0, created_at DATETIME);

CREATE TABLE sys_config (id INTEGER PK, config_key VARCHAR(100) UNIQUE,
  config_value TEXT, updated_at DATETIME);
```

## 6. 部署

Docker多阶段构建: node构建前端 → golang构建后端 → debian+chromium运行镜像

```bash
docker build -t monitor-server -f server/Dockerfile .
docker run -d -p 8080:8080 -v ./data:/app/data monitor-server
```
