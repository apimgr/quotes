# Makefile Rules (PART 25)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Do not use Makefile targets in CI/CD workflows.
- Do not run Go tools directly on the host.
- Do not write build or test artifacts into tracked source paths.

## CRITICAL - ALWAYS DO
- Provide Docker-backed local development, test, build, and release targets.
- Keep Makefile responsibilities local-development only.
- Use project-scoped OS temp paths for runtime test output.

## KEY DECISIONS (pre-answered)
| Question | Answer | Spec Reference |
|----------|--------|----------------|
| What governs implementation? | Makefile targets are for local development; CI runs explicit commands without Makefile targets. | PART 25 |

## TERMINOLOGY
| Term | Meaning |
|------|---------|
| make test | Docker-backed Go tests, vet, vulnerability scan, and coverage gate. |

## QUICK REFERENCE
- Keep Docker-backed build/test workflows, scoped temp outputs, required platforms, and accurate help/target behavior.

For complete details, see AI.md PART 25
