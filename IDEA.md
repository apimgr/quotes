## Project description

Quotes is a public, read-only quote service with a server-rendered browser, versioned REST API, GraphQL endpoint, and companion CLI. Its embedded collections cover general quotes, anime, Chuck Norris facts, dad jokes, and programming. Operators distribute collection updates with application releases.

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
owner_token: quotes_owner_token_qwhzkfta

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
- Standard server pages and preferences, help, security, privacy, contact, and integration surfaces required by AI.md PART 16
- CLI client (`quotes-cli`) for shell-pipeline use: `quotes-cli random --category programming`
- OpenAPI/Swagger docs at `/api/{api_version}/server/swagger`
- GraphQL queries at `GET` and `POST /api/graphql` (the canonical route is `/api/v1/server/graphql`; GET accepts query-only operations)

**Non-goals:**
- No user accounts, registration, or login. Template-required owner and linked credentials are non-account identities for token management and MCP access
- No admin web panel (server configured via `server.yml` only)
- No user-submitted or community quotes (curated dataset only, updated via releases)
- No paid tiers or customer API keys. Template-required abuse rate limiting and linked MCP credentials do not create access tiers
- No voting, favoriting, or social features

### Roles & permissions

There are no user accounts. Public browsing, REST, and GraphQL require no authentication. Template-required MCP and credential-management actions use revocable credentials without creating accounts.

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
| `source` | string — embedded provenance key (`original`); not returned by the API | Internal |

Quote records contain no PII. The template-required contact reporting flow handles reporter-submitted PII under its configured retention and delivery rules. Browser preference cookies contain no PII.

The owner, apimgr, declares the five embedded collections theirs. AI writes original content at the owner's direction, and the owner releases it under the project's MIT license. No third-party quote corpus is imported. Each collection currently ships 30 records; every embedded record has `source: "original"`. Future collection updates must preserve the same ownership and attribution boundary.

### Trust boundaries & external services

| Boundary | Trust level | Notes |
|----------|-------------|-------|
| Quote dataset (embedded at build) | Fully trusted | Static, compiled into binary |
| Incoming HTTP requests | **Untrusted** | Query parameters, credentials, and report submissions are validated |
| Contact report submitters | **Untrusted** | Abuse, DMCA, privacy, and security reports can contain untrusted text and links |

The quote retrieval path uses only the embedded catalog. Template-required operator-configured subsystems may contact their explicitly configured mail, certificate, GeoIP, Tor, or I2P services; each integration must document its trust assumptions and failure behavior before it is enabled.

### Threat model & abuse cases

**Primary assets:** service availability, quote-content provenance, operator credentials and configuration secrets, and any personal information submitted through the required report workflow.

**Attacker/abuser goals:**
- DoS via high-rate requests
- Bulk scraping of the full dataset

**Defenses:**
- Rate limiting on all endpoints
- Request size limits on all inputs
- Paginated list endpoints limit per-request data volume
- No user accounts removes account registration and account-based privilege escalation; credential actions remain protected by token validation
- Report submissions use the template-defined anti-spam, content, origin, and rate-limit checks and never add quotes to the dataset

### Security decisions & exceptions

- **Public read-only REST and GraphQL**: intentional. MCP and token-management surfaces follow AI.md credential requirements.
- **The quote API is cross-origin readable**: intentional. It serves public data. Other routes follow the configured CORS policy and credential rules in AI.md.

### API routes

Project-defined quote routes are public and read-only. Template-defined pages and protocols follow AI.md PART 16 and are included here so the product surface remains explicit.

| Method | Path | Auth | Purpose and request/response summary |
|---|---|---|---|
| GET | `/api/v1` | None | API identity, version, and category counts. |
| GET | `/api/autodiscover` | None or bearer | Unversioned discovery envelope with server/API versions, public origins, exposed overlay addresses, and integrity-verified CLI release metadata. |
| GET | `/cli/binaries/quotes-cli-{os}-{arch}[.exe]` | None, or bearer when configured | Streams the matching operator-provisioned, SHA-256-verified CLI release binary. |
| GET | `/api/v1/categories` | None | Category names and counts. |
| GET | `/api/v1/quotes` | None | Paginated quotes; optional category, author, page, and limit. |
| GET | `/api/v1/categories/{category}/quotes` | None | Paginated quotes within a category. |
| GET | `/api/v1/quotes/random` | None | One random quote, optionally filtered by category. |
| GET | `/api/v1/quotes/{id}` | None | One quote by global ID. |
| GET | `/api/v1/quotes.txt` | None | Plain-text quote list. |
| GET | `/api/v1/quotes/random.txt` | None | Plain-text random quote. |
| GET | `/api/v1/server/{about,help,privacy,terms}` | None | Machine-readable page content. |
| GET | `/api/v1/server/healthz` | None | Health status. |
| GET | `/api/v1/server/swagger`, `/api/swagger` | None | OpenAPI document and direct alias. |
| GET, POST | `/api/v1/server/graphql`, `/api/graphql` | None | Read-only GraphQL queries; GET accepts `query`, optional `variables`, and `operationName`. |
| GET | `/api/healthz` | None | Health compatibility route. |
| GET | `/api/v1/server/preferences` | None | Read/update browser preferences through cookies. |
| GET, POST | `/api/v1/server/preferences/sync` | None | Export/import browser preferences. |
| GET | `/api/v1/server/preferences/sync/import` | None | Apply validated shared preference values. |
| GET | `/api/v1/server/help/mcp` | None | MCP connection information. |
| POST | `/api/v1/server/preferences/tokens/primary` | None unless a valid owner token is already supplied | Issue a primary owner token; rejects callers already presenting a valid owner credential. |
| GET, POST | `/api/v1/server/preferences/tokens` | Primary owner token | List token metadata as `primary`, `linked`, and `max_linked`, or create a linked token with `name`, optional `scope` (`read` or `read-write`, default `read-write`), and optional `expires` or `expires_on`. |
| POST | `/api/v1/server/preferences/tokens/{id}/{reveal,rotate}` | Primary owner token | Reveal or rotate a linked token; scope remains unchanged on rotation. |
| PATCH, DELETE | `/api/v1/server/preferences/tokens/{id}` | Primary owner token | Rename/change expiry or revoke a linked token by its public UUIDv7 ID. Scope changes require revoke and recreation. |
| POST | `/api/v1/server/preferences/tokens/primary/{reveal,rotate}` | Primary owner token | Reveal or rotate the primary token. |
| DELETE | `/api/v1/server/preferences/tokens/primary` | Primary owner token | Revoke the primary token and delete the anonymous owner and its data. |
| POST | `/api/v1/server/mcp`, `/api/mcp`, `/mcp` | Linked credential | Streamable HTTP access to quote tools. |
| GET | `/api/v1/server/dpo` | None | Data-protection contact details. |
| GET | `/api/v1/server/security` | None | Security reporting channels. |
| GET | `/api/v1/server/security/policy` | None | Disclosure policy. |
| GET | `/api/v1/server/security/thanks` | None | Security acknowledgments. |
| GET | `/server/metrics`, `/server/metrics/prometheus` | Prometheus metrics token | Internal Prometheus exposition. |
| GET | `/server/metrics/grafana` | Grafana metrics token | Internal Grafana dashboard JSON. |
| GET | `/server/metrics/loki` | Loki metrics token | Internal recent structured log streams. |
| GET | `/api/v1/server/metrics`, `/api/v1/server/metrics/prometheus`, `/api/metrics`, `/api/metrics/prometheus` | Prometheus metrics token | Internal Prometheus aliases. |
| GET | `/api/v1/server/metrics/grafana`, `/api/metrics/grafana` | Grafana metrics token | Internal Grafana dashboard aliases. |
| GET | `/api/v1/server/metrics/loki`, `/api/metrics/loki` | Loki metrics token | Internal Loki stream aliases. |
| GET | `/metrics`, `/metrics/prometheus` | Prometheus metrics token | Optional internal root aliases when enabled. |
| GET | `/metrics/grafana` | Grafana metrics token | Optional internal root alias when enabled. |
| GET | `/metrics/loki` | Loki metrics token | Optional internal root alias when enabled. |
| GET, POST | `/api/v1/server/contact` | None | Report-form description and validated report submission. |
| GET | `/server/security/report/{tracking_id}` | Report status token | Read the status, maintainer comment, and disclosure date for a submitted security report. |

The browser surface includes `/`, `/random`, `/quotes/{id}`, `/server/about`, `/server/privacy`, `/server/contact`, `/server/help`, `/server/help/mcp`, `/server/terms`, `/server/dpo`, `/server/security`, `/server/healthz`, `/server/preferences`, `/server/preferences/sync`, `/server/preferences/tokens`, `/server/docs/swagger`, and `/server/docs/graphql`. Optional pages disabled by configuration are not registered and return 404.
