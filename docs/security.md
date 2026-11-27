# Security

Quotes is a public, read-only quote service with no user accounts or community quote submissions. MCP requires a valid Bearer token; public quote reads remain unauthenticated. The operator token (`server.token`) is accepted for the operator's own tooling and exposes operator-only tools, so it must not be shared with another person or a third-party client. Anonymous owners have no profile, password, or login. Their primary `global` token can manage revocable linked tokens. Linked tokens have `read` or `read-write` scope, and API routes accept Bearer credentials rather than browser cookies. Contact reports are validated, rate-limited, encrypted at rest, and handled separately from the quote catalog. Incoming query parameters are validated, request paths and bodies are bounded, and lists are paginated. CORS is open for public quote data.

The rate limiter, report spam controls, and GeoIP filter share a normalized client address. Forwarding headers are accepted only when the TCP peer matches `server.trusted_proxies.additional`; otherwise every forwarded value is ignored. Configure trusted proxy CIDRs only for proxies that overwrite the selected header. Responses include content-type sniffing, framing, referrer, content-security, and permissions policy headers. `/server/healthz` negotiates HTML, plain text, or JSON using `Accept`; `/api/v1/server/healthz` and `/api/healthz` return JSON. `/.well-known/security.txt` publishes the project security contact route.

Quote records are embedded in the executable. The current collections are original project content created at the project owner's direction. Any third-party material added later must have documented redistribution rights, source attribution, and license information before it is embedded.

The canonical machine-readable contact file is `/.well-known/security.txt`; `/security.txt` is intentionally not served. A security PGP keypair can be created with `quotes --maintenance pgp generate` after configuring `server.contact.security.email`. Its encrypted private key is stored separately from the public key, using a key derived from the database installation secret. Public-key publication is controlled by `web.security.publish_pgp_key`, and the public key is available at `/.well-known/pgp-key.asc` only while that setting is enabled and a keypair exists. The same conditional URL is shown in security.txt, the security page, and the security API response.

The coordinated disclosure page uses the built-in policy by default. Operators can replace its Markdown through `web.security.policy` and set the primary private reporting URL through `web.security.report_url` in the configuration files; the page sanitizes rendered content before serving it.

## Reporting a vulnerability

Start with `/.well-known/security.txt` served by the Quotes instance. Its `Contact` lines are ordered by preference: the configured private GitHub reporting URL, the instance's [`/server/contact?security_id=...`](</server/contact?security_id=...>) form when contact reporting is enabled, and the configured security email when present. Use the channel that matches your disclosure needs and follow the recipient's instructions. The `security_id` identifies the security report form; it is not an API token.

The human-readable [security page](/server/security) explains the reporting channels and links to the coordinated disclosure [policy](/server/security/policy). The report form also appears at [`/server/contact`](</server/contact>) when enabled. It validates and rate-limits reports, encrypts stored submissions, and does not add content to the quote collections. Do not include a secret or credential in a public issue.

## Well-known endpoints

Quotes implements these root discovery paths:

| Path | Availability | Purpose |
|---|---|---|
| `/.well-known/security.txt` | Always | RFC 9116 security contacts and policy |
| `/.well-known/llms.txt` | Always | Generated API and integration discovery |
| `/.well-known/pgp-key.asc` | Only when a keypair exists and `web.security.publish_pgp_key` is enabled | Public security encryption key |

`/llms.txt` is an alias for `/.well-known/llms.txt`. `/security.txt` is not an alias. Unknown `/.well-known/*` paths return `404`; the namespace is not a general static-file directory.


The operator Bearer token is generated into `secrets.yml` with mode `0600`; keep it private and never place it in browser storage, pages, or logs. Metrics bearer secrets are separate per-service tokens in the same protected file. Browser preference writes use a `csrf_token` double-submit cookie and require the same token in the form. Preference values are limited to the supported theme and language enums; sync imports are additive and ignore unsupported fields. Theme and language cookies are written only after that browser accepts preference cookies; otherwise the preference API reports the values as ignored with `consent_required: true`. The cookie-consent choice is device-local and is not included in preference exports. When `server.privacy.data.sold` is enabled, the CCPA opt-out form writes a same-site `ccpa_opt_out` cookie through a CSRF-protected route. Protected API routes do not use browser cookie credentials.
