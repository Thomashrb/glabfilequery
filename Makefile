all: test build

test:
	go test ./internal/gitlab/*

build:
	go build -o bin/glabfilequery ./cmd/main.go
