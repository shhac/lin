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

fmt:
	gofmt -w .
	goimports -w .

vet:
	go vet ./...

generate:
	go generate ./internal/linear/...

dev:
	go run ./cmd/lin $(ARGS)

clean:
	rm -f $(BINARY)

.PHONY: build test test-short lint fmt vet generate dev clean
