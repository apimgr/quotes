# Quotes

[![CI](https://github.com/apimgr/quotes/actions/workflows/ci.yml/badge.svg)](https://github.com/apimgr/quotes/actions/workflows/ci.yml) [![Release](https://img.shields.io/github/v/release/apimgr/quotes)](https://github.com/apimgr/quotes/releases) [![License](https://img.shields.io/github/license/apimgr/quotes)](LICENSE.md) [![Documentation](https://readthedocs.org/projects/apimgr-quotes/badge/?version=latest)](https://apimgr-quotes.readthedocs.io/)

Quotes serves five embedded quote collections through a public read-only web interface, REST API, GraphQL endpoint, and CLI. Collection data ships in the executable; no accounts or quote submissions are supported.

## Install

Download the latest release from [GitHub Releases](https://github.com/apimgr/quotes/releases/latest). Each release contains server and CLI binaries for Linux, macOS, Windows, and FreeBSD on amd64 and arm64.

| Platform | Server binary | CLI binary |
|---|---|---|
| Linux amd64 | `quotes-linux-amd64` | `quotes-cli-linux-amd64` |
| Linux arm64 | `quotes-linux-arm64` | `quotes-cli-linux-arm64` |
| macOS Intel | `quotes-darwin-amd64` | `quotes-cli-darwin-amd64` |
| macOS Apple Silicon | `quotes-darwin-arm64` | `quotes-cli-darwin-arm64` |
| Windows amd64 | `quotes-windows-amd64.exe` | `quotes-cli-windows-amd64.exe` |
| Windows arm64 | `quotes-windows-arm64.exe` | `quotes-cli-windows-arm64.exe` |
| FreeBSD amd64 | `quotes-freebsd-amd64` | `quotes-cli-freebsd-amd64` |
| FreeBSD arm64 | `quotes-freebsd-arm64` | `quotes-cli-freebsd-arm64` |

On Linux, install the matching server binary with:

```sh
arch=$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')
curl -LSsf "https://github.com/apimgr/quotes/releases/latest/download/quotes-linux-${arch}" -o quotes
chmod +x quotes
./quotes --port 8080
```

On macOS, choose the matching `darwin` binary, make it executable, then remove quarantine if needed:

```sh
chmod +x quotes
xattr -d com.apple.quarantine ./quotes 2>/dev/null || true
```

Windows users can download the matching `.exe` files and add their directory to `PATH`.

## Docker

```sh
docker compose -f docker/docker-compose.yml up -d
```

The container listens on port 80. Map it to a host port and persist `/config` so generated `server.yml` and `secrets.yml` remain available.

## quotes-cli

```sh
quotes-cli --server http://localhost:8080 random --category programming
quotes-cli --server http://localhost:8080 categories
quotes-cli --server http://localhost:8080 quote 42
```

The CLI accepts `--server URL`, `QUOTES_SERVER`, or its user config. It supports `random`, `categories`, and `quote ID`; `--output` accepts `plain` or `json`. `--lang` supports English, Spanish, French, Arabic, German, and Italian.

## Web and API

Open `http://localhost:8080/` to filter and browse quotes. REST uses `/api/v1/`; GraphQL is at `POST /api/v1/server/graphql`; OpenAPI JSON is at `/api/v1/server/swagger`. Browser API references are at `/server/docs/swagger` and `/server/docs/graphql`. The MCP Streamable HTTP endpoint is `/api/v1/server/mcp` and requires the operator Bearer token stored in `secrets.yml`.

See the [API documentation](docs/api.md), [configuration reference](docs/configuration.md), and [integration guide](docs/integrations.md).

## Development

Go development uses Docker-backed Makefile targets: `make dev`, `make local`, `make build`, `make test`, `make release`, and `make docker`. The release build targets all eight OS/architecture pairs listed above. See [development documentation](docs/development.md).

## License

MIT — see [LICENSE.md](LICENSE.md).
