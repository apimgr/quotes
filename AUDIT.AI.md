# Project Audit

Started: 2026-11-24

## Pass 1: Security
- [ ] PART 12 secrets contract (AI.md 8081–8090): `loadSecrets` accepts only operator token, encryption key, metrics tokens, SMTP password, and backup password; it cannot load several required project secret paths (captcha, DNS credentials, connection credentials, CVE key, contact webhooks). This also prevents moved credentials from being applied from `secrets.yml`.

## Pass 2: Code Quality
- [ ] No further findings recorded in this pass.

## Pass 3: Logic
- [ ] PART 12 init-only configuration (AI.md 7944–7958): startup does not seed `APPLICATION_NAME`/`APPLICATION_TAGLINE` or the documented database/cache variables; current configuration types omit the database/cache settings and their required precedence behavior.
- [ ] PART 16 security overview/policy (AI.md 15806 and API contract 29905): security policy HTML/API are fixed prose, and overview only reports GitHub, optional contact form, and optional email. There is no configured PGP key discovery/fingerprint, dynamic `Expires` sharing, or common security.txt-derived representation.
- [ ] PART 16 privacy contract (AI.md 29354–29445): no `server.tracking` config or analytics provider integration exists; privacy API hardcodes `tracking` disabled and data usage text, while the HTML policy does not fully derive disclosures from cookie/data-sharing/tracking settings.
- [ ] Project data rights (AGENTS.md project rules; IDEA.md dataset scope): all five JSON collections contain generated placeholder records, not curated/public-domain records. The requested redistribution terms remain unverified for anime and other collections. Do not label these as public domain or replace them with copyrighted material without cleared provenance.

## Pass 4: Documentation
- [ ] PART 12 environment documentation (AI.md 7910–7958): docs/configuration.md describes `PORT` and `LISTEN` but omits `APPLICATION_NAME`, `APPLICATION_TAGLINE`, and currently implemented `QUOTES_SERVER_ADDRESS`/`QUOTES_SERVER_PORT`; implementation and the documented variable set need reconciliation.
- [ ] PART 16 API and browser consistency (AI.md 29839–29910): privacy, security, About, and Terms HTML/API representations remain only partly aligned; the security API contract has had recent partial changes and needs one complete comparison after integration.

## Pass 5: Spec Compliance
- [ ] PARTS 5, 12, and 16 need a complete schema-to-runtime pass after the in-progress migration and server-page work is integrated. The configuration type surface is smaller than the keys required by the configuration tables.

## Pass 6: Code Flow Trace
- [ ] The full AI.md requirement walk is incomplete; verify all 33 PARTs before declaring project-wide alignment.

## Completed
- Corrected the initial audit hypothesis that the AI.md Markdown renderer was unspecified: PART 2 explicitly approves Goldmark and Bluemonday. No TODO.API.md issue was logged for that claim.
