# 构建阶段
FROM golang:1.24.5-alpine AS builder

# 接收服务参数
ARG SERVICE_TYPE
ARG SERVICE_NAME

WORKDIR /app

# 安装依赖工具
RUN apk add --no-cache git curl

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制项目代码
COPY . .

# 构建服务
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -installsuffix cgo \
    -o ${SERVICE_NAME} \
    ./apps/${SERVICE_TYPE}/${SERVICE_NAME}/${SERVICE_NAME}${SERVICE_TYPE}.go

# 运行阶段
FROM alpine:3.18

# 安装必要工具
RUN apk --no-cache add ca-certificates tzdata curl

ENV TZ=Asia/Shanghai

ARG SERVICE_NAME
ARG SERVICE_TYPE

WORKDIR /app

# 复制二进制文件
COPY --from=builder /app/${SERVICE_NAME} .

# 复制配置文件
COPY --from=builder /app/apps/${SERVICE_TYPE}/${SERVICE_NAME}/etc /app/etc

# 暴露端口
EXPOSE 8080 8081 8082

CMD ["./service"]
