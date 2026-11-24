# Quotes

Quotes is a public read-only application for five embedded collections: general quotes, anime, Chuck Norris facts, dad jokes, and programming quotes. Browse the server-rendered interface or use the REST and GraphQL APIs, the authenticated MCP tools, and the companion `quotes-cli` client.

The application does not provide user accounts or community quote submissions. Its report form accepts abuse, DMCA, privacy, and security reports; anonymous owner tokens support MCP and credential management without creating accounts. Quote content is embedded at build time and must retain documented source and redistribution rights.

## Interfaces

- Browser: `/` with category and author filters, pagination, and system/light/dark themes; `/random` for a random quote; `/quotes/{id}` for a detail page.
- Server pages: about, help, MCP help, DPO, security, security policy, security acknowledgments, health, privacy, terms, contact reports, preferences, preference sync, and owner-token management under `/server/`.
- REST: `/api/v1/` with quote, category, and health routes. See [API reference](api.md).
- OpenAPI: `/api/v1/server/swagger`; browser page at `/server/docs/swagger`.
- GraphQL: `POST /api/v1/server/graphql`; browser explorer at `/server/docs/graphql`.
- MCP: `POST /api/v1/server/mcp` (owner-linked Bearer token required); help at `/server/help/mcp`.
- CLI: `quotes-cli --server URL random --category programming`.

Start with [installation](installation.md), [configuration](configuration.md), or the [API reference](api.md).
