# Docker Rules (PART 26)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Never run containers as privileged or expose unnecessary ports, mounts, or capabilities.
- Never bake credentials or mutable runtime data into an image.
- Never use floating, unsafe container behavior where the spec requires pinned or constrained configuration.

## CRITICAL - ALWAYS DO
- Use the specified base image, build/runtime separation, non-root behavior, and health checks.
- Keep image entrypoints and compose services aligned with the server CLI and configuration.
- Use explicit volumes for persistent data and preserve safe defaults.

## KEY DECISIONS (pre-answered)
| Question | Answer | Spec Reference |
|---|---|---|
| Are credentials part of the image? | No. | PART 26 |
| Where does persistent application state live? | In explicitly configured volumes. | PART 26 |

## TERMINOLOGY
| Term | Meaning |
|---|---|
| Runtime volume | A mounted persistent location for mutable application state. | PART 26 |

## QUICK REFERENCE
Least privilege, explicit persistence, health checks, and no embedded secrets.

---
For complete details, see AI.md PART 26
