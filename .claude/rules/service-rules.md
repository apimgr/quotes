# Service Rules (PART 23, 24)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Do not run permanently privileged when the prescribed privilege drop is available.
- Do not install or modify host services as an implementation workaround.
- Do not claim unsupported init systems or service actions.

## CRITICAL - ALWAYS DO
- Implement the specified privilege-escalation and service lifecycle behavior.
- Ship service templates for supported platforms using the stable internal name.
- Keep install/uninstall operations explicit and safely reversible.

## KEY DECISIONS (pre-answered)
| Question | Answer | Spec Reference |
|----------|--------|----------------|
| What governs implementation? | Service support uses stable internal identifiers and explicit, reversible lifecycle operations. | PART 23, 24 |

## TERMINOLOGY
| Term | Meaning |
|------|---------|
| service install | platform-native service registration; privilege drop → least-privilege runtime identity. |

## QUICK REFERENCE
- Support only documented init systems and platform lifecycle behavior; never use host changes as a workaround.

For complete details, see AI.md PART 23, 24
