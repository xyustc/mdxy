# Go 后端 Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /build

# 安装依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制源码
COPY . .

# 构建二进制文件
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# 最终镜像
FROM alpine:latest

RUN apk --no-cache add ca-certificates sqlite-libs

WORKDIR /app

# 复制二进制文件
COPY --from=builder /build/server .
COPY --from=builder /build/config.yaml .

# 创建数据目录
RUN mkdir -p /app/data

EXPOSE 8080

CMD ["./server"]
