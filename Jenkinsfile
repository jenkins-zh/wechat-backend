pipeline {
    agent {
        label "golang"
    }

    environment {
        IMAGE_TAG = ""
        FOLDER = 'src/github.io/jenkins-zh'
    }

    stages{
        stage("clone") {
            steps {
                dir(FOLDER) {
                    checkout scm
                }
            }
        }

        stage("build") {
            environment {
                GOPATH = "${WORKSPACE}"
            }
            steps {
                container('golang'){
                    dir(FOLDER) {
                        sh '''
                        CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -s" -a -installsuffix cgo -o bin/wechat-backend
                        upx bin/wechat-backend
                        '''
                    }
                }
            }
        }

        stage("image") {
            steps {
                container('golang'){
                    dir(FOLDER) {
                        sh '''
                        IMAGE_TAG=$(git rev-parse --short HEAD)
                        docker build -t surenpi/jenkins-wechat:$IMAGE_TAG .
                        docker push surenpi/jenkins-wechat:$IMAGE_TAG
                        '''
                    }
                }
            }
        }

        stage("push-image") {
            environment {
                DOCKER_CREDS = credentials('docker-surenpi')
            }
            steps {
                container('golang') {
                    sh '''
                    docker login -u $DOCKER_CREDS_USR -p $DOCKER_CREDS_PSW
                    make push-image
                    docker logout
                    '''
                }
            }
        }
    }
}