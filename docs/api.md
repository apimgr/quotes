# API

All API routes are read-only and public. JSON responses use `{ "ok": true, "data": ... }`; list responses also include `pagination` with page, limit, total, and pages. Errors contain `ok: false`, an error code, and a human-readable message. CORS permits public cross-origin reads.

## REST routes

| Method | Route | Purpose |
|---|---|---|
| GET | `/api/v1/` | API name, version, category counts |
| GET | `/api/v1/categories` | Category-to-count map |
| GET | `/api/v1/quotes` | Paginated list; optional `category`, `author`, `page`, `limit` |
| GET | `/api/v1/quotes/random` | Random quote; optional `category` |
| GET | `/api/v1/quotes/{id}` | Retrieve by global integer ID |
| GET | `/api/v1/categories/{category}/quotes` | Paginated category list |
| GET | `/api/v1/server/{about,help,privacy,terms,healthz}` | Machine-readable server information |
| GET | `/api/v1/server/swagger` | OpenAPI document |
| POST | `/api/v1/server/graphql` | GraphQL query endpoint |

List endpoints default to 250 items and allow up to 1000. `/api/v1/quotes` also supports `Accept: text/plain`; `.txt` endpoints are available for list and random results. Query and request bodies have size limits; per-client rate limits are enabled by default.

## Examples

```sh
curl 'http://localhost:8080/api/v1/quotes/random?category=programming'
curl 'http://localhost:8080/api/v1/quotes?author=Grace%20Hopper&page=1&limit=20'
curl http://localhost:8080/api/v1/server/swagger
```

GraphQL supports `quote(id)`, `randomQuote(category)`, `quotes(category, author, page, limit)`, and `categories`.

Interactive documentation is available at `/server/docs/swagger`; GraphQL explorer at `/server/docs/graphql`.

## Browser pages

The server-rendered catalog is available at `/`. It supports category and author filters, pagination, and theme selection. Choose **Get a random quote** to open `/random` (optionally with a `category` query parameter), or select a quote ID to open `/quotes/{id}`. The quote detail page links back to the catalog and to another random quote.
