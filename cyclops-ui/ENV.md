# Environment Configuration

## Setup

1. Copy the environment template:

   ```bash
   cp .env.template .env
   ```

2. Update the `.env` file with your configuration values:
   - `REACT_APP_SUPPORT_EMAIL`: Email address for bug reports

## Environment Variables

- `NODE_ENV`: Application environment (development/production)
- `REACT_APP_CYCLOPS_CTRL_HOST`: Backend controller host URL
- `REACT_APP_ENABLE_STREAMING`: Enable streaming features (true/false)
- `REACT_APP_SUPPORT_EMAIL`: Email address for support and bug reports

## Security Note

The `.env` files are gitignored to prevent sensitive information like email addresses from being committed to the repository.
