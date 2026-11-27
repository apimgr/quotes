# Testing Rules (PART 28, 29, 30)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Never treat a successful build as proof that the application is complete.
- Never place coverage or test artifacts in the repository.
- Never claim browser, accessibility, localization, or documentation checks that were not run.

## CRITICAL - ALWAYS DO
- Cover public API, public web, server API, error responses, and edge cases.
- Collect coverage with `-test.gocoverdir` and convert using `go tool covdata textfmt`; do not use `-coverprofile`.
- Keep documentation buildable and strict, and test localization and accessibility requirements.
- Keep validation commands and results explicit in delivery notes.

## KEY DECISIONS (pre-answered)
| Question | Answer | Spec Reference |
|---|---|---|
| Does build success establish completion? | No; all specified product behavior must be implemented and verified. | PART 28 |
| How is coverage collected? | `-test.gocoverdir` followed by `go tool covdata textfmt`. | PART 28 |

## TERMINOLOGY
| Term | Meaning |
|---|---|
| Coverage data | Per-test-process files written to the project-scoped temp coverage directory. | PART 28 |

## QUICK REFERENCE
Test all required surfaces and error cases; keep test output out of the repo; build docs strictly; verify I18N and accessibility.

---
For complete details, see AI.md PART 28, 29, 30
