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

## Key rules summary
- Tests run in the prescribed Docker/Incus environment and follow host safety requirements.
- ReadTheDocs requires root MkDocs config, `.readthedocs.yaml`, docs pages, theme assets, and dependencies.
- Browser E2E scenarios cover SSR, no-JS, and full-JS behavior when beta testing is requested.

For complete details, see AI.md PART 28, 29, 30
