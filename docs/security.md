# Security

Quotes is a public, read-only service. It has no accounts, passwords, API keys, user submissions, or privileged web interface. Incoming parameters are validated, request URIs and bodies are size-limited, lists are paginated, and per-client rate limits are enabled by default. CORS is open because the dataset is public.

`/server/healthz` supports HTML, plain text, and JSON based on `Accept`. `/api/v1/server/healthz` and `/api/healthz` return machine-readable health details. `/.well-known/security.txt` and `/security.txt` provide security contact guidance. The service does not forward client-supplied proxy headers when determining the rate-limit address.

The bundled JSON files are part of the executable and should only contain content the project is authorized to distribute.
