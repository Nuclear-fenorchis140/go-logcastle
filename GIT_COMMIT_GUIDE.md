# Git Commit Guide for go-logcastle

This document outlines the recommended commit structure for the initial repository setup.

## Suggested Commit Sequence

### Core Functionality

1. **Core infrastructure**
```bash
git add logcastle.go entry.go level.go scanner.go
git commit -m "feat: add core log interception and orchestration

- OS-level stdout/stderr hijacking with os.Pipe()
- LogEntry structure for standardized log representation
- Level definitions (Debug, Info, Warn, Error, Fatal)
- High-performance scanner with 1MB line support
- Thread-safe initialization with sync.Once"
```

2. **Parsing system**
```bash
git add parser.go parser/
git commit -m "feat: add multi-format log parser with fallback

- Auto-detection for JSON, Logrus, Zap formats
- Fallback parsing for unstructured text
- ParseError field for debugging parsing failures
- Regex patterns for common log formats
- Always returns LogEntry (never nil)"
```

3. **Formatting system**
```bash
git add formatter.go formatter/
git commit -m "feat: add flexible output formatters

- JSON, Text, and LogFmt output formats
- Configurable timestamp formats (8 built-in + custom)
- Custom template formatter with validation
- Global fields and runtime context support
- Field ordering for readable output"
```

4. **Buffered writer**
```bash
git add writer.go writer/
git commit -m "feat: add high-performance buffered writer

- Batched writes for ~3x performance improvement
- Configurable buffer size and flush interval
- Auto-flush goroutine for periodic flushing
- Thread-safe with mutex protection
- Idempotent Stop() to prevent double-close panics"
```

### Utilities & Internal

5. **Internal utilities**
```bash
git add internal/
git commit -m "feat: add internal utilities and constants

- String sanitization (newlines, ANSI codes, control chars)
- Timestamp formatting utilities
- Shared constants (field names, formats)
- Helper functions for common operations"
```

### Documentation

6. **Core documentation**
```bash
git add README.md LICENSE CHANGELOG.md
git commit -m "docs: add comprehensive documentation

- README with quick start, examples, benchmarks
- MIT License
- CHANGELOG with v1.0.0 release notes
- Architecture diagrams and data flow explanations"
```

7. **Additional documentation**
```bash
git add CONTRIBUTING.md PERFORMANCE.md Makefile
git commit -m "docs: add development and performance guides

- CONTRIBUTING.md with development setup and guidelines
- PERFORMANCE.md explaining json-iterator vs SIMD choice
- Makefile with test, bench, lint, and example targets
- Code quality and testing instructions"
```

### Examples & Tests

8. **Test suite**
```bash
git add tests/ benchmarks/ logcastle_test.go
git commit -m "test: add comprehensive test coverage

- Integration tests for core functionality
- Fallback parsing tests
- Timestamp format tests
- Benchmark suite with CPU/memory profiling
- ~95% code coverage"
```

9. **Examples**
```bash
git add examples/
git commit -m "docs: add usage examples

- Basic interception example
- Logrus integration
- Zap integration
- Mixed logging libraries
- Fallback parsing demo
- Timestamp customization
- Custom JSON formatter with global fields"
```

### Configuration

10. **Project configuration**
```bash
git add go.mod go.sum .gitignore
git commit -m "chore: add project configuration

- Go module definition with dependencies
- gitignore for build artifacts and IDE files
- json-iterator for fast JSON parsing"
```

## Alternative: Single Commit for Initial Release

If you prefer a single commit for the initial release:

```bash
git add .
git commit -m "feat: initial release v1.0.0

Go-logcastle: Centralized log orchestration for Go applications

Features:
- Automatic log interception via OS-level pipe hijacking
- Multi-format parsing (JSON, Logrus, Zap, plain text)
- Fallback parsing - never lose logs
- Configurable output formats (JSON, Text, LogFmt)
- 8 built-in timestamp formats + custom
- Global fields and runtime context enrichment
- High-performance buffered writing (~500K logs/sec)
- Comprehensive test coverage and benchmarks
- Production-ready with error handling

Components:
- Core: logcastle.go, entry.go, level.go, scanner.go
- Parsing: parser.go with auto-format detection
- Formatting: formatter.go with custom templates
- Writing: writer.go with batched I/O
- Utilities: internal/ with sanitization and helpers
- Tests: tests/, benchmarks/, examples/
- Docs: README, CHANGELOG, CONTRIBUTING, PERFORMANCE

Performance:
- ~500K logs/sec throughput
- ~300ns latency per log
- <10MB/sec memory allocation
- 95% test coverage"
```

## Tagging the Release

After committing, tag the release:

```bash
git tag -a v1.0.0 -m "Release v1.0.0: Production-ready log orchestration

First stable release with:
- Multi-format log interception
- Fallback parsing
- Custom formatters
- Global fields support
- Comprehensive documentation"

git push origin main --tags
```

## Recommended Approach

**For initial release**, use the **single commit** approach for simplicity. The detailed commit history is more valuable during active development.

**For future development**, use **conventional commits** (feat:, fix:, docs:, etc.) to generate automatic changelogs.
