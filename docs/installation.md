# Installation

Quotes ships a server binary and a companion `quotes-cli` binary for Linux,
macOS, Windows, and FreeBSD on amd64 and arm64. Release assets are published
without a version in their asset name; use the latest-release URL shown below.

## Docker

Pull and start the published container image with Docker Compose:

```sh
docker compose -f docker/docker-compose.yml up -d
```

The server listens on port 80 inside the container. Persist `/config` to keep
the generated `server.yml` and `secrets.yml`. The image includes the embedded
quote collections and the Tor and I2P provider executables; no external quote
dataset download is required. I2P is enabled by default in the application
container, while Tor starts automatically when its executable is present.

## Binary downloads

Choose the server or CLI artifact for your operating system and architecture:

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

All artifacts are available from the [latest GitHub release](https://github.com/apimgr/quotes/releases/latest).

### Convenience installer scripts

From a source checkout, run the installer for the current platform:

| Platform | Command |
|---|---|
| Linux | `bash scripts/install-linux.sh` |
| macOS | `bash scripts/install-macos.sh` |
| FreeBSD | `sh scripts/install-bsd.sh` |
| Windows | `powershell -ExecutionPolicy Bypass -File scripts/install-windows.ps1` |

The Unix installers select the matching architecture, verify the downloaded
binary against the release `sha256.txt`, install it under `/usr/local/bin`
when run as root (otherwise under `~/.local/bin`), and call the binary's
`--service --install` command. Linux/macOS/FreeBSD installers accept `-h`/
`--help` and `-V`/`--version`. The Windows PowerShell installer also verifies
the release checksum, adds its binary directory to `PATH`, and installs the
Windows service when run as Administrator; without elevation it leaves the
binary installed for manual use.

### Linux

```sh
arch=$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')
curl -q -LSsf "https://github.com/apimgr/quotes/releases/latest/download/quotes-linux-${arch}" -o quotes
chmod +x quotes
./quotes --port 8080
```

To install and start a system service as root, or a user service as an
unprivileged account, run `quotes --service --install`. Native installations
choose and persist an available port from `64000`–`64999` when none is
configured; containers listen on port `80`. See
[service management](development.md#service-management) for lifecycle
commands.

### macOS

Download the matching `darwin` artifact. This command detects Intel or Apple
Silicon automatically:

```sh
arch=$(uname -m | sed 's/x86_64/amd64/;s/arm64/arm64/')
curl -q -LSsf "https://github.com/apimgr/quotes/releases/latest/download/quotes-darwin-${arch}" -o quotes
chmod +x quotes
xattr -d com.apple.quarantine ./quotes 2>/dev/null || true
./quotes --port 8080
```

Run `quotes --service --install` to install and start a LaunchAgent for your
user, or use `sudo` to install a system LaunchDaemon. Both use the configured
port; a native installation chooses and persists one from `64000`–`64999` if
none was configured.

### Windows

Download the matching `.exe` server and CLI assets from the
[latest release](https://github.com/apimgr/quotes/releases/latest), then place
them in a directory on your `PATH`. Start the server in PowerShell:

```powershell
.\quotes-windows-amd64.exe --port 8080
```

Run PowerShell as Administrator and invoke
`quotes-windows-amd64.exe --service --install` to register and start the
Windows service. The service runs with its Windows service virtual account and
uses the configured port; a native installation chooses and persists one from
`64000`–`64999` if none was configured.

### FreeBSD

Download the matching FreeBSD artifact, make it executable, and start it:

```sh
arch=$(uname -m | sed 's/amd64/amd64/;s/aarch64/arm64/')
curl -q -LSsf "https://github.com/apimgr/quotes/releases/latest/download/quotes-freebsd-${arch}" -o quotes
chmod +x quotes
./quotes --port 8080
```

Run `quotes --service --install` as root to install, enable, and start the
native `rc.d` service. It uses the configured port; a native installation
chooses and persists one from `64000`–`64999` if none was configured.

## CLI

The `quotes-cli` binary requires a server URL unless configured in `cli.yml` or
with `QUOTES_SERVER_PRIMARY`. Replace `YOUR_SERVER` with the host name of your Quotes instance:

```sh
./quotes-cli --server https://YOUR_SERVER categories
./quotes-cli --server https://YOUR_SERVER random --category programming
```

See the [CLI reference](cli.md) for commands, configuration, and token use.

## From source

Development and release builds use the Docker-backed Makefile. See
[development](development.md) for supported targets and validation commands.

## Overlay network setup

Tor support starts automatically when a `tor` executable is installed and
available to the server process. Host installations therefore need Tor
installed and discoverable on `PATH` (or its path configured in `server.tor`).
I2P is off by default on bare-metal and VM installs; enable it with
`server.i2p.enabled`, `I2P_ENABLED`, or the `--i2p` startup flag. I2P needs a
managed `i2pd` provider or an external SAM service, whose default address is
`127.0.0.1:7656`. The Docker image includes `i2pd` and enables I2P by default.

Both overlays serve HTTP only and use their own listener. Their resolved
addresses are not published in autodiscovery unless the matching
`server.tor.expose` or `server.i2p.expose` setting is enabled. See
[integrations](integrations.md#overlay-networks) and
[configuration](configuration.md) for details.
