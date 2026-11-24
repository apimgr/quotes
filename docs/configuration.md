# Configuration

The server reads `{config_dir}/server.yml`; `--config` selects the configuration directory. `--data`, `--cache`, `--log`, `--backup`, and `--pid` override filesystem locations. `--address`, `--port`, and `--baseurl` override bind and URL-prefix settings. Directory overrides preserve absolute paths, resolve relative paths against the startup working directory, and reject control characters and `..` segments. When the configuration file is absent, the server writes defaults, including an unused 64xxx port (port 80 in an application container). It also creates a neighboring `secrets.yml` with an operator token and encryption key, using owner-only file permissions. Keep that file private and back it up securely. Never commit it.

## Server settings

| Setting | Default | Purpose |
|---|---|---|
| `server.address` | `[::]` | Bind address (dual-stack when supported). |
| `server.port` | Random unused 64xxx port; `80` in an application container | HTTP listener port; the selected default is persisted. |
| `server.fqdn` | empty | Public host name when configured. |
| `server.mode` | `production` | Application mode. |
| `server.update.branch` | `stable` | Release channel: `stable`, `beta`, or `devel` (`daily`, `rolling`, `dev`, and `development` normalize to `devel`). |
| `server.update.auto_install` | `false` | Whether scheduled update checks may install eligible updates. |
| `server.update.defer_days` | `0` | Stable/beta release age threshold in days (0–365); devel ignores this setting. |
| `server.update.cli_min_version` | empty | Minimum CLI version accepted by the server. |
| `server.tokens.max_linked` | `25` | Maximum active linked tokens per anonymous owner. |
| `server.tokens.cookie_max_age` | `1y` | Sliding owner-token cookie lifetime. |
| `server.tokens.idle_purge_days` | `0` | Purge idle owners only when they have no resources; 0 disables purging. |
| `server.baseurl` | empty | Optional URL path prefix; trusted proxy headers take precedence. |
| `server.pidfile` | `true` | Write a process PID file outside application containers. |
| `server.daemonize` | `false` | Detach manual starts; service-manager starts remain foreground. |
| `server.trusted_proxies.additional` | empty list | Additional trusted proxy IPs or CIDRs. |
| `server.trusted_proxies.client_ip_header` | `auto` | Client IP header selection for trusted proxies. |
| `server.rate_limit.enabled` | `true` | Enable per-client request limits. |
| `server.rate_limit.requests` | `120` | Requests allowed per client per window. |
| `server.rate_limit.window_seconds` | `60` | Rate-limit window length. |
| `server.logging.access_format` | `apache` | Access-log format setting. |
| `server.logging.level` | `info` | Log verbosity setting. |
| `server.metrics.enabled` | `true` | Enable internal Prometheus, Grafana, and Loki handlers. |
| `server.metrics.root.enabled` | `true` | Enable authenticated `/metrics` aliases. |
| `server.metrics.auth.allow_unauthenticated` | `false` | Skip metrics tokens only on a firewalled internal network. |
| `server.metrics.include_system` | `true` | Collect CPU, memory, and disk metrics. |
| `server.metrics.include_runtime` | `true` | Collect Go runtime metrics. |
| `server.geoip.enabled` | `true` | Enable country lookup and configured allow/deny policy when the local database is available. |
| `server.scheduler.timezone` | `America/New_York` | IANA timezone for the persistent built-in scheduler. |
| `server.scheduler.catch_up_window` | `1h` | Maximum age for scheduled work to run after a restart. |
| `web-ui.theme` | `dark` | Default web theme. |
| `web-ui.logo`, `web-ui.favicon` | empty | Branding asset values. |
| `web-ui.notifications.enabled` | `true` | Announcement setting. |
| `web-ui.notifications.announcements` | empty list | Announcement text entries. |
| `web-robots.allow` | `/`, `/api` | Allowed robots paths. |
| `web-robots.deny` | `/debug` | Disallowed robots paths. |
| `web-security.cors` | `*` | CORS allow-origin value. |

`secrets.yml` stores `server.token`, `server.security.encryption_key`, and the independent `server.metrics.auth.tokens` credentials; it is generated automatically and must not be added to source control. Metrics tokens are generated individually for Prometheus, Grafana, and Loki. `server.token` is reserved for protected maintenance operations; MCP requires an owner-linked token. SQLite stores anonymous owner-token and scheduler state in `<data_dir>/db/server.db`; quote collections remain embedded in the binary. The CLI accepts `--server`, `QUOTES_SERVER`, or a user configuration file.

## Email and reports

SMTP options live below `server.notifications.email.smtp` and accept `SMTP_HOST`, `SMTP_PORT`, `SMTP_USERNAME`, `SMTP_PASSWORD`, `SMTP_TLS`, `SMTP_FROM_NAME`, and `SMTP_FROM_EMAIL` environment overrides. TLS defaults to `auto`; port 465 uses implicit TLS and other ports use EHLO followed by STARTTLS when offered. Email-dependent features stay disabled when SMTP is unavailable. Default email templates are embedded; operator overrides belong in `{config_dir}/template/email/`.

## Scheduled tasks

The scheduler is always active and persists task state in SQLite. Configure IANA timezone and catch-up window under `server.scheduler`; individual tasks use `schedule`, `enabled`, `retry_on_fail`, and `retry_delay`. The server does not use cron or an external scheduler. `quotes --maintenance scheduler list|run|enable|disable` is the task-control interface.

## Storage locations

The server creates its directories at startup. On Linux, privileged runs use `/etc/apimgr/quotes`, `/var/lib/apimgr/quotes`, `/var/cache/apimgr/quotes`, and `/var/log/apimgr/quotes`; user runs follow XDG configuration, data, cache, and log locations. macOS uses `~/Library/Application Support/apimgr/quotes`, `~/Library/Caches/apimgr/quotes`, and `~/Library/Logs/apimgr/quotes`. Windows uses `%APPDATA%` for configuration and `%LOCALAPPDATA%` for data, cache, and logs. BSD uses `/usr/local/etc`, `/var/db`, `/var/cache`, and `/var/log`.

Inside a container, configuration is under `/config/quotes`, data and cache are under `/data/quotes`, logs are under `/data/log/quotes`, and backups are under `/data/backups/quotes`. Mount `/config` and `/data` as persistent volumes.

The mode is selected from `--mode`, `MODE`, `server.mode`, then `production`. Supported values are `production`, `development`, and `debug`; `prod`, `dev`, and `devel` are aliases. `--maintenance mode <mode>` persists the server mode after authorization. Debug endpoints require the separate `--debug` flag or a truthy `DEBUG` environment value. The default theme is dark; selecting a theme submits a same-origin form protected by the `csrf_token` double-submit cookie.

Service commands use the detected host service manager: `quotes --service start|stop|restart|reload`, `quotes --service --install`, `quotes --service --disable`, and `quotes --service --uninstall`. Install enables and starts the service; disable stops it and removes it from automatic startup while leaving its service files and data in place. Uninstall requires the operator to type `yes` before removing service files, the managed account, and application directories. A regular user gets a systemd user unit or macOS LaunchAgent; system installations run with a dedicated service account after privileged listener binding. `--update branch <stable|beta|devel>` changes the persisted update channel; `--update check` checks GitHub Releases and `--update yes` verifies the release checksum before replacing the executable.
