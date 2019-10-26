def scmObj
pipeline {
    agent {
        label "golang"
    }

    environment {
        FOLDER = 'src/github.com/jenkins-zh/wechat-backend'
    }

    stages{
        stage("clone") {
            steps {
                dir(FOLDER) {
                    script {
                        scmObj = checkout scm
                    }
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
            environment {
                IMAGE_TAG = getCurrentCommit(scmObj)
            }
            steps {
                container('tools'){
                    dir(FOLDER) {
                        sh '''
                        docker build -t surenpi/jenkins-wechat:$IMAGE_TAG .
                        docker build -t surenpi/jenkins-wechat:$IMAGE_TAG-ubuntu -f Dockerfile.ubuntu .
                        '''
                    }
                }
            }
        }

        stage("push-image") {
            environment {
                DOCKER_CREDS = credentials('docker-surenpi')
                IMAGE_TAG = getCurrentCommit(scmObj)
            }
            steps {
                container('tools') {
                    sh '''
                    docker login -u $DOCKER_CREDS_USR -p $DOCKER_CREDS_PSW
                    docker push surenpi/jenkins-wechat:$IMAGE_TAG
                    docker push surenpi/jenkins-wechat:$IMAGE_TAG-ubuntu
                    docker logout
                    '''
                }
            }
        }
    }
}

def getCurrentCommit(scmObj) {
    return scmObj.GIT_COMMIT
}
