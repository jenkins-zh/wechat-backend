pipeline {
    agent {
        label "golang"
    }

    stages{
        stage("build") {
            steps {
                sh 'make build'
            }
        }

        stage("image") {
            steps {
                sh 'make image'
            }
        }

        stage("push-image") {
            steps {
                withCredentials([usernamePassword(credentialsId: '', passwordVariable: 'passwd', usernameVariable: 'user')]) {
                    sh '''
                    docker login -u $user -p $passwd
                    make push-image
                    docker logout
                    '''
                }
            }
        }
    }
}