# Quotes Project Loader

## Source of Truth

- `AI.md` is the read-only implementation specification. Read the applicable PART before changing code, pages, workflows, or documentation.
- `IDEA.md` defines Quotes product behavior and project variables; it does not override `AI.md`.
- Read the matching `.claude/rules/*.md` and project TODO instructions before starting related work.

## Before Any Code Change

1. Read the relevant AI.md PART and existing source/tests.
2. Check the implementation against the spec; do not guess when the spec is silent.
3. Preserve existing worktree changes and search for reusable code before adding more.
4. Verify the affected behavior and keep browser, API, CLI, and documentation aligned.

## Critical Rules

### NEVER

- Modify `AI.md`, invent quote provenance, add user accounts or community quote submissions, or leave stubs and unfinished TODOs in production code.
- Run Go builds, tests, or binaries directly on the host; use the Docker-backed Makefile targets or specified Docker/Incus test scripts.
- Use Makefile targets in CI; workflows use explicit commands.

### ALWAYS

- Keep the quote corpus original and project-owned, follow the eight-platform static Go build matrix, and preserve the public read-only API behavior.
- Keep server-rendered pages functional without JavaScript and synchronize all public documentation and API descriptions.
- Run required test and lint gates before committing; use `gitcommit --dir {project_dir} all` with an accurate `.git/COMMIT_MESS`.

Full rules are in `.claude/rules/`; current tasks are tracked in `TODO.AI.md`. `AI.md` remains authoritative.
