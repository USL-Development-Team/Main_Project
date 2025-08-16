#!/bin/bash

# Environment management script
# Usage: ./scripts/env.sh [develop|production|staging]

ENV=${1:-develop}

case $ENV in
  "develop")
    echo "Loading development environment..."
    cp .env.develop .env
    ;;
  "production")
    echo "Loading production environment..."
    if [ ! -f .env.production ]; then
      echo "Error: .env.production not found!"
      echo "Create it first and fill in your production values."
      exit 1
    fi
    cp .env.production .env
    ;;
  "staging")
    echo "Loading staging environment..."
    if [ ! -f .env.staging ]; then
      echo "Error: .env.staging not found!"
      exit 1
    fi
    cp .env.staging .env
    ;;
  *)
    echo "Usage: $0 [develop|production|staging]"
    echo "Current environment files:"
    ls -la .env* 2>/dev/null || echo "No environment files found"
    exit 1
    ;;
esac

echo "Environment set to: $ENV"
echo "Current .env file:"
echo "=================="
head -5 .env