# 构建阶段
FROM golang:1.24-alpine AS builder

WORKDIR /app

# 安装必要的构建工具
RUN apk add --no-cache git ca-certificates tzdata

# 复制 go.mod 和 go.sum
COPY go.mod go.sum* ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建二进制文件
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /tidalcore ./cmd/server/

# 运行阶段
FROM alpine:3.19

WORKDIR /app

# 安装时区数据和 CA 证书
RUN apk add --no-cache ca-certificates tzdata

# 从构建阶段复制二进制文件
COPY --from=builder /tidalcore .

# 创建非 root 用户和备份目录
RUN adduser -D -g '' appuser && \
    mkdir -p /app/backups && \
    chown -R appuser:appuser /app
USER appuser

# 暴露端口
EXPOSE 8080

# 启动命令 (配置通过环境变量传入)
CMD ["./tidalcore"]
