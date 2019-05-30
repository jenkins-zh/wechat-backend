pipeline {
    agent {
        label "golang"
    }

    environment {
        IMAGE_TAG = ""
        FOLDER = 'src/github.com/jenkins-zh/wechat-backend'
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
                dir(FOLDER) {
                    container('golang'){
                            sh '''
                            CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -s" -a -installsuffix cgo -o bin/wechat-backend
                            '''
                    }
                    container('tools') {
                        sh 'upx bin/wechat-backend'
                    }
                }
            }
        }

        stage("image") {
            steps {
                container('tools'){
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
                container('tools') {
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