# Integrations

Quotes exposes public REST and GraphQL interfaces, OpenAPI documentation, plain-text quote responses, and a Streamable HTTP MCP endpoint. REST and GraphQL quote reads need no credentials. MCP requests require an operator Bearer token; the server does not issue linked tokens or per-client credentials yet.

A web app manifest and service worker provide install metadata and basic offline caching. Security contact information is published at `/.well-known/security.txt` and `/security.txt`. The project does not currently enable OAuth/OIDC, webhooks, federation, native app associations, or external runtime data APIs. There is no configured official server URL; clients select the server they use.
