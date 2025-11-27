#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Functions
print_success() { echo -e "${GREEN}✅ $1${NC}"; }
print_error() { echo -e "${RED}❌ $1${NC}"; }
print_warning() { echo -e "${YELLOW}⚠️  $1${NC}"; }
print_info() { echo -e "${BLUE}ℹ️  $1${NC}"; }

APP_NAME="event-campus-api"
PREVIOUS_IMAGE="${APP_NAME}:previous"

echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${YELLOW}🔄 Event Campus Rollback Script${NC}"
echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

# Check if previous image exists
if ! docker images ${PREVIOUS_IMAGE} -q | grep -q .; then
    print_error "No previous version found!"
    print_info "Previous image tag '${PREVIOUS_IMAGE}' does not exist"
    print_info "Available images:"
    docker images ${APP_NAME}
    exit 1
fi

print_warning "This will rollback to the previous version"
read -p "Continue? (y/N) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    print_info "Rollback cancelled"
    exit 0
fi

# Get current container info for logging
print_info "Current container status:"
docker ps -a | grep ${APP_NAME} || echo "No container running"

# Stop and remove current container
print_info "Stopping current container..."
docker stop ${APP_NAME} 2>/dev/null || true
docker rm ${APP_NAME} 2>/dev/null || true
print_success "Container stopped"

# Start previous version
print_info "Starting previous version..."
DEPLOY_PATH="/opt/event-campus"

docker run -d \
    --name ${APP_NAME} \
    -p 3000:8080 \
    --env-file ${DEPLOY_PATH}/.env \
    -v ${DEPLOY_PATH}/storage:/app/storage \
    --restart unless-stopped \
    ${PREVIOUS_IMAGE}

# Wait and health check
print_info "Waiting for application to start..."
sleep 5

MAX_RETRIES=6
RETRY_COUNT=0
while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    if curl -f http://localhost:3000/health > /dev/null 2>&1; then
        print_success "Application is healthy!"
        break
    fi
    RETRY_COUNT=$((RETRY_COUNT + 1))
    if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
        print_error "Health check failed!"
        print_info "Check logs with: docker logs ${APP_NAME}"
        exit 1
    fi
    echo -n "."
    sleep 5
done

echo ""
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}✅ ROLLBACK SUCCESSFUL!${NC}"
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}API Endpoint:${NC} http://localhost:3000"
echo -e "${BLUE}Health Check:${NC} http://localhost:3000/health"
echo -e "${BLUE}View Logs:${NC} docker logs -f ${APP_NAME}"
echo ""
