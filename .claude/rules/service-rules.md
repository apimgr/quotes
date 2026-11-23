# Service Rules (PART 23, 24)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Do not run permanently privileged when the prescribed privilege drop is available.
- Do not install or modify host services as an implementation workaround.
- Do not claim unsupported init systems or service actions.

## CRITICAL - ALWAYS DO
- Implement the specified privilege-escalation and service lifecycle behavior.
- Ship service templates for supported platforms using the stable internal name.
- Keep install/uninstall operations explicit and safely reversible.

## Key rules summary
- Service files belong in the prescribed repository paths and use the expected runtime binary.
- PID, signal, user/group, permissions, and shutdown behavior follow PART 23.
- Systemd, runit, rc.d, launchd, and Windows support follow PART 24.

For complete details, see AI.md PART 23, 24
