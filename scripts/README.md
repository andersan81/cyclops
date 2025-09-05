# Development Scripts

This directory contains scripts to help with Cyclops development.

## Scripts

- **`dev-setup.sh`** - First-time setup script that installs dependencies and creates environment files
- **`dev-start.sh`** - Starts both backend and frontend services concurrently
- **`dev-backend.sh`** - Starts only the backend controller
- **`dev-frontend.sh`** - Starts only the frontend UI

## Usage

All scripts can be run from anywhere in the project (they automatically detect the project root):

```bash
# From project root
./scripts/dev-setup.sh     # First time only
./scripts/dev-start.sh      # Start both services

# Or from anywhere in the project
cd cyclops-ui/
../scripts/dev-start.sh

# Or from the scripts directory
cd scripts/
./dev-setup.sh
./dev-start.sh
```

For more detailed information, see the main [DEVELOPMENT.md](../DEVELOPMENT.md) file.
