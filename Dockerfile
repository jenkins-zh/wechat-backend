FROM alpine:3.3

USER root

RUN sed -i 's|dl-cdn.alpinelinux.org|mirrors.aliyun.com|g' /etc/apk/repositories
RUN apk add --no-cache ca-certificates curl

COPY bin/wechat-backend wechat-backend

RUN chmod u+x wechat-backend

CMD ["./wechat-backend"]
