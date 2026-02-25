#!/bin/bash

# Configuration
MAX_RETRIES=50
RETRY_DELAY=5

echo "Starting robust Git Push process..."
echo "This script will automatically retry if network connection fails during push."

# Function to push with retries
git_push_retry() {
    local attempt=1
    while [ $attempt -le $MAX_RETRIES ]; do
        echo "--------------------------------------------------"
        echo "Attempt $attempt of $MAX_RETRIES: Pushing to GitHub..."
        
        # Try to push
        # Using http.postBuffer to handle large commits if necessary
        git config http.postBuffer 524288000
        
        if git push; then
            echo "✅ Git push successful!"
            return 0
        else
            echo "❌ Push failed (likely network issue or timeout)."
            echo "Retrying in $RETRY_DELAY seconds..."
            sleep $RETRY_DELAY
            attempt=$((attempt+1))
        fi
    done
    
    echo "❌ Failed to push after $MAX_RETRIES attempts."
    return 1
}

# Execute
git_push_retry
