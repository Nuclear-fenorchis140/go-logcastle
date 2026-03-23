# Changelog

All notable changes to go-logcastle will be documented in this file.

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
