pipeline {
    agent any
    
    environment {
        APP_NAME = 'event-campus-api'
        DOCKER_IMAGE = "${APP_NAME}:${BUILD_NUMBER}"
        DEPLOY_PATH = '/opt/event-campus'
    }
    
    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }
        
        stage('Test') {
            steps {
                sh 'go test -v ./...'
            }
        }
        
        stage('Build Docker Image') {
            steps {
                sh "docker build -t ${DOCKER_IMAGE} ."
                sh "docker tag ${DOCKER_IMAGE} ${APP_NAME}:latest"
            }
        }
        
        stage('Deploy to VPS') {
            steps {
                script {
                    // Stop existing container
                    sh """
                        docker stop ${APP_NAME} || true
                        docker rm ${APP_NAME} || true
                    """
                    
                    // Run new container
                    sh """
                        docker run -d \
                            --name ${APP_NAME} \
                            -p 8080:8080 \
                            --env-file ${DEPLOY_PATH}/.env \
                            -v ${DEPLOY_PATH}/storage:/root/storage \
                            --restart unless-stopped \
                            ${APP_NAME}:latest
                    """
                }
            }
        }
        
        stage('Health Check') {
            steps {
                script {
                    sleep 5
                    sh 'curl -f http://localhost:8080/health || exit 1'
                }
            }
        }
        
        stage('Cleanup') {
            steps {
                sh "docker image prune -f"
            }
        }
    }
    
    post {
        success {
            echo 'Deployment successful!'
        }
        failure {
            echo 'Deployment failed!'
            // Rollback to previous version
            sh "docker stop ${APP_NAME} || true"
            sh "docker rm ${APP_NAME} || true"
        }
    }
}
