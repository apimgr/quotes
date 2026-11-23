# Security

Quotes is a public, read-only quote service with no user accounts or quote submissions. The MCP route requires the operator Bearer token; other quote reads are public. Incoming query parameters are validated, request paths and bodies are bounded, lists are paginated, and per-client rate limits are enabled by default. CORS is open because the data is public.

The rate limiter uses the direct remote address and does not trust caller-supplied forwarding headers. Responses include content-type sniffing, framing, referrer, content-security, and permissions policy headers. `/server/healthz` negotiates HTML, plain text, or JSON using `Accept`; `/api/v1/server/healthz` and `/api/healthz` return JSON. `/.well-known/security.txt` and `/security.txt` publish the project security contact route.

Quote records are embedded in the executable. Only include material with documented redistribution rights and retain source attribution and license information. The current sample data is illustrative and must not be presented as an authoritative or licensed corpus.


The operator Bearer token is generated into `secrets.yml` with mode `0600`; keep it private and never place it in browser storage, pages, or logs. Browser preference writes use a `csrf_token` double-submit cookie and require the same token in the form. Preference values are limited to the supported theme and language enums; sync imports are additive and ignore unsupported fields. The preference API only writes browser preference cookies; protected API routes do not use browser cookie credentials.
