# Binary Rules (PART 7, 8, 32)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Do not use CGO or build only one platform.
- Do not shell out with user-controlled values.
- Do not omit the required companion CLI client.

## CRITICAL - ALWAYS DO
- Build server and client binaries as static, stripped artifacts for all eight OS/architecture pairs.
- Implement the exact CLI flags, exit codes, NO_COLOR behavior, and version output.
- Keep client requests bounded and errors useful on stderr.

## Key rules summary
- Use Docker-backed Go Makefile targets; never run Go builds directly on the host.
- Preserve binary names `quotes` and `quotes-cli`.
- Implement PART 32 client behavior, including smart context and first-run behavior where applicable.

For complete details, see AI.md PART 7, 8, 32
