# Contributing

Quotes is a public, read-only quote service. It does not accept user-submitted
quotes or create accounts. Please keep changes aligned with [AI.md](../AI.md),
[IDEA.md](../IDEA.md), and the project documentation.

## Local setup and checks

Use the Docker-backed Makefile targets for Go work:

```sh
make build
make test
```

The manual integration suites are separate from `make test`:

```sh
bash tests/run_tests.sh
bash tests/incus.sh
./tests/e2e.sh
```

The documentation source is in `docs/`; its MkDocs configuration is at the
repository root in `mkdocs.yml`. See the [development guide](https://apimgr-quotes.readthedocs.io/development/)
for project workflows.

## Changes and review

Contribute through a fork and pull request targeting `main`. Include the reason
for the change, evidence from relevant checks, and any compatibility impact.
Update tests and documentation when behavior or configuration changes. Do not
add generic placeholder content or quote material without documented
redistribution rights and provenance.

Do not include credentials or private data in issues, logs, fixtures, or pull
requests. For a source-code vulnerability, use the private reporting path in
[SECURITY.md](SECURITY.md); do not open a public issue.
