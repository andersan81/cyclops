#!/bin/bash

# Development server startup script for Cyclops
# This script starts both the backend controller and the UI server

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[CYCLOPS DEV]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[CYCLOPS DEV]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[CYCLOPS DEV]${NC} $1"
}

print_error() {
    echo -e "${RED}[CYCLOPS DEV]${NC} $1"
}

# Find the project root directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Check if we found the correct project root
if [ ! -f "$PROJECT_ROOT/cyclops-ctrl/Makefile" ] || [ ! -f "$PROJECT_ROOT/cyclops-ui/package.json" ]; then
    print_error "Could not find cyclops project root. Make sure this script is in the scripts/ directory of the cyclops project."
    exit 1
fi

print_status "Project root: $PROJECT_ROOT"

# Check environment files
print_status "Checking environment configuration..."
if [ ! -f "$PROJECT_ROOT/cyclops-ui/.env" ]; then
    print_warning "No .env file found in cyclops-ui/"
    if [ -f "$PROJECT_ROOT/cyclops-ui/.env.template" ]; then
        print_status "Copying .env.template to .env"
        cp "$PROJECT_ROOT/cyclops-ui/.env.template" "$PROJECT_ROOT/cyclops-ui/.env"
        print_warning "Please review and update cyclops-ui/.env with your configuration"
    else
        print_error "No .env.template found in cyclops-ui/"
        exit 1
    fi
fi

if [ ! -f "$PROJECT_ROOT/cyclops-ctrl/.env" ]; then
    print_warning "No .env file found in cyclops-ctrl/"
    if [ -f "$PROJECT_ROOT/cyclops-ctrl/.env.template" ]; then
        print_status "Copying .env.template to .env"
        cp "$PROJECT_ROOT/cyclops-ctrl/.env.template" "$PROJECT_ROOT/cyclops-ctrl/.env"
        print_warning "Please review and update cyclops-ctrl/.env with your configuration"
    fi
fi

# Function to cleanup background processes
cleanup() {
    print_status "Shutting down services..."
    jobs -p | xargs -r kill
    wait
    print_success "Services stopped"
}

# Set trap to cleanup on script exit
trap cleanup EXIT INT TERM

# Start the backend controller
print_status "Starting Cyclops Controller (Backend)..."
cd "$PROJECT_ROOT/cyclops-ctrl"
make start &
BACKEND_PID=$!
cd "$PROJECT_ROOT"

# Wait a moment for the backend to start
sleep 3

# Check if backend is still running
if ! kill -0 $BACKEND_PID 2>/dev/null; then
    print_error "Backend failed to start"
    exit 1
fi

print_success "Backend started (PID: $BACKEND_PID)"

# Start the UI server
print_status "Starting Cyclops UI (Frontend)..."
cd "$PROJECT_ROOT/cyclops-ui"

# Check if node_modules exists
if [ ! -d "node_modules" ]; then
    print_status "Installing UI dependencies..."
    npm install
fi

npm start &
FRONTEND_PID=$!
cd "$PROJECT_ROOT"

# Wait a moment for the frontend to start
sleep 3

# Check if frontend is still running
if ! kill -0 $FRONTEND_PID 2>/dev/null; then
    print_error "Frontend failed to start"
    exit 1
fi

print_success "Frontend started (PID: $FRONTEND_PID)"

print_success "🚀 Cyclops development environment is running!"
print_status "Backend: http://localhost:8888 (Controller API)"
print_status "Frontend: http://localhost:3000 (React UI)"
print_status ""
print_status "Press Ctrl+C to stop all services"

# Wait for both processes
wait
