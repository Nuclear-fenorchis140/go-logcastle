# Makefile for go-logcastle

.PHONY: all test bench lint fmt clean examples install

all: test

# Run all tests with coverage
test:
	@echo "Running tests with race detector and coverage..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Run tests without race detector (faster)
test-fast:
	@go test -v ./...

# Run specific test
test-one:
	@go test -v -run $(TEST) ./...

# Run benchmarks with memory profiling
bench:
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem -cpuprofile=cpu.prof -memprofile=mem.prof ./benchmarks/
	@echo "Profiles: cpu.prof, mem.prof"

# View benchmark results
bench-compare:
	@go test -bench=. -benchmem ./benchmarks/ | tee bench-new.txt

# Lint code
lint:
	@echo "Running linters..."
	@golangci-lint run --timeout=5m

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@goimports -w .

# Vet code
vet:
	@go vet ./...

# Build all examples
examples:
	@echo "Building examples..."
	@mkdir -p bin
	@go build -o bin/basic ./examples/basic
	@go build -o bin/logrus ./examples/logrus
	@go build -o bin/zap ./examples/zap
	@go build -o bin/mixed ./examples/mixed
	@go build -o bin/fallback-parsing ./examples/fallback-parsing
	@go build -o bin/timestamp-formats ./examples/timestamp-formats
	@go build -o bin/json-custom ./examples/json-custom
	@echo "Examples built in bin/"

# Run examples
run-basic:
	@go run ./examples/basic

run-logrus:
	@go run ./examples/logrus

run-zap:
	@go run ./examples/zap

run-mixed:
	@go run ./examples/mixed

run-fallback:
	@go run ./examples/fallback-parsing

run-timestamp:
	@go run ./examples/timestamp-formats

run-json-custom:
	@go run ./examples/json-custom

# Install dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

# Clean generated files
clean:
	@echo "Cleaning..."
	@rm -f coverage.out coverage.html cpu.prof mem.prof
	@rm -rf bin/
	@rm -f bench-new.txt bench-old.txt

# View coverage in browser
cover: test
	@go tool cover -html=coverage.out

# View CPU profile
profile-cpu: bench
	@go tool pprof -http=:8080 cpu.prof

# View memory profile
profile-mem: bench
	@go tool pprof -http=:8080 mem.prof

# Install golangci-lint (if not installed)
install-lint:
	@command -v golangci-lint >/dev/null 2>&1 || { \
		echo "Installing golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin; \
	}

# Install goimports (if not installed)
install-tools:
	@go install golang.org/x/tools/cmd/goimports@latest

# Full check before commit
check: fmt vet lint test
	@echo "✅ All checks passed!"

# Quick validation
validate: fmt vet test-fast
	@echo "✅ Quick validation passed!"

# Help
help:
	@echo "Available targets:"
	@echo "  make test          - Run tests with coverage"
	@echo "  make test-fast     - Run tests without race detector"
	@echo "  make bench         - Run benchmarks"
	@echo "  make lint          - Run linters"
	@echo "  make fmt           - Format code"
	@echo "  make examples      - Build all examples"
	@echo "  make run-*         - Run specific example"
	@echo "  make check         - Run all checks (pre-commit)"
	@echo "  make clean         - Remove generated files"
	@echo "  make cover         - View coverage in browser"
	@echo "  make help          - Show this help"
