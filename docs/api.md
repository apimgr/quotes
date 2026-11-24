# API

Quote reads are public and read-only. MCP and owner-token management require the credentials described below. REST success responses use `{ "ok": true, "data": ... }`; list responses also include `pagination` (`page`, `limit`, `total`, `pages`). Errors use `{ "ok": false, "error": "CODE", "message": "..." }`. CORS allows public cross-origin quote access.

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
| GET | `/api/v1/server/healthz.txt` | Plain-text health status |
| GET | `/api/v1/server/help/mcp` | MCP connection details and tools |
| GET | `/api/v1/server/dpo` | Data-protection contact details |
| GET | `/api/v1/server/security` | Security contact information |
| GET | `/api/v1/server/security/policy` | Vulnerability disclosure policy |
| GET | `/api/v1/server/security/thanks` | Security acknowledgments |
| GET | `/api/v1/server/metrics[/{service}]` | Internal metrics; each service requires its own Bearer token. |
| GET | `/api/metrics[/{service}]` | Unversioned internal metrics aliases. |
| GET | `/metrics[/{service}]` | Root internal aliases when `server.metrics.root.enabled` is true. |
| GET, POST | `/api/v1/server/contact` | Contact report form definition and validated report submission |
| GET | `/api/v1/server/preferences/tokens` | List masked owner and linked token metadata (primary credential required) |
| POST | `/api/v1/server/preferences/tokens/primary` | Issue a primary owner token |
| POST | `/api/v1/server/preferences/tokens/primary/{reveal,rotate}` | Reveal or rotate the primary token |
| DELETE | `/api/v1/server/preferences/tokens/primary` | Revoke the primary token and its owner |
| POST | `/api/v1/server/preferences/tokens` | Create a linked token |
| POST | `/api/v1/server/preferences/tokens/{id}/{reveal,rotate}` | Reveal or rotate a linked token |
| PATCH, DELETE | `/api/v1/server/preferences/tokens/{id}` | Update or revoke a linked token |
| GET | `/api/v1/server/preferences` | Read the browser's theme and language preferences |
| POST | `/api/v1/server/preferences` | Apply valid theme or language preferences as browser cookies |
| GET | `/api/v1/server/preferences/sync` | Export preferences as a URL and base64url code |
| GET | `/api/v1/server/preferences/sync/import` | Import allowed preferences from direct parameters or a code |
| POST | `/api/v1/server/mcp` | MCP Streamable HTTP endpoint; owner or linked Bearer token required |
| POST | `/api/mcp`, `/mcp` | Direct MCP aliases |
| GET | `/api/v1/server/swagger` | OpenAPI document |
| GET | `/api/swagger` | OpenAPI compatibility alias |
| POST | `/api/v1/server/graphql` | GraphQL query endpoint |
| POST | `/api/graphql` | GraphQL compatibility alias |
| GET | `/api/healthz` | Health compatibility route |

Lists default to 250 items and allow up to 1000. Invalid categories, IDs, and pagination values return an error. The server applies configured per-client rate limits and request-size limits. `Accept: text/plain` on quote JSON routes requests text output.

Metrics services are `prometheus`, `grafana`, and `loki`. `/server/metrics` is the Prometheus exposition handler; `/server/metrics/grafana` returns the importable dashboard and `/server/metrics/loki` returns bounded recent log streams. Their generated bearer tokens are stored separately under `server.metrics.auth.tokens` in `secrets.yml`. Query-string credentials are not accepted. Keep these internal endpoints behind network firewall rules and do not proxy them publicly.

## Examples

```sh
curl 'http://localhost:8080/api/v1/quotes/random?category=programming'
curl 'http://localhost:8080/api/v1/quotes?author=Grace%20Hopper&page=1&limit=20'
curl http://localhost:8080/api/v1/server/swagger
```

GraphQL supports `quote(id)`, `randomQuote(category)`, `quotes(category, author, page, limit)`, and `categories`. Send JSON with a `query` field to either GraphQL route. MCP uses JSON-RPC over POST and requires an owner or linked Bearer token. Create a separate linked token for each MCP client; the operator token is not accepted for MCP. Its tools are `quotes_list`, `quotes_random`, `quote_get`, and `categories_list`; use `tools/list` before `tools/call`.

## Browser pages

The server-rendered catalog at `/` supports category and author filters and pagination. `/server/preferences` changes appearance and language; `/server/preferences/sync` exports or imports browser preferences. `/server/preferences/tokens` issues, restores, and manages anonymous owner credentials without accounts. `/random` selects a random quote and `/quotes/{id}` shows a quote by ID. Standard pages include `/server/about`, `/server/help`, `/server/help/mcp`, `/server/healthz`, `/server/privacy`, `/server/terms`, `/server/dpo`, `/server/security`, `/server/security/policy`, `/server/security/thanks`, and `/server/contact`. The contact form accepts only abuse, DMCA, privacy, and security reports and is reached through the security page or security.txt. Interactive references are at `/server/docs/swagger` and `/server/docs/graphql`.
