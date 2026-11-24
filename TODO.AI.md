# Remaining blocked alignment work

The implementation and documentation work is substantially aligned, but these items need external input or an accepted test-gate invocation before the repository can be declared complete.

- **Approve the quote corpus.** Provide a redistribution-cleared source and attribution terms for the general, anime, Chuck Norris, dad-joke, and programming collections, or explicitly change the product scope. The current JSON files contain generated placeholder records and must not be represented as curated quotes.
- **Replace and document the approved data.** After approval, integrate the corpus with its provenance and license records, preserve stable IDs and category counts, and update the affected fixtures and public data documentation.
- **Complete full-app validation and delivery.** Obtain a test-gate-accepted invocation for `tests/incus.sh` and `tests/e2e.sh`, run both after corpus integration, resolve any failures, then prepare `COMMIT_MESS` and complete the repository gitcommit workflow.
