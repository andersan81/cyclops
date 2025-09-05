# Environment Configuration

## Setup

1. Copy the environment template:
   ```bash
   cp .env.template .env
   ```

2. Update the `.env` file with your configuration values as needed.

## Environment Variables

- `DISABLE_TELEMETRY`: Disable telemetry collection (true/false)
- `PORT`: Port for the controller to listen on
- `WATCH_NAMESPACE`: Kubernetes namespace to watch for Cyclops resources
- `WATCH_NAMESPACE_HELM`: Namespace to watch for Helm releases (optional)
- `CYCLOPS_VERSION`: Version of Cyclops
- `MODULE_TARGET_NAMESPACE`: Default namespace for deploying modules
- `MAX_CONCURRENT_RECONCILES`: Maximum concurrent reconciliations (optional)

## Security Note

The `.env` files are gitignored to prevent sensitive configuration from being committed to the repository.
