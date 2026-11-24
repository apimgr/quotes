# Integrations

Quotes exposes public REST and GraphQL interfaces, OpenAPI documentation, plain-text quote responses, and a Streamable HTTP MCP endpoint. REST and GraphQL quote reads need no credentials. MCP requests require an owner-linked Bearer token. Issue a primary owner token at `/server/preferences/tokens`, then create a separate linked token for each MCP client; do not use the operator token.

A web app manifest and service worker provide install metadata and basic offline caching. Security contact information is published at `/.well-known/security.txt` and `/security.txt`. The project does not currently enable OAuth/OIDC, webhooks, federation, native app associations, or external runtime data APIs. There is no configured official server URL; clients select the server they use.

## SMTP notifications

When a working SMTP relay is available, the server enables the notification queue and report delivery. `server.notifications.email.smtp.host` may be left empty for local relay detection; configured relays are tested at startup. `tls: auto` uses implicit TLS on port 465 and negotiates STARTTLS after EHLO on other ports. Authentication credentials are sent only after TLS. Without working SMTP, email-dependent features are disabled. Templates ship in the binary and operator overrides are loaded from `{config_dir}/template/email/`.

## Internal metrics

Metrics are internal endpoints and require independent bearer credentials stored in `secrets.yml`: `/server/metrics/prometheus`, `/server/metrics/grafana`, and `/server/metrics/loki`. The Prometheus service is also available at `/server/metrics` and, when enabled, `/metrics`. Do not expose these routes through a public reverse proxy; use the configured network firewall as an additional control.

## Overlay networks

Tor onion service and I2P eepsite support follow the project's overlay configuration. I2P is opt-in on host installs and defaults on in a container. Each overlay endpoint uses its dedicated listener and publishes its resolved URL through the supported server information surfaces when enabled.
