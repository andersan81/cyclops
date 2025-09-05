# 🚀 Development Guide

This guide helps you set up and run Cyclops in development mode.

## 📋 Prerequisites

- **Go 1.22+** - Backend development
- **Node.js 16+** - Frontend development
- **npm** - Package management
- **kubectl** - Kubernetes CLI (for local K8s development)

## ⚡ Quick Start

### 1. First-time Setup
```bash
./scripts/dev-setup.sh
```
This script will:
- Check for required dependencies (Go, Node.js, npm)
- Create `.env` files from templates
- Install Go and npm dependencies
- Make all scripts executable

### 2. Start Both Services
```bash
./scripts/dev-start.sh
```
This will start both the backend controller and frontend UI concurrently.

- **Backend API**: http://localhost:8888
- **Frontend UI**: http://localhost:3000

### 3. Individual Services

Start only the backend:
```bash
./scripts/dev-backend.sh
```

Start only the frontend:
```bash
./scripts/dev-frontend.sh
```

## 🔧 Environment Configuration

The scripts automatically create `.env` files from templates if they don't exist:

### Backend (`cyclops-ctrl/.env`)
- `DISABLE_TELEMETRY=true` - Disable telemetry in development
- `PORT=8888` - Controller API port
- `WATCH_NAMESPACE=cyclops` - Kubernetes namespace to watch
- `MODULE_TARGET_NAMESPACE=vision` - Default deployment namespace

### Frontend (`cyclops-ui/.env`)
- `REACT_APP_CYCLOPS_CTRL_HOST=http://localhost:8888` - Backend URL
- `REACT_APP_ENABLE_STREAMING=true` - Enable real-time features
- `REACT_APP_SUPPORT_EMAIL=your-email@example.com` - Support contact

## 🛠️ Manual Commands

If you prefer to run commands manually:

### Backend
```bash
cd cyclops-ctrl
make start
# OR
go run cmd/main/main.go
```

### Frontend
```bash
cd cyclops-ui
npm start
```

## 🐛 Troubleshooting

### Port Already in Use
If you get port errors:
- **Port 8888**: Another instance of the controller is running
- **Port 3000**: Another React app is running

Stop the conflicting processes or change ports in the `.env` files.

### Missing Dependencies
Run the setup script again:
```bash
./scripts/dev-setup.sh
```

### Go Module Issues
If you see Go module errors:
```bash
cd cyclops-ctrl
go mod download
go mod tidy
```

### npm Issues
If you see npm errors:
```bash
cd cyclops-ui
rm -rf node_modules package-lock.json
npm install
```

## 🔄 Development Workflow

1. **Start services**: `./scripts/dev-start.sh`
2. **Make changes** to code
3. **Backend changes**: The Go server will need manual restart
4. **Frontend changes**: React will auto-reload
5. **Stop services**: Press `Ctrl+C` in the terminal running `scripts/dev-start.sh`

## 📂 Project Structure

```
cyclops/
├── cyclops-ctrl/          # Go backend controller
│   ├── .env.template      # Backend environment template
│   ├── Makefile          # Build and run commands
│   └── cmd/main/         # Main entry point
├── cyclops-ui/           # React frontend
│   ├── .env.template     # Frontend environment template
│   ├── package.json      # npm dependencies
│   └── src/              # React source code
└── scripts/              # Development scripts
    ├── dev-setup.sh      # First-time setup script
    ├── dev-start.sh      # Start both services
    ├── dev-backend.sh    # Start backend only
    └── dev-frontend.sh   # Start frontend only
```

## 🔒 Security Notes

- `.env` files are gitignored and contain sensitive configuration
- Always use `.env.template` files as templates
- Review environment variables before running in production
