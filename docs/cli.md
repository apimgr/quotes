# Command-line client

`quotes-cli` reads from the Quotes API. Public quote retrieval does not require a token. The client has no compiled public-server default; configure a server with `--server`, `QUOTES_SERVER_PRIMARY`, or `server.primary` in its configuration. Starting `quotes-cli` without a command selects an interface from the environment and `display.mode` configuration. Static default builds provide the TUI and CLI; optional GUI builds provide the native GUI on local Linux, macOS, and Windows desktops.

```sh
quotes-cli --server http://localhost:8080 random --category programming
quotes-cli --server http://localhost:8080 categories
quotes-cli --server http://localhost:8080 quote 42
quotes-cli --server http://localhost:8080 list --category programming --author "Ada Lovelace"
quotes-cli --server http://localhost:8080 --output json random
quotes-cli --update check
quotes-cli --update yes
```

Commands are `random [--category CATEGORY]`, `categories`, `quote ID`, and `list [--category CATEGORY] [--author NAME] [--page N] [--limit N]`. `--output` accepts `table`, `json`, or `plain`; the default is `table`. Choose defaults with `defaults.output`, `output.format`, or `QUOTES_OUTPUT_FORMAT`; the flag takes precedence. `--color` accepts `auto`, `yes`, or `no`; `--pager` accepts `auto`, `always`, or `never`; `--quiet` suppresses non-essential output; and `--verbose` writes request paths to stderr. Quiet mode suppresses verbose details. `--debug` enables client debug mode. `--lang` accepts `en`, `es`, `fr`, `ar`, `de`, or `it` and sends the corresponding `Accept-Language` request header; unsupported values fall back to English. `--token` supplies a Bearer token for protected endpoints; `--token-file FILE` reads the token bytes verbatim from a file. API errors are written to stderr. With no command, the TUI offers random quotes, category counts and browsing, quote lookup by ID, and quote listing with category and author filters. Use arrows or `j`/`k`, Enter, Tab, and `q` to navigate. First run asks for a server URL and optional token, tests the connection, and lets you save the settings.

Both binaries have built-in shell integration. Use `quotes-cli --shell completions [bash|zsh|fish|powershell|pwsh]` to print completion code, or `quotes-cli --shell init [SHELL]` to print an eval-ready command. Omit the shell name to use `$SHELL`; the server binary supports the same commands. `sh`, `dash`, and `ksh` have no completions, but `--shell init` safely prints a comment for them. `--shell --help` shows shell-specific help.

`--update check` checks the configured server for an available CLI release. `--update yes` downloads the matching platform binary, checks its SHA-256 against `/api/autodiscover`, replaces the current executable, and restarts it. Updates are downloaded only from HTTPS servers. On normal invocations, the CLI checks for updates; interactive sessions ask before installing, while non-interactive sessions install silently only when `update.auto: true`. The server can enforce a minimum CLI version; an outdated CLI then requires `--update yes` before it can make API requests.

## Display modes and GUI build

The static `quotes-cli-{os}-{arch}` artifact is built without GUI dependencies and works across distributions. It launches the TUI for interactive runs. The optional `quotes-cli-gui-{os}-{arch}` artifact is built with the `gui` build tag for Linux, macOS, and Windows; it uses `github.com/gogpu/ui` and `github.com/gogpu/gogpu`, requires the platform system loader, and is not a cross-distribution static binary. The static updater does not replace GUI artifacts.

The `display.mode` setting accepts `auto`, `gui`, or `tui`; there are no UI mode flags. `auto` selects the native GUI on a local desktop in GUI builds, otherwise the TUI when an interactive terminal is available. SSH and Mosh sessions use the TUI. A forced GUI request on a static build reports `gui_unavailable`; it never silently starts the TUI. Help, version, and command arguments remain CLI operations. On first GUI launch without a configured server, the setup screen accepts a server URL and optional token, tests discovery and authentication, and offers configuration, default-server, and per-user login-start options.

## Server selection

The client resolves its server in this order: `--server`, `QUOTES_SERVER_PRIMARY`, `server.primary` in the config file, then a compiled default if one is provided by the build. Quotes has no compiled default. URLs must be absolute HTTP or HTTPS URLs.

The default configuration path is `~/.config/apimgr/quotes/cli.yml` on Unix and `%APPDATA%\apimgr\quotes\cli.yml` on Windows. The file is created with user-only permissions. A minimal configuration is:

```yaml
server:
  primary: https://YOUR_SERVER
auth:
  token: ""
output:
  format: table
  color: auto
  pager: auto
  # Optional: true or false overrides NO_COLOR for emoji only; TERM=dumb still disables it
  # emoji: true
  quiet: false
  verbose: false
update:
  auto: false
  check_interval: per_invocation
  channel: stable
display:
  mode: auto
theme:
  mode: auto
tui:
  enabled: true
  theme: dark
  mouse: true
  unicode: true
debug: false
defaults:
  lang: auto
  limit: 20
  output: table
```

`--config NAME` selects `NAME.yml` (or an existing `.yaml`) under the same configuration directory. Paths beginning with `~/` and absolute paths are also accepted. A valid `--server` value is used for that invocation; it is saved only when the configured server is empty or invalid. Set `QUOTES_TOKEN` or use `--token` to provide a token for protected operations.

`output.color` accepts `auto`, `yes`, or `no`; precedence is CLI flag, config, `NO_COLOR`, then TTY/`TERM` auto-detection. `--color=yes` overrides `NO_COLOR` and `TERM=dumb` for color sequences; `--color=no` disables colors. `output.emoji` is optional and overrides `NO_COLOR` for emoji, while `TERM=dumb` always disables emoji. `output.pager` accepts `auto`, `always`, or `never`; `auto` invokes `less` or `more` only when output exceeds the terminal height. `QUOTES_OUTPUT_PAGER`, `QUOTES_OUTPUT_QUIET`, and `QUOTES_OUTPUT_VERBOSE` can set the corresponding defaults.

`defaults.lang` accepts `auto`, `en`, `es`, `fr`, `ar`, `de`, or `it`; `QUOTES_LANG` overrides it and `LANG` supplies the automatic locale. `defaults.limit` sets the list limit (1–1000; default 20). `debug` is overridden by `DEBUG` or `MODE=debug`, then by `--debug`. `tui.enabled: false` disables interactive TUI mode. `tui.theme` accepts `dark` or `light`; `system` follows `theme.mode` and terminal background detection. `tui.mouse` controls mouse input, and `tui.unicode: false` uses ASCII prompts and masks.

JSON preserves structured quote fields. Table output includes column headings. Plain output keeps the compact quote text suitable for shell pipelines. Color styling applies only to table headings, so JSON output remains machine-readable.

CLI configuration, data, cache, and log directories are scoped to the invoking user. The CLI creates them with private permissions; updater audit events are recorded in `cli.log` and written to stderr.
