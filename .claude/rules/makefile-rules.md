# Makefile Rules (PART 25)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Do not run Go toolchain commands directly on the host.
- Do not use Makefile targets from CI workflows.

## CRITICAL - ALWAYS DO
- Use Docker-backed dev, local, build, and test targets.
- Keep release and clean targets scoped to project artifacts.

For complete details, see AI.md PART 25
