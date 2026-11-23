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

## Key rules summary
- Include required cache, concurrency, timeout, cleanup, and coverage behavior.
- Keep test, lint, build, and CI commands explicit and auditable.
- Preserve every required target and platform in the project matrix.

For complete details, see AI.md PART 25
