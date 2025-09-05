#!/bin/bash

# Start only the Cyclops UI (Frontend)

set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${BLUE}[FRONTEND]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[FRONTEND]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[FRONTEND]${NC} $1"
}

# Check if we're in the right directory
if [ ! -f "../cyclops-ui/package.json" ]; then
    echo "Please run this script from the scripts directory in the cyclops project"
    exit 1
fi

print_status "Starting Cyclops UI..."
cd ../cyclops-ui

# Check environment file
if [ ! -f ".env" ]; then
    if [ -f ".env.template" ]; then
        print_status "Copying .env.template to .env"
        cp .env.template .env
        print_warning "Please review and update .env with your configuration"
    fi
fi

# Install dependencies if needed
if [ ! -d "node_modules" ]; then
    print_status "Installing dependencies..."
    npm install
fi

# Start the development server
print_status "Starting React development server..."
npm start
