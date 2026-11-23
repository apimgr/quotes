# API Rules (PART 13, 14, 15)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Do not retain superseded project routes as compatibility shims.
- Do not add authentication to public quote reads; protect only routes whose spec requires credentials.
- Do not let REST, GraphQL, OpenAPI, MCP, browser, and CLI surfaces drift.

## CRITICAL - ALWAYS DO
- Version REST routes and return consistent JSON/error envelopes.
- Keep documented aliases attached to the same handler as canonical routes.
- Implement health/versioning and built-in TLS requirements from their PARTs.

## Key rules summary
- API inventory in IDEA.md must cover each project-defined route.
- OpenAPI and GraphQL descriptions must match actual handlers and response shapes.
- Plain text and HTML negotiation follows the exact frontend/API rules.

For complete details, see AI.md PART 13, 14, 15
