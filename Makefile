build:
	env GOOS=linux go build -o bin/wechat-backend
	upx bin/wechat-backend

image: build
	docker build -t surenpi/jenkins-wechat .
	docker push surenpi/jenkins-wechat