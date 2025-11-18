.PHONY: build clean test run install deps lint coverage

BINARY_NAME=scaffold
BUILD_DIR=bin
CMD_DIR=./cmd/scaffold

all: test build

build:
	@echo "Building..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)

clean:
	@echo "Cleaning..."
	@go clean
	@rm -rf $(BUILD_DIR)

test:
	@echo "Testing..."
	@go test -v ./...

run: build
	@./$(BUILD_DIR)/$(BINARY_NAME)

install:
	@echo "Installing..."
	@go install $(CMD_DIR)

deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

lint:
	@echo "Running linters..."
	@golangci-lint run ./...

coverage:
	@echo "Generating coverage report..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
