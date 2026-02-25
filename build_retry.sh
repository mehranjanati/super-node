#!/bin/bash

# Configuration
IMAGE_NAME="nexus-super-node:latest"
MAX_RETRIES=100
RETRY_DELAY=5

echo "Starting robust build process for Super Node Backend..."
echo "This script will use cache to avoid starting from scratch after network failures."

# Function to build docker image with retries and cache
build_image() {
    local attempt=1
    while [ $attempt -le $MAX_RETRIES ]; do
        echo "--------------------------------------------------"
        echo "Attempt $attempt of $MAX_RETRIES: Building Docker image..."
        
        # Using --build-arg if needed, but primarily relying on Docker's layer caching
        # The key is that Docker naturally caches successful layers (like go mod download)
        if docker build --cache-from $IMAGE_NAME -t $IMAGE_NAME .; then
            echo "✅ Docker build successful!"
            return 0
        else
            echo "❌ Build failed (likely network issue)."
            echo "Docker cached previous successful layers. Retrying in $RETRY_DELAY seconds..."
            sleep $RETRY_DELAY
            attempt=$((attempt+1))
        fi
    done
    
    echo "❌ Failed to build image after $MAX_RETRIES attempts."
    return 1
}

# Function to start services
start_services() {
    echo "--------------------------------------------------"
    echo "Starting all services (Infra + App)..."
    
    # Try to bring up services
    if docker-compose -f docker-compose.infra.yml -f docker-compose.app.yml up -d; then
        echo "✅ All backend services are up and running!"
        return 0
    else
        echo "❌ Failed to start services. Check docker-compose files."
        return 1
    fi
}

# Main execution
if build_image; then
    start_services
else
    echo "Exiting due to build failure."
    exit 1
fi
