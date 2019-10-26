CGO_ENABLED = 0

TAG=dev-$(shell cat .version)-$(shell git config --get user.email | sed -e "s/@/-/")

build:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -s" -a -installsuffix cgo -o bin/wechat-backend
	upx bin/wechat-backend

run:
	CGO_ENABLED=0 go build -ldflags "-w -s" -a -installsuffix cgo -o bin/wechat-backend
	./bin/wechat-backend

build-local:
	env go build -o bin/wechat-backend
	upx bin/wechat-backend

image: build
	docker build -t surenpi/jenkins-wechat:${TAG} .

image-alauda: build
	docker build -t index.alauda.cn/alaudak8s/jenkins-wechat .

push-image: image
	docker push surenpi/jenkins-wechat:${TAG}

push-image-alauda: image-alauda
	docker push index.alauda.cn/alaudak8s/jenkins-wechat

image-ubuntu: build
	docker build -t surenpi/jenkins-wechat:ubuntu . -f Dockerfile.ubuntu
	docker push surenpi/jenkins-wechat:ubuntu

init-mock-dep:
	go get github.com/golang/mock/gomock
	go install github.com/golang/mock/mockgen

update:
	kubectl set image deploy/wechat wechat=surenpi/jenkins-wechat:${TAG}
	make restart

update-alauda:
	kubectl set image deploy/wechat wechat=index.alauda.cn/alaudak8s/jenkins-wechat
	make restart

restart:
	kubectl scale deploy/wechat --replicas=0
	kubectl scale deploy/wechat --replicas=1

test:
	go test ./... -v -coverprofile coverage.out
