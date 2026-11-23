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

## Key rules summary
- Every registered `/server/*` page has skip links, one main, one h1, shared chrome, footer, and page shell.
- Preferences and help hubs are generated from their registries; all forms and links work without JavaScript.
- Use external embedded CSS/JS, theme tokens, semantic HTML, and accessible labels/focus states.

For complete details, see AI.md PART 16
