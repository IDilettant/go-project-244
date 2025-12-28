test:
	go test -v ./... -race

lint:
	golangci-lint run ./...

build:
	go build -o bin/gendiff ./cmd/gendiff

cover:
	go test ./... -count=1 -coverprofile=coverage.out && go tool cover -func=coverage.out

.PHONY: test lint build cover
