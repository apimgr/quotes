# Backend Rules (PART 9, 10, 11, 31)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Do not persist quote data or visitor PII outside the defined report and operational workflows.
- Do not trust forwarded client IP headers without a trusted-proxy policy.
- Do not expose secrets, raw stack traces, or unsafe user-controlled content.

## CRITICAL - ALWAYS DO
- Bound requests, apply configured per-client rate limits, and use canonical errors.
- Use parameterized database operations, connection pooling, and timeouts.
- Implement security logging, token checks, and overlay-network rules as prescribed.

## KEY DECISIONS (pre-answered)
| Question | Answer | Spec Reference |
|----------|--------|----------------|
| What governs implementation? | Public quote reads stay unauthenticated; credentialed surfaces follow their specific token rules. | PART 9, 10, 11, 31 |

## TERMINOLOGY
| Term | Meaning |
|------|---------|
| owner token | anonymous credential-management identity; overlay → Tor hidden service or I2P eepsite. |

## QUICK REFERENCE
- Parameterize database queries, bound work and inputs, protect secrets, and apply the prescribed overlay lifecycle.

For complete details, see AI.md PART 9, 10, 11, 31
