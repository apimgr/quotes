PROJECTNAME := $(shell git remote get-url origin 2>/dev/null | sed -E 's|.*/([^/]+)$$|\1|' | sed -E 's|\.git$$||' || basename "$$(pwd)")
PROJECTORG := $(shell git remote get-url origin 2>/dev/null | sed -E 's|.*/([^/]+)/[^/]+$$|\1|' || basename "$$(dirname "$$(pwd)")")
INTERNAL_NAME := $(PROJECTNAME)
VERSION := $(shell cat release.txt 2>/dev/null || echo "0.1.0")
BUILD_EPOCH := $(shell date -u +%s)
BUILD_DATE := $(shell date -u -d @$(BUILD_EPOCH) +"%Y-%m-%dT%H:%M:%SZ")
COMMIT_ID := $(shell git rev-parse --short=7 HEAD 2>/dev/null || echo N/A)
OFFICIAL_SITE := $(shell [ -f site.txt ] && cat site.txt || echo "$${OFFICIAL_SITE:-}")
LDFLAGS := -s -w -X 'main.Version=$(VERSION)' -X 'main.CommitID=$(COMMIT_ID)' -X 'main.BuildDate=$(BUILD_DATE)' -X 'main.BuildEpoch=$(BUILD_EPOCH)' -X 'main.OfficialSite=$(OFFICIAL_SITE)'
BINDIR := binaries
RELDIR := releases
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64 windows/arm64 freebsd/amd64 freebsd/arm64
REGISTRY ?= ghcr.io/$(PROJECTORG)/$(PROJECTNAME)
GO_CACHE ?= $(HOME)/go/pkg/mod
GO_BUILD ?= $(HOME)/.cache/go-build/$(PROJECTNAME)
DOCKER_MEM ?= 4g
DOCKER_CPUS ?= 2
GO_DOCKER_RUN := docker run --rm --name $(PROJECTNAME)-$$(tr -dc 'a-z0-9' </dev/urandom | head -c8) --memory=$(DOCKER_MEM) --cpus=$(DOCKER_CPUS) -v $(PWD):/app -v $(GO_CACHE):/usr/local/share/go/pkg/mod -v $(GO_BUILD):/usr/local/share/go/cache -w /app -e CGO_ENABLED=0 -e GOFLAGS=-buildvcs=false
GO_DOCKER := $(GO_DOCKER_RUN) casjaysdev/go:latest

.PHONY: dev local build test release docker clean

dev:
	@mkdir -p "$(GO_CACHE)" "$(GO_BUILD)" "$${TMPDIR:-/tmp}/$(PROJECTORG)"
	@$(GO_DOCKER) go mod tidy
	@BUILD_DIR=$$(mktemp -d "$${TMPDIR:-/tmp}/$(PROJECTORG)/$(INTERNAL_NAME)-XXXXXX"); mkdir -p "$$BUILD_DIR"; docker run --rm --name $(PROJECTNAME)-$$(tr -dc 'a-z0-9' </dev/urandom | head -c8) --memory=$(DOCKER_MEM) --cpus=$(DOCKER_CPUS) -v "$(PWD):/app" -v "$$BUILD_DIR:/out" -v "$(GO_CACHE):/usr/local/share/go/pkg/mod" -v "$(GO_BUILD):/usr/local/share/go/cache" -w /app -e CGO_ENABLED=0 -e GOFLAGS=-buildvcs=false casjaysdev/go:latest sh -c "go build -buildvcs=false -o /out/$(PROJECTNAME) ./src && go build -buildvcs=false -o /out/$(PROJECTNAME)-cli ./src/client"

local: clean
	@mkdir -p $(GO_CACHE) $(GO_BUILD)
	@mkdir -p $(BINDIR)
	@$(GO_DOCKER) go mod tidy
	@$(GO_DOCKER) sh -c "go build -buildvcs=false -trimpath -ldflags \"$(LDFLAGS)\" -o $(BINDIR)/$(PROJECTNAME) ./src && go build -buildvcs=false -trimpath -ldflags \"-X 'main.Version=$(VERSION)' -X 'main.CommitID=$(COMMIT_ID)' -X 'main.BuildDate=$(BUILD_DATE)' -X 'main.BuildEpoch=$(BUILD_EPOCH)'\" -o $(BINDIR)/$(PROJECTNAME)-cli ./src/client"

build: clean
	@mkdir -p $(GO_CACHE) $(GO_BUILD)
	@mkdir -p $(BINDIR)
	@$(GO_DOCKER) go mod tidy
	@$(GO_DOCKER) go mod download
	@for platform in $(PLATFORMS); do os=$${platform%/*}; arch=$${platform#*/}; ext=; [ "$$os" = windows ] && ext=.exe; $(GO_DOCKER_RUN) -e GOOS=$$os -e GOARCH=$$arch casjaysdev/go:latest go build -buildvcs=false -trimpath -ldflags "$(LDFLAGS)" -o $(BINDIR)/$(PROJECTNAME)-$$os-$$arch$$ext ./src || exit 1; $(GO_DOCKER_RUN) -e GOOS=$$os -e GOARCH=$$arch casjaysdev/go:latest go build -buildvcs=false -trimpath -ldflags "-X 'main.Version=$(VERSION)' -X 'main.CommitID=$(COMMIT_ID)' -X 'main.BuildDate=$(BUILD_DATE)' -X 'main.BuildEpoch=$(BUILD_EPOCH)' -X 'main.OfficialSite=$(OFFICIAL_SITE)'" -o $(BINDIR)/$(PROJECTNAME)-cli-$$os-$$arch$$ext ./src/client || exit 1; done

test:
	@mkdir -p "$(GO_CACHE)" "$(GO_BUILD)" "$${TMPDIR:-/tmp}/$(PROJECTORG)"
	@$(GO_DOCKER) sh -c 'mkdir -p "$${TMPDIR:-/tmp}/$(PROJECTORG)" && covdir=$$(mktemp -d "$${TMPDIR:-/tmp}/$(PROJECTORG)/$(PROJECTNAME)-XXXXXX") && trap "rm -rf $$covdir" EXIT && go vet ./... && govulncheck ./... && go test -coverprofile=$$covdir/coverage.out ./... && coverage=$$(go tool cover -func=$$covdir/coverage.out | awk "/^total:/ {gsub(/%/, \"\", \$$3); print \$$3}") && echo "Coverage: $$coverage%" && awk -v coverage="$$coverage" "BEGIN { if (coverage+0 < 60) exit 1 }"'

release: build
	@mkdir -p $(RELDIR)
	@printf '%s\n' '$(VERSION)' > $(RELDIR)/version.txt
	@tar --exclude=.git --exclude=.github --exclude=binaries --exclude=releases -czf $(RELDIR)/$(PROJECTNAME)-$(VERSION)-source.tar.gz .
	@cd $(RELDIR) && sha256sum $(PROJECTNAME)-$(VERSION)-source.tar.gz version.txt > sha256.txt && sha512sum $(PROJECTNAME)-$(VERSION)-source.tar.gz version.txt > sha512.txt

docker:
	@docker buildx build -f docker/Dockerfile --platform linux/amd64,linux/arm64 --build-arg VERSION=$(VERSION) --build-arg COMMIT_ID=$(COMMIT_ID) --build-arg BUILD_EPOCH=$(BUILD_EPOCH) --build-arg BUILD_DATE=$(BUILD_DATE) --build-arg OFFICIAL_SITE='$(OFFICIAL_SITE)' -t $(REGISTRY):$(VERSION) -t $(REGISTRY):latest --push .

clean:
	@rm -rf $(BINDIR) $(RELDIR)
