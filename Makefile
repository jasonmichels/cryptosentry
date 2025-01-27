.PHONY: all build run test clean

# Default target: run both tests and build
all: test build

test:
	@echo "Running tests..."
	go test -v ./... -bench . -cover

build:
	@echo "Building binary..."
	go build -o bin/cryptosentry ./cmd/main.go

run:
	@echo "Running application..."
	go run ./cmd/main.go

clean:
	@echo "Cleaning up..."
	rm -rf bin
