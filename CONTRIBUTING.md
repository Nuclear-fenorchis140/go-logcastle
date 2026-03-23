# Contributing to go-logcastle 🏰

Thank you for your interest in contributing! This document provides guidelines for contributing to the project.

## Code of Conduct

- Be respectful and inclusive
- Focus on constructive feedback
- Help make go-logcastle better for everyone

## How to Contribute

### Reporting Bugs

1. Check if the bug is already reported in Issues
2. Create a new issue with:
   - Clear title and description
   - Steps to reproduce
   - Expected vs actual behavior
   - Go version and OS
   - Relevant code snippets

### Suggesting Features

1. Open an issue with the `enhancement` label
2. Describe the feature and use case
3. Explain why it's valuable
4. Provide examples if possible

### Submitting Pull Requests

1. **Fork the repository**

2. **Create a feature branch**
   ```bash
   git checkout -b feat/your-feature-name
   ```

3. **Follow conventional commits**
   - `feat:` new feature
   - `fix:` bug fix
   - `docs:` documentation changes
   - `perf:` performance improvements
   - `refactor:` code refactoring
   - `test:` adding tests
   - `chore:` maintenance tasks

   Examples:
   ```
   feat: add slog format detection
   fix: race condition in pipe initialization
   docs: update README with new examples
   perf: optimize JSON parsing with pooling
   ```

4. **Write tests**
   - Add unit tests for new features
   - Ensure existing tests pass
   - Aim for high coverage

5. **Run quality checks**
   ```bash
   go test ./...
   go test -race ./...
   go vet ./...
   golangci-lint run
   ```

6. **Update documentation**
   - Update README.md if needed
   - Add godoc comments
   - Update CHANGELOG.md

7. **Submit the PR**
   - Clear title and description
   - Reference related issues
   - Explain what and why, not just how

## Development Setup

### Prerequisites
- Go 1.21 or higher
- git
- golangci-lint (for linting)
- goimports (for formatting)

### Clone and setup
```bash
git clone https://github.com/yourusername/go-logcastle.git
cd go-logcastle
go mod download
```

### Running tests
```bash
# All tests with coverage
make test

# Fast tests (no race detector)
make test-fast

# Specific test
TEST=TestFallbackParsing make test-one

# With coverage report in browser
make cover
```

### Running benchmarks
```bash
# Run benchmarks with profiling
make bench

# View CPU profile in browser
make profile-cpu

# View memory profile
make profile-mem
```

### Code quality
```bash
# Format code
make fmt

# Run linters
make lint

# Vet code
make vet

# Full pre-commit checks
make check
```

### Project Structure
```
go-logcastle/
├── logcastle.go          # Main API and orchestration
├── entry.go              # LogEntry data structure
├── level.go              # Log level definitions
├── parser.go             # Multi-format log parser
├── formatter.go          # Output formatters (JSON/Text/LogFmt)
├── writer.go             # Buffered writer
├── scanner.go            # High-performance line scanner
├── formatter/            # Extended formatters
│   ├── formatter.go      # Base formatter types
│   ├── custom.go         # Custom template formatter
│   └── json_custom.go    # JSON formatter with global fields
├── parser/               # Extended parsers
│   └── parser.go         # Additional parsing logic
├── writer/               # Extended writers
│   └── writer.go         # Additional writer logic
├── internal/             # Internal utilities
│   ├── constants/        # Shared constants
│   └── utils/            # Helper functions
├── tests/                # Integration tests
├── benchmarks/           # Performance benchmarks
└── examples/             # Usage examples
```

## Code Style

### General Guidelines
- Follow standard Go conventions
- Use `gofmt` for formatting
- Keep functions focused and small
- Write clear comments
- Use meaningful variable names

### Performance Considerations
- Minimize allocations in hot paths
- Use buffering for I/O operations
- Avoid unnecessary locks
- Profile before optimizing
- Add benchmarks for critical paths

### Error Handling
- Return errors, don't panic (unless truly exceptional)
- Wrap errors with context using `fmt.Errorf`
- Log errors to original stderr for debuggability

### Testing
- Write table-driven tests where appropriate
- Test edge cases and error paths
- Use subtests for related scenarios
- Keep tests fast and deterministic
- Mock external dependencies

## Performance Requirements

go-logcastle is a high-performance library. Contributions should:
- Not significantly degrade performance
- Include benchmarks for new features
- Profile for allocations and CPU usage
- Maintain low latency overhead (<1ms p99)

## Documentation

- Use godoc-style comments
- Include examples in godoc
- Keep README.md up to date
- Document breaking changes
- Explain "why" not just "what"

## Questions?

- Open a GitHub Discussion for questions
- Join our community chat (if available)
- Tag maintainers in issues for assistance

## License

By contributing, you agree that your contributions will be licensed under the project's license (typically MIT or Apache 2.0).

---

**Thank you for contributing to go-logcastle! 🚀**
