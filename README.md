# Quotes API

A simple and elegant REST API for serving inspirational quotes with a web UI.

## Features

- Random quote generation via REST API
- 20+ pre-loaded inspirational quotes
- Multiple query options (by ID, category, author)
- Web UI with interactive quote display
- Admin panel with authentication
- SQLite database for settings and admin users
- Docker support with health checks
- Multi-platform binaries (Linux, Windows, macOS, BSD)

## Quick Start

### Using Docker

```bash
# Pull and run the latest image
docker run -d \
  --name quotes \
  -p 8080:80 \
  -v ./data:/data \
  ghcr.io/apimgr/quotes:latest
```

Access the API at `http://localhost:8080`

### Using Docker Compose

```bash
# Production deployment
docker-compose up -d

# Development/testing
docker-compose -f docker-compose.test.yml up -d
```

### Using Binary

```bash
# Download the latest release
wget https://github.com/apimgr/quotes/releases/latest/download/quotes-linux-amd64.tar.gz

# Extract
tar -xzf quotes-linux-amd64.tar.gz

# Run
./quotes --port 8080
```

## API Endpoints

### Public Endpoints

- `GET /api/v1/random` - Get a random quote
- `GET /api/v1/quotes` - Get all quotes
- `GET /api/v1/quotes/{id}` - Get a specific quote by ID
- `GET /api/v1/quotes/category/{category}` - Get quotes by category
- `GET /api/v1/quotes/author/{author}` - Get quotes by author
- `GET /api/v1/status` - Get API status and version

### Admin Endpoints (Authentication Required)

- `GET /api/v1/admin/settings` - Get all settings
- `POST /api/v1/admin/settings` - Set or update a setting
- `DELETE /api/v1/admin/settings/{key}` - Delete a setting

### Example Request

```bash
# Get a random quote
curl http://localhost:8080/api/v1/random

# Response
{
  "success": true,
  "data": {
    "id": 1,
    "quote": "The only way to do great work is to love what you do.",
    "author": "Steve Jobs",
    "category": "inspiration"
  }
}
```

## Configuration

### Environment Variables

- `PORT` - Server port (default: 8080)
- `ADDRESS` - Server address (default: 0.0.0.0)
- `CONFIG_DIR` - Configuration directory
- `DATA_DIR` - Data directory
- `LOGS_DIR` - Logs directory
- `DB_PATH` - Database file path
- `ADMIN_USER` - Admin username (first run only)
- `ADMIN_PASSWORD` - Admin password (first run only)
- `ADMIN_TOKEN` - Admin API token (first run only)

### Admin Credentials

On first run, admin credentials are automatically generated and saved to:
- `$CONFIG_DIR/admin-credentials.txt`

Or set them manually using environment variables.

## Building from Source

### Prerequisites

- Go 1.23 or later
- Make (optional)

### Build Commands

```bash
# Build for all platforms
make build

# Build for current platform only
go build -o quotes ./src

# Run tests
make test

# Create release artifacts
make release

# Build Docker image (development)
make docker-dev
```

## Development

### Project Structure

```
quotes/
├── src/
│   ├── data/
│   │   └── quotes.json       # Quote data
│   ├── quotes/               # Quote service
│   ├── database/             # Database layer
│   ├── paths/                # OS-specific paths
│   ├── server/               # HTTP server
│   └── main.go              # Entry point
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── README.md
```

### Adding New Quotes

Edit `src/data/quotes.json` and add your quote:

```json
{
  "id": 21,
  "quote": "Your inspirational quote here",
  "author": "Author Name",
  "category": "category-name"
}
```

## Docker Deployment

### Production

Port: `172.17.0.1:64180:80`
Storage: `./rootfs/`

```bash
docker-compose up -d
```

### Development

Port: `64181:80`
Storage: `/tmp/quotes/rootfs/`

```bash
docker-compose -f docker-compose.test.yml up -d
```

## Web UI

Access the web interface at:
- Home: `http://localhost:8080/`
- Admin Panel: `http://localhost:8080/admin`

## License

MIT License - see LICENSE.md for details

## Credits

Created by apimgr
