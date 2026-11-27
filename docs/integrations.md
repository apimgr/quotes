# Integrations

Quotes exposes public REST and GraphQL interfaces, OpenAPI documentation, plain-text quote responses, and a Streamable HTTP MCP endpoint. REST and GraphQL quote reads need no credentials. MCP requests require a valid Bearer token. The operator may use `server.token` for their own MCP tooling; it has full operator access and exposes operator-only tools, so do not share it or give it to a third-party client. For each standing MCP client, issue a separate linked token from `/server/preferences/tokens` with the minimum required scope. Owner and linked tokens expose only the tools their scope may call.

A web app manifest and service worker provide install metadata and basic offline caching. Security contact information is published at `/.well-known/security.txt`; AI clients can discover the public API at `/.well-known/llms.txt` and `/llms.txt`. `/api/autodiscover` returns the standard JSON envelope with the running server/API versions, request and configured public origins, trusted proxy origins learned by the server, and checksum-verified CLI release metadata. Exposed Tor and I2P addresses appear only in their own fields when the corresponding overlay and `expose` setting are enabled. The endpoint does not publish server settings or configuration schemas.

CLI setup may validate a supplied Bearer token with the machine-only `GET /api/v1/server/auth/check` endpoint. A valid response returns `data.kind` (`operator` or `owner`); an invalid token returns 401. The endpoint does not echo the credential.

Operators provision CLI releases in `{data_dir}/cli/` (or `cli.binary_download.dir`) using the release binaries, `version.txt`, and `sha256.txt`. The server recomputes each checksum and advertises only matching files. `GET /cli/binaries/{filename}` streams an advertised binary; downloads are public by default and can require a Bearer token with `cli.binary_download.require_auth: true`. See the [API reference](api.md) and [CLI reference](cli.md) for response and update behavior. The discovery documents use the request's resolved public origin and do not advertise internal metrics. The project does not currently enable OAuth/OIDC, webhooks, federation, native app associations, or external runtime data APIs. There is no configured official server URL; clients select the server they use.

Search engines can use the dynamically generated `/robots.txt` and `/sitemap.xml`. The sitemap includes public HTML pages and quote detail pages allowed by the robots configuration; it excludes API paths, contact forms, and private token pages. The canonical security contact file is `/.well-known/security.txt`; a generated PGP public key is served at `/.well-known/pgp-key.asc` only when a keypair exists and `web.security.publish_pgp_key` is enabled. `/security.txt` is not an alias.

Feed readers can subscribe to `/feed.xml` (Atom 1.0) or `/rss.xml` (RSS 2.0). Both feeds contain the same latest public quote records and support conditional caching.

## SMTP notifications

When a working SMTP relay is available, the server enables the notification queue and report delivery. `server.notifications.email.smtp.host` may be left empty for local relay detection; configured relays are tested at startup. `tls: auto` uses implicit TLS on port 465 and negotiates STARTTLS after EHLO on other ports. Authentication credentials are sent only after TLS. Without working SMTP, email-dependent features are disabled. Templates ship in the binary and operator overrides are loaded from `{config_dir}/template/email/`.

## Internal metrics

Metrics are internal endpoints and require independent bearer credentials stored in `secrets.yml`: `/server/metrics/prometheus`, `/server/metrics/grafana`, and `/server/metrics/loki`. The Prometheus service is also available at `/server/metrics` and, when enabled, `/metrics`. Do not expose these routes through a public reverse proxy; use the configured network firewall as an additional control.

## Overlay networks

Quotes can serve its HTTP interface as a Tor v3 onion service and as an I2P eepsite. Tor starts automatically when the `tor` executable is available; there is no enable/disable switch. I2P is disabled by default on bare-metal, VM, and system-container installs, and enabled by default in an application container. On Linux, the application container includes `i2pd`; otherwise I2P uses the configured external SAM endpoint (default `127.0.0.1:7656`). If no provider is available, that overlay remains unavailable.

Each active overlay uses a dedicated loopback backend listener, separate from the clearnet listener. Both overlays serve plain HTTP only: their onion or eepsite address is the endpoint identity, and Quotes does not add HTTPS, HSTS, or an HTTP-to-HTTPS redirect there. Tor and I2P addresses are private by default and are omitted from autodiscovery. Set the relevant `server.tor.expose` or `server.i2p.expose` option to publish an address in the supported discovery and server-information responses. See [configuration](configuration.md) for the full overlay settings and defaults.
