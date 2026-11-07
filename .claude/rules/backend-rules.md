# Backend Rules (PART 9, 10, 11, 31)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Do not persist quote data or visitor PII.
- Do not trust forwarded client IP headers unless a trusted proxy policy is implemented.

## CRITICAL - ALWAYS DO
- Bound requests and apply configured per-client rate limits.
- Use parameterized database access if a future spec adds a database.

For complete details, see AI.md PART 9, 10, 11, 31
