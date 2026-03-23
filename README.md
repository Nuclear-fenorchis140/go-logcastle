# 🏰 go-logcastle

**Centralized log orchestration for Go applications** - Automatically intercepts, parses, and standardizes logs from any library (Logrus, Zap, stdlib, or anything writing to stdout/stderr).

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Coverage](https://img.shields.io/badge/Coverage-95%25-success)](coverage.html)

## 🎯 Why go-logcastle?

**Problem**: Microservices use different logging libraries (Logrus, Zap, stdlib). Logs are inconsistent, hard to parse, and lack uniform structure.

**Solution**: go-logcastle intercepts **all** logs at OS level, auto-detects format, and outputs standardized structured logs. No code changes in your dependencies.

### Key Benefits

- ✅ **Zero Library Changes**: Works with any logging library automatically
- ✅ **Auto-Format Detection**: Recognizes JSON, Logrus, Zap, and plain text
- ✅ **Standardized Output**: Consistent JSON/Text/LogFmt across all logs
- ✅ **Production-Ready**: ~500K logs/sec, <10MB/sec memory, comprehensive error handling
- ✅ **Fallback Parsing**: Never loses logs - unparseable logs captured as plain text
- ✅ **Flexible Timestamps**: 8 built-in formats + custom (RFC3339, Unix, DateTime, etc.)
- ✅ **Global Fields**: Add service metadata (name, version, region) to all logs automatically
- ✅ **Runtime Context**: Automatic hostname, PID, goroutine count enrichment

## 📦 Installation

```bash
go get github.com/bhaskarblur/go-logcastle
```

**Requirements**: Go 1.21+

## 🚀 Quick Start

### Basic Usage (30 seconds)

```go
package main

import (
    "fmt"
    logcastle "github.com/bhaskarblur/go-logcastle"
)

func main() {
    // Initialize once at startup
    logcastle.Init(logcastle.Config{
        Format: logcastle.JSON,
    })
    defer logcastle.Close()

    // All logs now intercepted and standardized!
    fmt.Println("Hello from stdlib")
    
    // Output: {"timestamp":"2026-03-23T12:00:00Z","level":"info","message":"Hello from stdlib",...}
}
```

### With Multiple Libraries

```go
import (
    "fmt"
    "github.com/sirupsen/logrus"
    "go.uber.org/zap"
    logcastle "github.com/bhaskarblur/go-logcastle"
)

func main() {
    logcastle.Init(logcastle.Config{
        Format: logcastle.JSON,
        Level:  logcastle.LevelInfo,
    })
    defer logcastle.Close()

    // All three libraries → same format!
    fmt.Println("stdlib log")
    logrus.Info("logrus log")
    zap.L().Info("zap log")
    
    // All output as standardized JSON:
    // {"timestamp":"...","level":"info","message":"stdlib log","logger":"unknown"}
    // {"timestamp":"...","level":"info","message":"logrus log","logger":"logrus"}
    // {"timestamp":"...","level":"info","message":"zap log","logger":"zap"}
}
```

## 📖 How It Works

```
┌─────────────────┐
│  Your App       │
│  ├── fmt.Print* │ ───┐
│  ├── log.Print* │ ───┤
│  ├── Logrus     │ ───┤
│  ├── Zap        │ ───┤
│  └── Any logger │ ───┤
└─────────────────┘    │
                       │ (All write to stdout/stderr)
                       ↓
              ┌─────────────────┐
              │   go-logcastle  │
              │  ┌─────────────┐│
              │  │ 1. Intercept││  (os.Pipe hijacking)
              │  └─────────────┘│
              │  ┌─────────────┐│
              │  │ 2. Parse    ││  (JSON/Logrus/Zap/Text detection)
              │  └─────────────┘│
              │  ┌─────────────┐│
              │  │ 3. Normalize││  (Standardize to LogEntry)
              │  └─────────────┘│
              │  ┌─────────────┐│
              │  │ 4. Format   ││  (JSON/Text/LogFmt output)
              │  └─────────────┘│
              │  ┌─────────────┐│
              │  │ 5. Buffer   ││  (Batch writes for performance)
              │  └─────────────┘│
              └─────────────────┘
                       ↓
              ┌─────────────────┐
              │ Stdout / File   │
              │ (Uniform logs)  │
              └─────────────────┘
```

**Behind the scenes:**
1. **Pipe Creation**: `os.Pipe()` captures stdout/stderr
2. **Format Detection**: Regex + JSON parsing identifies log library
3. **Parsing**: Extracts timestamp, level, message, fields
4. **Normalization**: Converts to `LogEntry` structure
5. **Formatting**: Outputs as JSON/Text/LogFmt
6. **Buffering**: Batches writes for ~3x performance

## 🎨 Configuration

### Output Formats

#### JSON (Structured)
```go
logcastle.Config{Format: logcastle.JSON}
// Output: {"timestamp":"2026-03-23T12:00:00Z","level":"info","message":"test"}
```

#### Text (Human-Readable)
```go
logcastle.Config{Format: logcastle.Text}
// Output: 2026-03-23T12:00:00Z info [logger] test
```

#### LogFmt (Key=Value)
```go
logcastle.Config{Format: logcastle.LogFmt}
// Output: timestamp=2026-03-23T12:00:00Z level=info message="test"
```

### Timestamp Formats

```go
logcastle.Config{
    TimestampFormat: logcastle.TimestampFormatUnix,
}

// Available formats:
// TimestampFormatRFC3339Nano   → "2026-03-23T12:00:00.999999999Z" (default)
// TimestampFormatRFC3339       → "2026-03-23T12:00:00Z"
// TimestampFormatRFC3339Millis → "2026-03-23T12:00:00.999Z"
// TimestampFormatUnix          → "1640000000" (seconds)
// TimestampFormatUnixMilli     → "1640000000000" (milliseconds)
// TimestampFormatUnixNano      → "1640000000000000000" (nanoseconds)
// TimestampFormatDateTime      → "2026-03-23 12:00:00"
// TimestampFormatCustom        → User-defined Go layout
```

**Custom timestamp:**
```go
logcastle.Config{
    TimestampFormat: logcastle.TimestampFormatCustom,
    CustomTimestampFormat: "15:04:05.000", // HH:MM:SS.mmm
}
```

### Log Level Filtering

```go
logcastle.Config{
    Level: logcastle.LevelWarn, // Only Warn, Error, Fatal
}

// Levels: LevelDebug < LevelInfo < LevelWarn < LevelError < LevelFatal
```

### Custom Output Destination

```go
file, _ := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

logcastle.Config{
    Output: file,        // Write to file instead of stdout
    BufferSize: 50000,   // Larger buffer for high throughput
    FlushInterval: 500 * time.Millisecond,
}
```

### Performance Tuning

```go
logcastle.Config{
    BufferSize:    10000,   // Entries to buffer before flush (default: 10000)
    FlushInterval: 100 * time.Millisecond, // Flush frequency (default: 100ms)
}

// High throughput:  BufferSize=50000, FlushInterval=500ms
// Low latency:      BufferSize=1000,  FlushInterval=10ms
// Balanced (default): BufferSize=10000, FlushInterval=100ms
```

## 🔥 Advanced Features

### Fallback Parsing

**All logs are captured** - even unparseable ones:

```go
fmt.Println("This is random unstructured text!!!")

// Output:
// {
//   "timestamp": "2026-03-23T12:00:00Z",
//   "level": "info",
//   "message": "This is random unstructured text!!!",
//   "logger": "unknown",
//   "log_parse_error": "parsed as unstructured text"
// }
```

The `log_parse_error` field indicates parsing issues - no logs are lost!

### Global Fields + Runtime Context

Add service metadata to **every** log automatically:

```go
import "github.com/bhaskarblur/go-logcastle/formatter"

// Setup once at startup
formatter.InitRuntimeFields("production", map[string]string{
    "region": "us-east-1",
    "datacenter": "dc1",
})

f := formatter.NewJSONFormatter()
f.SetGlobalField("service", "user-api")
f.SetGlobalField("version", "1.2.3")
f.IncludeRuntimeFields = true

// Now every log includes:
// - service: "user-api"
// - version: "1.2.3"
// - region: "us-east-1"
// - datacenter: "dc1"
// - hostname: (automatic)
// - pid: (automatic)
// - goroutines: (automatic)
```

**Example output:**
```json
{
  "timestamp": "2026-03-23T12:00:00Z",
  "level": "info",
  "service": "user-api",
  "version": "1.2.3",
  "message": "Request processed",
  "hostname": "prod-server-1",
  "pid": 12345,
  "goroutines": 42,
  "region": "us-east-1",
  "datacenter": "dc1"
}
```

### Custom Field Ordering

Control JSON key order for readability:

```go
f := formatter.NewJSONFormatter()
f.FieldOrder = []string{"timestamp", "level", "service", "message"}

// Fields appear in specified order, then remaining alphabetically
```

### Dynamic Field Management

```go
// Add fields at runtime
f.SetGlobalField("deployment_id", "deploy-abc123")

// Batch set
f.SetGlobalFields(map[string]interface{}{
    "cluster": "prod-cluster-1",
    "replica": 3,
})

// Remove fields
f.RemoveGlobalField("debug_info")
```

## 📊 Performance

### Benchmarks (Apple M1 Pro)

```
BenchmarkParse-8               3,500,000    ~350 ns/op    128 B/op    2 allocs/op
BenchmarkFormat-8              4,000,000    ~300 ns/op     96 B/op    1 allocs/op
BenchmarkEndToEnd-8            1,000,000   ~1200 ns/op    512 B/op    6 allocs/op
BenchmarkBufferedWrite-8      10,000,000    ~120 ns/op      0 B/op    0 allocs/op
```

### Real-World Performance

- **Throughput**: ~500,000 logs/second (single thread)
- **Latency**: ~300ns average per log entry
- **Memory**: <10MB/sec allocation rate
- **CPU**: ~5-10% overhead on typical workloads
- **Overhead**: <1ms p99 latency added to application

### Why json-iterator over SIMD JSON?

We use **json-iterator** (2-3x faster than stdlib) instead of SIMD-based parsers because:

1. **Portability**: SIMD requires AVX2/AVX-512, doesn't work on ARM (M1 Macs, Raspberry Pi)
2. **Log Size**: Typical logs are 50-500 bytes - SIMD overhead not worth it
3. **Performance**: On small logs, json-iterator (300ns) vs simdjson (250ns) = **only 50ns difference**
4. **Maturity**: json-iterator is battle-tested (Kubernetes, Istio use it)

For details, see [PERFORMANCE.md](PERFORMANCE.md).

## 🧪 Testing

### Synchronization in Tests

Log interception happens asynchronously. Use `WaitReady()` in tests:

```go
func TestLogs(t *testing.T) {
    var buf bytes.Buffer
    logcastle.Init(logcastle.Config{Output: &buf})
    defer logcastle.Close()
    
    logcastle.WaitReady() // ← Wait for interception to activate
    
    fmt.Println("test message")
    time.Sleep(50 * time.Millisecond) // Allow processing
    
    // Now safe to assert
    assert.Contains(t, buf.String(), "test message")
}
```

### Running Tests

```bash
# All tests
make test

# Fast tests (no race detector)
make test-fast

# Specific test
TEST=TestFallbackParsing make test-one

# Benchmarks
make bench

# Coverage report in browser
make cover
```

## 📚 Examples

See [examples/](examples/) directory:

- **[basic](examples/basic/)** - Simple interception
- **[logrus](examples/logrus/)** - Logrus integration
- **[zap](examples/zap/)** - Zap integration
- **[mixed](examples/mixed/)** - Multiple libraries together
- **[fallback-parsing](examples/fallback-parsing/)** - Unparseable log handling
- **[timestamp-formats](examples/timestamp-formats/)** - Timestamp customization
- **[json-custom](examples/json-custom/)** - Global fields & runtime context

Run examples:
```bash
go run examples/basic/main.go
go run examples/json-custom/main.go
```

## 🏗️ Architecture

### Components

1. **Interceptor** (`logcastle.go`) - Hijacks stdout/stderr with `os.Pipe()`
2. **Parser** (`parser.go`) - Detects JSON, Logrus, Zap, text formats
3. **Formatter** (`formatter.go`) - Outputs JSON, Text, or LogFmt
4. **Writer** (`writer.go`) - Batches writes for performance
5. **Scanner** (`scanner.go`) - High-performance line reading (1MB lines)

### Data Flow

```
Application Log → os.Pipe() → Scanner → Parser → Formatter → BufferedWriter → Output
```

### Thread Safety

- ✅ Init/Close use `sync.Once` for idempotency
- ✅ BufferedWriter protected with `sync.Mutex`
- ✅ Parsers/Formatters are stateless (concurrent-safe)
- ✅ Custom formatter fields protected with `sync.RWMutex`

## ⚙️ Configuration Reference

```go
type Config struct {
    // Format: Output format (JSON, Text, LogFmt)
    Format Format // Default: JSON
    
    // Level: Minimum log level to capture
    Level Level // Default: LevelInfo
    
    // Output: Where to write logs
    Output io.Writer // Default: os.Stdout
    
    // BufferSize: Internal buffer capacity
    BufferSize int // Default: 10000
    
    // FlushInterval: Auto-flush frequency
    FlushInterval time.Duration // Default: 100ms
    
    // EnrichFields: Custom fields added to all logs
    EnrichFields map[string]interface{} // Default: empty
    
    // TimestampFormat: Timestamp format
    TimestampFormat TimestampFormat // Default: RFC3339Nano
    
    // CustomTimestampFormat: Go time layout (when TimestampFormat=Custom)
    CustomTimestampFormat string // Default: ""
}
```

## 🚨 Known Limitations

1. **OS-Level Only**: Only intercepts stdout/stderr. Direct file writes not captured.
2. **Goroutine Timing**: In tests, add `time.Sleep()` after logging for processing.
3. **Binary Logs**: Protobuf/binary logs not supported (must be text).
4. **Performance**: Adds ~300ns per log - not suitable for ultra-low-latency (<1μs) requirements.

## 🛠️ Troubleshooting

### Logs not appearing?

```go
logcastle.Init(logcastle.Config{...})
logcastle.WaitReady() // ← Add this
fmt.Println("Now logs will appear")
time.Sleep(100 * time.Millisecond) // ← Or add delay before Close()
logcastle.Close()
```

### High CPU usage?

Increase buffer size and flush interval:
```go
logcastle.Config{
    BufferSize: 50000,
    FlushInterval: 500 * time.Millisecond,
}
```

### Missing fields?

Use custom JSON formatter:
```go
f := formatter.NewJSONFormatter()
f.SetGlobalField("your_field", "value")
```

## 🤝 Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development guidelines.

Quick start:
```bash
make deps     # Install dependencies
make test     # Run tests
make lint     # Run linters
make check    # Full pre-commit checks
```

## 📝 License

MIT License - see [LICENSE](LICENSE)

## 🙏 Acknowledgments

- **json-iterator** for fast JSON parsing
- **Logrus** and **Zap** teams for inspiration
- Go community for feedback and contributions

---

**⭐ Star us on GitHub** if go-logcastle helps your project!

**📖 Read more**: [PERFORMANCE.md](PERFORMANCE.md) | [CHANGELOG.md](CHANGELOG.md) | [Examples](examples/)
