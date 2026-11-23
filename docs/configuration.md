# Configuration

The server reads `server.yml`; use `--config` to select a file and `--address` or `--port` to override bind settings. When the configuration file is absent, the server writes defaults. It also creates a neighboring `secrets.yml` with an operator token and encryption key, using owner-only file permissions. Keep that file private and back it up securely. Never commit it.

## Server settings

| Setting | Default | Purpose |
|---|---|---|
| `server.address` | `0.0.0.0` | Bind address. |
| `server.port` | `8080` at runtime | HTTP listener port. |
| `server.fqdn` | empty | Public host name when configured. |
| `server.mode` | `production` | Application mode. |
| `server.rate_limit.enabled` | `true` | Enable per-client request limits. |
| `server.rate_limit.requests` | `120` | Requests allowed per client per window. |
| `server.rate_limit.window_seconds` | `60` | Rate-limit window length. |
| `server.logging.access_format` | `apache` | Access-log format setting. |
| `server.logging.level` | `info` | Log verbosity setting. |
| `server.metrics.enabled` | `false` | Metrics setting; metrics are not exposed by the current HTTP handler. |
| `server.metrics.endpoint` | `/metrics` | Configured metrics path. |
| `web-ui.theme` | `dark` | Default web theme. |
| `web-ui.logo`, `web-ui.favicon` | empty | Branding asset values. |
| `web-ui.notifications.enabled` | `true` | Announcement setting. |
| `web-ui.notifications.announcements` | empty list | Announcement text entries. |
| `web-robots.allow` | `/`, `/api` | Allowed robots paths. |
| `web-robots.deny` | `/debug` | Disallowed robots paths. |
| `web-security.cors` | `*` | CORS allow-origin value. |

`secrets.yml` stores `server.token` and `server.security.encryption_key`; it is generated automatically and must not be added to source control. `server.token` is required for MCP requests. The CLI accepts `--server`, `QUOTES_SERVER`, or a user configuration file. Quote content is embedded in the binary; there is no runtime database.
