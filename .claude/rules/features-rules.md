# Feature Rules (PART 17–22)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Never deliver notification or email content containing secrets or unsanitized user input.
- Never run scheduled work without bounded execution, observable status, and safe shutdown behavior.
- Never enable GeoIP, metrics, backup, or update behavior beyond its configured scope.

## CRITICAL - ALWAYS DO
- Follow the specified notification-channel configuration and delivery results.
- Make scheduled tasks safe, bounded, and reportable.
- Protect metrics, backups, and update artifacts according to the relevant feature rules.
- Preserve backup integrity and verify updates before applying them.

## KEY DECISIONS (pre-answered)
| Question | Answer | Spec Reference |
|---|---|---|
| Are optional feature integrations implicitly enabled? | No; use the specified configuration defaults. | PART 17–22 |
| Can an unverified update be applied? | No. | PART 22 |

## TERMINOLOGY
| Term | Meaning |
|---|---|
| Feature integration | An optional subsystem controlled by its documented settings. | PART 17–22 |

## QUICK REFERENCE
Keep feature behavior configured, bounded, observable, verified, and safe under failure.

---
For complete details, see AI.md PART 17–22
