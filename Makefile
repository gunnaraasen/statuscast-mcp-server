.PHONY: build test lint run-stdio run-http

build:
	go build -o bin/statuscast-mcp-server ./cmd/server

test:
	go test ./...

lint:
	golangci-lint run

run-stdio:
	go run ./cmd/server

run-http:
	TRANSPORT=http go run ./cmd/server
