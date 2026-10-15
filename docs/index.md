# Quotes

Quotes is a read-only service for five embedded collections: general quotes, anime quotes, Chuck Norris facts, dad jokes, and programming quotes. Browse the server rendered interface or use the public REST and GraphQL APIs and companion `quotes-cli` client.

The application does not store user data or call external services. Collection contents ship inside the binary. The currently bundled records are generated sample entries and should not be represented as curated or authoritative quotations until licensed source data is selected.

## Interfaces

- Browser: `/` with category and author filters, pagination, and system/light/dark themes; `/random` for a random quote and `/quotes/{id}` for a quote detail page.
- REST: `/api/v1/quotes`, `/api/v1/quotes/random`, `/api/v1/categories`, and `/api/v1/quotes/{id}`.
- OpenAPI: `/api/v1/server/swagger`; browser documentation at `/server/docs/swagger`.
- GraphQL: `POST /api/v1/server/graphql`; browser explorer at `/server/docs/graphql`.
- Server information: `/server/about`, `/server/help`, `/server/healthz`, `/server/privacy`, and `/server/terms`.
- CLI: `quotes-cli --server URL random --category programming`.

Start with [installation](installation.md), [configuration](configuration.md), or the [API reference](api.md).
