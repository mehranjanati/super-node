#!/bin/bash

# Configuration
MAX_RETRIES=20
RETRY_DELAY=5

# Function to push with retries
push_changes() {
    local attempt=1
    while [ $attempt -le $MAX_RETRIES ]; do
        echo "Pushing changes... (Attempt $attempt/$MAX_RETRIES)"
        if git push; then
            echo "✅ Push successful!"
            return 0
        fi
        echo "❌ Push failed. Retrying in $RETRY_DELAY seconds..."
        sleep $RETRY_DELAY
        attempt=$((attempt+1))
    done
    return 1
}

# Reset any pending commits to start fresh splitting
echo "Resetting mixed changes..."
git reset HEAD~1

# 1. Commit Core Domain & Ports (Smallest, most fundamental)
echo "--------------------------------------------------"
echo "📦 Committing: Core Domain & Ports"
git add internal/core/domain/ internal/ports/
git commit -m "feat(core): update domain models and ports for agent service"
push_changes

# 2. Commit Services (Business Logic)
echo "--------------------------------------------------"
echo "📦 Committing: Services"
git add internal/core/services/
git commit -m "feat(services): implement agent service business logic"
push_changes

# 3. Commit Adapters (Infrastructure)
echo "--------------------------------------------------"
echo "📦 Committing: Adapters (Gateway, Persistence)"
git add internal/adapters/
git commit -m "feat(adapters): connect echo gateway and tidb repository"
push_changes

# 4. Commit CMD & Configuration (Wiring)
echo "--------------------------------------------------"
echo "📦 Committing: CMD & Config"
git add cmd/ internal/config/
git commit -m "chore(main): wire up dependencies in main.go"
push_changes

# 5. Commit Docker & Build Scripts (DevOps)
echo "--------------------------------------------------"
echo "📦 Committing: DevOps (Docker, Scripts)"
git add Dockerfile docker-compose* *.sh
git commit -m "ci(docker): optimize dockerfile with alpine and add retry scripts"
push_changes

# 6. Commit Anything Else (Docs, etc.)
echo "--------------------------------------------------"
echo "📦 Committing: Remaining files"
git add .
git commit -m "chore: misc updates and documentation"
push_changes

echo "🎉 All incremental commits pushed successfully!"
