# CI/CD Rules (PART 27)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Do not use actions/setup-go or Makefile targets in CI.
- Do not use floating action references or pull_request_target.

## CRITICAL - ALWAYS DO
- Pin GitHub Actions by full commit SHA.
- Set version, commit ID, and build epoch explicitly; run lint, tests, build, and security scans.

For complete details, see AI.md PART 27
