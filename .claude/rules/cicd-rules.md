# CI/CD Rules (PART 27)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Do not invoke Makefile targets from CI.
- Do not use floating action references or unverified third-party actions.
- Do not add a custom Go build-image workflow.

## CRITICAL - ALWAYS DO
- Use the repository provider's supported workflows and required security gates.
- Pin actions to verified full commit SHAs and keep Renovate annotations intact.
- Build and test in the prescribed Go container across required platforms.

## KEY DECISIONS (pre-answered)
| Question | Answer | Spec Reference |
|----------|--------|----------------|
| What governs implementation? | All five supported providers ship equivalent security, test, build, and release gates. | PART 27 |

## TERMINOLOGY
| Term | Meaning |
|------|---------|
| workflow pin | verified full commit SHA or provider-required immutable digest; Renovate → only supported dependency updater. |

## QUICK REFERENCE
- Keep provider workflows fork-portable, least-privileged, secret-safe, and aligned with local checks; verify real CI after push.

For complete details, see AI.md PART 27
