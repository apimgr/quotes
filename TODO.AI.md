# Remaining alignment work

- **Complete OpenRC integration verification.** The systemd Incus lifecycle, routes, Tor/I2P behavior, persistence, and clean uninstall passed. OpenRC verification could not complete because Alpine package-index fetches failed with transient DNS errors in the Incus guest. Retry when guest package networking works.
- **Finish the full AI.md alignment walk.** Trace every applicable requirement in PARTS 0–32 against code, server pages, docs, and CI; implement and verify each remaining gap. Previous coverage checks do not establish that this requirement-by-requirement pass is complete.
- **Finish delivery after the remaining blockers.** Repeat the applicable validation after corpus integration and OpenRC verification, then use the repository gitcommit workflow.
