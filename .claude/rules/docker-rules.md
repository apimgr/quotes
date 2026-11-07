# Docker Rules (PART 26)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Do not add unnecessary packages or services to the runtime image.
- Do not publish the production port on every host interface.

## CRITICAL - ALWAYS DO
- Build static binaries in casjaysdev/go and use a minimal runtime image.
- Keep compose files under docker/ with explicit health checks, logging, and project networks.

For complete details, see AI.md PART 26
