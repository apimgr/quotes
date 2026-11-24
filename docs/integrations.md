# Integrations

Quotes exposes public REST and GraphQL interfaces, OpenAPI documentation, plain-text quote responses, and a Streamable HTTP MCP endpoint. REST and GraphQL quote reads need no credentials. MCP requests require an owner-linked Bearer token. Issue a primary owner token at `/server/preferences/tokens`, then create a separate linked token for each MCP client; do not use the operator token.

A web app manifest and service worker provide install metadata and basic offline caching. Security contact information is published at `/.well-known/security.txt` and `/security.txt`. The project does not currently enable OAuth/OIDC, webhooks, federation, native app associations, or external runtime data APIs. There is no configured official server URL; clients select the server they use.
