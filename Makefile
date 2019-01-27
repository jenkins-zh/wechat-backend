build:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -s" -a -installsuffix cgo -o bin/wechat-backend
	upx bin/wechat-backend

build-local:
	env go build -o bin/wechat-backend
	upx bin/wechat-backend

image: build
	docker build -t surenpi/jenkins-wechat .

image-alauda: build
	docker build -t index.alauda.cn/alaudak8s/jenkins-wechat .

push-image: image
	docker push surenpi/jenkins-wechat

push-image-alauda: image-alauda
	docker push index.alauda.cn/alaudak8s/jenkins-wechat

image-ubuntu: build
	docker build -t surenpi/jenkins-wechat:ubuntu . -f Dockerfile.ubuntu
	docker push surenpi/jenkins-wechat:ubuntu

init-mock-dep:
	go get github.com/golang/mock/gomock
	go install github.com/golang/mock/mockgen

update:
	kubectl set image deploy/wechat wechat=surenpi/jenkins-wechat
	make restart

update-alauda:
	kubectl set image deploy/wechat wechat=index.alauda.cn/alaudak8s/jenkins-wechat
	make restart

restart:
	kubectl scale deploy/wechat --replicas=0
	kubectl scale deploy/wechat --replicas=1