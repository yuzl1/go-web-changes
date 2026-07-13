# 阶段1: 前端构建
FROM node:22-alpine AS web-builder
WORKDIR /web
COPY web/package.json web/package-lock.json ./
RUN npm ci 2>/dev/null || npm install
COPY web/ ./
RUN npm run build

# 阶段2: 后端构建
FROM golang:1.22-alpine AS server-builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /server
COPY server/go.mod server/go.sum ./
RUN go mod download
COPY server/ ./
RUN CGO_ENABLED=1 go build -o monitor-server .

# 阶段3: 运行镜像
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    chromium ca-certificates fonts-liberation fonts-noto-cjk \
    && rm -rf /var/lib/apt/lists/*

ENV CHROME_BIN=/usr/bin/chromium
ENV CHROME_PATH=/usr/bin/chromium

WORKDIR /app
COPY --from=server-builder /server/monitor-server .
COPY --from=server-builder /server/config/config.yaml ./config/
COPY --from=web-builder /web/dist ./web/dist/

RUN mkdir -p /app/data /app/logs

EXPOSE 8080
CMD ["./monitor-server"]
