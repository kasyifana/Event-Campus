#!/bin/bash
set -e

echo "🚀 Starting deployment..."

# Navigate to project root (assuming script is run from deployment/ or root)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
cd "$PROJECT_ROOT"

# Pull latest changes
echo "📥 Pulling latest code..."
git pull origin main

# Build and start containers
echo "🐳 Building and starting containers..."
docker-compose -f deployment/docker-compose.yml up -d --build

# Prune unused images
echo "🧹 Cleaning up..."
docker image prune -f

echo "✅ Deployment successful!"
