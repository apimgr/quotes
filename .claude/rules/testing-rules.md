# Testing Rules (PART 28, 29, 30)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Do not lower coverage gates to hide missing tests.
- Do not claim validation that was not run.
- Do not put runtime artifacts or coverage output in the project tree.

## CRITICAL - ALWAYS DO
- Keep Go tests beside the source and exercise API edge cases.
- Keep ReadTheDocs navigation, browser content, CLI, API, and accessibility documentation aligned.
- Keep locale files in parity and follow the prescribed translation workflow.

## KEY DECISIONS (pre-answered)
| Question | Answer | Spec Reference |
|----------|--------|----------------|
| What governs implementation? | The Go coverage gate is 60%; browser E2E covers SSR, no-JS, and full-JS behavior. | PART 28, 29, 30 |

## TERMINOLOGY
| Term | Meaning |
|------|---------|
| Phase 1 | Docker toolchain unit-test gate; Phase 2 → explicit binary validation in Docker or Incus. |

## QUICK REFERENCE
- Keep tests beside Go source, docs complete and buildable, locales synchronized, and test outputs in scoped temp paths.

For complete details, see AI.md PART 28, 29, 30
