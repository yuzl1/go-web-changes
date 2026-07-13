# ============ 阶段1: 前端构建 ============
FROM node:20-alpine AS web-builder
WORKDIR /web
COPY web/package.json web/package-lock.json ./
RUN npm ci 2>/dev/null || npm install
COPY web/ ./
RUN npm run build
# 输出: /web/dist/

# ============ 阶段2: 后端构建 ============
FROM golang:1.21-alpine AS server-builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /server
COPY server/go.mod server/go.sum ./
RUN go mod download
COPY server/ ./
# 将前端构建产物拷贝到后端目录（供 embed 使用）
COPY --from=web-builder /web/dist ./web/dist/
RUN CGO_ENABLED=1 go build -o monitor-server .

# ============ 阶段3: 运行镜像 ============
FROM debian:bookworm-slim

# 安装 Chromium 及运行依赖
RUN apt-get update && apt-get install -y --no-install-recommends \
    chromium \
    ca-certificates \
    fonts-liberation \
    fonts-noto-cjk \
    && rm -rf /var/lib/apt/lists/*

# 设置 Chromium 环境变量
ENV CHROME_BIN=/usr/bin/chromium
ENV CHROME_PATH=/usr/bin/chromium

# 创建应用目录
WORKDIR /app
COPY --from=server-builder /server/monitor-server .
COPY --from=server-builder /server/config/config.yaml ./config/

# 创建数据和日志目录
RUN mkdir -p /app/data /app/logs

# 暴露端口
EXPOSE 8080

# 启动
CMD ["./monitor-server"]
