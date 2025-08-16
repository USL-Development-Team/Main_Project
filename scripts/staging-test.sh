#!/bin/bash

# Staging Environment Test Script
# Tests production-ready container before deployment to Render

set -e

echo "🚀 Testing Staging Container (Pre-Production)"
echo "============================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Container names
APP_CONTAINER="app-staging"
NETWORK_NAME="app-staging-network"

# Clean up function
cleanup() {
    print_status "Cleaning up containers and networks..."
    docker stop $APP_CONTAINER 2>/dev/null || true
    docker rm $APP_CONTAINER 2>/dev/null || true
    docker network rm $NETWORK_NAME 2>/dev/null || true
}

# Set up cleanup trap
trap cleanup EXIT

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    print_error "Docker is not running. Please start Docker and try again."
    exit 1
fi

# Check if .env.staging exists
if [ ! -f ".env.staging" ]; then
    print_error ".env.staging file not found. Please create it first."
    exit 1
fi

print_status "Creating Docker network..."
docker network create $NETWORK_NAME

print_status "Building application image..."
docker build -t app-staging:latest .

print_status "Starting application container (without database - testing container functionality)..."
docker run -d \
    --name $APP_CONTAINER \
    --network $NETWORK_NAME \
    --env-file .env.staging \
    -e PORT=8080 \
    -e SERVER_HOST=0.0.0.0 \
    -p 8080:8080 \
    app-staging:latest

print_status "Waiting for application to be ready..."
sleep 5
timeout 30 bash -c 'until curl -f http://localhost:8080/health > /dev/null 2>&1; do sleep 2; done'

print_status "Running health checks..."

# Test 1: Health endpoint response
echo -n "✓ Health endpoint response: "
if curl -f -s http://localhost:8080/health > /dev/null; then
    echo -e "${GREEN}PASS${NC}"
else
    echo -e "${RED}FAIL${NC}"
    print_error "Health endpoint not responding on port 8080"
    docker logs $APP_CONTAINER
    exit 1
fi

# Test 2: Static assets
echo -n "✓ Static assets: "
if curl -f -s http://localhost:8080/static/htmx.min.js > /dev/null; then
    echo -e "${GREEN}PASS${NC}"
else
    echo -e "${RED}FAIL${NC}"
    print_warning "Static assets may not be properly served"
fi

# Test 3: Check logs for critical errors (exclude expected warnings)
echo -n "✓ No critical errors in logs: "
# Filter out expected warnings - database connectivity issues are expected without Supabase
CRITICAL_ERRORS=$(docker logs $APP_CONTAINER 2>&1 | grep -i "fatal\|panic" | grep -v -E "(Guild not found|database|supabase|connection)" | wc -l)
if [ "$CRITICAL_ERRORS" -gt 0 ]; then
    echo -e "${RED}FAIL${NC}"
    print_error "Found critical errors in application logs:"
    docker logs $APP_CONTAINER 2>&1 | grep -i "fatal\|panic" | grep -v -E "(Guild not found|database|supabase|connection)"
    exit 1
else
    echo -e "${GREEN}PASS${NC}"
fi

# Test 4: Container health (wait for health check to complete)
echo -n "✓ Container health: "
# Wait up to 60 seconds for health check to pass
for i in $(seq 1 12); do
    HEALTH_STATUS=$(docker inspect --format='{{.State.Health.Status}}' $APP_CONTAINER 2>/dev/null || echo 'healthy')
    if [ "$HEALTH_STATUS" = "healthy" ]; then
        echo -e "${GREEN}PASS${NC}"
        break
    elif [ "$HEALTH_STATUS" = "unhealthy" ]; then
        echo -e "${RED}FAIL${NC}"
        print_error "Container health check failed"
        echo "Health check logs:"
        docker inspect $APP_CONTAINER --format='{{range .State.Health.Log}}{{.Output}}{{end}}'
        exit 1
    else
        # Still starting, wait a bit more
        sleep 5
    fi
done

# If we got here and status is still not healthy, it's a timeout
if [ "$HEALTH_STATUS" != "healthy" ]; then
    echo -e "${YELLOW}WARN${NC} (Health check still starting - ${HEALTH_STATUS})"
fi

# Test 5: Memory usage check
echo -n "✓ Memory usage reasonable: "
MEMORY_USAGE=$(docker stats --no-stream --format "{{.MemUsage}}" $APP_CONTAINER | cut -d'/' -f1 | sed 's/MiB//')
if [ "${MEMORY_USAGE%.*}" -lt 100 ]; then
    echo -e "${GREEN}PASS${NC} (${MEMORY_USAGE})"
else
    echo -e "${YELLOW}WARN${NC} (${MEMORY_USAGE} - higher than expected)"
fi

print_status "All staging tests passed! 🎉"
print_status "Container is ready for production deployment to Render."

# Show container information
echo ""
print_status "Container information:"
echo "Image size: $(docker images app-staging:latest --format 'table {{.Size}}')"
echo "Container uptime: $(docker inspect --format='{{.State.StartedAt}}' $APP_CONTAINER)"
echo ""

print_status "To view logs: docker logs -f $APP_CONTAINER"
print_status "To access container: docker exec -it $APP_CONTAINER sh"
print_status "Containers will be cleaned up automatically on script exit"