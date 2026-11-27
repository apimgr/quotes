# Command-line client

`quotes-cli` reads from the Quotes API. Public quote retrieval does not require a token. The client has no compiled public-server default; configure a server with `--server`, `QUOTES_SERVER_PRIMARY`, or `server.primary` in its configuration. Starting `quotes-cli` without a command opens its terminal interface when an interactive terminal is available.

```sh
quotes-cli --server http://localhost:8080 random --category programming
quotes-cli --server http://localhost:8080 categories
quotes-cli --server http://localhost:8080 quote 42
quotes-cli --server http://localhost:8080 list --category programming --author "Ada Lovelace"
quotes-cli --server http://localhost:8080 --output json random
quotes-cli --update check
quotes-cli --update yes
```

Commands are `random [--category CATEGORY]`, `categories`, `quote ID`, and `list [--category CATEGORY] [--author NAME] [--page N] [--limit N]`. `--output` accepts `plain` or `json`. `--lang` accepts `en`, `es`, `fr`, `ar`, `de`, or `it` and sends the corresponding `Accept-Language` request header; unsupported values fall back to English. `--token` supplies a Bearer token for protected endpoints. API errors are written to stderr. The no-command TUI supports a random quote, category counts, and author search; use arrow keys or `j`/`k`, Enter, and `q` to navigate.

`--update check` checks the configured server for an available CLI release. `--update yes` downloads the matching platform binary, checks its SHA-256 against `/api/autodiscover`, replaces the current executable, and restarts it. Updates are downloaded only from HTTPS servers. On normal invocations, the CLI checks for updates; interactive sessions ask before installing, while non-interactive sessions install silently only when `update.auto: true`. The server can enforce a minimum CLI version; an outdated CLI then requires `--update yes` before it can make API requests.

## Server selection

The client resolves its server in this order: `--server`, `QUOTES_SERVER_PRIMARY`, `server.primary` in the config file, then a compiled default if one is provided by the build. Quotes has no compiled default. URLs must be absolute HTTP or HTTPS URLs.

The default configuration path is `~/.config/apimgr/quotes/cli.yml` on Unix and `%APPDATA%\apimgr\quotes\cli.yml` on Windows. The file is created with user-only permissions. A minimal configuration is:

```yaml
server:
  primary: https://YOUR_SERVER
auth:
  token: ""
update:
  auto: false
  check_interval: per_invocation
  channel: stable
display:
  mode: auto
theme:
  mode: auto
```

`--config NAME` selects `NAME.yml` (or an existing `.yaml`) under the same configuration directory. Paths beginning with `~/` and absolute paths are also accepted. A valid `--server` value is used for that invocation; it is saved only when the configured server is empty or invalid. Set `QUOTES_TOKEN` or use `--token` to provide a token for protected operations.

CLI configuration, data, cache, and log directories are scoped to the invoking user. The CLI creates them with private permissions; updater audit events are recorded in `cli.log` and written to stderr.
