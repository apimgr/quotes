# Binary Rules (PART 7, 8, 32)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Never build project Go binaries directly on the host.
- Never emit secrets, uncontrolled diagnostics, or inconsistent output from CLI commands.
- Never omit a required platform or CLI mode from the prescribed build matrix.

## CRITICAL - ALWAYS DO
- Apply the required static binary, naming, version, and linker conventions.
- Keep server and client CLI flags, help, shell completion, and output behavior synchronized.
- Run builds and binary checks in the required Docker-backed environment.

## KEY DECISIONS (pre-answered)
| Question | Answer | Spec Reference |
|---|---|---|
| Where are Go binaries built? | In the Docker-backed toolchain. | PART 7 |
| Which client formats are supported? | Only the formats and behaviors specified for the CLI. | PART 32 |

## TERMINOLOGY
| Term | Meaning |
|---|---|
| Server binary | The Quotes API service executable. | PART 8 |
| Client binary | The companion command-line client. | PART 32 |

## QUICK REFERENCE
Use the prescribed names, version metadata, platform matrix, flags, output modes, and Docker build targets.

---
For complete details, see AI.md PART 7, 8, 32
