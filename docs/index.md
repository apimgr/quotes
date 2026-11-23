# Quotes

Quotes is a public read-only application for five embedded collections: general quotes, anime, Chuck Norris facts, dad jokes, and programming quotes. Browse the server-rendered interface or use the REST and GraphQL APIs, the authenticated MCP tools, and the companion `quotes-cli` client.

The application does not provide accounts or submissions. Quote content is embedded at build time; the current bundled records are illustrative sample entries and should not be described as curated or authoritative quotations. Verify rights before redistributing quote content.

## Interfaces

- Browser: `/` with category and author filters, pagination, and system/light/dark themes; `/random` for a random quote; `/quotes/{id}` for a detail page.
- Server pages: about, help, health, privacy, and terms under `/server/`.
- REST: `/api/v1/` with quote, category, and health routes. See [API reference](api.md).
- OpenAPI: `/api/v1/server/swagger`; browser page at `/server/docs/swagger`.
- GraphQL: `POST /api/v1/server/graphql`; browser explorer at `/server/docs/graphql`.
- MCP: `POST /api/v1/server/mcp` (operator Bearer token required); help at `/server/help/mcp`.
- CLI: `quotes-cli --server URL random --category programming`.

Start with [installation](installation.md), [configuration](configuration.md), or the [API reference](api.md).
