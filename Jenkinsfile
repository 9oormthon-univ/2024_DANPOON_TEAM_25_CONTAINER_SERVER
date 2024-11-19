pipeline {
    agent any

    environment {
        DOCKER_IMAGE = "milkymilky0116/flake-ide-container-server"
        DOCKER_TAG = "latest"
    }

    stages {
        stage('Clone Repository') {
            steps {
                // Git 리포지토리를 클론
                git url: 'https://github.com/9oormthon-univ/2024_DANPOON_TEAM_25_CONTAINER_SERVER.git', branch: 'dev'
            }
        }

        stage('Build Docker Image') {
            steps {
                // Docker 이미지를 빌드
                script {
                    docker.build("${DOCKER_IMAGE}:${DOCKER_TAG}")
                }
            }
        } 

        stage('Push Docker Image') {
            steps {
                script {
                    docker.withRegistry('https://registry.hub.docker.com', 'dockerhub') {
                        docker.image("${DOCKER_IMAGE}:${DOCKER_TAG}").push()
                    }
                }
            }
        }
    }

    post {
        success {
            echo "Docker image pushed to Docker Hub successfully."
        }
        failure {
            echo "Failed to build or push Docker image."
        }
    }
}
