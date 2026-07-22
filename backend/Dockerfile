# backend/Dockerfile
# 多阶段构建：编译 Go 后端服务

# 阶段 1：编译
FROM golang:1.25-alpine AS builder

WORKDIR /app

# 先复制依赖文件并下载，利用 Docker 缓存层
COPY go.mod go.sum ./
RUN go mod download

# 复制源码
COPY . .

# 编译（静态链接，无 CGO）
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./cmd/server

# 阶段 2：运行（最小镜像）
FROM alpine:latest

WORKDIR /app

# 复制编译产物与配置
COPY --from=builder /app/server .
COPY config/ ./config/

# 暴露服务端口（与 config/config.yaml 中的 server.port 保持一致）
EXPOSE 8080

CMD ["./server"]
