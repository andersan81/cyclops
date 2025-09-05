# Development Scripts

This directory contains scripts to help with Cyclops development.

## Scripts

- **`dev-setup.sh`** - First-time setup script that installs dependencies and creates environment files
- **`dev-start.sh`** - Starts both backend and frontend services concurrently
- **`dev-backend.sh`** - Starts only the backend controller
- **`dev-frontend.sh`** - Starts only the frontend UI

## Usage

All scripts should be run from this `scripts/` directory:

```bash
cd scripts/
./dev-setup.sh     # First time only
./dev-start.sh      # Start both services
```

For more detailed information, see the main [DEVELOPMENT.md](../DEVELOPMENT.md) file.
