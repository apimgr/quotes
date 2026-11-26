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
| GET | `/api/v1/server/security` | Security contact information and conditional PGP key fingerprint |
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
| POST | `/api/v1/server/preferences` | Apply valid theme or language preferences only after this device accepts preference cookies; otherwise return `consent_required` and list the values as ignored |
| GET | `/api/v1/server/preferences/sync` | Export preferences as a URL and base64url code |
| GET | `/api/v1/server/preferences/sync/import` | Import allowed preferences from direct parameters or a code |
| POST | `/api/v1/server/mcp` | MCP Streamable HTTP endpoint; owner or linked Bearer token required |
| POST | `/api/mcp`, `/mcp` | Direct MCP aliases |
| GET | `/api/v1/server/swagger` | OpenAPI document |
| GET | `/api/swagger` | OpenAPI compatibility alias |
| GET | `/api/v1/server/graphql` | GraphQL endpoint description and query example |
| GET, POST | `/api/v1/server/graphql` | Read-only GraphQL query endpoint; GET accepts `query`, optional `variables`, and `operationName` |
| GET, POST | `/api/graphql` | GraphQL compatibility alias for the current version |
| GET | `/api/healthz` | Health compatibility route |
| GET | `/api/v1/server/security/report/{tracking_id}?token=...` | Read an authorized security-report status; responses are not cacheable |

Lists default to 250 items and allow up to 1000. Invalid categories, IDs, and pagination values return an error. The server applies configured per-client rate limits and request-size limits. API routes return JSON to browser and API clients; `Accept: text/plain` requests readable text output across JSON endpoints. The health API also has an explicit `.txt` endpoint.

## Public feeds and discovery

`GET /feed.xml` returns an Atom 1.0 feed and `GET /rss.xml` returns an RSS 2.0 feed with the same most recent public quote items. Both include absolute links built from the request's validated public origin, XML-escaped content, and `ETag`, `Last-Modified`, and `Cache-Control: public, max-age=300` headers. Conditional requests may return `304 Not Modified`. `GET /robots.txt` and `GET /sitemap.xml` provide search crawler rules and the public-page sitemap; `/.well-known/llms.txt` and `/llms.txt` describe the public API for AI clients.

Metrics services are `prometheus`, `grafana`, and `loki`. `/server/metrics` is the Prometheus exposition handler; `/server/metrics/grafana` returns the importable dashboard and `/server/metrics/loki` returns bounded recent log streams. Their generated bearer tokens are stored separately under `server.metrics.auth.tokens` in `secrets.yml`. Query-string credentials are not accepted. Keep these internal endpoints behind network firewall rules and do not proxy them publicly.

## Examples

```sh
curl -q -LSsf 'http://localhost:8080/api/v1/quotes/random?category=programming'
curl -q -LSsf 'http://localhost:8080/api/v1/quotes?author=Grace%20Hopper&page=1&limit=20'
curl -q -LSsf http://localhost:8080/api/v1/server/swagger
```

GraphQL supports `quote(id)`, `randomQuote(category)`, `quotes(category, author, page, limit)`, and `categories`. POST accepts JSON with a `query` field; GET accepts the query and optional `variables` and `operationName` parameters. MCP uses JSON-RPC over POST and requires an owner or linked Bearer token. Create a separate linked token for each MCP client; the operator token is not accepted for MCP. Its tools are `quotes_list`, `quotes_random`, `quote_get`, and `categories_list`; use `tools/list` before `tools/call`.

## Browser pages

The server-rendered catalog at `/` supports category and author filters and pagination. `/server/preferences` changes appearance and language; `/server/preferences/sync` exports or imports browser preferences. Theme and language cookies are stored only after the visitor accepts the preference-cookie category on that device; a save or import does not count as consent. When `server.privacy.data.sold` is enabled, `/server/privacy` and `/server/preferences` show the CCPA opt-out and `POST /server/ccpa` saves the device choice. `/server/preferences/tokens` issues, restores, and manages anonymous owner credentials without accounts. `/random` selects a random quote and `/quotes/{id}` shows a quote by ID. Standard pages include `/server/about`, `/server/help`, `/server/help/mcp`, `/server/healthz`, `/server/privacy`, `/server/terms`, `/server/dpo`, `/server/security`, `/server/security/policy`, `/server/security/thanks`, and `/server/contact`. The contact form accepts only abuse, DMCA, privacy, and security reports and is reached through the security page or security.txt. Interactive references are at `/server/docs/swagger` and `/server/docs/graphql`.

Frontend routes return HTML for browser requests. `Accept: text/html` selects HTML, `Accept: text/plain` returns a readable text rendering, and `Accept: application/json` returns the paired API representation where one exists. HTTP tools such as curl receive readable text when they omit an explicit browser format. Text browsers such as Lynx and w3m receive server-rendered HTML with executable scripts removed. The preference and reporting pages remain standard server-rendered forms. `/server/docs/swagger` serves its version-pinned Swagger UI assets locally and includes an OpenAPI JSON fallback when JavaScript is unavailable; `/server/docs/graphql` serves its version-pinned GraphiQL assets locally and links directly to the GraphQL endpoint for no-JavaScript clients.
