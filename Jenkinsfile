pipeline {
    agent any
    
    environment {
        APP_NAME = 'event-campus-api'
        BUILD_VERSION = sh(script: "git rev-parse --short HEAD", returnStdout: true).trim()
        DOCKER_IMAGE = "${APP_NAME}:${BUILD_VERSION}"
        DEPLOY_PATH = '/opt/event-campus'
        PREVIOUS_IMAGE = "${APP_NAME}:previous"
    }
    
    stages {
        stage('Checkout') {
            steps {
                echo '📥 Checking out code...'
                checkout scm
            }
        }
        
        stage('Environment Check') {
            steps {
                echo '🔍 Checking deployment environment...'
                script {
                    // Check if .env file exists
                    sh """
                        if [ ! -f ${DEPLOY_PATH}/.env ]; then
                            echo "❌ ERROR: ${DEPLOY_PATH}/.env not found!"
                            echo "Please create .env file using .env.example as template"
                            exit 1
                        fi
                    """
                    
                    // Check disk space (at least 1GB free)
                    sh """
                        available=\$(df ${DEPLOY_PATH} | tail -1 | awk '{print \$4}')
                        if [ \$available -lt 1048576 ]; then
                            echo "⚠️  WARNING: Low disk space! Available: \${available}KB"
                        fi
                    """
                }
            }
        }
        
        // Test stage skipped - no test files in project
        // Uncomment below if you add unit tests in the future
        /*
        stage('Test') {
            steps {
                echo '🧪 Running tests in Docker...'
                script {
                    timeout(time: 5, unit: 'MINUTES') {
                        sh """
                            docker run --rm \
                                -v \$(pwd):/app \
                                -w /app \
                                golang:1.24-alpine \
                                sh -c 'go mod download && go test -v ./... || exit 0'
                        """
                    }
                }
            }
        }
        */
        
        stage('Backup Current Version') {
            steps {
                echo '💾 Backing up current version...'
                script {
                    // Tag current running image as 'previous' for rollback
                    sh """
                        if docker images ${APP_NAME}:latest -q | grep -q .; then
                            docker tag ${APP_NAME}:latest ${PREVIOUS_IMAGE} || true
                            echo "✅ Backed up current version"
                        else
                            echo "ℹ️  No previous version to backup (first deployment)"
                        fi
                    """
                }
            }
        }
        
        stage('Build Docker Image') {
            steps {
                echo '🐳 Building Docker image...'
                sh """
                    docker build \
                        --build-arg BUILD_VERSION=${BUILD_VERSION} \
                        -t ${DOCKER_IMAGE} \
                        -t ${APP_NAME}:latest .
                """
                sh "echo '✅ Image built: ${DOCKER_IMAGE}'"
            }
        }
        
        stage('Deploy to VPS') {
            steps {
                echo '🚀 Deploying new version...'
                script {
                    // Stop and remove existing container
                    sh """
                        docker stop ${APP_NAME} 2>/dev/null || echo 'No container to stop'
                        docker rm ${APP_NAME} 2>/dev/null || echo 'No container to remove'
                    """
                    
                    // Run new container
                    sh """
                        docker run -d \
                            --name ${APP_NAME} \
                            -p 3000:8080 \
                            --env-file ${DEPLOY_PATH}/.env \
                            -v ${DEPLOY_PATH}/storage:/app/storage \
                            --restart unless-stopped \
                            --health-cmd="curl -f http://localhost:8080/health || exit 1" \
                            --health-interval=30s \
                            --health-timeout=5s \
                            --health-retries=3 \
                            ${APP_NAME}:latest
                    """
                }
            }
        }
        
        stage('Health Check') {
            steps {
                echo '🏥 Performing health check...'
                script {
                    // Wait for container to start
                    sleep 5
                    
                    // Retry health check up to 6 times (30 seconds total)
                    retry(6) {
                        sleep 5
                        sh """
                            curl -f http://localhost:8080/health || {
                                echo "❌ Health check failed!"
                                exit 1
                            }
                        """
                    }
                    echo '✅ Application is healthy!'
                }
            }
        }
        
        stage('Cleanup') {
            steps {
                echo '🧹 Cleaning up old images...'
                script {
                    // Keep only last 3 builds + latest + previous
                    sh """
                        # Remove dangling images
                        docker image prune -f
                        
                        # Keep only recent images
                        docker images ${APP_NAME} --format '{{.Tag}}' | \
                            grep -v 'latest' | \
                            grep -v 'previous' | \
                            tail -n +4 | \
                            xargs -r -I {} docker rmi ${APP_NAME}:{} 2>/dev/null || true
                    """
                    echo '✅ Cleanup complete'
                }
            }
        }
    }
    
    post {
        success {
            echo '✅ =========================================='
            echo '✅ DEPLOYMENT SUCCESSFUL!'
            echo '✅ =========================================='
            echo "Version: ${BUILD_VERSION}"
            echo "Container: ${APP_NAME}"
            echo "Health: http://localhost:3000/health"
        }
        
        failure {
            echo '❌ =========================================='
            echo '❌ DEPLOYMENT FAILED!'
            echo '❌ =========================================='
            echo '🔄 Attempting rollback to previous version...'
            
            script {
                sh """
                    # Stop failed container
                    docker stop ${APP_NAME} 2>/dev/null || true
                    docker rm ${APP_NAME} 2>/dev/null || true
                    
                    # Start previous version if exists
                    if docker images ${PREVIOUS_IMAGE} -q | grep -q .; then
                        docker run -d \
                            --name ${APP_NAME} \
                            -p 3000:8080 \
                            --env-file ${DEPLOY_PATH}/.env \
                            -v ${DEPLOY_PATH}/storage:/app/storage \
                            --restart unless-stopped \
                            ${PREVIOUS_IMAGE}
                        echo "✅ Rolled back to previous version"
                    else
                        echo "❌ No previous version available for rollback!"
                    fi
                """
            }
        }
        
        always {
            echo '📊 Deployment Summary:'
            sh 'docker ps -a | grep ${APP_NAME} || echo "No container found"'
        }
    }
}

