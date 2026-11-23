# Features Rules (PART 17-22)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Do not add email, notifications, scheduler, GeoIP, metrics, backups, or updater behavior contrary to its PART.
- Do not call external services without a documented trust boundary and configured enablement.
- Do not silently claim disabled features are active.

## CRITICAL - ALWAYS DO
- Implement mandatory PART behavior and document disabled optional subsystems accurately.
- Keep scheduled work internal and bounded.
- Protect backups and sensitive operational output as specified.

## Key rules summary
- SMTP configuration, automatic task scheduling, metrics exposure, backup encryption, and update verification follow their individual contracts.
- GeoIP data and attribution must follow the specified license and refresh rules.
- IDEA.md documents every integration trust assumption and failure mode.

For complete details, see AI.md PART 17-22
