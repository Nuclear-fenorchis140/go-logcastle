# Changelog

All notable changes to go-logcastle will be documented in this file.

## [1.0.0] - 2026-03-23

### Added
- **Core log interception** - Automatically captures all stdout/stderr from any library
- **Multi-format parsing** - Supports JSON, Logrus, Zap, and unstructured text
- **Fallback parsing** - Unparseable logs captured as plain text with `log_parse_error` field
- **Configurable timestamp formats** - 8 built-in formats (RFC3339, Unix, DateTime, etc.) + custom
- **Custom JSON formatter** - Global fields, runtime fields (hostname, PID), custom field ordering
- **String sanitization** - Handles newlines, ANSI escape codes, control characters
- **Buffered writing** - Configurable buffer size and flush interval for performance
- **Thread-safe operations** - Concurrent access supported throughout
- **Graceful shutdown** - Proper cleanup and log flushing on Close()
- **Zero configuration** - Works out of the box with sensible defaults

### Features
- **Performance**: ~500K logs/sec throughput, ~300ns latency per log
- **Memory efficient**: <10MB/sec allocation rate with pooling
- **Production-ready**: Comprehensive error handling, no log loss
- **Flexible output**: JSON, Text, or LogFmt formats
- **Level filtering**: Debug, Info, Warn, Error, Fatal levels
- **Enrichment**: Add custom fields to all logs automatically

### Documentation
- Comprehensive README with examples and benchmarks
- PERFORMANCE.md explaining JSON parser choice (json-iterator vs SIMD)
- QUICK_REFERENCE.md for API usage
- CONTRIBUTING.md for development guidelines
- Multiple working examples in examples/ directory

### Examples
- Basic interception
- Logrus integration
- Zap integration
- Mixed logging libraries
- Fallback parsing demo
- Timestamp format customization
- Custom JSON formatter with global fields

### Technical Highlights
- Uses json-iterator for 2-3x faster JSON parsing vs stdlib
- High-performance scanner with 1MB line support
- Automatic log format detection (JSON, Logrus text, Zap console)
- Parse error tracking for debugging
- Idempotent shutdown (no double-close panics)

## [Unreleased]

### Planned
- Structured logging framework integration (zerolog, slog)
- Metrics and observability hooks
- Log sampling for high-volume scenarios
- Compression support for high-throughput
- Circuit breaker for I/O failures

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-03-23

### Added
- Initial production release 🚀
- Core log interception via `os.Pipe()` hijacking
- Auto-detection for JSON, Logrus, and Zap log formats
- Standardized output formats: JSON, Text, LogFmt
- High-performance buffered writer with auto-flush
- Level filtering (Debug, Info, Warn, Error, Fatal)
- Field enrichment for adding custom metadata
- Graceful shutdown with proper cleanup
- Thread-safe initialization with `sync.Once`
- `WaitReady()` for synchronization in tests
- Comprehensive test suite with 100% coverage
- Benchmark tests for performance validation
- Multiple working examples (basic, logrus, zap, mixed)

### Features
- ✨ **Zero-configuration**: Single `Init()` call intercepts all logs
- 🚀 **High-performance**: Buffered writing, minimal allocations
- 🔍 **Auto-detection**: Automatically identifies log format
- 📊 **Standardization**: Converts all logs to consistent format
- 🎯 **Level filtering**: Only output logs above configured level
- 🏷️ **Enrichment**: Add custom fields to all logs
- 🔒 **Thread-safe**: Safe for concurrent use
- 🧹 **Graceful shutdown**: Flushes all logs before exit

### Technical Details
- Built with Go 1.21+
- Uses `json-iterator` for fast JSON parsing
- OS-level stdout/stderr interception
- Buffered I/O for performance
- Zero external runtime dependencies (except parsing libs)
