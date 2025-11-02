# Makefile for portctl
.PHONY: all build test lint integration clean release

all: build

build:
	go build -o cmd/portctl/portctl ./cmd/portctl

test:
	go test ./...

lint:
	# optional: run golangci-lint if available
	if command -v golangci-lint >/dev/null 2>&1; then golangci-lint run; else echo "golangci-lint not installed, skipping"; fi

integration:
	chmod +x tests/integration/*.sh
	cd tests/integration && ./run_tests.sh

clean:
	rm -f cmd/portctl/portctl
	rm -rf dist

release:
	goreleaser release --rm-dist
