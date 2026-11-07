# Project SPEC

Project: quotes
Role: Efficient loader for AI.md

## Source of truth
- Read `IDEA.md` for Quotes product behavior and project variables.
- Read `AI.md` for implementation rules; it is read-only.
- Read `SPEC.md` first when present; a non-empty SPEC.md overrides AI.md.
- Read relevant `.claude/rules/*.md` before changing matching code.

## Project rules
- This is a public, read-only quote API with no accounts, submissions, or operator web panel.
- Preserve the rights and provenance boundary for bundled quote content. Do not import datasets without clear redistribution permission.
- Keep API, GraphQL, browser, CLI, and documentation behavior aligned.
- Use Docker-backed Makefile targets for Go development and leave validation to the primary session.
- Keep `TODO.AI.md` limited to genuinely unfinished or blocked work.

## Session workflow
1. Read relevant source and tests before editing.
2. Follow the matching AI.md PART and any linked requirements.
3. Preserve existing worktree changes.
4. Do not commit or push from an implementation subtask.
