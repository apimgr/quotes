# Configuration Rules (PART 5, 6, 12)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Never make invalid configuration booleans stop server startup when the spec requires defaulting.
- Never expose secrets through configuration output, diagnostics, or logs.
- Never silently invent configuration keys or precedence.

## CRITICAL - ALWAYS DO
- Use the shared boolean parsers: configuration input uses `ParseBoolOrDefault`; request input uses `ParseBool` and validation errors.
- Apply the specified configuration precedence and safe defaults.
- Keep debug routes absent unless debug is enabled, and require the operator token on every debug route.

## KEY DECISIONS (pre-answered)
| Question | Answer | Spec Reference |
|---|---|---|
| What happens to invalid configuration booleans? | Warn once and use the configured default. | PART 5 |
| Does debug mode bypass debug authentication? | No; every debug endpoint requires the operator token. | PART 6 |

## TERMINOLOGY
| Term | Meaning |
|---|---|
| Configuration input | Environment, config file, or startup flag. | PART 5 |
| Request input | User-provided API, form, or JSON value. | PART 5 |

## QUICK REFERENCE
Use shared config parsing, keep secrets sanitized, and mount diagnostics only behind the debug switch plus operator authentication.

---
For complete details, see AI.md PART 5, 6, 12
