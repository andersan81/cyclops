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

# Find the project root directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Check if we found the correct project root
if [ ! -f "$PROJECT_ROOT/cyclops-ctrl/Makefile" ]; then
    echo "Could not find cyclops project root. Make sure this script is in the scripts/ directory of the cyclops project."
    exit 1
fi

print_status "Starting Cyclops Controller..."
cd "$PROJECT_ROOT/cyclops-ctrl"

# Check environment file
if [ ! -f ".env" ] && [ -f ".env.template" ]; then
    print_status "Copying .env.template to .env"
    cp .env.template .env
fi

# Start the controller
make start
