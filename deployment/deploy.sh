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

# Navigate to project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
cd "$PROJECT_ROOT"

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}🚀 Event Campus Deployment Script${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

# Pre-deployment checks
print_info "Running pre-deployment checks..."

# Check if .env exists
if [ ! -f .env ]; then
    print_error ".env file not found!"
    print_info "Please create .env file from .env.example"
    exit 1
fi
print_success ".env file found"

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    print_error "Docker is not running!"
    exit 1
fi
print_success "Docker is running"

# Check disk space (minimum 1GB)
AVAILABLE_KB=$(df . | tail -1 | awk '{print $4}')
AVAILABLE_GB=$((AVAILABLE_KB / 1024 / 1024))
if [ $AVAILABLE_KB -lt 1048576 ]; then
    print_warning "Low disk space! Available: ${AVAILABLE_GB}GB"
    read -p "Continue anyway? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_error "Deployment cancelled"
        exit 1
    fi
else
    print_success "Disk space OK (${AVAILABLE_GB}GB available)"
fi

# Pull latest changes
print_info "Pulling latest code from repository..."
if git pull origin main; then
    print_success "Code updated"
else
    print_warning "Git pull failed or no changes"
fi

# Build and start containers
print_info "Building and starting containers..."
docker-compose -f deployment/docker-compose.yml down 2>/dev/null || true
docker-compose -f deployment/docker-compose.yml up -d --build

# Wait for health check
print_info "Waiting for application to be healthy..."
sleep 5

MAX_RETRIES=12
RETRY_COUNT=0
while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    if curl -f http://localhost:8080/health > /dev/null 2>&1; then
        print_success "Application is healthy!"
        break
    fi
    RETRY_COUNT=$((RETRY_COUNT + 1))
    if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
        print_error "Health check failed after ${MAX_RETRIES} attempts"
        print_info "Check logs with: docker logs event-campus-api"
        exit 1
    fi
    echo -n "."
    sleep 5
done

# Cleanup old images
print_info "Cleaning up old images..."
docker image prune -f > /dev/null 2>&1
print_success "Cleanup complete"

# Show container status
echo ""
print_info "Container status:"
docker ps | grep event-campus-api || print_error "Container not found!"

echo ""
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}✅ DEPLOYMENT SUCCESSFUL!${NC}"
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}API Endpoint:${NC} http://localhost:8080"
echo -e "${BLUE}Health Check:${NC} http://localhost:8080/health"
echo -e "${BLUE}View Logs:${NC} docker logs -f event-campus-api"
echo ""

