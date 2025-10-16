# API Reference

Complete API documentation for the Quotes API server.

## Base URL

```
http://localhost:8080/api/v1
```

## Response Format

All API responses follow a consistent JSON structure:

### Success Response

```json
{
  "success": true,
  "data": {
    // Response data here
  },
  "timestamp": "2025-10-14T12:00:00Z"
}
```

### Error Response

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message"
  },
  "timestamp": "2025-10-14T12:00:00Z"
}
```

## Quote Object

All quote collections return objects with this structure:

```json
{
  "id": 1,
  "quote": "Quote text here",
  "author": "Author Name",
  "category": "category-name",
  "tags": ["tag1", "tag2", "tag3"]
}
```

## Health & Status Endpoints

### GET /health

Check server health status.

**Request:**
```bash
curl http://localhost:8080/health
```

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2025-10-14T12:00:00Z"
}
```

### GET /status

Get detailed server status with collection statistics.

**Request:**
```bash
curl http://localhost:8080/status
```

**Response:**
```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "version": "0.0.1",
    "uptime": "2h30m15s",
    "collections": {
      "quotes": 5500,
      "anime": 5500,
      "chucknorris": 5500,
      "dadjokes": 5500,
      "programming": 5500
    },
    "total_quotes": 27500
  }
}
```

## General Quotes Collection

### GET /api/v1/quotes

Get all general quotes.

**Request:**
```bash
curl http://localhost:8080/api/v1/quotes
```

**Query Parameters:**
- `limit` (integer, optional): Limit number of results (default: 100, max: 1000)
- `offset` (integer, optional): Offset for pagination (default: 0)

**Response:**
```json
{
  "success": true,
  "data": {
    "quotes": [
      {
        "id": 1,
        "quote": "The only way to do great work is to love what you do.",
        "author": "Steve Jobs",
        "category": "motivation",
        "tags": ["work", "passion", "success"]
      }
    ],
    "total": 5500,
    "limit": 100,
    "offset": 0
  }
}
```

### GET /api/v1/quotes/random

Get a random general quote.

**Request:**
```bash
curl http://localhost:8080/api/v1/quotes/random
```

**Response:**
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

### GET /api/v1/quotes/:id

Get a specific quote by ID.

**Request:**
```bash
curl http://localhost:8080/api/v1/quotes/42
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 42,
    "quote": "Life is what happens when you're busy making other plans.",
    "author": "John Lennon",
    "category": "life",
    "tags": ["life", "wisdom", "philosophy"]
  }
}
```

**Error Response (Not Found):**
```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "Quote with ID 9999 not found"
  }
}
```

### GET /api/v1/quotes/search

Search general quotes by keyword.

**Request:**
```bash
curl "http://localhost:8080/api/v1/quotes/search?q=success&limit=10"
```

**Query Parameters:**
- `q` (string, required): Search query
- `limit` (integer, optional): Limit results (default: 50, max: 500)

**Response:**
```json
{
  "success": true,
  "data": {
    "quotes": [
      {
        "id": 15,
        "quote": "Success is not final, failure is not fatal.",
        "author": "Winston Churchill",
        "category": "success",
        "tags": ["success", "failure", "perseverance"]
      }
    ],
    "total": 42,
    "query": "success"
  }
}
```

## Anime Quotes Collection

### GET /api/v1/anime

Get all anime quotes.

**Request:**
```bash
curl http://localhost:8080/api/v1/anime?limit=50
```

**Query Parameters:** Same as `/api/v1/quotes`

### GET /api/v1/anime/random

Get a random anime quote.

**Request:**
```bash
curl http://localhost:8080/api/v1/anime/random
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 789,
    "quote": "The world isn't perfect, but it's there for us trying the best it can.",
    "author": "Roy Mustang",
    "category": "philosophy",
    "tags": ["fullmetal-alchemist", "wisdom", "life"]
  }
}
```

### GET /api/v1/anime/:id

Get a specific anime quote by ID.

**Request:**
```bash
curl http://localhost:8080/api/v1/anime/100
```

### GET /api/v1/anime/search

Search anime quotes.

**Request:**
```bash
curl "http://localhost:8080/api/v1/anime/search?q=strength"
```

## Chuck Norris Collection

### GET /api/v1/chucknorris

Get all Chuck Norris facts.

**Request:**
```bash
curl http://localhost:8080/api/v1/chucknorris
```

### GET /api/v1/chucknorris/random

Get a random Chuck Norris fact.

**Request:**
```bash
curl http://localhost:8080/api/v1/chucknorris/random
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 321,
    "quote": "Chuck Norris can divide by zero.",
    "author": "Chuck Norris",
    "category": "humor",
    "tags": ["math", "impossible", "superhuman"]
  }
}
```

### GET /api/v1/chucknorris/:id

Get a specific Chuck Norris fact by ID.

**Request:**
```bash
curl http://localhost:8080/api/v1/chucknorris/50
```

### GET /api/v1/chucknorris/search

Search Chuck Norris facts.

**Request:**
```bash
curl "http://localhost:8080/api/v1/chucknorris/search?q=roundhouse"
```

## Dad Jokes Collection

### GET /api/v1/dadjokes

Get all dad jokes.

**Request:**
```bash
curl http://localhost:8080/api/v1/dadjokes
```

### GET /api/v1/dadjokes/random

Get a random dad joke.

**Request:**
```bash
curl http://localhost:8080/api/v1/dadjokes/random
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 456,
    "quote": "Why don't scientists trust atoms? Because they make up everything!",
    "author": "Anonymous",
    "category": "science",
    "tags": ["science", "pun", "wordplay"]
  }
}
```

### GET /api/v1/dadjokes/:id

Get a specific dad joke by ID.

**Request:**
```bash
curl http://localhost:8080/api/v1/dadjokes/200
```

### GET /api/v1/dadjokes/search

Search dad jokes.

**Request:**
```bash
curl "http://localhost:8080/api/v1/dadjokes/search?q=chicken"
```

## Programming Quotes Collection

### GET /api/v1/programming

Get all programming quotes.

**Request:**
```bash
curl http://localhost:8080/api/v1/programming
```

### GET /api/v1/programming/random

Get a random programming quote.

**Request:**
```bash
curl http://localhost:8080/api/v1/programming/random
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 987,
    "quote": "There are only two hard things in Computer Science: cache invalidation and naming things.",
    "author": "Phil Karlton",
    "category": "wisdom",
    "tags": ["computer-science", "humor", "truth"]
  }
}
```

### GET /api/v1/programming/:id

Get a specific programming quote by ID.

**Request:**
```bash
curl http://localhost:8080/api/v1/programming/150
```

### GET /api/v1/programming/search

Search programming quotes.

**Request:**
```bash
curl "http://localhost:8080/api/v1/programming/search?q=debugging"
```

## Admin Endpoints

All admin endpoints require authentication via Bearer token.

### Authentication

Include the admin token in the `Authorization` header:

```bash
curl -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  http://localhost:8080/api/v1/admin/stats
```

### GET /api/v1/admin/stats

Get server statistics (requires authentication).

**Request:**
```bash
curl -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  http://localhost:8080/api/v1/admin/stats
```

**Response:**
```json
{
  "success": true,
  "data": {
    "server": {
      "version": "0.0.1",
      "uptime": "2h30m15s",
      "start_time": "2025-10-14T10:00:00Z"
    },
    "collections": {
      "quotes": 5500,
      "anime": 5500,
      "chucknorris": 5500,
      "dadjokes": 5500,
      "programming": 5500,
      "total": 27500
    },
    "requests": {
      "total": 12345,
      "success": 12100,
      "errors": 245,
      "rate": "124.5 req/s"
    },
    "memory": {
      "allocated": "52.4 MB",
      "sys": "68.2 MB",
      "gc_runs": 23
    }
  }
}
```

### GET /api/v1/admin/settings

Get all server settings (requires authentication).

**Request:**
```bash
curl -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  http://localhost:8080/api/v1/admin/settings
```

**Response:**
```json
{
  "success": true,
  "data": {
    "server_name": "Quotes API",
    "max_results": 1000,
    "enable_cors": true,
    "log_level": "info",
    "cache_enabled": true
  }
}
```

### POST /api/v1/admin/settings

Update server settings (requires authentication).

**Request:**
```bash
curl -X POST \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{"max_results": 500, "log_level": "debug"}' \
  http://localhost:8080/api/v1/admin/settings
```

**Response:**
```json
{
  "success": true,
  "data": {
    "updated": ["max_results", "log_level"],
    "settings": {
      "server_name": "Quotes API",
      "max_results": 500,
      "enable_cors": true,
      "log_level": "debug",
      "cache_enabled": true
    }
  }
}
```

## Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `SUCCESS` | 200 | Request successful |
| `BAD_REQUEST` | 400 | Invalid request parameters |
| `UNAUTHORIZED` | 401 | Missing or invalid authentication |
| `FORBIDDEN` | 403 | Insufficient permissions |
| `NOT_FOUND` | 404 | Resource not found |
| `INTERNAL_ERROR` | 500 | Internal server error |
| `SERVICE_UNAVAILABLE` | 503 | Service temporarily unavailable |

## Rate Limiting

Currently, there are no rate limits enforced. However, for production deployments, it's recommended to implement rate limiting at the reverse proxy level (nginx, Caddy, etc.).

## CORS

CORS is enabled by default for all origins. For production, configure specific origins in the admin settings.

## Pagination

Endpoints that return lists support pagination via `limit` and `offset` parameters:

```bash
# Get quotes 100-199
curl "http://localhost:8080/api/v1/quotes?limit=100&offset=100"
```

## Best Practices

### Caching

The API responses are cacheable. Implement client-side caching for better performance:

```bash
curl -H "Cache-Control: max-age=3600" \
  http://localhost:8080/api/v1/quotes/random
```

### Error Handling

Always check the `success` field in responses:

```javascript
const response = await fetch('http://localhost:8080/api/v1/quotes/random');
const data = await response.json();

if (data.success) {
  console.log(data.data.quote);
} else {
  console.error(data.error.message);
}
```

### Search Performance

For better search performance, use specific queries and limit results:

```bash
# Good: Specific query with limit
curl "http://localhost:8080/api/v1/quotes/search?q=success&limit=10"

# Avoid: Very broad queries without limit
curl "http://localhost:8080/api/v1/quotes/search?q=a"
```

## Examples

### JavaScript (Fetch API)

```javascript
// Get random quote
async function getRandomQuote() {
  const response = await fetch('http://localhost:8080/api/v1/quotes/random');
  const data = await response.json();

  if (data.success) {
    return data.data;
  }
  throw new Error(data.error.message);
}

// Usage
const quote = await getRandomQuote();
console.log(`${quote.quote} - ${quote.author}`);
```

### Python (requests)

```python
import requests

# Get random anime quote
def get_random_anime_quote():
    response = requests.get('http://localhost:8080/api/v1/anime/random')
    data = response.json()

    if data['success']:
        return data['data']
    raise Exception(data['error']['message'])

# Usage
quote = get_random_anime_quote()
print(f"{quote['quote']} - {quote['author']}")
```

### Go

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

type Response struct {
    Success bool `json:"success"`
    Data    struct {
        Quote  string `json:"quote"`
        Author string `json:"author"`
    } `json:"data"`
}

func getRandomQuote() (*Response, error) {
    resp, err := http.Get("http://localhost:8080/api/v1/quotes/random")
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var result Response
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }

    return &result, nil
}

func main() {
    quote, _ := getRandomQuote()
    fmt.Printf("%s - %s\n", quote.Data.Quote, quote.Data.Author)
}
```

### cURL

```bash
# Get random programming quote and format with jq
curl -s http://localhost:8080/api/v1/programming/random | \
  jq -r '"\(.data.quote) - \(.data.author)"'

# Search dad jokes about dogs
curl -s "http://localhost:8080/api/v1/dadjokes/search?q=dog&limit=5" | \
  jq '.data.quotes[]'

# Get server statistics (requires auth)
curl -s -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8080/api/v1/admin/stats | \
  jq '.data'
```

## Changelog

### Version 0.0.1 (2025-10-14)
- Initial API release
- 5 quote collections with 27,500 total quotes
- 25+ endpoints
- Admin authentication
- Search functionality
- Pagination support

---

**Version**: 0.0.1
**Last Updated**: 2025-10-14
