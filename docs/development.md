# Development

The application is implemented in Go and keeps its collection data in `src/data/`. The `src/catalog` package loads and validates the embedded datasets; `src/httpapp` serves the web, REST, OpenAPI, and GraphQL surfaces; `src/client` contains the companion CLI.

Use the Docker-backed Makefile targets `dev`, `local`, `build`, `test`, `release`, and `docker`. Tests live beside their packages. Do not add collection content without a source and distribution license; update documentation whenever public routes or configuration change.
