#!/bin/bash

# Start only the Cyclops Controller (Backend)

set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${BLUE}[BACKEND]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[BACKEND]${NC} $1"
}

# Check if we're in the right directory
if [ ! -f "../cyclops-ctrl/Makefile" ]; then
    echo "Please run this script from the scripts directory in the cyclops project"
    exit 1
fi

print_status "Starting Cyclops Controller..."
cd ../cyclops-ctrl

# Check environment file
if [ ! -f ".env" ] && [ -f ".env.template" ]; then
    print_status "Copying .env.template to .env"
    cp .env.template .env
fi

# Start the controller
make start
