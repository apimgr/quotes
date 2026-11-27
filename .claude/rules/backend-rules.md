# Backend Rules (PART 9, 10, 11, 31)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Never leak internal errors, credentials, personal data, or cryptographic material.
- Never use unbounded request bodies, unsafe SQL construction, or unvalidated external input.
- Never enable an overlay service or privileged operation without its documented controls.

## CRITICAL - ALWAYS DO
- Return the required status, response envelope, cache behavior, and safe client-facing error.
- Use database transactions and migrations according to the persistence rules.
- Apply security headers, structured logging, secret redaction, and rate limiting as specified.
- Treat Tor and I2P as optional, explicitly configured integrations with safe failure behavior.

## KEY DECISIONS (pre-answered)
| Question | Answer | Spec Reference |
|---|---|---|
| May internal errors reach clients? | No; log safely and return a sanitized response. | PART 9, 11 |
| Are overlay networks enabled by default? | No; follow the explicit configuration rules. | PART 31 |

## TERMINOLOGY
| Term | Meaning |
|---|---|
| Cache | Bounded, expiring storage governed by configured policy. | PART 9 |
| Secret | Any credential or key that must be redacted from output. | PART 11 |

## QUICK REFERENCE
Validate inputs, use safe persistence, sanitize errors and logs, and keep optional integrations opt-in.

---
For complete details, see AI.md PART 9, 10, 11, 31
