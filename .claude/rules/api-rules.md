# API Rules (PART 13, 14, 15)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Never register a route outside the central route contract or duplicate route metadata in API discovery documents.
- Never expose private or internal routes in public discovery output.
- Never treat TLS errors or invalid request values as successful responses.

## CRITICAL - ALWAYS DO
- Keep API versions, health/version responses, route auth, and response formats consistent.
- Define each route once; use the registry as the source for routing, OpenAPI, llms.txt, and text serialization.
- Apply the specified TLS and certificate lifecycle behavior.

## KEY DECISIONS (pre-answered)
| Question | Answer | Spec Reference |
|---|---|---|
| Where does route metadata come from? | One shared `Route` registry entry. | PART 14 |
| What does health report? | The specified health contract and status, not a guessed summary. | PART 13 |

## TERMINOLOGY
| Term | Meaning |
|---|---|
| Listed route | A route eligible for public API documentation output. | PART 14 |
| Internal route | A mounted route excluded from public discovery surfaces. | PART 14 |

## QUICK REFERENCE
Register once; derive every public description from the route definition; preserve version, auth, TLS, and health contracts.

---
For complete details, see AI.md PART 13, 14, 15
