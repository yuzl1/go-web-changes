# Web 网页变更监听与推送系统

定时监听网页内容变更，通过 jQuery 脚本提取数据，变更时邮件推送。

## 技术栈

| 层 | 技术 |
|----|------|
| 后端 | Go 1.21+ / Gin / chromedp / robfig/cron |
| 前端 | Vue 3 / TypeScript / Element Plus / Vite |
| 数据库 | SQLite 3 + GORM |
| 浏览器 | Chromium (headless) |
| 部署 | Docker 单容器 |

## 快速开始

```bash
# 本地开发
cd server && go run main.go     # :8080
cd web && npm install && npm run dev  # :5173

# Docker 部署
docker run -d -p 8080:8080 -v $(pwd)/data:/app/data yu1998/go-web-changes:latest
```

访问 `http://localhost:8080`

## Docker 镜像构建与推送

### 多平台构建（linux/amd64 + linux/arm64）

```bash
# 1. 创建 buildx 构建器
docker buildx create --name multiarch --use

# 2. 构建并推送多平台镜像到 Docker Hub
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -t yu1998/go-web-changes:latest \
  -f server/Dockerfile \
  --push .

# 3. 单独构建本地平台（不推送）
docker build -t yu1998/go-web-changes:latest -f server/Dockerfile .
```

### 手动构建（网络受限时）

```bash
# 后端
cd server && CGO_ENABLED=1 go build -o monitor-server .

# 前端
cd web && npm install && npm run build

# 打包
tar czf go-web-changes.tar.gz -C server monitor-server config/config.yaml -C ../web dist
```

### 镜像说明

| 标签 | 说明 |
|------|------|
| `yu1998/go-web-changes:latest` | 最新版本（多架构） |
| 基础镜像 | debian:bookworm-slim + Chromium |
| 端口 | 8080 |
| 数据卷 | /app/data (SQLite), /app/logs |

## 核心功能

- **jQuery脚本规则** — 编写jQuery代码提取网页数据，支持多步骤链式执行
- **链式引用** — 上一步结果通过 `window.__prev_result` / `window.__prev_html` 传递
- **元素计数** — 测试时自动显示选择器匹配元素数量
- **定时扫描** — 5分钟~每天 7档频率，cron调度
- **变更通知** — SMTP邮件(PLAIN/TLS)，每任务独立开关
- **扫描历史** — 内容预览 + Diff对比

## 规则示例

```javascript
// 步骤1: 提取标题
$('h1').text()

// 步骤2: 提取价格
$('div.price > span.current').text()

// 步骤3: 从上一步HTML提取链接
$(window.__prev_html).find('a').attr('href')
```

## 项目结构

```
server/                # Go后端
├── internal/
│   ├── engine/        # chromedp + jQuery执行
│   ├── scheduler/     # 定时调度
│   ├── notifier/      # 邮件通知
│   └── handler/       # API控制器
web/                   # Vue3前端
└── src/views/         # 页面组件
```

## API

| 方法 | 路径 | 说明 |
|------|------|------|
| GET/POST | /api/tasks | 任务列表/新增 |
| POST | /api/tasks/:id/execute | 手动扫描 |
| POST | /api/rules/test | 测试jQuery脚本 |
| GET | /api/tasks/:id/records | 扫描历史 |
| GET/PUT | /api/config/smtp | 邮件配置 |
