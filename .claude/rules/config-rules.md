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

## Key rules summary
- `server.yml`, drop-ins, and `secrets.yml` are the configuration source of truth.
- Separate secret values from ordinary configuration; first-run credentials are generated securely.
- Apply path safety, boolean parsing, environment precedence, and mode behavior from their PARTs.

For complete details, see AI.md PART 5, 6, 12
