BINARY := lin
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -s -w -X main.version=$(VERSION)

build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) ./cmd/lin

test:
	go test ./... -count=1

test-short:
	go test ./... -count=1 -short

lint:
	golangci-lint run ./...

# Scoped to tracked files: this repo keeps a module cache under .cache/, which
# the go tool skips (dot-directory) but gofmt and goimports walk into, so a bare
# `-w .` rewrites vendored dependencies and makes `gofmt -l .` report noise.
fmt:
	gofmt -w $$(git ls-files '*.go')
	@command -v goimports >/dev/null && goimports -w $$(git ls-files '*.go') || echo "goimports not installed (optional; install: go install golang.org/x/tools/cmd/goimports@latest)"

vet:
	go vet ./...

generate:
	go generate ./internal/linear/...
	go run scripts/fix-omitempty.go -- internal/linear/generated.go

dev:
	go run ./cmd/lin $(ARGS)

clean:
	rm -f $(BINARY)

.PHONY: build test test-short lint fmt vet generate dev clean
