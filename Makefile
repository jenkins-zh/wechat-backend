build:
	env GOOS=linux go build

push:
	scp wechat-backend root@surenpi.com:/root