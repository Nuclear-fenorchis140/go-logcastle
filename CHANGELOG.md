# Changelog

All notable changes to go-logcastle will be documented in this file.

## [1.0.2] - 2026-03-23

### Added
- **Configurable metadata fields** - `IncludeLoggerField` and `IncludeParseError` config options (both default: false)
- **Smart logger detection** - Distinguishes stdlib `log` vs direct `stdout` writes
- **Clean output by default** - No metadata clutter unless explicitly enabled

### Changed
- Removed 150+ lines of hardcoded library pattern matching (MongoDB, Redis, GIN, etc.)
- Simplified logger detection using timestamp pattern recognition
- Default output now excludes `logger` and `log_parse_error` fields

### Improved
- More maintainable parser code
- Faster generic log parsing

## [1.0.1] - 2026-03-23

### Fixed
- **CRITICAL BUG**: stdlib `log.Println()` and `log.Printf()` calls were not being captured
- Added `log.SetOutput()` reconfiguration after stderr replacement
- Root cause: stdlib log package caches `os.Stderr` at import time and must be explicitly updated

### Changed
- Increased interception delay from 5ms to 10ms for more reliable initialization

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