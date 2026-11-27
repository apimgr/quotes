# Service Rules (PART 23, 24)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Never escalate privileges for ordinary requests or broader work than the single authorized operation.
- Never claim service-manager support without implementing the specified platform behavior.
- Never leave a service in an ambiguous state after install, upgrade, or removal.

## CRITICAL - ALWAYS DO
- Validate privileged operations, minimize their scope, and report failures accurately.
- Implement service lifecycle behavior and supported managers exactly as specified.
- Keep service operations idempotent and safe to retry where required.

## KEY DECISIONS (pre-answered)
| Question | Answer | Spec Reference |
|---|---|---|
| May the app elevate arbitrary commands? | No; only the narrowly defined service operations. | PART 23 |
| What managers are supported? | Only the managers explicitly listed by the spec. | PART 24 |

## TERMINOLOGY
| Term | Meaning |
|---|---|
| Privileged operation | A narrowly scoped action requiring elevated OS permissions. | PART 23 |

## QUICK REFERENCE
Minimize privilege, support only specified managers, and verify service lifecycle outcomes.

---
For complete details, see AI.md PART 23, 24
