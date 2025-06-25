# Muse Backend Makefile

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
BINARY_NAME=muse-backend

# Add Go bin to PATH
export PATH := $(PATH):$(shell go env GOPATH)/bin

.PHONY: all build clean test test-coverage lint fmt vet deps help

all: deps lint test build

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) -v .

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Run tests
test:
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	$(GOTEST) -race -coverprofile=coverage.out -covermode=atomic ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Run unit tests only (no database)
test-unit:
	$(GOTEST) -short -race -coverprofile=coverage.out -covermode=atomic \
		$$(go list ./... | grep -v '/internal/repository/postgres')

# Run linting
lint:
	@echo "Running golangci-lint..."
	golangci-lint run --timeout=5m

# Format code
fmt:
	$(GOFMT) -s -w .

# Run go vet
vet:
	$(GOCMD) vet ./...

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) verify

# Install development tools
install-tools:
	@echo "Installing development tools..."
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.64.8

# Run all checks (what CI runs)
ci-checks: deps lint vet test-unit

# Run integration tests (requires database)
test-integration:
	$(GOTEST) -v -race -coverprofile=integration-coverage.out ./...

# Run benchmarks
bench:
	$(GOTEST) -bench=. -benchmem -run=^$$ ./...

# Help
help:
	@echo "Available commands:"
	@echo "  build           - Build the application"
	@echo "  clean           - Clean build artifacts"
	@echo "  test            - Run all tests"
	@echo "  test-unit       - Run unit tests only"
	@echo "  test-coverage   - Run tests with coverage report"
	@echo "  test-integration- Run integration tests"
	@echo "  lint            - Run golangci-lint"
	@echo "  fmt             - Format code"
	@echo "  vet             - Run go vet"
	@echo "  deps            - Download dependencies"
	@echo "  install-tools   - Install development tools"
	@echo "  ci-checks       - Run all CI checks locally"
	@echo "  bench           - Run benchmarks"
	@echo "  help            - Show this help" 