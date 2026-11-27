# Development

The application is implemented in Go. `src/catalog` loads the embedded collections, `src/database` stores anonymous owner-token, report, and scheduler state in SQLite, `src/server/metrics` implements the metrics registry and collectors, `src/server` serves server-rendered pages, REST, OpenAPI, GraphQL, MCP, metrics, and report workflows, and `src/client` contains the companion CLI.

The Makefile has exactly six Docker-backed targets: `build`, `test`, `release`, `docker`, `clean`, and `help`. `build` creates unoptimized debug binaries for all eight platform targets; `release` creates optimized local production artifacts, a source archive, and checksums without tagging or publishing. `test` runs vet, static analysis, vulnerability checks, shell lint, Go unit and GUI tests with coverage, then `tests/run_tests.sh` for basic binary checks. Go unit tests live beside their packages. Run the SSR, no-JavaScript, and full browser suite on demand with `./tests/e2e.sh`. Run systemd or OpenRC service checks with `./tests/incus.sh [systemd|openrc|all]`.

## Service management

Install, enable, and start the native service with `quotes --service --install`.
An elevated install creates the system service; an unprivileged install uses a
user service where the platform supports it. Native installations select and
persist an available port from `64000`–`64999` when none is configured;
containers use port `80`. Use `quotes --service start`,
`stop`, `restart`, or `reload` to control it. `quotes --service --disable`
stops and disables the service while retaining its files and data.

`quotes --service --uninstall` asks for confirmation before removing service
data and the managed Unix service account. On Unix, the service binds its
listener before dropping to the dedicated `quotes` account; Windows uses its
service virtual account. The in-process scheduler starts with the server and
is controlled with `quotes --maintenance scheduler list|run|enable|disable`.

`--maintenance backup [file]` writes a verified `.tar.gz` archive with a consistent SQLite snapshot; `--maintenance restore <file>` verifies it before replacing data and requires root, an empty database, or the service account plus the operator token. `--maintenance mode <production|development|debug>` updates `server.mode`; the running server applies it after restart, and an invalid value exits with usage code 2. `--maintenance setup` resets server configuration with its required authorization checks. PGP key management supports `--maintenance pgp generate`, `--maintenance pgp export public [path]`, and `--maintenance pgp delete`. `--update check` and `--update yes` use the project GitHub Releases API, and `--update branch <name>` persists the channel under `server.update.branch`.

Do not add third-party quote content without a documented source and redistribution license. The owner declares the bundled collections as original material created at the owner's direction; records carry `source: "original"`, and the license/about page provide the attribution. Update the API and user documentation when public behavior or configuration changes. Run the complete local quality gate with `make test`; it includes Go tests, vulnerability checks, and shell lint.

## License attribution check

Run the dependency license check inside the supported Go container:

```sh
docker run --rm --name "quotes-$(od -An -N4 -tx1 /dev/urandom | tr -d ' \n')" \
  -v "$PWD:/app" -w /app casjaysdev/go:latest ./scripts/verify-licenses.sh
```

The checker accepts `--help` and `--version`; both can be passed after the
script path to inspect its usage without scanning dependencies.

The script scans Go package licenses and temporarily generates a protected
report under `${TMPDIR:-/tmp}/{project_org}/{internal_name}-XXXXXX`; it removes
the report directory when it exits. It downloads the pinned `mathutil` module
into the Go module cache to verify its license file. It fails closed if the
scanner cannot identify another dependency license or detects a copyleft
license. The pinned `modernc.org/mathutil` license is checked against its
tagged upstream `LICENSE` file and added as BSD-3-Clause because `go-licenses`
cannot classify its nonstandard license metadata.
