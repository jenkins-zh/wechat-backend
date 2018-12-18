build:
	env GOOS=linux go build -o bin/wechat-backend
	upx bin/wechat-backend

image: build
	docker build -t surenpi/jenkins-wechat .

push-image:
	docker push surenpi/jenkins-wechat

image-ubuntu: build
	docker build -t surenpi/jenkins-wechat:ubuntu . -f Dockerfile.ubuntu
	docker push surenpi/jenkins-wechat:ubuntu

init-mock-dep:
	go get github.com/golang/mock/gomock
	go install github.com/golang/mock/mockgen

update:
	kubectl set image deploy/wechat wechat=surenpi/jenkins-wechat
	make restart

restart:
	kubectl scale deploy/wechat --replicas=0
	kubectl scale deploy/wechat --replicas=1