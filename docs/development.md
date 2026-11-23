# Development

The application is implemented in Go. `src/catalog` loads the embedded collections, `src/httpapp` serves the browser, REST, OpenAPI, and GraphQL surfaces, and `src/client` contains the companion CLI.

Use the Docker-backed Makefile targets `dev`, `local`, `build`, `test`, `release`, and `docker`. Go unit tests live beside their packages. The executable basic binary-validation suite is `./tests/run_tests.sh`; it requires Docker and delegates to `tests/docker.sh`. AI.md also specifies separate full-app Incus and browser E2E suites.

Do not add quote content without a documented source and redistribution license. Update the API and user documentation when public behavior or configuration changes. Run Go tests and vulnerability checks with `make test`; run shell lint through the repository's lint wrapper.
