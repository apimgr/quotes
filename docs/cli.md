# CLI

`quotes-cli` retrieves quotes from a Quotes API server. It has no built-in server URL, so supply `--server URL`, set `QUOTES_SERVER`, or configure the CLI file.

```sh
quotes-cli --server http://localhost:8080 random --category programming
quotes-cli --server http://localhost:8080 categories
quotes-cli --server http://localhost:8080 quote 42
quotes-cli --server http://localhost:8080 --output json random
```

Commands are `random [--category CATEGORY]`, `categories`, and `quote ID`. `--output` accepts `plain` or `json`. The CLI makes bounded HTTP requests, reports API failures on stderr, and does not require authentication.
