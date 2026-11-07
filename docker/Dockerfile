FROM casjaysdev/go:latest AS builder

RUN apk add --no-cache bash git

ARG TARGETARCH
ARG VERSION=dev
ARG BUILD_EPOCH=0
ARG COMMIT_ID=unknown
ARG OFFICIAL_SITE=

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN mkdir -p /app/binary && \
    CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -buildvcs=false -trimpath \
      -ldflags "-s -w -X 'main.Version=${VERSION}' -X 'main.CommitID=${COMMIT_ID}' -X 'main.BuildEpoch=${BUILD_EPOCH}' -X 'main.OfficialSite=${OFFICIAL_SITE}'" \
      -o /app/binary/quotes ./src && \
    CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -buildvcs=false -trimpath \
      -ldflags "-s -w -X 'main.Version=${VERSION}' -X 'main.CommitID=${COMMIT_ID}' -X 'main.BuildEpoch=${BUILD_EPOCH}'" \
      -o /app/binary/quotes-cli ./src/client

FROM alpine:latest

ARG BUILD_EPOCH=0
ARG COMMIT_ID=unknown
ARG LICENSE=MIT

RUN apk add --no-cache tini
COPY --from=builder /app/binary/quotes /usr/local/bin/quotes
COPY --from=builder /app/binary/quotes-cli /usr/local/bin/quotes-cli
COPY docker/rootfs/ /
COPY docker/Dockerfile /root/Dockerfile
RUN chmod 755 /usr/local/bin/*

ENV TZ=America/New_York

EXPOSE 80
STOPSIGNAL SIGRTMIN+3
HEALTHCHECK --start-period=90s --interval=10s --timeout=5s --retries=3 \
    CMD wget -q -O /dev/null http://127.0.0.1:80/server/healthz || exit 1
ENTRYPOINT [ "tini", "-p", "SIGTERM", "--", "/usr/local/bin/entrypoint.sh" ]
