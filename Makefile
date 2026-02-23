.PHONY: all build test lint clean

APP_NAME := dockit
BUILD_DIR := bin

all: lint test build

build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) ./main.go

test:
	@echo "Running tests..."
	go test -v ./...

lint:
	@echo "Running linter..."
	golangci-lint run ./... || go vet ./...

clean:
	@echo "Cleaning artifacts..."
	rm -rf $(BUILD_DIR)

