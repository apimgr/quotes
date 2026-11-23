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

## Key rules summary
- CI, release, beta, daily, and Docker workflows have distinct triggers and responsibilities.
- Keep permissions minimal, use concurrency controls, and avoid exposing secrets in logs.
- CI commands are explicit; container images follow the project-declared toolchain.

For complete details, see AI.md PART 27
