# Changelog

All notable changes to go-logcastle will be documented in this file.

## [1.0.3] - 2026-03-23

### Added
- **FlattenFields config option** (default: true) - Merges enrichment fields to root level instead of nested "fields" object
  - Critical for Grafana/Loki label extraction and query performance
  - Example: `{"env":"prod","service":"api"}` vs `{"fields":{"env":"prod","service":"api"}}`
- **PrettyPrint config option** (default: false) - Multi-line JSON with indentation for development/debugging
  - Improves terminal readability during development
  - Single-line output remains default for production log aggregation
- **ColorOutput config option** (default: false) - ANSI color codes for Text format terminal output
  - ERROR appears in bold red, WARN in yellow, INFO in green, DEBUG in gray
  - Only applies to Text format (ignored in JSON/LogFmt)
- **FieldOrder config option** - Custom field ordering in JSON output
  - Optimizes log readability in ELK/Logstash/Splunk where field order matters
  - Example: `FieldOrder: []string{"timestamp", "level", "service", "env", "message"}`
  - Unspecified fields appear after ordered fields, in alphabetical order
- **Comprehensive examples** in `examples/formatting/main.go` demonstrating:
  - Terminal development mode (pretty print + flattened)
  - Production Grafana/Loki mode (single-line + flattened)
  - ELK/Logstash mode (custom field ordering)
  - Terminal with ANSI colors (Text format)
- **Complete test coverage** for all new features:
  - `TestFlattenFields` and `TestNestedFields` - Field flattening behavior
  - `TestPrettyPrint` and `TestSingleLineJSON` - Output formatting
  - `TestColorOutput` and `TestNoColorOutput` - ANSI color codes
  - `TestFieldOrder` - Custom field ordering

### Improved
- **Package documentation** with comprehensive usage examples for each format type
- **Config struct documentation** with before/after examples showing exactly what each option does
- **Formatter documentation** explaining all format types, advanced features, and use cases
- **DefaultConfig documentation** listing all defaults with explanations
- Format type constants now include detailed descriptions and example outputs

### Features
- **Production observability optimizations** for modern log aggregation platforms:
  - Grafana/Loki: Flattened fields enable label extraction from root-level fields
  - ELK/Logstash: Custom field ordering improves parsing performance and readability
  - Terminal/Console: ANSI colors and pretty printing for better developer experience
- **Format flexibility**: Choose between flat vs nested structure based on your observability stack
- **Development/Production modes**: Easy switching between readable (pretty, colored) and optimized (single-line, flat) outputs

### Technical Highlights
- All 21 tests passing including 7 new formatting tests
- Zero breaking changes - all new features are opt-in with sensible defaults
- FlattenFields defaults to `true` for optimal Grafana/Loki experience out of the box

## [1.0.2] - 2026-03-23

### Added
- **Configurable metadata fields** - `IncludeLoggerField` and `IncludeParseError` config options (both default: false)
- **Smart logger detection** - Distinguishes stdlib `log` vs direct `stdout` writes based on timestamp patterns
- **Clean output by default** - No metadata clutter unless explicitly enabled
- **Reset() function** - Enables proper test isolation (test-only, not for production use)
- **Comprehensive test suite** - Tests for logger detection, clean output, stdlib log capture, and fallback parsing

### Fixed
- **CRITICAL BUG**: stdlib `log.Println()` and `log.Printf()` calls were not being captured
- Added `log.SetOutput()` reconfiguration after stderr replacement
- Root cause: stdlib log package caches `os.Stderr` at import time and must be explicitly updated

### Changed
- Simplified logger detection using timestamp pattern recognition (stdlib log format: `YYYY/MM/DD HH:MM:SS`)
- Default output now excludes `logger` and `log_parse_error` fields
- Increased interception delay from 5ms to 10ms for more reliable initialization
- Updated all tests to use Reset() for proper isolation between test runs
- Increased interception delay from 5ms to 10ms for more reliable initialization

### Improved
- More maintainable parser code
- Faster generic log parsing without pattern matching overhead

## [1.0.1] - 2026-03-23

### Fixed
- Module path corrected to `github.com/bhaskarblur/go-logcastle`

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