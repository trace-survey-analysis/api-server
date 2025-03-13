.DEFAULT_GOAL := build

.PHONY: build fmt vet run

fmt:
	echo "Running go fmt"
	go fmt ./...

vet:
	echo "Running go vet"
	go vet ./...

build:
	echo "Building the binary"
	go build -v -o bin/ ./...

all: 
	echo "Running checks and Building the binary"
	make fmt vet build

# for development
run: fmt vet
	echo "Running the server"
	go run cmd/api-server/main.go