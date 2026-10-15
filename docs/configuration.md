# Configuration

The server accepts `--address`, `--port`, and `--config`. The default configuration path is selected by the application configuration package; pass `--config` to use a specific `server.yml` file. The initial configuration contains server bind settings, web theme, logging, and rate-limit settings.

The `server.rate_limit` settings are `enabled`, `requests`, and `window_seconds`. Defaults are enabled, 120 requests, and 60 seconds per client address. Set `server.address` and `server.port` in the configuration or use the corresponding command-line flags.

The CLI stores its server URL in the invoking user's configuration directory with user-only file permissions. Use `quotes-cli --server https://example.invalid ...` for a one-off URL, or set `QUOTES_SERVER`. The project does not define a default public server; no official site URL is assumed.

All quote content is compiled into the binary. There is no runtime database, account configuration, admin panel, or external service integration.
