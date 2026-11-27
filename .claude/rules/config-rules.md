# Configuration Rules (PART 5, 6, 12)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Do not add web administration or runtime configuration mutation APIs.
- Do not silently ignore invalid configuration or expose secrets in `server.yml`.
- Do not trust paths, forwarded headers, or config values without validation.

## CRITICAL - ALWAYS DO
- Keep server configuration file-based and document defaults.
- Store secrets in `secrets.yml` with owner-only permissions.
- Validate configuration before starting services and hot-reload only as specified.

## KEY DECISIONS (pre-answered)
| Question | Answer | Spec Reference |
|----------|--------|----------------|
| What governs implementation? | Server settings come from server.yml, yaml.d area files, and secrets.yml. | PART 5, 6, 12 |

## TERMINOLOGY
| Term | Meaning |
|------|---------|
| secrets.yml | generated credentials; hot reload → supported configuration changes without restart. |

## QUICK REFERENCE
- Validate before startup, keep secrets out of server.yml, and document actual defaults and reload behavior.

For complete details, see AI.md PART 5, 6, 12
