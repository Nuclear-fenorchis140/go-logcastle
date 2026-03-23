# Makefile for go-logcastle

.PHONY: all test bench lint fmt clean examples

all: test

# Run tests
test:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run benchmarks
bench:
	go test -bench=. -benchmem -cpuprofile=cpu.prof -memprofile=mem.prof

# Lint code
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...
	goimports -w .

# Build examples
examples:
	go build -o bin/basic ./examples/basic
	go build -o bin/logrus ./examples/logrus
	go build -o bin/zap ./examples/zap
	go build -o bin/mixed ./examples/mixed

# Run example
run-basic:
	go run ./examples/basic

run-logrus:
	go run ./examples/logrus

run-zap:
	go run ./examples/zap

run-mixed:
	go run ./examples/mixed

# Install dependencies
deps:
	go mod download
	go mod tidy

# Clean generated files
clean:
	rm -f coverage.out coverage.html cpu.prof mem.prof
	rm -rf bin/

# View coverage
cover: test
	go tool cover -html=coverage.out
