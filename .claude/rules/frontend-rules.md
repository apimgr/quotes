# Frontend Rules (PART 16)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Do not make API-backed browsing dependent on client-side rendering.
- Do not insert unescaped user input into HTML.
- Do not use inline CSS, inline event handlers, or JavaScript for native HTML behavior.

## CRITICAL - ALWAYS DO
- Use server-rendered, mobile-friendly, accessible pages with shared shell and components.
- Keep theme selection, pagination, and filtering usable without JavaScript.
- Register pages centrally and make every error page conform to the shell.

## KEY DECISIONS (pre-answered)
| Question | Answer | Spec Reference |
|----------|--------|----------------|
| What governs implementation? | Core browsing and forms must work without JavaScript using server-rendered pages. | PART 16 |

## TERMINOLOGY
| Term | Meaning |
|------|---------|
| SSR | server-side rendering; page shell → shared accessible layout and navigation. |

## QUICK REFERENCE
- Use mobile-first accessible HTML, shared page chrome, external assets, and JavaScript only for capabilities HTML/CSS cannot provide.

For complete details, see AI.md PART 16
