# Security Policy

## Supported versions

Security fixes are published for the latest stable release. Use the newest
release available from [GitHub Releases](https://github.com/apimgr/quotes/releases).

## Report a source-code vulnerability

Use [GitHub private vulnerability reporting](https://github.com/apimgr/quotes/security/advisories/new),
available from the repository's **Security** tab under **Report a vulnerability**.
Do not report a vulnerability in a public issue or pull request. Mirrors should
direct researchers to the GitHub origin repository.

For a deployed Quotes API instance, use its `/.well-known/security.txt` or
`/server/contact?security_id=...` security-report flow. The [security guide](https://apimgr-quotes.readthedocs.io/security/)
describes the `/server/security` page and the deployed-service reporting paths.

## Response and disclosure

Maintainers aim to acknowledge reports within 48 hours and coordinate a fix and
release with the reporter. Critical issues are prioritized for a fix within 14
days when practical. Please allow maintainers to investigate and publish a fix
before disclosing technical details. If a report is out of scope or cannot be
reproduced, the reporter will be told why.
