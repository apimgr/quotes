## Project description

Quotes is a full-stack Go web application serving a curated multi-category collection of quotes through a versioned REST API, GraphQL endpoint, and a server-side rendered web UI. Categories include general inspirational quotes, anime quotes, Chuck Norris facts, dad jokes, and programming quotes. All quote data is embedded in the binary at build time. A companion CLI tool lets users pull random or category-filtered quotes directly from the terminal. Deployed as a single self-contained static binary.

## Project variables

project_name: quotes
project_org: apimgr
internal_name: quotes
internal_org: apimgr
app_name: Quotes API
repo: https://github.com/apimgr/quotes
license: MIT
binary: quotes
client_binary: quotes-cli

## Business logic

### Product scope & non-goals

**In scope:**
- Multi-category quote dataset: inspirational (quotes), anime, chucknorris, dadjokes, programming
- Random quote retrieval (any category or filtered by category)
- Category-filtered retrieval
- Quote lookup by ID
- Author-filtered retrieval
- List all available categories with quote counts
- Full web frontend (server-side Go templates, dark/light/auto theme, PWA, mobile-first)
- Server pages: `/server/about`, `/server/help`, `/server/healthz`, `/server/privacy`, `/server/terms`
- CLI client (`quotes-cli`) for shell-pipeline use: `quotes-cli random --category programming`
- OpenAPI/Swagger docs at `/api/{api_version}/server/swagger`
- GraphQL at `/graphql`

**Non-goals:**
- No user accounts, registration, or login of any kind
- No admin web panel (server configured via `server.yml` only)
- No user-submitted or community quotes (curated dataset only, updated via releases)
- No paid tiers, no API keys, no rate-limited access tiers
- No voting, favoriting, or social features

### Roles & permissions

There are no user roles. All endpoints are public and require no authentication.

| Actor | Access |
|-------|--------|
| **Anonymous visitor (browser)** | Full read access to all web pages and API endpoints |
| **Anonymous API client (curl/CLI)** | Full read access to all API endpoints |
| **Server operator** | Configures server via `server.yml` only; no web management interface |

### Data model & sensitivity

**Quote record** (embedded at build time, no PII):

| Field | Type | Sensitivity |
|-------|------|-------------|
| `id` | integer — record identifier | Public |
| `quote` | string — quote text | Public |
| `author` | string — attributed author or character name | Public |
| `category` | string — category (anime, chucknorris, dadjokes, programming, quotes) | Public |

No PII stored or served.

### Trust boundaries & external services

| Boundary | Trust level | Notes |
|----------|-------------|-------|
| Quote dataset (embedded at build) | Fully trusted | Static, compiled into binary |
| Incoming HTTP requests | **Untrusted** | All query parameters validated |

No external services called at runtime.

### Threat model & abuse cases

**Primary assets:** service availability.

**Attacker/abuser goals:**
- DoS via high-rate requests
- Bulk scraping of the full dataset

**Defenses:**
- Rate limiting on all endpoints
- Request size limits on all inputs
- Paginated list endpoints limit per-request data volume
- No user accounts eliminates credential stuffing and privilege escalation entirely

### Security decisions & exceptions

- **No authentication on any endpoint**: intentional. Public read-only reference API.
- **All responses include `Access-Control-Allow-Origin: *`**: intentional. Public data API designed for cross-origin browser use.
