# Features Rules (PART 17-22)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Do not add email, notifications, scheduler, GeoIP, metrics, backups, or updater behavior contrary to its PART.
- Do not call external services without a documented trust boundary and configured enablement.
- Do not silently claim disabled features are active.

## CRITICAL - ALWAYS DO
- Implement mandatory PART behavior and document disabled optional subsystems accurately.
- Keep scheduled work internal and bounded.
- Protect backups and sensitive operational output as specified.

## KEY DECISIONS (pre-answered)
| Question | Answer | Spec Reference |
|----------|--------|----------------|
| What governs implementation? | Required built-in features follow their individual PART contracts and documented failure behavior. | PART 17, 18, 19, 20, 21, 22 |

## TERMINOLOGY
| Term | Meaning |
|------|---------|
| scheduler | internal persisted task runner; GeoIP → deny-only risk signal; backup → verified archive. |

## QUICK REFERENCE
- Keep tasks bounded and persisted, protect backup data, and reflect feature health and configuration accurately.

For complete details, see AI.md PART 17-22
