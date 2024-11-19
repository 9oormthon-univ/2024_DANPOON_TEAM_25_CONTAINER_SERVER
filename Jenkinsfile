pipeline {
    agent any
    
    environment {
        DOCKER_IMAGE = "milkymilky0116/flake-ide-container-server"
        DOCKER_TAG = "latest"
        // Docker Hub 크리덴셜 ID를 설정합니다
        DOCKER_CREDENTIALS_ID = 'jenkins-credentials'
    }
    
    stages {
        stage('Clone Repository') {
            steps {
                git branch: 'dev',
                    url: 'https://github.com/9oormthon-univ/2024_DANPOON_TEAM_25_CONTAINER_SERVER.git'
            }
        }
        
        stage('Build Docker Image') {
            steps {
                script {
                    // Jenkins 사용자에게 sudo 권한을 부여하여 Docker 커맨드 실행
                    sh '''
                        sudo docker build -t $DOCKER_IMAGE:$DOCKER_TAG .
                    '''
                }
            }
        }
        
        stage('Push Docker Image') {
            steps {
                script {
                    // Docker Hub 로그인을 위한 크리덴셜 사용
                    withCredentials([usernamePassword(credentialsId: DOCKER_CREDENTIALS_ID, passwordVariable: 'DOCKER_PASSWORD', usernameVariable: 'DOCKER_USERNAME')]) {
                        sh '''
                            echo $DOCKER_PASSWORD | sudo docker login -u $DOCKER_USERNAME --password-stdin
                            sudo docker push $DOCKER_IMAGE:$DOCKER_TAG
                            sudo docker rmi $DOCKER_IMAGE:$DOCKER_TAG
                        '''
                    }
                }
            }
        }
    }
    
    post {
        success {
            echo 'Successfully built and pushed Docker image!'
        }
        failure {
            echo 'Failed to build or push Docker image.'
        }
        always {
            // 클린업 작업
            sh 'sudo docker logout'
        }
    }
}
