# Docker Rules (PART 26)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- Do not put Dockerfiles at repository root.
- Do not copy or symlink host-built binaries into runtime images.
- Do not run the service as root or persist runtime data in the image layer.

## CRITICAL - ALWAYS DO
- Use prescribed multi-stage Dockerfiles, OCI labels, tini entrypoint, health checks, and rootfs overlay.
- Keep Compose files and host/container paths aligned with the conventions.
- Use the required image variants and never invent a custom Go toolchain image.

## Key rules summary
- Production, development, and test compose services follow their separate rules.
- Bind ports and networks using the specified project conventions.
- Container config/data volumes persist operator configuration and secrets outside the container.

For complete details, see AI.md PART 26
