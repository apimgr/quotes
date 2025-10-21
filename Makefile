# Makefile for Quotes API

# Project variables
PROJECTNAME := quotes
PROJECTORG := apimgr
VERSION := $(shell cat release.txt 2>/dev/null || echo "0.1.0")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

# Build variables
LDFLAGS := -X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildDate=$(BUILD_DATE) -w -s
BUILD_FLAGS := -ldflags "$(LDFLAGS)" -a -installsuffix cgo

# Output directories
BINDIR := binaries
RELEASEDIR := releases

# Platforms
PLATFORMS := \
	linux/amd64 \
	linux/arm64 \
	windows/amd64 \
	windows/arm64 \
	darwin/amd64 \
	darwin/arm64 \
	freebsd/amd64 \
	freebsd/arm64

.PHONY: all build test clean release docker docker-dev help

# Default target
all: build

# Build for all platforms using Docker (SPEC compliant)
build:
	@echo "Building $(PROJECTNAME) v$(VERSION) for all platforms..."
	@mkdir -p $(BINDIR) $(RELEASEDIR)
	@docker run --rm -v $$(pwd):/workspace -w /workspace golang:alpine sh -c ' \
		apk add --no-cache git make && \
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o binaries/$(PROJECTNAME)-linux-amd64 ./src && \
		CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(BUILD_FLAGS) -o binaries/$(PROJECTNAME)-linux-arm64 ./src && \
		CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o binaries/$(PROJECTNAME)-macos-amd64 ./src && \
		CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build $(BUILD_FLAGS) -o binaries/$(PROJECTNAME)-macos-arm64 ./src && \
		CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) -o binaries/$(PROJECTNAME)-windows-amd64.exe ./src && \
		CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build $(BUILD_FLAGS) -o binaries/$(PROJECTNAME)-windows-arm64.exe ./src && \
		CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build $(BUILD_FLAGS) -o binaries/$(PROJECTNAME)-bsd-amd64 ./src && \
		CGO_ENABLED=0 GOOS=freebsd GOARCH=arm64 go build $(BUILD_FLAGS) -o binaries/$(PROJECTNAME)-bsd-arm64 ./src'
	@echo "✓ Build complete: $(PROJECTNAME) v$(VERSION)"

# Run tests in Docker (SPEC compliant)
test:
	@echo "Running tests..."
	@docker run --rm -v $$(pwd):/workspace -w /workspace golang:alpine sh -c 'go test -v -race -timeout 5m ./src/...'
	@echo "✓ Tests passed"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BINDIR) $(RELEASEDIR)
	@rm -f $(PROJECTNAME) $(PROJECTNAME).exe
	@echo "✓ Clean complete"

# Create release artifacts
release: build
	@echo "Creating release artifacts..."
	@mkdir -p $(RELEASEDIR)
	@$(foreach platform,$(PLATFORMS), \
		tar -czf $(RELEASEDIR)/$(PROJECTNAME)-$(VERSION)-$(word 1,$(subst /, ,$(platform)))-$(word 2,$(subst /, ,$(platform))).tar.gz \
		-C $(BINDIR) $(PROJECTNAME)-$(word 1,$(subst /, ,$(platform)))-$(word 2,$(subst /, ,$(platform)))$(if $(findstring windows,$(platform)),.exe,) && \
		echo "✓ Created $(PROJECTNAME)-$(VERSION)-$(word 1,$(subst /, ,$(platform)))-$(word 2,$(subst /, ,$(platform))).tar.gz" ; \
	)
	@echo "✓ Release artifacts created in $(RELEASEDIR)/"

# Build and push multi-platform Docker images (release)
docker:
	@echo "Building multi-platform Docker images..."
	@docker buildx build \
		--platform linux/amd64,linux/arm64 \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		-t ghcr.io/$(PROJECTORG)/$(PROJECTNAME):latest \
		-t ghcr.io/$(PROJECTORG)/$(PROJECTNAME):$(VERSION) \
		--push \
		.
	@echo "✓ Docker images pushed to ghcr.io/$(PROJECTORG)/$(PROJECTNAME):$(VERSION)"

# Build Docker image for development (local only, not pushed)
docker-dev:
	@echo "Building development Docker image..."
	@docker build \
		--build-arg VERSION=$(VERSION)-dev \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		-t $(PROJECTNAME):dev \
		.
	@echo "✓ Docker development image built: $(PROJECTNAME):dev"

# Show help
help:
	@echo "Quotes API - Build System"
	@echo ""
	@echo "Usage:"
	@echo "  make build       - Build binaries for all platforms"
	@echo "  make test        - Run tests"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make release     - Create release artifacts"
	@echo "  make docker      - Build and push multi-platform Docker images"
	@echo "  make docker-dev  - Build development Docker image (local only)"
	@echo "  make help        - Show this help message"
	@echo ""
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Build Date: $(BUILD_DATE)"
