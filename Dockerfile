FROM golang:1.13
WORKDIR /workspace
COPY . .
ENV GOPROXY=https://mirrors.aliyun.com/goproxy/
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -s" -a -installsuffix cgo -o wechat-backend
RUN chmod u+x wechat-backend

FROM alpine:3.3
USER root
RUN sed -i 's|dl-cdn.alpinelinux.org|mirrors.aliyun.com|g' /etc/apk/repositories
RUN apk add --no-cache ca-certificates curl
COPY --from=0 /workspace/wechat-backend .
CMD ["./wechat-backend"]
