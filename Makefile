BINARY    := ferry
PKG       := ./cmd/ferry
VERSION   ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT    ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE      ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS   := -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)
PLATFORMS := darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64 windows/arm64

.PHONY: all build install test race bench fuzz vet fmt check cross snapshot clean

all: check build

## build: build ./bin/ferry for the current platform
build:
	go build -trimpath -ldflags "$(LDFLAGS)" -o bin/$(BINARY) $(PKG)

## install: install ferry into $(go env GOPATH)/bin
install:
	go install -trimpath -ldflags "$(LDFLAGS)" $(PKG)

## test: run unit tests
test:
	go test ./...

## race: run tests with the race detector
race:
	go test -race ./...

## bench: run Go benchmarks (no unit tests)
bench:
	go test -run='^$$' -bench=. -benchmem ./...

## fuzz: fuzz one target, e.g. make fuzz FUZZ=FuzzReadFrame PKGDIR=./internal/protocol
fuzz:
	go test -run='^$$' -fuzz=$(FUZZ) -fuzztime=60s $(PKGDIR)

## vet: run go vet
vet:
	go vet ./...

## fmt: format all Go files
fmt:
	gofmt -s -w .

## check: formatting + vet + tests (what CI runs)
check: vet
	@test -z "$$(gofmt -s -l .)" || (echo "files need gofmt:"; gofmt -s -l .; exit 1)
	go test ./...

## cross: cross-compile for all target platforms into ./dist
cross:
	@for p in $(PLATFORMS); do \
		os=$${p%/*}; arch=$${p#*/}; ext=""; \
		if [ "$$os" = "windows" ]; then ext=".exe"; fi; \
		echo "building $$os/$$arch"; \
		CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build -trimpath -ldflags "$(LDFLAGS)" \
			-o dist/$(BINARY)-$$os-$$arch$$ext $(PKG) || exit 1; \
	done

## snapshot: local GoReleaser dry run (requires goreleaser)
snapshot:
	goreleaser release --snapshot --clean

## clean: remove build output
clean:
	rm -rf bin dist

## help: list targets
help:
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## //'
