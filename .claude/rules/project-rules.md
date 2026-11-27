# Project Rules (PART 2, 3, 4)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Do not copy quote datasets without clear provenance and redistribution permission.
- Do not put runtime data, credentials, binaries, or logs in the project tree.
- Do not change stable internal names without an explicit rename and migration.

## CRITICAL - ALWAYS DO
- Keep assets embedded and licenses/attribution explicit.
- Use the prescribed project structure and OS-specific paths.
- Keep source, config, data, cache, log, and backup paths distinct.

## KEY DECISIONS (pre-answered)
| Question | Answer | Spec Reference |
|----------|--------|----------------|
| What governs implementation? | Original quote content is project-owned; preserve its provenance and attribution boundary. | PART 2, 3, 4 |

## TERMINOLOGY
| Term | Meaning |
|------|---------|
| internal_name | stable on-disk identifier; embedded data → shipped static assets; runtime data → OS data directory. |

## QUICK REFERENCE
- Keep source and runtime paths separate, preserve stable names, and include accurate third-party license attribution.

For complete details, see AI.md PART 2, 3, 4
