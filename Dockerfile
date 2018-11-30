FROM alpine:3.8

USER root

RUN apk -U upgrade && \
    apk -U add ca-certificates && \
    update-ca-certificates

COPY bin/wechat-backend wechat-backend

RUN chmod u+x wechat-backend

CMD ["./wechat-backend"]
