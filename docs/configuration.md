# Configuration

Configuration is file-backed. The merged files are the only source of operator-editable settings. Database records hold application state, not configuration. There is no web administration panel or runtime configuration mutation API.

## Files and precedence

The configuration directory contains `server.yml`, `yaml.d/` area files, and `secrets.yml`. The merge order is built-in defaults, `server.yml`, `yaml.d/*.yml` and `yaml.d/*.yaml` in byte-wise filename order, then `secrets.yml`. Later values override earlier values. Maps merge recursively; lists and scalars replace the earlier value. CLI flags and environment variables have the precedence specified for each setting.

`server.yml` holds node settings: `server.port`, `server.address`, `server.fqdn`, `server.mode`, `server.api_version`, `server.baseurl`, `server.pidfile`, `server.daemonize`, `server.debug.*`, `server.database`, `server.cache`, `server.ssl`, and `server.trusted_proxies`. Other sections belong in area files:

| File | Sections |
|---|---|
| `yaml.d/branding.yml` | `server.branding`, `server.seo`, `web.ui` |
| `yaml.d/web.yml` | `web.*` except `web.ui`; `server.cors`, `server.csrf`, `server.url_detection`, `server.compression`, `server.healthz`, `server.i18n` |
| `yaml.d/pages.yml` | `server.pages` (report form, antispam, help pages) |
| `yaml.d/contact.yml` | `server.contact` |
| `yaml.d/privacy.yml` | `server.privacy`, `server.compliance`, `server.tracking` |
| `yaml.d/security.yml` | `server.security`, `server.rate_limit`, `server.geoip`, `server.limits`, `server.tokens`, `data.*` |
| `yaml.d/logs.yml` | `server.logs` |
| `yaml.d/metrics.yml` | `server.metrics` |
| `yaml.d/notifications.yml` | `server.notifications` |
| `yaml.d/tor.yml` | `server.tor`, top-level `tor` |
| `yaml.d/i2p.yml` | `server.i2p` |
| `yaml.d/maintenance.yml` | `server.maintenance`, `server.scheduler`, `server.update`, `server.auto_update`, `server.backup`, `server.restore`, top-level `cli` |
| `yaml.d/app.yml` | Other project-defined sections |

The merge accepts these sections in any YAML file, but the table is the canonical location. Existing area sections in `server.yml` are honored and automatically migrated on startup. Before moving values, migration creates a `server.yml.bak`; it writes area files atomically and then rewrites `server.yml`. If the configuration directory is read-only, migration is skipped with a warning and the merged settings remain usable.

Startup also migrates the legacy setting names supported by the current release before splitting area files. Canonical keys win when both names are present. Explicitly retired keys are removed after the backup is made; unlisted keys in defined sections are retained for the normal unknown-key warning. This includes the old logging, scheduler, web, database, branding, SEO, page content, contact, admin, SSL, and security names. The HTTPS port is folded into a single numeric `server.port` as a dual-port value; an already dual port is preserved. Legacy scheduler task rows retain only `enabled`, `schedule`, `retry_on_fail`, and `retry_delay`.

Legacy installation and previous-key values are first moved from `server.yml` and drop-ins into owner-only `secrets.yml`. On startup, they are imported into `app_secrets` only when the canonical database row is missing; the legacy file keys are removed only after the database import succeeds. A database failure leaves the legacy values available for the next startup.

The web area file can configure the generated AI discovery document and conditional PGP public-key publication:

```yaml
web:
  llms:
    include_endpoints: true
    include_schemas: false
    custom_sections: []
  security:
    report_url: https://github.com/apimgr/quotes/security/advisories/new
    publish_pgp_key: false
    expires: ""
    policy: ""
```

`/.well-known/llms.txt` and `/llms.txt` are always available. An operator may add up to ten Markdown sections; unsafe links are reduced to plain link text and HTML tags are escaped. `web.security.report_url` sets the primary private vulnerability-reporting URL, `web.security.policy` replaces the built-in coordinated disclosure policy with operator-supplied Markdown, and `web.security.expires` sets the RFC 9116 expiry (RFC 3339 or date; blank calculates one year ahead). `web.security.publish_pgp_key` only publishes an existing security keypair. `quotes --maintenance pgp generate` creates a keypair after operator authorization and enables publication; `quotes --maintenance pgp export public [path]` exports its public key; `quotes --maintenance pgp delete` requires a typed confirmation and disables publication.

Only `*.yml` and `*.yaml` regular files are loaded (internal symlinks are allowed). Hidden files, backups, other extensions, and subdirectories are ignored. Files `server.yml`, `secrets.yml`, and every drop-in use mode `0600`; `yaml.d/` uses `0700`. The server corrects modes at startup. A drop-in that remains group/world writable is ignored and logged. The watcher reloads `server.yml`, `yaml.d/`, `secrets.yml`, and `assets/`; invalid individual values fall back to defaults, and an unparseable file keeps its last good copy.

## Secrets and validation

`secrets.yml` is the at-rest home for credentials, including database credentials, operator and encryption keys, SMTP passwords, and metrics credentials. A hand-written secret in `server.yml` or an area file is moved to `secrets.yml` during startup migration; avoid placing new credentials in those files. The server creates missing operator and encryption keys securely. Keep `secrets.yml` private, include it in protected backups, and never commit it. Secret values are never logged; the newly generated operator token is disclosed once at creation.

Every setting is validated. An invalid value is replaced by its documented default with one warning naming the file, line, and key; configuration errors do not stop startup. A YAML file that cannot be parsed is ignored, retaining its last good copy or built-in defaults. Duplicate scalar keys and explicit `null` values are invalid. Unknown keys in defined sections are warned about and ignored.

`secrets.yml` follows the same per-key fallback rule. A malformed secrets document retains the last good secret values for that configuration path (or defaults when no last-good state exists), and warnings identify only the file and key; secret contents are never included in diagnostics.

YAML comments go above the setting, never inline. For example:

```yaml
server:
  # Listen on all supported interfaces.
  address: "[::]"
```

Booleans accept case-insensitive true/false forms including `1`/`0`, `yes`/`no`, `on`/`off`, `enable`/`disable`, `y`/`n`, and the documented aliases `enabled`/`disabled`. Empty values use the setting default; invalid values warn and use the default.

## Node settings and database

```yaml
server:
  address: "[::]"
  port: 64000
  fqdn: ""
  mode: production
  baseurl: /
  pidfile: true
  daemonize: false
  database:
    driver: sqlite
    url: "{data_dir}/db/server.db"
  cache:
    type: memory
  ssl:
    enabled: false
```

The server chooses and persists an available port in the 64000–64999 range when no port is configured or supplied at runtime; application containers default to port 80. At every startup, `--port` overrides `PORT`, which overrides the saved `server.port`; `--address` overrides `LISTEN`, which overrides the saved `server.address`. These runtime listener values are not written to `server.yml`. `server.address` is the bind address; `server.fqdn` is the public hostname. `server.baseurl` is the URL path prefix, default `/`. Trusted reverse-proxy headers are honored only from trusted peers. Configure additional proxy IPs/CIDRs under `server.trusted_proxies.additional`; do not trust forwarded headers from arbitrary clients.

SQLite is local and creates its database under the data directory. For remote Turso/libsql, configure `server.database.driver: libsql` and the connection URL; put its auth token in `secrets.yml` as `server.database.token`. Cache supports the configured local or remote backend. Database and cache connection environment variables seed first-run configuration only; later changes belong in files.

TLS listener settings are `server.ssl.enabled`, `server.ssl.cert`, `server.ssl.key`, and `server.ssl.min_version`; certificate and key paths refer to operator-managed files. ACME settings live under `server.ssl.letsencrypt`; the ACME client reads updated values at its next issuance or renewal. Request limits and response compression are configured in `yaml.d/security.yml` and `yaml.d/web.yml`, respectively. The settings index below lists every supported key and built-in default. Proxy trust applies only to explicitly trusted peers; changes to listener settings require a restart, while reloadable settings are watched.

## Environment and command-line overrides

Directory options are available as flags and environment variables. A flag wins over its environment variable, which wins over the platform default. Directory paths preserve absolute paths; relative paths resolve against the startup working directory. Paths containing control characters or a `..` segment are rejected.

| Environment | Flag | Purpose |
|---|---|---|
| `CONFIG_DIR` | `--config` | Configuration directory |
| `DATA_DIR` | `--data` | Data directory |
| `LOG_DIR` | `--log` | Log directory |
| `DATABASE_DIR` | — | SQLite database directory |
| `BACKUP_DIR` | `--backup` | Backup directory |

`APPLICATION_NAME` and `APPLICATION_TAGLINE` seed `server.branding` once on first run. `DATABASE_DRIVER`, `DATABASE_URL`, and `CACHE_URL` seed the initial database and cache connections only; later changes belong in files. Credential-bearing connection URLs are stored in `secrets.yml`. `PORT` and `LISTEN` are runtime variables read on every startup, with their corresponding CLI flags taking precedence; neither variable is first-run-only.

Other filesystem overrides include `--cache`, `--log`, and `--pid`. Listener and URL overrides include `--address`, `--port`, and `--baseurl`. `DOMAIN` overrides hostname discovery. `MODE` selects `production`, `development`, or `debug`; `DEBUG` independently enables debug behavior. SMTP environment variables are `SMTP_HOST`, `SMTP_PORT`, `SMTP_USERNAME`, `SMTP_PASSWORD`, `SMTP_FROM_NAME`, `SMTP_FROM_EMAIL`, and `SMTP_TLS`. Flags take precedence over environment variables, and runtime variables are read on every startup; first-run variables seed configuration only.

Debug endpoints exist only when `--debug` or a truthy `DEBUG` value enables them. `MODE=debug` enables them by default, while an explicit `DEBUG` value or `--debug` flag can turn them off. Every `/debug/*` request still requires the operator token as a Bearer credential. Diagnostics live only at root paths such as `/debug/pprof` and `/debug/memory`; they have no versioned API aliases. The pprof command-line endpoint and expvar output redact process arguments.

## Area configuration

The following examples show the area-file shape; omit settings that should use their defaults. The complete key, type, and built-in default index follows the examples.

### Languages

`server.i18n` is stored in `yaml.d/web.yml`. `default_language` selects the server's language when a request has no supported `?lang=`, language cookie, or `Accept-Language` preference. `available_languages` is the allowlist used for request negotiation and the language selector; it must include the six bundled languages. Unsupported languages and missing translation keys fall back to English.

```yaml
server:
  i18n:
    default_language: en
    available_languages: [en, es, fr, ar, de, it]
```

### Logging

Log files are created in the selected log directory. The configurable file kinds are `access.log`, `server.log`, `error.log`, `audit.log`, `security.log`, and `debug.log`; the logging subsystem also writes `app.log` and `auth.log` using their built-in formats. Audit records use JSON Lines and a hash chain. Files use mode `0640`, and rotated files follow each kind's retention setting.

```yaml
server:
  logs:
    level: info
    defaults:
      custom: ""
      rotate: weekly,50MB
      keep: none
    access:
      enabled: true
      format: apache
      rotate: monthly
    server:
      enabled: true
      format: text
    error:
      enabled: true
      format: text
    audit:
      enabled: true
      format: json
      rotate: daily
    security:
      enabled: true
      format: fail2ban
    debug:
      enabled: false
      format: text
```

Formats depend on the log kind: general logs support `text` or `json`, access logs support `apache`, `nginx`, `json`, or a custom format, audit logs require `json`, and security logs use `fail2ban`. Rotation supports daily, weekly, monthly, or size limits (the default is weekly or 50 MB, whichever comes first). Retention accepts `none`, a count, or a duration in days, weeks, or months, or `forever`. Secret values and raw owner tokens are never written to logs. Successful health probes are omitted from the access log; failed probes remain visible.

### Analytics tracking

Analytics is disabled by default. Configure a supported provider under `server.tracking`; the server validates its ID and HTTPS URL requirements at load time and disables an invalid provider. Browser scripts load only after analytics consent. They are omitted on onion and I2P requests, and the same-origin `/server/tracking.js` endpoint returns `404` without consent or when tracking is disabled. CSP adds only the configured provider's sources.

```yaml
server:
  tracking:
    # google, matomo, piwik, owa, fathom, plausible, umami, simple, cloudflare
    type: ""
    id: ""
    url: ""
```

Self-hosted Matomo, Piwik, OWA, and Umami require an HTTPS `url`. Fathom and Plausible accept an optional HTTPS URL for self-hosting. Google, Simple Analytics, and Cloudflare do not use `url`. Use the provider-issued site or measurement identifier for `id`; invalid identifiers disable the provider during configuration validation.

```yaml
# yaml.d/security.yml
server:
  rate_limit:
    enabled: true
  tokens:
    max_linked: 25
```

```yaml
# yaml.d/contact.yml
server:
  contact:
    general:
      email: "reports@<your-domain>"
```

```yaml
# yaml.d/notifications.yml
server:
  notifications:
    webui:
      position: top-right
      duration: 3
      warning_duration: 5
    email:
      smtp:
        host: "<your-smtp-host>"
        port: 587
        tls: auto
```

Use `yaml.d/pages.yml` for report page behavior and antispam, `yaml.d/branding.yml` for branding and SEO, `yaml.d/privacy.yml` for privacy and compliance disclosures, `yaml.d/metrics.yml` for metrics, `yaml.d/tor.yml` and `yaml.d/i2p.yml` for overlay networking, and `yaml.d/maintenance.yml` for scheduling, backups, restores, and updates. Contact reports do not create accounts or accept quote contributions. Owner tokens are anonymous credentials; they are not user accounts.

Email behavior and templates, scheduler tasks, GeoIP filtering, metrics endpoints, backup and restore safeguards, and update channels are covered in their respective sections of this guide and the [installation](installation.md), [security](security.md), and [integrations](integrations.md) guides. Backups include configuration area files and `secrets.yml` when present. Treat backup files as sensitive.

## Platform paths and operations

Use the platform configuration, data, cache, log, and backup directories selected by the application. On Linux, privileged configuration is under `/etc/apimgr/quotes/` and user configuration under `~/.config/apimgr/quotes/`; on macOS, use `/Library/Application Support/apimgr/quotes/` or `~/Library/Application Support/apimgr/quotes/`; on Windows, use `%ProgramData%\apimgr\quotes\` or `%AppData%\apimgr\quotes\`; on BSD, use `/usr/local/etc/apimgr/quotes/` or `~/.config/apimgr/quotes/`. Container deployments use persistent mounts at `/config/quotes/` and `/data/quotes/`; do not assume container paths exist on a native host.

The built-in scheduler persists state in the configured database; do not use cron or an external scheduler. Maintenance commands require operator authorization where specified. The server has no web panel for changing configuration: edit the files and let the watcher reload supported settings. Listener address, port, daemonization, and PID-file changes require a restart.

## Related guides

- [Installation](installation.md) — platform paths, services, and container setup.
- [Security](security.md) — reporting, public endpoints, and credential handling.
- [Integrations](integrations.md) — API discovery and external services.

## Complete settings index

This index covers every configuration key, its YAML value type, and its built-in default. `empty` means unset; secret values are intentionally omitted. Structured map and list defaults are summarized by item count. Platform-dependent paths are described by name rather than expanded for a particular host. The legacy `web-ui`, `web-robots`, and `web-security` prefixes are shown for migration reference; use the canonical `server.branding` and `web.*` sections listed above for new configuration. For configuration examples and behavior, see the sections above.

| Setting | Type | Built-in default |
|---|---|---|
| `server.port` | `string` | `empty` |
| `server.fqdn` | `string` | `empty` |
| `server.address` | `string` | `::` |
| `server.baseurl` | `string` | `empty` |
| `server.pidfile` | `bool` | `true` |
| `server.daemonize` | `bool` | `false` |
| `server.mode` | `string` | `production` |
| `server.update.branch` | `string` | `stable` |
| `server.update.auto_install` | `bool` | `false` |
| `server.update.defer_days` | `int` | `0` |
| `server.update.cli_min_version` | `string` | `empty` |
| `server.tokens.max_linked` | `int` | `25` |
| `server.tokens.cookie_max_age` | `string` | `1y` |
| `server.tokens.idle_purge_days` | `int` | `0` |
| `server.trusted_proxies.additional` | `[]string` | `empty` |
| `server.trusted_proxies.client_ip_header` | `string` | `auto` |
| `server.backup.encryption.enabled` | `bool` | `false` |
| `server.backup.encryption_password` | `string` | `empty` |
| `server.backup.retention.max_backups` | `int` | `1` |
| `server.backup.retention.keep_weekly` | `int` | `0` |
| `server.backup.retention.keep_monthly` | `int` | `0` |
| `server.backup.retention.keep_yearly` | `int` | `0` |
| `server.backup.retention.max_total_size` | `string` | `10%` |
| `server.restore.max_archive_size` | `string` | `1G` |
| `server.restore.max_expanded_size` | `string` | `4G` |
| `server.restore.max_files` | `int` | `10000` |
| `server.restore.max_file_size` | `string` | `512M` |
| `server.restore.max_ratio` | `int` | `100` |
| `server.compliance.gdpr` | `bool` | `false` |
| `server.compliance.ccpa` | `bool` | `false` |
| `server.compliance.hipaa` | `bool` | `false` |
| `server.compliance.soc2` | `bool` | `false` |
| `server.compliance.pci_dss` | `bool` | `false` |
| `server.compliance.iso27001` | `bool` | `false` |
| `server.compliance.fedramp` | `bool` | `false` |
| `server.compliance.lgpd` | `bool` | `false` |
| `server.compliance.pipeda` | `bool` | `false` |
| `server.compliance.appi` | `bool` | `false` |
| `server.compliance.pdpa` | `bool` | `false` |
| `server.privacy.data.sold` | `bool` | `false` |
| `server.privacy.data.stored_on_server` | `bool` | `true` |
| `server.privacy.data.sharing` | `[]config.PrivacySharingConfig` | `1 item(s)` |
| `server.privacy.retention.export_available` | `bool` | `false` |
| `server.privacy.retention.deletion_available` | `bool` | `false` |
| `server.privacy.consent.show_until_acknowledged` | `bool` | `true` |
| `server.privacy.consent.message` | `string` | `built-in text` |
| `server.privacy.consent.message_if_sold` | `string` | `built-in text` |
| `server.privacy.consent.policy.text` | `string` | `Privacy Policy` |
| `server.privacy.consent.policy.url` | `string` | `/server/privacy` |
| `server.privacy.consent.policy.version` | `string` | `1` |
| `server.privacy.consent.buttons.decline` | `string` | `Decline` |
| `server.privacy.consent.buttons.accept` | `string` | `Allow preferences` |
| `server.privacy.consent.show_preferences` | `bool` | `true` |
| `server.privacy.consent.preferences_text` | `string` | `Manage Preferences` |
| `server.privacy.cookies.essential.enabled` | `bool` | `true` |
| `server.privacy.cookies.essential.description` | `string` | `built-in text` |
| `server.privacy.cookies.preferences.enabled` | `bool` | `true` |
| `server.privacy.cookies.preferences.description` | `string` | `Remember this browser's theme and language choices after you accept the preferences category.` |
| `server.privacy.cookies.analytics.enabled` | `bool` | `true` |
| `server.privacy.cookies.analytics.description` | `string` | `No analytics cookies are set unless an analytics provider is configured and separately accepted.` |
| `server.privacy.cookies.analytics.description_suffix_not_sold` | `string` | `Analytics data is not sold.` |
| `server.privacy.cookies.analytics.description_suffix_sold` | `string` | `Analytics data may be shared or sold as described in this policy.` |
| `server.privacy.third_party.services` | `[]config.PrivacyThirdPartyServiceConfig` | `empty` |
| `server.tracking.type` | `string` | `empty` |
| `server.tracking.id` | `string` | `empty` |
| `server.tracking.url` | `string` | `empty` |
| `server.branding.title` | `string` | `empty` |
| `server.branding.tagline` | `string` | `empty` |
| `server.branding.description` | `string` | `empty` |
| `server.branding.content.header` | `*config.PageContentConfig` | `unset` |
| `server.branding.content.about` | `*config.PageContentConfig` | `unset` |
| `server.branding.content.privacy` | `*config.PageContentConfig` | `unset` |
| `server.branding.content.footer` | `*config.PageContentConfig` | `unset` |
| `server.branding.content.terms` | `*config.PageContentConfig` | `unset` |
| `server.pages.help.content` | `string` | `built-in text` |
| `server.pages.help.mcp` | `string` | `built-in text` |
| `server.pages.contact.captcha_secret_key` | `string` | `empty` |
| `server.ssl.cert` | `string` | `empty` |
| `server.ssl.key` | `string` | `empty` |
| `server.ssl.letsencrypt.dns_credentials` | `string` | `empty` |
| `server.metrics.enabled` | `bool` | `true` |
| `server.metrics.root.enabled` | `bool` | `true` |
| `server.metrics.auth.allow_unauthenticated` | `bool` | `false` |
| `server.metrics.auth.tokens` | `map[string]string` | `empty` |
| `server.metrics.include_system` | `bool` | `true` |
| `server.metrics.include_runtime` | `bool` | `true` |
| `server.metrics.loki.max_entries` | `int` | `1000` |
| `server.metrics.loki.max_age` | `string` | `1h` |
| `server.metrics.duration_buckets` | `[]float64` | `empty` |
| `server.metrics.size_buckets` | `[]float64` | `empty` |
| `server.notifications.webui.position` | `string` | `top-right` |
| `server.notifications.webui.duration` | `int` | `3` |
| `server.notifications.webui.warning_duration` | `int` | `5` |
| `server.notifications.email.smtp.host` | `string` | `empty` |
| `server.notifications.email.smtp.port` | `string` | `587` |
| `server.notifications.email.smtp.username` | `string` | `empty` |
| `server.notifications.email.smtp.password` | `string` | `empty` |
| `server.notifications.email.smtp.tls` | `string` | `auto` |
| `server.notifications.email.from.name` | `string` | `empty` |
| `server.notifications.email.from.email` | `string` | `empty` |
| `server.notifications.email.reply_to` | `string` | `empty` |
| `server.contact.admin.email` | `string` | `empty` |
| `server.contact.admin.name` | `string` | `empty` |
| `server.contact.admin.webhooks` | `map[string]string` | `empty` |
| `server.contact.security.email` | `string` | `empty` |
| `server.contact.security.name` | `string` | `empty` |
| `server.contact.security.webhooks` | `map[string]string` | `empty` |
| `server.contact.abuse.email` | `string` | `empty` |
| `server.contact.abuse.name` | `string` | `empty` |
| `server.contact.abuse.webhooks` | `map[string]string` | `empty` |
| `server.contact.general.email` | `string` | `empty` |
| `server.contact.general.name` | `string` | `empty` |
| `server.contact.general.webhooks` | `map[string]string` | `empty` |
| `server.contact.dpo.email` | `string` | `empty` |
| `server.contact.dpo.name` | `string` | `empty` |
| `server.contact.dpo.webhooks` | `map[string]string` | `empty` |
| `server.i18n.default_language` | `string` | `en` |
| `server.i18n.available_languages` | `[]string` | `6 item(s)` |
| `server.healthz.root.enabled` | `bool` | `false` |
| `server.geoip.enabled` | `bool` | `true` |
| `server.geoip.deny_countries` | `[]string` | `empty` |
| `server.geoip.allow_countries` | `[]string` | `empty` |
| `server.geoip.presets` | `map[string][]string` | `empty` |
| `server.tor.binary` | `string` | `empty` |
| `server.tor.expose` | `bool` | `false` |
| `server.tor.use_network` | `bool` | `false` |
| `server.tor.max_circuits` | `int` | `32` |
| `server.tor.startup_timeout` | `int` | `60` |
| `server.tor.circuit_timeout` | `int` | `60` |
| `server.tor.bootstrap_timeout` | `int` | `180` |
| `server.tor.safe_logging` | `bool` | `true` |
| `server.tor.max_streams_per_circuit` | `int` | `100` |
| `server.tor.close_circuit_on_stream_limit` | `bool` | `true` |
| `server.tor.bandwidth_rate` | `string` | `1 MB` |
| `server.tor.bandwidth_burst` | `string` | `2 MB` |
| `server.tor.max_monthly_bandwidth` | `string` | `100 GB` |
| `server.tor.num_intro_points` | `int` | `3` |
| `server.tor.virtual_port` | `int` | `80` |
| `server.i2p.enabled` | `bool` | `true` |
| `server.i2p.expose` | `bool` | `false` |
| `server.i2p.binary` | `string` | `empty` |
| `server.i2p.sam_address` | `string` | `(loopback SAM service)` |
| `server.i2p.virtual_port` | `int` | `80` |
| `server.i2p.inbound_length` | `int` | `3` |
| `server.i2p.outbound_length` | `int` | `3` |
| `server.i2p.inbound_quantity` | `int` | `5` |
| `server.i2p.outbound_quantity` | `int` | `5` |
| `server.i2p.signature_type` | `int` | `7` |
| `server.i2p.bootstrap_timeout` | `int` | `300` |
| `server.scheduler.timezone` | `string` | `America/New_York` |
| `server.scheduler.catch_up_window` | `string` | `1h` |
| `server.scheduler.tasks` | `map[string]config.SchedulerTaskConfig` | `15 item(s)` |
| `server.logs.level` | `string` | `empty` |
| `server.logs.defaults.custom` | `string` | `empty` |
| `server.logs.defaults.rotate` | `string` | `weekly,50MB` |
| `server.logs.defaults.keep` | `string` | `none` |
| `server.logs.access.enabled` | `bool` | `true` |
| `server.logs.access.format` | `string` | `apache` |
| `server.logs.access.custom` | `string` | `empty` |
| `server.logs.access.filename` | `string` | `empty` |
| `server.logs.access.rotate` | `string` | `monthly` |
| `server.logs.access.keep` | `string` | `empty` |
| `server.logs.access.compress` | `bool` | `false` |
| `server.logs.server.enabled` | `bool` | `true` |
| `server.logs.server.format` | `string` | `text` |
| `server.logs.server.custom` | `string` | `empty` |
| `server.logs.server.filename` | `string` | `empty` |
| `server.logs.server.rotate` | `string` | `empty` |
| `server.logs.server.keep` | `string` | `empty` |
| `server.logs.server.compress` | `bool` | `false` |
| `server.logs.error.enabled` | `bool` | `true` |
| `server.logs.error.format` | `string` | `text` |
| `server.logs.error.custom` | `string` | `empty` |
| `server.logs.error.filename` | `string` | `empty` |
| `server.logs.error.rotate` | `string` | `empty` |
| `server.logs.error.keep` | `string` | `empty` |
| `server.logs.error.compress` | `bool` | `false` |
| `server.logs.audit.enabled` | `bool` | `true` |
| `server.logs.audit.format` | `string` | `json` |
| `server.logs.audit.custom` | `string` | `empty` |
| `server.logs.audit.filename` | `string` | `empty` |
| `server.logs.audit.rotate` | `string` | `daily` |
| `server.logs.audit.keep` | `string` | `empty` |
| `server.logs.audit.compress` | `bool` | `false` |
| `server.logs.security.enabled` | `bool` | `true` |
| `server.logs.security.format` | `string` | `fail2ban` |
| `server.logs.security.custom` | `string` | `empty` |
| `server.logs.security.filename` | `string` | `empty` |
| `server.logs.security.rotate` | `string` | `empty` |
| `server.logs.security.keep` | `string` | `empty` |
| `server.logs.security.compress` | `bool` | `false` |
| `server.logs.debug.enabled` | `bool` | `false` |
| `server.logs.debug.format` | `string` | `text` |
| `server.logs.debug.custom` | `string` | `empty` |
| `server.logs.debug.filename` | `string` | `empty` |
| `server.logs.debug.rotate` | `string` | `empty` |
| `server.logs.debug.keep` | `string` | `empty` |
| `server.logs.debug.compress` | `bool` | `false` |
| `server.rate_limit.enabled` | `bool` | `true` |
| `server.rate_limit.requests` | `int` | `120` |
| `server.rate_limit.window_seconds` | `int` | `60` |
| `server.database.driver` | `string` | `sqlite` |
| `server.database.url` | `string` | `platform-specific SQLite path` |
| `server.database.token` | `string` | `empty` |
| `server.cache.type` | `string` | `memory` |
| `server.cache.url` | `string` | `empty` |
| `server.cache.host` | `string` | `(local cache host)` |
| `server.cache.port` | `int` | `6379` |
| `server.cache.username` | `string` | `empty` |
| `server.cache.password` | `string` | `empty` |
| `server.cache.db` | `int` | `0` |
| `server.cache.tls` | `bool` | `false` |
| `server.cache.tls_skip_verify` | `bool` | `true` |
| `server.cache.tls_ca` | `string` | `empty` |
| `server.cache.tls_cert` | `string` | `empty` |
| `server.cache.tls_key` | `string` | `empty` |
| `server.cache.pool_size` | `int` | `10` |
| `server.cache.min_idle` | `int` | `2` |
| `server.cache.timeout` | `string` | `5s` |
| `server.cache.prefix` | `string` | `quotes:` |
| `server.cache.ttl` | `string` | `1h` |
| `server.security.encryption_key` | `string` | `empty` |
| `server.security.encryption_key_version` | `int` | `0` |
| `server.token` | `string` | `empty` |
| `cli.binary_download.dir` | `string` | `{data_dir}/cli` |
| `cli.binary_download.require_auth` | `bool` | `false` |
| `data.cve.api_key` | `string` | `empty` |
| `web.llms.include_endpoints` | `bool` | `true` |
| `web.llms.include_schemas` | `bool` | `false` |
| `web.llms.custom_sections` | `[]config.LLMSCustomSection` | `empty` |
| `web.security.admin` | `string` | `empty` |
| `web.security.cors` | `string` | `empty` |
| `web.security.report_url` | `string` | `https://github.com/apimgr/quotes/security/advisories/new` |
| `web.security.publish_pgp_key` | `bool` | `false` |
| `web.security.expires` | `string` | `empty` |
| `web.security.policy` | `string` | `empty` |
| `web.sitemap.max_urls` | `int` | `50000` |
| `web.sitemap.include_images` | `bool` | `false` |
| `web-ui.theme` | `string` | `dark` |
| `web-ui.logo` | `string` | `empty` |
| `web-ui.favicon` | `string` | `empty` |
| `web-ui.notifications.enabled` | `bool` | `true` |
| `web-ui.notifications.announcements` | `[]string` | `empty` |
| `web-robots.allow` | `[]string` | `2 item(s)` |
| `web-robots.deny` | `[]string` | `empty` |
| `web-security.admin` | `string` | `empty` |
| `web-security.cors` | `string` | `*` |
| `web-security.report_url` | `string` | `empty` |
| `web-security.publish_pgp_key` | `bool` | `false` |
| `web-security.expires` | `string` | `empty` |
| `web-security.policy` | `string` | `empty` |
