# 启动编译环境  builder 阶段
FROM golang:1.17-alpine AS builder

# 配置编译环境
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct

# 拷贝源代码到镜像中
COPY .  /go/src/httpserver
WORKDIR /go/src/httpserver/

# 编译
RUN go build -o /bin/httpserver

## alpine 阶段
FROM alpine
COPY --from=builder /go/src/httpserver/. /bin/httpserver
# --from=<name> 将从 from 指定的构建阶段中寻找源文件；因为文件只存在于前一个 Docker 阶段
ENV VERSION = 1.0.0

# 设置服务入口
ENTRYPOINT ["/bin/httpserver"]
# EXPOSE 80
