# Project Rules (PART 2, 3, 4)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Never copy third-party data or assets without verified redistribution rights and required attribution.
- Never put runtime data, generated output, secrets, or host-specific workarounds in source control.
- Never rename or add project paths contrary to the prescribed directory layout.

## CRITICAL - ALWAYS DO
- Preserve the project license and provenance for all bundled quote content.
- Keep code in the prescribed singular package directories; `src/paths/` is an explicitly named layout exception.
- Follow the OS-specific path and platform rules for runtime files.

## KEY DECISIONS (pre-answered)
| Question | Answer | Spec Reference |
|---|---|---|
| May bundled material have unknown rights? | No; verify ownership and redistribution terms first. | PART 2 |
| Where do runtime files belong? | In the OS-appropriate application directories, not the repository. | PART 4 |

## TERMINOLOGY
| Term | Meaning |
|---|---|
| Project file | Source-controlled code, configuration, or documentation. | PART 3 |
| Runtime data | Mutable data created or consumed while the application runs. | PART 4 |

## QUICK REFERENCE
Keep project layout exact, keep the quote corpus original and attributable, and use platform-specific runtime paths.

---
For complete details, see AI.md PART 2, 3, 4
