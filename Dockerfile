FROM golang:1.26-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制 go mod 文件并下载依赖（利用层缓存）
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 安装 git 并构建应用
RUN apk add --no-cache git && \
    CGO_ENABLED=0 GOOS=linux go build -tags netgo -ldflags "-s -w" -o ./cloudgram-go ./main.go

# 生产阶段
FROM alpine:latest

# 安装 ca-certificates 以支持 HTTPS 请求
RUN apk --no-cache add ca-certificates tzdata && \
    # 创建非 root 用户
    addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件和入口脚本
COPY --from=builder /app/cloudgram-go /app/cloudgram-go
COPY --from=builder /app/entrypoint.sh /app/entrypoint.sh
COPY --from=builder /app/dist /app/dist

# 设置环境变量（保持原有环境变量）
ENV LISTEN=
ENV AUTH_USER=
ENV AUTH_PASSWORD=
ENV JWT_SECRET_KEY=
ENV TELEGRAM_BOT_TOKEN=
ENV DB_TYPE=
ENV DB_DSN=
ENV LOG_PATH=
ENV DEBUG=

# 更改文件所有者并设置执行权限
RUN chown -R appuser:appgroup /app && \
    chmod +x /app/cloudgram-go && \
    chmod +x /app/entrypoint.sh

# 切换到非 root 用户
USER appuser

# 健康检查（可选，根据实际端口调整）
# HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
#   CMD wget --quiet --tries=1 --spider http://localhost:5244/health || exit 1

ENTRYPOINT ["/app/entrypoint.sh"]
CMD []