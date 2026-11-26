# Configuration

Configuration is file-backed. `AI.md` PART 5 defines the configuration contract; `IDEA.md` defines product behavior. The merged files are the only source of operator-editable settings. Database records hold application state, not configuration. There is no web administration panel or runtime configuration mutation API.

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

Startup also migrates the legacy setting names listed in AI.md PART 5 before splitting area files. Canonical keys win when both names are present. Explicitly retired keys are removed after the backup is made; unlisted keys in defined sections are retained for the normal unknown-key warning. This includes the old logging, scheduler, web, database, branding, SEO, page content, contact, admin, SSL, and security names. The HTTPS port is folded into a single numeric `server.port` as a dual-port value; an already dual port is preserved. Legacy scheduler task rows retain only `enabled`, `schedule`, `retry_on_fail`, and `retry_delay`.

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

Booleans accept case-insensitive true/false forms including `1`/`0`, `yes`/`no`, `on`/`off`, `enable`/`disable`, `y`/`n`, and the additional forms defined in AI.md PART 5. Empty values use the setting default; invalid values warn and use the default.

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

SSL settings are `server.ssl.enabled`, `server.ssl.cert`, `server.ssl.key`, and the Let's Encrypt settings defined in PART 12. The certificate and key paths refer to operator-managed files. Request limits and response compression are configured in `yaml.d/security.yml` and `yaml.d/web.yml`, respectively. See AI.md PART 12 for complete keys, defaults, proxy trust behavior, and graceful-reload rules.

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

Other filesystem overrides include `--cache`, `--log`, and `--pid`. Listener and URL overrides include `--address`, `--port`, and `--baseurl`. `DOMAIN` overrides hostname discovery. `MODE` selects `production`, `development`, or `debug`; `DEBUG` independently enables debug behavior. SMTP environment variables are `SMTP_HOST`, `SMTP_PORT`, `SMTP_USERNAME`, `SMTP_PASSWORD`, `SMTP_FROM_NAME`, `SMTP_FROM_EMAIL`, and `SMTP_TLS`. See AI.md PART 12 for exact precedence and the full runtime/init-only variable lists.

## Area configuration

All feature settings and defaults are specified in AI.md. The following examples show the area-file shape; omit settings that should use their defaults.

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

Supported formats and rotation/retention values are listed in AI.md PART 12. Secret values and raw owner tokens are never written to logs. Successful health probes are omitted from the access log; failed probes remain visible.

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

Self-hosted Matomo, Piwik, OWA, and Umami require an HTTPS `url`. Fathom and Plausible accept an optional HTTPS URL for self-hosting. Google, Simple Analytics, and Cloudflare do not use `url`. See AI.md PART 12 for each provider's ID format.

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
      email: reports@example.org
```

```yaml
# yaml.d/notifications.yml
server:
  notifications:
    email:
      smtp:
        host: mail.example.org
        port: 587
        tls: auto
```

Use `yaml.d/pages.yml` for report page behavior and antispam, `yaml.d/branding.yml` for branding and SEO, `yaml.d/privacy.yml` for privacy and compliance disclosures, `yaml.d/metrics.yml` for metrics, `yaml.d/tor.yml` and `yaml.d/i2p.yml` for overlay networking, and `yaml.d/maintenance.yml` for scheduling, backups, restores, and updates. Contact reports do not create accounts or accept quote contributions. Owner tokens are anonymous credentials; they are not user accounts.

Email behavior, templates, scheduler tasks, GeoIP, metrics endpoints, backup/restore safeguards, and update channels are documented in AI.md PARTS 17–22. Backups include configuration area files and `secrets.yml` when present. Treat backup files as sensitive.

## Platform paths and operations

Use the platform configuration, data, cache, log, and backup directories selected by the application. Container deployments use persistent mounts for configuration and data. The exact platform paths and service-manager behavior are defined in AI.md PARTS 4, 8, and 12; do not assume a host path when using a container.

The built-in scheduler persists state in the configured database; do not use cron or an external scheduler. Maintenance commands require operator authorization where specified. The server has no web panel for changing configuration: edit the files and let the watcher reload supported settings. Listener address, port, daemonization, and PID-file changes require a restart.

## Related requirements

- AI.md PART 5: file layout, merge, migration, permissions, environment variables, paths, and boolean parsing.
- AI.md PART 12: server settings, validation, modes, database/cache, TLS, request limits, and proxy handling.
- AI.md PART 16: browser behavior, web/security settings, CSRF, and pages.
- AI.md PARTS 17–22: email, scheduler, GeoIP, metrics, backup/restore, and updates.
- AI.md PART 28: Docker-backed development and verification workflow.
- AI.md PART 29: Read the Docs documentation requirements and navigation.
