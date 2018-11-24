FROM alpine:3.8

COPY bin/wechat-backend wechat-backend

RUN chmod u+x wechat-backend

CMD ["./wechat-backend"]
