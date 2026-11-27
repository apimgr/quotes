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

## KEY DECISIONS (pre-answered)
| Question | Answer | Spec Reference |
|----------|--------|----------------|
| What governs implementation? | Quote retrieval is public, read-only, versioned, and consistent across REST, GraphQL, CLI, and browser. | PART 13, 14, 15 |

## TERMINOLOGY
| Term | Meaning |
|------|---------|
| canonical route | versioned server path; alias → documented direct path using the same behavior. |

## QUICK REFERENCE
- Keep OpenAPI, GraphQL, response envelopes, content negotiation, and docs aligned with registered handlers.
- Owner token JSON includes scope: primary tokens are `global`; linked tokens are `read` or `read-write` (default `read-write`). Token management remains primary-only.

For complete details, see AI.md PART 13, 14, 15
