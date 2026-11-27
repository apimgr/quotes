# Makefile Rules (PART 25)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Never use Makefile targets in CI/CD workflows.
- Never build Go binaries with a host toolchain.
- Never leave generated build output in source-controlled directories.

## CRITICAL - ALWAYS DO
- Keep local development, build, test, release, and clean targets consistent with project conventions.
- Run Go work through Docker-backed targets and keep temporary output in project-scoped OS temp directories.
- Make clean operations target only known project build artifacts.

## KEY DECISIONS (pre-answered)
| Question | Answer | Spec Reference |
|---|---|---|
| Should CI call `make test`? | No; CI runs explicit provider steps. | PART 25, 27 |
| Where does temporary output belong? | A project-scoped OS temp directory. | PART 25 |

## TERMINOLOGY
| Term | Meaning |
|---|---|
| Local target | A developer-facing Makefile command, not a CI job abstraction. | PART 25 |

## QUICK REFERENCE
Docker-backed local targets; explicit CI commands; scoped temp output; safe cleanup.

---
For complete details, see AI.md PART 25
