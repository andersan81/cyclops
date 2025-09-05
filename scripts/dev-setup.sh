#!/bin/bash

# Development environment setup script for Cyclops

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${BLUE}[SETUP]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SETUP]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[SETUP]${NC} $1"
}

print_error() {
    echo -e "${RED}[SETUP]${NC} $1"
}

print_status "🔧 Setting up Cyclops development environment..."

# Check if we're in the right directory
if [ ! -f "../cyclops-ctrl/Makefile" ] || [ ! -f "../cyclops-ui/package.json" ]; then
    print_error "Please run this script from the scripts directory in the cyclops project"
    exit 1
fi

# Check Go installation
print_status "Checking Go installation..."
if ! command -v go &> /dev/null; then
    print_error "Go is not installed. Please install Go 1.22+ and try again."
    exit 1
fi

GO_VERSION=$(go version | cut -d' ' -f3 | sed 's/go//')
print_success "Go $GO_VERSION found"

# Check Node.js installation
print_status "Checking Node.js installation..."
if ! command -v node &> /dev/null; then
    print_error "Node.js is not installed. Please install Node.js and try again."
    exit 1
fi

NODE_VERSION=$(node --version)
print_success "Node.js $NODE_VERSION found"

# Check npm installation
if ! command -v npm &> /dev/null; then
    print_error "npm is not installed. Please install npm and try again."
    exit 1
fi

NPM_VERSION=$(npm --version)
print_success "npm $NPM_VERSION found"

# Setup backend environment
print_status "Setting up backend environment..."
cd ../cyclops-ctrl

if [ ! -f ".env" ]; then
    if [ -f ".env.template" ]; then
        print_status "Creating .env from template"
        cp .env.template .env
        print_warning "Please review cyclops-ctrl/.env and update as needed"
    else
        print_warning "No .env.template found in cyclops-ctrl/"
    fi
else
    print_success "Backend .env file already exists"
fi

# Download Go dependencies
print_status "Downloading Go dependencies..."
go mod download
print_success "Go dependencies downloaded"

cd ../scripts

# Setup frontend environment
print_status "Setting up frontend environment..."
cd ../cyclops-ui

if [ ! -f ".env" ]; then
    if [ -f ".env.template" ]; then
        print_status "Creating .env from template"
        cp .env.template .env
        print_warning "Please review cyclops-ui/.env and update as needed"
    else
        print_warning "No .env.template found in cyclops-ui/"
    fi
else
    print_success "Frontend .env file already exists"
fi

# Install npm dependencies
print_status "Installing npm dependencies..."
npm install
print_success "npm dependencies installed"

cd ../scripts

# Make scripts executable
print_status "Making development scripts executable..."
chmod +x dev-start.sh dev-backend.sh dev-frontend.sh dev-setup.sh

print_success "✅ Development environment setup complete!"
print_status ""
print_status "Available commands:"
print_status "  ./dev-start.sh     - Start both backend and frontend"
print_status "  ./dev-backend.sh   - Start only the backend controller"
print_status "  ./dev-frontend.sh  - Start only the frontend UI"
print_status ""
print_status "Environment files created:"
if [ -f "../cyclops-ctrl/.env" ]; then
    print_status "  ✅ cyclops-ctrl/.env"
else
    print_warning "  ❌ cyclops-ctrl/.env (missing)"
fi
if [ -f "../cyclops-ui/.env" ]; then
    print_status "  ✅ cyclops-ui/.env"
else
    print_warning "  ❌ cyclops-ui/.env (missing)"
fi
print_status ""
print_warning "🔍 Please review the .env files and update with your configuration before running the services."
