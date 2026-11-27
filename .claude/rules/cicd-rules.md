# CI/CD Rules (PART 27)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Never call Makefile targets from CI workflows.
- Never expose secrets in workflow output or grant broader permissions than a job needs.
- Never bypass the prescribed workflow validation, artifact checks, or security gates.

## CRITICAL - ALWAYS DO
- Keep all required provider workflows synchronized with the same build, test, lint, security, and release requirements.
- Pin third-party workflow actions as specified and grant least privilege.
- Verify checksums and signatures for downloaded or published artifacts where required.

## KEY DECISIONS (pre-answered)
| Question | Answer | Spec Reference |
|---|---|---|
| How does CI build/test? | Explicit provider commands, never `make` targets. | PART 27 |
| How are actions referenced? | Use the required immutable references. | PART 27 |

## TERMINOLOGY
| Term | Meaning |
|---|---|
| Workflow gate | A required check that must pass before release or merge. | PART 27 |

## QUICK REFERENCE
Use least privilege, pinned actions, explicit commands, complete provider parity, and verified release artifacts.

---
For complete details, see AI.md PART 27
