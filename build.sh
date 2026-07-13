#!/bin/bash
set -e

echo "=== 1. 构建前端 ==="
cd web
npm install --silent
npm run build
cd ..

echo "=== 2. 拷贝前端产物 ==="
cp -r web/dist server/web/dist

echo "=== 3. 构建后端 ==="
cd server
go mod tidy
go build -o monitor-server .
cd ..

echo "=== 4. 构建完成 ==="
echo "后端二进制: server/monitor-server"
echo "启动命令: cd server && ./monitor-server"
echo ""
echo "=== Docker 构建（可选）==="
echo "docker build -t monitor-server:latest -f server/Dockerfile ."
