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

## KEY DECISIONS (pre-answered)
| Question | Answer | Spec Reference |
|----------|--------|----------------|
| What governs implementation? | The local docker target loads a host-platform image; CI owns multi-platform image publication. | PART 26 |

## TERMINOLOGY
| Term | Meaning |
|------|---------|
| runtime image | shipped service container; OCI metadata → build-time labels/annotations required by the platform. |

## QUICK REFERENCE
- Use the prescribed Dockerfiles, non-root runtime, tini, health checks, Compose conventions, and image contract checks.

For complete details, see AI.md PART 26
