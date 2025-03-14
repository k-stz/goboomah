.DEFAULT_GOAL := run

.PHONY: fmt vet build
fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build: vet
	go build -o bin/goboomah

wsl: vet
	GOOS=windows go build -o bin/goboomah
	./bin/goboomah

test:
	go test

run: build
	./bin/goboomah