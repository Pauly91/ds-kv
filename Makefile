# Name of the main binary (change as needed)
BINARY_NAME=ds-kv

# Directory containing main.go
CMD_DIR=cmd/$(BINARY_NAME)

.PHONY: all build run test clean fmt vet

all: build

build:
	go build -o bin/$(BINARY_NAME) $(CMD_DIR)/cli.go

run: build
	./bin/$(BINARY_NAME)

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	rm -rf bin/$(BINARY_NAME)
