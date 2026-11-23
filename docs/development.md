# Development

The application is implemented in Go. `src/catalog` loads the embedded collections, `src/httpapp` serves the browser, REST, OpenAPI, and GraphQL surfaces, and `src/client` contains the companion CLI.

Use the Docker-backed Makefile targets `dev`, `local`, `build`, `test`, `release`, and `docker`. Tests live beside their packages. Do not add quote content without a documented source and redistribution license. Update the API and user documentation when public behavior or configuration changes.
