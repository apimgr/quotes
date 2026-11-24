# Configuration

The server reads `server.yml`; use `--config` to select a file and `--address` or `--port` to override bind settings. When the configuration file is absent, the server writes defaults. It also creates a neighboring `secrets.yml` with an operator token and encryption key, using owner-only file permissions. Keep that file private and back it up securely. Never commit it.

## Server settings

| Setting | Default | Purpose |
|---|---|---|
| `server.address` | `0.0.0.0` | Bind address. |
| `server.port` | `8080` at runtime | HTTP listener port. |
| `server.fqdn` | empty | Public host name when configured. |
| `server.mode` | `production` | Application mode. |
| `server.update.branch` | `stable` | Release channel: `stable`, `beta`, or `devel` (`daily`, `rolling`, `dev`, and `development` normalize to `devel`). |
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

`secrets.yml` stores `server.token` and `server.security.encryption_key`; it is generated automatically and must not be added to source control. `server.token` is required for MCP requests and protected maintenance operations. SQLite stores anonymous owner-token state in `<data_dir>/db/server.db`; quote collections remain embedded in the binary. The CLI accepts `--server`, `QUOTES_SERVER`, or a user configuration file.

## Storage locations

The server creates its directories at startup. On Linux, privileged runs use `/etc/apimgr/quotes`, `/var/lib/apimgr/quotes`, `/var/cache/apimgr/quotes`, and `/var/log/apimgr/quotes`; user runs follow XDG configuration, data, cache, and log locations. macOS uses `~/Library/Application Support/apimgr/quotes`, `~/Library/Caches/apimgr/quotes`, and `~/Library/Logs/apimgr/quotes`. Windows uses `%APPDATA%` for configuration and `%LOCALAPPDATA%` for data, cache, and logs. BSD uses `/usr/local/etc`, `/var/db`, `/var/cache`, and `/var/log`.

Inside a container, configuration is under `/config/quotes`, data and cache are under `/data/quotes`, logs are under `/data/log/quotes`, and backups are under `/data/backups/quotes`. Mount `/config` and `/data` as persistent volumes.

The mode is selected from `--mode`, `MODE`, `server.mode`, then `production`. Supported values are `production`, `development`, and `debug`; `prod`, `dev`, and `devel` are aliases. `--maintenance mode <mode>` persists the server mode after authorization. Debug endpoints require the separate `--debug` flag or a truthy `DEBUG` environment value. The default theme is dark; selecting a theme submits a same-origin form protected by the `csrf_token` double-submit cookie.

Service commands use the detected host service manager: `quotes --service start|stop|restart|reload`, `quotes --service --install`, `quotes --service --disable`, and `quotes --service --uninstall`. Install enables and starts the service; disable stops it and removes it from automatic startup while leaving its service files and data in place. Uninstall requires the operator to type `yes` before removing service files, the managed account, and application directories. A regular user gets a systemd user unit or macOS LaunchAgent; system installations run with a dedicated service account after privileged listener binding. `--update branch <stable|beta|devel>` changes the persisted update channel; `--update check` checks GitHub Releases and `--update yes` verifies the release checksum before replacing the executable.
