# Quotes

A public, read-only web application and API for five embedded quote collections: general quotes, anime quotes, Chuck Norris facts, dad jokes, and programming quotes. A companion `quotes-cli` client supports shell use. No accounts, submissions, admin panel, authentication, or runtime external services are part of the product.

## Run

```sh
docker build -f docker/Dockerfile -t quotes .
docker run --rm -p 8080:80 quotes
```

Open `http://localhost:8080/` to browse and filter the catalog. Quote IDs open individual detail pages at `/quotes/{id}`, and `/random` displays a random quote (optionally filter with `?category=programming`). The server listens on port 80 in the container. A direct binary can be started with `./quotes --port 8080`.

## API

```sh
curl http://localhost:8080/api/v1/quotes/random
curl 'http://localhost:8080/api/v1/quotes?category=programming&page=1&limit=20'
curl http://localhost:8080/api/v1/categories
```

REST uses `/api/v1/`; GraphQL is at `POST /api/v1/server/graphql`, and OpenAPI JSON is at `/api/v1/server/swagger`. Browser documentation is served at `/server/docs/swagger` and `/server/docs/graphql`. The user-facing pages include `/server/about`, `/server/help`, `/server/healthz`, `/server/privacy`, and `/server/terms`.

## CLI

```sh
quotes-cli --server http://localhost:8080 random --category programming
quotes-cli --server http://localhost:8080 categories
quotes-cli --server http://localhost:8080 quote 42
```

The CLI also reads `QUOTES_SERVER` or its user configuration file. There is no default public server URL.

## Data and configuration

The data is embedded at build time from `src/data/*.json`. The current bundled entries are sample generated records; they should not be described as curated quotations until authorized source collections and licenses are selected. Server configuration is YAML (`server.yml`); rate limiting is enabled by default at 120 requests per client per 60 seconds. See [the documentation](docs/index.md).

## Development

The project uses Docker-backed Makefile targets (`dev`, `local`, `build`, `test`, `release`, and `docker`). See [development docs](docs/development.md).
