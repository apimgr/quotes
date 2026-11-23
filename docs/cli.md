# CLI

`quotes-cli` retrieves quotes from a Quotes API server. Supply `--server URL`, set `QUOTES_SERVER`, or configure the CLI file; the project does not assume a default public server.

```sh
quotes-cli --server http://localhost:8080 random --category programming
quotes-cli --server http://localhost:8080 categories
quotes-cli --server http://localhost:8080 quote 42
quotes-cli --server http://localhost:8080 --output json random
```

Commands are `random [--category CATEGORY]`, `categories`, and `quote ID`. `--output` accepts `plain` or `json`. The persistent `--lang` option sends an `Accept-Language` preference; supported languages are English (`en`), Spanish (`es`), French (`fr`), Arabic (`ar`), German (`de`), and Italian (`it`). Unsupported values fall back to English. Requests are bounded, API failures are written to stderr, and no authentication is required.
