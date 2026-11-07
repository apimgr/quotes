# Binary Rules (PART 7, 8, 32)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Do not enable CGO or introduce runtime dependencies.
- Do not use hand-written argument loops for primary CLI parsing.

## CRITICAL - ALWAYS DO
- Build static server and CLI binaries with shared version metadata.
- Respect NO_COLOR and display the actual executable name.

For complete details, see AI.md PART 7, 8, 32
