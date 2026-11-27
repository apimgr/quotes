# AI Assistant Rules (PART 0, 1)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Never treat existing code, docs, TODOs, or audit notes as the specification; `AI.md` is authoritative for implementation.
- Never edit `AI.md`; report contradictions in the designated API-template TODO file.
- Never leave unfinished stubs, fabricated product decisions, or generated-author attribution in project files.

## CRITICAL - ALWAYS DO
- Read each relevant PART completely before editing and follow every referenced PART before returning.
- Preserve user work and search for existing implementations before adding new ones.
- Keep durable project decisions in `IDEA.md`, not in assistant memory or duplicate specifications.

## KEY DECISIONS (pre-answered)
| Question | Answer | Spec Reference |
|---|---|---|
| What is authoritative? | `AI.md` defines implementation; `IDEA.md` defines product behavior. | PART 0 |
| How are spec contradictions handled? | Leave `AI.md` unchanged and record the exact conflict for template maintenance. | PART 0 |

## TERMINOLOGY
| Term | Meaning |
|---|---|
| Requirement | A normative `AI.md` rule that must be implemented. | PART 0 |

## QUICK REFERENCE
Read → trace references → implement → compare against the complete requirement. Never infer a product decision from code.

---
For complete details, see AI.md PART 0, 1
