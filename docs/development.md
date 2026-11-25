# Development

The application is implemented in Go. `src/catalog` loads the embedded collections, `src/database` stores anonymous owner-token, report, and scheduler state in SQLite, `src/httpapp` serves server-rendered pages, REST, OpenAPI, GraphQL, MCP, metrics, and report workflows, and `src/client` contains the companion CLI.

Use the Docker-backed Makefile targets `dev`, `local`, `build`, `test`, `release`, and `docker`. Go unit tests live beside their packages. Run the SSR, no-JavaScript, and browser suite with `./tests/e2e.sh`; run systemd and OpenRC service checks with `./tests/incus.sh`. The executable basic binary-validation suite is `./tests/run_tests.sh`; it requires Docker and delegates to `tests/docker.sh`.

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

`--maintenance backup [file]` writes a verified `.tar.gz` archive with a consistent SQLite snapshot; `--maintenance restore <file>` verifies it before replacing data and requires root, an empty database, or the service account plus the operator token. `--maintenance mode <mode>` and `--maintenance setup` update/reset server configuration with their required authorization checks. PGP key management supports `--maintenance pgp generate`, `--maintenance pgp export public [path]`, and `--maintenance pgp delete`. `--update check` and `--update yes` use the project GitHub Releases API, and `--update branch <name>` persists the channel under `server.update.branch`.

Do not add quote content without a documented source and redistribution license. Update the API and user documentation when public behavior or configuration changes. Run Go tests and vulnerability checks with `make test`; run shell lint through the repository's lint wrapper.

## License attribution check

Run the dependency license check inside the supported Go container:

```sh
docker run --rm --name "quotes-$(od -An -N4 -tx1 /dev/urandom | tr -d ' \n')" \
  -v "$PWD:/app" -w /app casjaysdev/go:latest ./scripts/verify-licenses.sh
```

The checker accepts `--help` and `--version`; both can be passed after the
script path to inspect its usage without scanning dependencies.

The script scans Go package licenses and temporarily generates a report under
`/tmp/apimgr/`; it removes the report directory when it exits. It fails closed
if the scanner cannot identify a dependency license or detects a copyleft
license. At the current dependency versions, `go-licenses` reports
`modernc.org/mathutil` v1.7.1 as unknown despite its BSD-3-Clause attribution
in `LICENSE.md`, so the check exits non-zero until the scanner can classify it.
