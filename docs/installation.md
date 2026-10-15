# Installation

## Docker

Build from the repository with `docker build -f docker/Dockerfile -t quotes .`, then run `docker run --rm -p 8080:80 quotes`. The server listens on port 80 inside the container. For Compose, use `docker compose -f docker/docker-compose.yml up -d`.

## Binary

Download the server binary for your OS and architecture from the project release. Run `./quotes --port 8080`. The companion client is distributed as `quotes-cli` and requires an explicit server URL unless `site.txt` or a supported environment setting provides one.

## From source

The supported development workflow uses Docker and the repository Makefile. See [development](development.md) for the available targets.
