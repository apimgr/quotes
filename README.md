# Quotes

[![CI](https://github.com/apimgr/quotes/actions/workflows/ci.yml/badge.svg)](https://github.com/apimgr/quotes/actions/workflows/ci.yml) [![Release](https://img.shields.io/github/v/release/apimgr/quotes)](https://github.com/apimgr/quotes/releases) [![License](https://img.shields.io/github/license/apimgr/quotes)](LICENSE.md) [![Documentation](https://readthedocs.org/projects/apimgr-quotes/badge/?version=latest)](https://apimgr-quotes.readthedocs.io/)

## About

Quotes is a public, read-only quote service with a server-rendered browser, REST API, GraphQL endpoint, and companion CLI. Its five embedded collections cover general quotes, anime, Chuck Norris facts, dad jokes, and programming. The service has no user accounts or community quote submissions.

The owner declares the bundled collections as original material created at the owner's direction. Each collection currently contains 30 records, each record is marked `source: "original"`, and the content is released under the project MIT license. No third-party quote corpus is included.

## Features

The browser and CLI support six languages: English, Spanish, French, Arabic, German, and Italian. See the [translation guide](docs/translating.md) for locale and contribution details.

- Browse and filter quotes in a server-rendered web interface.
- Read quotes through the versioned REST API and GraphQL endpoint.
- Use interactive, locally hosted Swagger UI and GraphiQL pages.
- Connect an MCP client with its own linked Bearer token; the operator can use `server.token` for their own MCP tooling.
- Submit abuse, DMCA, privacy, or security reports without creating an account.
- Discover Atom and RSS feeds, the sitemap, robots policy, and `llms.txt`.
- Run optional Tor and I2P overlay services.

## Production

Download the latest release from [GitHub Releases](https://github.com/apimgr/quotes/releases/latest). Server and CLI binaries are published for all eight platform and architecture combinations:

| Platform | Server | CLI |
|---|---|---|
| Linux amd64 | `quotes-linux-amd64` | `quotes-cli-linux-amd64` |
| Linux arm64 | `quotes-linux-arm64` | `quotes-cli-linux-arm64` |
| macOS Intel | `quotes-darwin-amd64` | `quotes-cli-darwin-amd64` |
| macOS Apple Silicon | `quotes-darwin-arm64` | `quotes-cli-darwin-arm64` |
| Windows amd64 | `quotes-windows-amd64.exe` | `quotes-cli-windows-amd64.exe` |
| Windows arm64 | `quotes-windows-arm64.exe` | `quotes-cli-windows-arm64.exe` |
| FreeBSD amd64 | `quotes-freebsd-amd64` | `quotes-cli-freebsd-amd64` |
| FreeBSD arm64 | `quotes-freebsd-arm64` | `quotes-cli-freebsd-arm64` |

Optional native GUI CLI artifacts (`quotes-cli-gui-{os}-{arch}`) are also built for Linux, macOS, and Windows. GUI artifacts use the operating system's dynamic loader and are desktop-specific; the default CLI artifacts remain static and distribution-independent.

Platform installer scripts verify release checksums and use the binary's native service installer. See the [installation guide](docs/installation.md) for commands and options.

### Linux

```sh
arch=$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')
curl -q -LSsf "https://github.com/apimgr/quotes/releases/latest/download/quotes-linux-${arch}" -o quotes
chmod +x quotes
./quotes --port 8080
```

### macOS

```sh
arch=$(uname -m | sed 's/x86_64/amd64/;s/arm64/arm64/')
curl -q -LSsf "https://github.com/apimgr/quotes/releases/latest/download/quotes-darwin-${arch}" -o quotes
chmod +x quotes
xattr -d com.apple.quarantine ./quotes 2>/dev/null || true
./quotes --port 8080
```

### Windows

Download the matching `.exe` from the releases page and add its directory to `PATH`. Start the server in PowerShell:

```powershell
.\quotes-windows-amd64.exe --port 8080
```

### FreeBSD

Download the matching `freebsd` artifact, make it executable, and run it:

```sh
arch=$(uname -m | sed 's/amd64/amd64/;s/aarch64/arm64/')
curl -q -LSsf "https://github.com/apimgr/quotes/releases/latest/download/quotes-freebsd-${arch}" -o quotes
chmod +x quotes
./quotes --port 8080
```

On Linux and macOS, `quotes --service --install` installs a user service without elevated privileges or a system service when elevated. Native installs choose and persist an available port in the `64000`–`64999` range when none is configured; containers use port `80`. Windows and FreeBSD system service installation requires administrator or root privileges.

### Docker

```sh
docker compose -f docker/docker-compose.yml up -d
```

The container listens on port 80. Persist `/config` to retain the generated server configuration and secrets.

## Client

Set the server URL to your Quotes instance with `--server URL`, `QUOTES_SERVER_PRIMARY`, or `server.primary` in the CLI configuration. Quotes has no compiled default server; replace `YOUR_SERVER` below with your server host.

```sh
quotes-cli --server https://YOUR_SERVER random --category programming
quotes-cli --server https://YOUR_SERVER categories
quotes-cli --server https://YOUR_SERVER quote 42
quotes-cli --server https://YOUR_SERVER list --author "Ada Lovelace"
```

The client supports random quotes, category listing, quote lookup, and filtered, paginated listing. `--output` accepts `table`, `json`, or `plain`; `QUOTES_OUTPUT_FORMAT` and `output.format` set the default. See the [CLI guide](docs/cli.md) for display modes, shell integration, GUI artifacts, and configuration details.

## Configuration

The server stores configuration under the platform-specific configuration directory. Containers persist `/config`. The CLI configuration is `~/.config/apimgr/quotes/cli.yml` on Unix and `%APPDATA%\apimgr\quotes\cli.yml` on Windows. See the [configuration guide](docs/configuration.md) for configuration areas, defaults, permissions, and secrets handling.

## API

REST uses `/api/v1/`; GraphQL is at `/api/v1/server/graphql`; OpenAPI JSON is at `/api/v1/server/swagger`. Browser API references are at `/server/docs/swagger` and `/server/docs/graphql`. MCP is at `/api/v1/server/mcp` and requires a valid Bearer token. The operator may use `server.token` for their own tooling; each standing MCP client should use a separate linked token. Linked tokens have `read` or `read-write` scope; token management is available at `/server/preferences/tokens`.

See the [API reference](docs/api.md) for routes, formats, and authentication behavior.

## Other

See the [security model](docs/security.md), [integrations guide](docs/integrations.md), and [documentation site](https://apimgr-quotes.readthedocs.io/) for security reporting, service integrations, and additional guides. Third-party content is included only when its source and redistribution rights are documented.

## Development

Go development uses the six Docker-backed Makefile targets: `make build`, `make test`, `make release`, `make docker`, `make clean`, and `make help`. `make test` runs the lint, vulnerability, unit, coverage, and basic binary checks. Run browser E2E checks separately with `./tests/e2e.sh`, and service integration checks with `./tests/incus.sh`. Release builds create local artifacts only; CI publishes releases. Build and release artifacts target Linux, macOS, Windows, and FreeBSD on amd64 and arm64. See the [development guide](docs/development.md).

## Disclaimer

This software is provided "as is" without warranty of any kind. Use and operate it at your own risk.

- **No warranty:** The authors are not responsible for damages, data loss, or other issues arising from use of this software.
- **No professional advice:** The service and its content do not provide legal, financial, medical, or other professional advice.
- **Third-party services:** External services used for overlays, data, or integrations may have separate terms of service.
- **Security:** The project follows security practices, but no software can be guaranteed free of vulnerabilities.
- **Production use:** Evaluate the service and its configuration before deploying it in production.

By using this software, you acknowledge that you have read and understood this disclaimer.

## License

MIT — see [LICENSE.md](LICENSE.md).
