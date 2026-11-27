# Frontend Rules (PART 16)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Never create independent page shells or duplicate shared layout markup.
- Never put business logic in templates or render configured Markdown as trusted HTML.
- Never make required server pages depend on JavaScript for core behavior.

## CRITICAL - ALWAYS DO
- Use `layout/public.tmpl`, shared head/scripts/header/navigation/footer partials, and one page template per route.
- Keep all `/server/*` pages consistent, accessible, responsive, localized, and usable without JavaScript.
- Render Markdown with the prescribed Goldmark and bluemonday pipeline; raw HTML remains disabled.
- Keep browser pages, browser API calls, and their documented behavior aligned.

## KEY DECISIONS (pre-answered)
| Question | Answer | Spec Reference |
|---|---|---|
| Where does page structure live? | Shared layouts and partials; page templates contain page content. | PART 16 |
| Is raw HTML allowed in configured Markdown? | No. | PART 16 |

## TERMINOLOGY
| Term | Meaning |
|---|---|
| Page template | Route-specific content rendered inside the shared shell. | PART 16 |
| Partial | Reusable self-contained UI template component. | PART 16 |

## QUICK REFERENCE
Use the common page shell, safe Markdown renderer, minimal template logic, and accessible no-JS behavior.

---
For complete details, see AI.md PART 16
