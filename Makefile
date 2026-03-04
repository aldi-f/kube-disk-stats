.PHONY: build clean test lint install release

APP_NAME=kube-disk-stats
DOCKER_IMAGE=ghcr.io/aldi-f/kube-disk-stats
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags="-s -w -X main.Version=$(VERSION)"

GO=go
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)

DIST_DIR=dist

build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(DIST_DIR)
	$(GO) build $(LDFLAGS) -o $(DIST_DIR)/$(APP_NAME) .

build-all:
	@echo "Building $(APP_NAME) for all platforms..."
	@mkdir -p $(DIST_DIR)
	@for os in linux darwin windows; do \
		for arch in amd64 arm64; do \
			if [ "$$os" = "windows" ] && [ "$$arch" = "arm64" ]; then \
				continue; \
			fi; \
			output=$(DIST_DIR)/$(APP_NAME)-$$os-$$arch; \
			if [ "$$os" = "windows" ]; then \
				output=$$output.exe; \
			fi; \
			GOOS=$$os GOARCH=$$arch $(GO) build $(LDFLAGS) -o $$output .; \
			echo "Built $$output"; \
		done; \
	done

clean:
	@echo "Cleaning..."
	@rm -rf $(DIST_DIR)

test:
	@echo "Running tests..."
	$(GO) test -v -race -cover ./...

lint:
	@echo "Running linters..."
	golangci-lint run

install: build
	@echo "Installing $(APP_NAME)..."
	$(GO) install $(LDFLAGS) .

release: clean build-all
	@echo "Creating release packages..."
	@for file in $(DIST_DIR)/*; do \
		if [ -f "$$file" ]; then \
			extension="$${file##*.}"; \
			if [ "$$extension" = "exe" ]; then \
				base=$$(basename "$$file" .exe); \
				zip $(DIST_DIR)/$$base.zip "$$file"; \
			else \
				tar -czf $(DIST_DIR)/$$(basename "$$file").tar.gz "$$file"; \
			fi; \
		fi; \
	done

help:
	@echo "Available targets:"
	@echo "  build        - Build the binary for current platform"
	@echo "  build-all    - Build binaries for all platforms"
	@echo "  clean        - Remove build artifacts"
	@echo "  test         - Run tests"
	@echo "  lint         - Run linters"
	@echo "  install      - Install binary to GOPATH/bin"
	@echo "  release      - Create release packages"
