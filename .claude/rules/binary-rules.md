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

## KEY DECISIONS (pre-answered)
| Question | Answer | Spec Reference |
|----------|--------|----------------|
| What governs implementation? | Build static CGO-disabled server and CLI binaries for all eight OS/architecture pairs. | PART 7, 8, 32 |

## TERMINOLOGY
| Term | Meaning |
|------|---------|
| server | quotes binary; client → quotes-cli companion; NO_COLOR → disable terminal color output. |

## QUICK REFERENCE
- Use Docker-backed Makefile targets and keep CLI flags, exit codes, and version behavior in sync with the spec.

For complete details, see AI.md PART 7, 8, 32
