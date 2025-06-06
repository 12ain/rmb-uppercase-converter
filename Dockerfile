FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制go.mod和go.sum文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制项目文件
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 使用轻量级镜像运行应用
FROM alpine:latest

# 安装必要的证书
RUN apk --no-cache add ca-certificates

# 设置工作目录
WORKDIR /root/

# 从builder阶段复制构建好的应用
COPY --from=builder /app/main .
COPY --from=builder /app/static ./static

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./main"]