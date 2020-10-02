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
	docker build -t jenkinszh/jenkins-wechat:${TAG} .

push-image: image
	docker push jenkinszh/jenkins-wechat:${TAG}

image-ubuntu: build
	docker build -t jenkinszh/jenkins-wechat:ubuntu . -f Dockerfile.ubuntu
	docker push jenkinszh/jenkins-wechat:ubuntu

init-mock-dep:
	go get github.com/golang/mock/gomock
	go install github.com/golang/mock/mockgen

update:
	kubectl set image deploy/wechat wechat=jenkinszh/jenkins-wechat:${TAG}
	make restart

restart:
	kubectl scale deploy/wechat --replicas=0
	kubectl scale deploy/wechat --replicas=1

test:
	go test ./... -v -coverprofile coverage.out
