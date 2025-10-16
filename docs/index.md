# Quotes API Documentation

Welcome to the **Quotes API** documentation. This is a high-performance REST API server that provides access to multiple quote collections, built as a single static binary with embedded assets and data.

## Overview

Quotes API is designed to deliver quotes from various categories through a simple, efficient REST API. With 27,500 quotes across 5 collections, it's perfect for applications, websites, or services that need inspirational, humorous, or educational content.

### Key Features

- **5 Quote Collections**: quotes, anime, chucknorris, dadjokes, programming
- **27,500 Total Quotes**: 5,500 entries per collection
- **Single Binary**: All assets, templates, and data embedded via go:embed
- **High Performance**: ~10,000 requests/sec on a single core
- **Low Memory**: ~50MB with all collections loaded
- **Multi-platform**: Linux, Windows, macOS, BSD (amd64, arm64)
- **Docker Support**: Multi-arch images for easy deployment
- **Modern Web UI**: Responsive dark-themed interface
- **REST API**: Clean JSON responses with comprehensive error handling

## Quick Start

### Using Docker

The fastest way to get started is with Docker:

```bash
# Pull the latest image
docker pull ghcr.io/apimgr/quotes:latest

# Run the container
docker run -d \
  --name quotes \
  -p 8080:80 \
  -e ADMIN_USER=admin \
  -e ADMIN_PASSWORD=changeme \
  ghcr.io/apimgr/quotes:latest

# Access the API
curl http://localhost:8080/api/v1/quotes/random
```

### Using Binary

Download and run the binary directly:

```bash
# Download latest release (Linux AMD64)
wget https://github.com/apimgr/quotes/releases/latest/download/quotes-linux-amd64

# Make executable
chmod +x quotes-linux-amd64

# Run
./quotes-linux-amd64 --port 8080
```

### Using Docker Compose

For production deployments:

```bash
# Clone repository
git clone https://github.com/apimgr/quotes.git
cd quotes

# Start services
docker compose up -d

# Access at http://172.17.0.1:64180
```

## API Usage

### Get Random Quote

```bash
curl http://localhost:8080/api/v1/quotes/random
```

Response:
```json
{
  "success": true,
  "data": {
    "id": 1234,
    "quote": "The only way to do great work is to love what you do.",
    "author": "Steve Jobs",
    "category": "motivation",
    "tags": ["work", "passion", "success"]
  }
}
```

### Get Random Anime Quote

```bash
curl http://localhost:8080/api/v1/anime/random
```

### Get Random Chuck Norris Fact

```bash
curl http://localhost:8080/api/v1/chucknorris/random
```

### Get Random Dad Joke

```bash
curl http://localhost:8080/api/v1/dadjokes/random
```

### Get Random Programming Quote

```bash
curl http://localhost:8080/api/v1/programming/random
```

### Search Quotes

```bash
curl "http://localhost:8080/api/v1/quotes/search?q=success"
```

## Collections

### General Quotes (5,500)
Inspirational, motivational, and wisdom quotes from historical figures, leaders, and thinkers.

**Endpoint**: `/api/v1/quotes`

**Categories**: motivation, wisdom, life, success, happiness, courage

### Anime Quotes (5,500)
Memorable quotes from anime characters and series.

**Endpoint**: `/api/v1/anime`

**Categories**: action, drama, comedy, shonen, seinen

### Chuck Norris Facts (5,500)
Classic Chuck Norris facts and jokes.

**Endpoint**: `/api/v1/chucknorris`

**Categories**: humor, action, facts

### Dad Jokes (5,500)
Clean, family-friendly dad jokes and puns.

**Endpoint**: `/api/v1/dadjokes`

**Categories**: humor, puns, family-friendly

### Programming Quotes (5,500)
Programming-related humor, wisdom, and observations.

**Endpoint**: `/api/v1/programming`

**Categories**: humor, wisdom, software-development, debugging

## Web Interface

Access the modern web interface at `http://localhost:8080/`:

- Browse all collections
- Get random quotes
- Search functionality
- Dark/light theme toggle
- Mobile-responsive design

## Admin Panel

Access the admin panel at `http://localhost:8080/admin`:

- Server statistics
- Collection management
- Settings configuration
- Health monitoring
- Log viewer

Default credentials (change on first run):
- Username: `administrator`
- Password: Set via `ADMIN_PASSWORD` environment variable

## Architecture

### Technology Stack

**Backend**:
- Language: Go 1.23+
- Framework: Standard library (net/http)
- Database: SQLite3
- Templates: html/template

**Frontend**:
- Templates: Go html/template
- CSS: Vanilla CSS3 (~900 lines)
- JavaScript: Vanilla ES6+ (~130 lines)
- Theme: Dark mode default

**Infrastructure**:
- Container: Alpine 3.19
- Registry: GitHub Container Registry
- CI/CD: Jenkins + GitHub Actions

### Performance

- **Requests/sec**: ~10,000 (single core)
- **Memory**: ~50MB (all collections loaded)
- **Startup**: <100ms
- **Binary Size**: ~15MB (includes all data)

### Data Format

All quotes are stored in JSON format and embedded in the binary:

```json
{
  "id": 1,
  "quote": "Quote text here",
  "author": "Author Name",
  "category": "category-name",
  "tags": ["tag1", "tag2"]
}
```

## Documentation

### API Reference
Complete API documentation with all endpoints, parameters, and response formats.

[View API Reference](API.md)

### Server Administration
Deployment guides, configuration options, and troubleshooting.

[View Server Guide](SERVER.md)

### GitHub Repository
Source code, issues, and contributions.

[View on GitHub](https://github.com/apimgr/quotes)

## Support

- **Issues**: [GitHub Issues](https://github.com/apimgr/quotes/issues)
- **Repository**: [github.com/apimgr/quotes](https://github.com/apimgr/quotes)
- **Organization**: [github.com/apimgr](https://github.com/apimgr)

## License

MIT License - see [LICENSE.md](https://github.com/apimgr/quotes/blob/main/LICENSE.md)

---

**Version**: 0.0.1
**Last Updated**: 2025-10-14
