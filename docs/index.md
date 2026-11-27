# Quotes

Quotes is a public read-only application for five embedded collections: general quotes, anime, Chuck Norris facts, dad jokes, and programming quotes. Browse the server-rendered interface or use the REST and GraphQL APIs, the authenticated MCP tools, and the companion `quotes-cli` client.

The application does not provide user accounts or community quote submissions. Its report form accepts abuse, DMCA, privacy, and security reports; anonymous owner tokens support MCP and credential management without creating accounts. The five embedded collections contain original project content created at the project owner's direction, with 30 records per collection. No third-party quote corpus is imported.

## Quick start

Pull and start the published container image with Docker Compose:

```sh
docker compose -f docker/docker-compose.yml up -d
```

Or run the downloaded Linux amd64 release asset with `./quotes-linux-amd64 --port 8080`. See the
[installation guide](installation.md) for platform-specific downloads and
service installation.

## Features

- Browse and filter quotes in the server-rendered web interface.
- Read data through REST, GraphQL, and interactive OpenAPI/GraphQL explorers.
- Use the CLI, RSS/Atom feeds, and MCP. MCP clients use their own linked token; the operator can use `server.token` for their own tooling.
- Submit abuse, DMCA, privacy, and security reports without an account.
- Discover service metadata, sitemap, robots policy, and `llms.txt`.

## Interfaces

- Browser: `/` with category and author filters, pagination, and system/light/dark themes; `/random` for a random quote; `/quotes/{id}` for a detail page.
- Server pages: about, help, MCP help, DPO, security, security policy, security acknowledgments, health, privacy, terms, contact reports, preferences, preference sync, and owner-token management under `/server/`.
- REST: `/api/v1/` with quote, category, and health routes. See [API reference](api.md).
- OpenAPI: `/api/v1/server/swagger`; browser page at `/server/docs/swagger`.
- GraphQL: `POST /api/v1/server/graphql`; browser explorer at `/server/docs/graphql`.
- MCP: `POST /api/v1/server/mcp` (valid Bearer token required; operator token accepted for the operator's own tooling); help at `/server/help/mcp`.
- CLI: `quotes-cli --server URL random --category programming`.

## Documentation

- [Installation](installation.md) — binaries, Docker, and service setup.
- [Configuration](configuration.md) — server settings and runtime environment.
- [API reference](api.md) — REST, OpenAPI, GraphQL, feeds, and MCP.
- [CLI reference](cli.md) — commands, profiles, output formats, and tokens.
- [Security](security.md) — public endpoints, owner tokens, reporting, and privacy.
- [Integrations](integrations.md) — discovery, identity, and overlay services.
- [Development](development.md) — Docker-backed development and verification.
- [Translation](translating.md) — translation format, glossary, and review.

## Links

- [Repository](https://github.com/apimgr/quotes)
- Swagger UI: `/server/docs/swagger`
- GraphiQL: `/server/docs/graphql`

## License

MIT — see the [project license](https://github.com/apimgr/quotes/blob/main/LICENSE.md).
