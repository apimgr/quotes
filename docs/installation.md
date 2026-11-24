# Installation

## Docker

Build from the repository with `docker build -f docker/Dockerfile -t quotes .`, then run `docker run --rm -p 8080:80 quotes`. The server listens on port 80 inside the container. For Compose, use `docker compose -f docker/docker-compose.yml up -d`.

## Binary

Download the server binary for your OS and architecture from the [GitHub releases page](https://github.com/apimgr/quotes/releases/latest). For example, download the Linux amd64 build with:

```sh
curl -q -LSsf -o quotes https://github.com/apimgr/quotes/releases/latest/download/quotes-linux-amd64
chmod +x quotes
./quotes --port 8080
```

The companion client is distributed as `quotes-cli` and requires an explicit server URL unless `site.txt` or a supported environment setting provides one.

## From source

The supported development workflow uses Docker and the repository Makefile. See [development](development.md) for the available targets.
