#!/bin/bash

# Build all Docker images for the test-workflow project

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo "Building all Docker images..."

# Build API Gateway
echo "Building API Gateway (Go)..."
docker build -t test-workflow/api-gateway:latest "$PROJECT_ROOT/services/api-gateway"

# Build Order Service
echo "Building Order Service (Java)..."
docker build -t test-workflow/order-service:latest "$PROJECT_ROOT/services/order-service"

# Build Inventory Service
echo "Building Inventory Service (C#)..."
docker build -t test-workflow/inventory-service:latest "$PROJECT_ROOT/services/inventory-service"

# Build Notification Service
echo "Building Notification Service (Ruby)..."
docker build -t test-workflow/notification-service:latest "$PROJECT_ROOT/services/notification-service"

# Build Frontend
echo "Building Frontend (React)..."
docker build -t test-workflow/frontend:latest "$PROJECT_ROOT/frontend"

echo ""
echo "All images built successfully!"
echo ""
echo "Images:"
docker images | grep test-workflow
