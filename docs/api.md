# API

All API routes are public and read-only. REST success responses use `{ "ok": true, "data": ... }`; list responses also include `pagination` (`page`, `limit`, `total`, `pages`). Errors use `{ "ok": false, "error": "CODE", "message": "..." }`. CORS allows public cross-origin access.

## REST routes

| Method | Route | Purpose |
|---|---|---|
| GET | `/api/v1` | API name, version, and category counts |
| GET | `/api/v1/categories` | Category-to-count map |
| GET | `/api/v1/quotes` | Paginated list; optional `category`, `author`, `page`, `limit` |
| GET | `/api/v1/categories/{category}/quotes` | Paginated category list |
| GET | `/api/v1/quotes/random` | Random quote; optional `category` |
| GET | `/api/v1/quotes/{id}` | Retrieve by global integer ID |
| GET | `/api/v1/quotes.txt` | Plain-text list |
| GET | `/api/v1/quotes/random.txt` | Plain-text random quote |
| GET | `/api/v1/server/{about,help,privacy,terms}` | Machine-readable page information |
| GET | `/api/v1/server/healthz` | JSON health status |
| GET | `/api/v1/server/help/mcp` | MCP connection details and tools |
| GET | `/api/v1/server/dpo` | Data-protection contact details |
| GET | `/api/v1/server/security` | Security contact information |
| POST | `/api/v1/server/mcp` | Bearer-authenticated MCP Streamable HTTP endpoint |
| POST | `/api/mcp`, `/mcp` | Direct MCP aliases |
| GET | `/api/v1/server/swagger` | OpenAPI document |
| GET | `/api/swagger` | OpenAPI compatibility alias |
| POST | `/api/v1/server/graphql` | GraphQL query endpoint |
| POST | `/api/graphql` | GraphQL compatibility alias |
| GET | `/api/healthz` | Health compatibility route |

Lists default to 250 items and allow up to 1000. Invalid categories, IDs, and pagination values return an error. The server applies configured per-client rate limits and request-size limits. `Accept: text/plain` on quote JSON routes requests text output.

## Examples

```sh
curl 'http://localhost:8080/api/v1/quotes/random?category=programming'
curl 'http://localhost:8080/api/v1/quotes?author=Grace%20Hopper&page=1&limit=20'
curl http://localhost:8080/api/v1/server/swagger
```

GraphQL supports `quote(id)`, `randomQuote(category)`, `quotes(category, author, page, limit)`, and `categories`. Send JSON with a `query` field to either GraphQL route. MCP uses JSON-RPC over POST and requires the operator Bearer token from the protected `secrets.yml` file. Its tools are `quotes_list`, `quotes_random`, `quote_get`, and `categories_list`; use `tools/list` before `tools/call`.

## Browser pages

The server-rendered catalog at `/` supports category and author filters, pagination, and system/light/dark themes. `/random` selects a random quote and `/quotes/{id}` shows a quote by ID. The about, help, help/MCP, health, privacy, terms, data-protection, and security pages are at `/server/about`, `/server/help`, `/server/help/mcp`, `/server/healthz`, `/server/privacy`, `/server/terms`, `/server/dpo`, and `/server/security`. Interactive references are at `/server/docs/swagger` and `/server/docs/graphql`.
