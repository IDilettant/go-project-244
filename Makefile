test:
	go test -v ./... -race

lint:
	golangci-lint run ./...

build:
	go build -o bin/gendiff ./cmd/gendiff

.PHONY: test lint build
