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
- ✅ **Advanced Formatting** (v1.0.3+): 
  - **FlattenFields**: Grafana/Loki label extraction optimization
  - **PrettyPrint**: Multi-line JSON for terminal readability
  - **ColorOutput**: ANSI colors for Text format (ERROR=red, WARN=yellow, etc.)
  - **FieldOrder**: Custom field ordering for ELK/Logstash pipelines

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

#### JSON (Structured - Default)
```go
logcastle.Config{Format: logcastle.JSON}
// Output: {"timestamp":"2026-03-23T12:00:00Z","level":"info","message":"test"}
```

#### Text (Human-Readable)
```go
logcastle.Config{Format: logcastle.Text}
// Output: 2026-03-23T12:00:00Z INFO test
```

#### LogFmt (Key=Value)
```go
logcastle.Config{Format: logcastle.LogFmt}
// Output: timestamp=2026-03-23T12:00:00Z level=info message="test"
```

### Advanced Formatting Options (v1.0.3+)

#### FlattenFields - Grafana/Loki Optimization

**Critical for production observability!** Merges enrichment fields to root level.

```go
// Flattened (default: true) - RECOMMENDED for Grafana/Loki
logcastle.Config{
    FlattenFields: true,
    EnrichFields: map[string]interface{}{
        "env":     "production",
        "service": "payment-service",
    },
}
// Output: {"timestamp":"...","level":"info","env":"production","service":"payment-service",...}

// Nested (false) - Fields grouped under "fields" key
logcastle.Config{
    FlattenFields: false,
    EnrichFields: map[string]interface{}{
        "env": "prod",
    },
}
// Output: {"timestamp":"...","level":"info","fields":{"env":"prod"},...}
```

**Why flatten?** Grafana/Loki can extract labels from root-level fields for filtering: `{service="payment-service", env="production"}`. Nested fields cannot be used as labels.

#### PrettyPrint - Development Readability

Multi-line JSON with indentation for terminal viewing.

```go
// Pretty (true) - Development/Debugging
logcastle.Config{
    Format:      logcastle.JSON,
    PrettyPrint: true,
}
// Output:
// {
//   "timestamp": "2026-03-23T12:00:00Z",
//   "level": "info",
//   "message": "server started"
// }

// Single-line (default: false) - Production
logcastle.Config{
    PrettyPrint: false,
}
// Output: {"timestamp":"2026-03-23T12:00:00Z","level":"info","message":"server started"}
```

#### ColorOutput - Terminal Colors

ANSI color codes for Text format (ignored in JSON/LogFmt).

```go
logcastle.Config{
    Format:      logcastle.Text,
    ColorOutput: true,
}
// Output (with colors):
// 2026-03-23T12:00:00Z \033[31mERROR\033[0m Failed to connect  (red)
// 2026-03-23T12:00:00Z \033[33mWARN\033[0m High memory usage   (yellow)
// 2026-03-23T12:00:00Z \033[32mINFO\033[0m Server started      (green)
// 2026-03-23T12:00:00Z \033[90mDEBUG\033[0m Cache hit          (gray)
```

#### FieldOrder - ELK/Logstash Optimization

Specify which fields appear first in JSON output.

```go
logcastle.Config{
    Format:        logcastle.JSON,
    FlattenFields: true,
    FieldOrder:    []string{"timestamp", "level", "service", "env", "message"},
    EnrichFields: map[string]interface{}{
        "service": "api-gateway",
        "env":     "staging",
    },
}
// Output: {"timestamp":"...","level":"info","service":"api-gateway","env":"staging","message":"...","caller":"..."}
// Fields appear in specified order, remaining fields alphabetically after
```

### Common Configuration Patterns

#### Development Mode (Terminal)
```go
logcastle.Config{
    Format:        logcastle.JSON,
    Level:         logcastle.LevelDebug,    // See all logs
    PrettyPrint:   true,                    // Readable multi-line
    FlattenFields: true,                    // Clean structure
    EnrichFields: map[string]interface{}{
        "env":     "development",
        "service": "my-service",
    },
}
```

#### Development with Colors (Text Format)
```go
logcastle.Config{
    Format:             logcastle.Text,
    Level:              logcastle.LevelDebug,
    ColorOutput:        true,                    // ANSI colors
    IncludeLoggerField: true,                    // Show log source
    EnrichFields: map[string]interface{}{
        "service": "my-service",
    },
}
```

#### Production - Grafana/Loki
```go
logcastle.Config{
    Format:        logcastle.JSON,
    Level:         logcastle.LevelInfo,
    FlattenFields: true,                    // CRITICAL for Loki labels
    PrettyPrint:   false,                   // Single-line for aggregation
    EnrichFields: map[string]interface{}{
        "env":       "production",
        "service":   "payment-service",
        "region":    "us-east-1",
        "pod":       os.Getenv("POD_NAME"),
    },
}
```

#### Production - ELK/Logstash
```go
logcastle.Config{
    Format:        logcastle.JSON,
    Level:         logcastle.LevelInfo,
    FlattenFields: true,
    FieldOrder:    []string{"timestamp", "level", "service", "message"},
    EnrichFields: map[string]interface{}{
        "service":  "user-api",
        "cluster":  "k8s-prod",
        "hostname": os.Getenv("HOSTNAME"),
    },
}
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

### Benchmarks (Apple M2)

```
BenchmarkParse-8               3,500,000    ~350 ns/op    128 B/op    2 allocs/op
BenchmarkFormat-8              4,000,000    ~300 ns/op     96 B/op    1 allocs/op
BenchmarkEndToEnd-8            1,000,000   ~1200 ns/op    512 B/op    6 allocs/op
BenchmarkBufferedWrite-8      10,000,000    ~120 ns/op      0 B/op    0 allocs/op
```

### Real-World Performance

**Baseline (Default Config):**
- **Throughput**: ~500,000 logs/second (single thread)
- **Latency**: ~300ns average per log entry
- **Memory**: <10MB/sec allocation rate
- **CPU**: ~5-10% overhead on typical workloads
- **Overhead**: <1ms p99 latency added to application

### Performance by Configuration

Different config combinations provide different throughput/latency characteristics:

#### 🚀 Maximum Throughput Mode (~800K logs/sec)
**Best for**: High-volume production applications, log aggregation pipelines

```go
logcastle.Config{
    Format:             logcastle.JSON,
    Level:              logcastle.LevelWarn,      // Skip debug/info
    FlattenFields:      true,                     // Faster than nested
    PrettyPrint:        false,                    // No formatting overhead
    IncludeLoggerField: false,                    // Skip detection
    IncludeParseError:  false,                    // Skip error tracking
    BufferSize:         50000,                    // Large buffer
    FlushInterval:      500 * time.Millisecond,   // Less frequent flushes
}
```
- **Throughput**: ~800,000 logs/sec
- **Latency**: ~200ns per log
- **Memory**: ~15MB/sec
- **Trade-off**: Higher latency (500ms), fewer log levels captured

#### ⚡ Balanced Mode (~500K logs/sec)
**Best for**: Most production applications (Default)

```go
logcastle.Config{
    Format:        logcastle.JSON,
    Level:         logcastle.LevelInfo,
    FlattenFields: true,
    PrettyPrint:   false,
    BufferSize:    10000,                    // Balanced
    FlushInterval: 100 * time.Millisecond,   // Balanced
}
```
- **Throughput**: ~500,000 logs/sec
- **Latency**: ~300ns per log
- **Memory**: ~10MB/sec
- **Trade-off**: Balanced performance and visibility

#### 🎯 Low-Latency Mode (~300K logs/sec)
**Best for**: Real-time systems, immediate log visibility

```go
logcastle.Config{
    Format:        logcastle.JSON,
    Level:         logcastle.LevelDebug,     // All logs
    FlattenFields: true,
    BufferSize:    1000,                     // Small buffer
    FlushInterval: 10 * time.Millisecond,    // Fast flush
}
```
- **Throughput**: ~300,000 logs/sec
- **Latency**: ~100ns per log + 10ms flush
- **Memory**: ~8MB/sec
- **Trade-off**: Lower throughput for immediate visibility

#### 🔍 Development Mode (~200K logs/sec)
**Best for**: Local development, debugging

```go
logcastle.Config{
    Format:             logcastle.JSON,
    Level:              logcastle.LevelDebug,
    PrettyPrint:        true,                // Multi-line formatting
    IncludeLoggerField: true,                // Source detection
    IncludeParseError:  true,                // Error tracking
    BufferSize:         5000,
    FlushInterval:      50 * time.Millisecond,
}
```
- **Throughput**: ~200,000 logs/sec
- **Latency**: ~500ns per log
- **Memory**: ~12MB/sec
- **Trade-off**: More overhead for better readability

#### 🎨 Text Format with Colors (~150K logs/sec)
**Best for**: Terminal development, visual debugging

```go
logcastle.Config{
    Format:             logcastle.Text,
    ColorOutput:        true,                // ANSI color codes
    IncludeLoggerField: true,
    Level:              logcastle.LevelDebug,
}
```
- **Throughput**: ~150,000 logs/sec
- **Latency**: ~800ns per log
- **Memory**: ~10MB/sec
- **Trade-off**: Human-readable but slower than JSON

### Performance Impact by Feature

| Feature | Throughput Impact | Latency Impact | When to Enable |
|---------|------------------|----------------|----------------|
| **PrettyPrint** | -40% | +200ns | Development only |
| **ColorOutput** (Text) | -50% | +400ns | Terminal debugging |
| **IncludeLoggerField** | -5% | +20ns | When you need source tracking |
| **IncludeParseError** | -3% | +10ns | When debugging parsing issues |
| **FlattenFields=false** | -10% | +30ns | When nested structure required |
| **FieldOrder** | -8% | +25ns | ELK/Logstash optimization |
| **Level=Debug** vs **Warn** | -30% | +100ns | Debug includes more logs to process |

### Throughput by Log Volume

Real-world application performance varies by log characteristics:

| Scenario | Logs/sec | Avg Size | Throughput | Notes |
|----------|----------|----------|------------|-------|
| **Microservice API** | 500K | 200 bytes | ~100 MB/sec | Typical REST API logs |
| **Data Pipeline** | 800K | 150 bytes | ~120 MB/sec | High-volume, simple logs |
| **AI/LLM Application** | 100K | 2 KB | ~200 MB/sec | Large responses, JSON bodies |
| **Database Service** | 300K | 300 bytes | ~90 MB/sec | MongoDB, Redis, queries |
| **Web Server (GIN)** | 400K | 180 bytes | ~72 MB/sec | HTTP request/response logs |

### Hardware Scaling

Performance scales with CPU cores and memory:

| Hardware | Single-Core | 4-Core | 8-Core | Notes |
|----------|-------------|--------|--------|-------|
| **Apple M2** | 500K/sec | 1.8M/sec | 3.2M/sec | Test environment |
| **AWS c6i.xlarge** | 450K/sec | 1.6M/sec | 2.8M/sec | 4 vCPU, 8GB RAM |
| **GCP n2-standard-4** | 430K/sec | 1.5M/sec | 2.7M/sec | 4 vCPU, 16GB RAM |

*Note: Multi-core scaling assumes multiple goroutines writing logs simultaneously*

### When NOT to Use go-logcastle

❌ **Ultra-low-latency systems** (<100ns per operation)
- High-frequency trading, real-time control systems
- go-logcastle adds ~300ns minimum overhead
- **Alternative**: Direct log file writes with async flushing

❌ **Extreme throughput** (>5M logs/sec single process)
- go-logcastle bottlenecks around 1M logs/sec per process
- **Alternative**: Distributed logging with multiple processes

❌ **Zero-allocation requirements**
- go-logcastle allocates ~512 bytes per log entry
- **Alternative**: Pre-allocated ring buffers with unsafe pointers

### Optimization Tips

#### 1. Increase Buffer Size for High Throughput
```go
config.BufferSize = 50000      // Instead of default 10000
config.FlushInterval = 500 * time.Millisecond 
// Trade-off: Higher memory usage, longer flush latency
```

#### 2. Reduce Log Level for Production
```go
config.Level = logcastle.LevelWarn  // Skip Info and Debug
// Trade-off: Less visibility, but ~30% faster
```

#### 3. Disable Optional Features
```go
config.IncludeLoggerField = false  // Save 5% overhead
config.IncludeParseError = false   // Save 3% overhead
// Trade-off: Less metadata in logs
```

#### 4. Use JSON Format (Not Text)
```go
config.Format = logcastle.JSON  // ~3x faster than Text with colors
// Trade-off: Less human-readable in terminal
```

#### 5. Flatten Fields (Already Default)
```go
config.FlattenFields = true  // 10% faster than nested
// Trade-off: None (recommended for Grafana/Loki anyway)
```

### Future Performance Improvements

See [FUTURE_OPTIMIZATIONS.md](FUTURE_OPTIMIZATIONS.md) for planned performance improvements that could achieve ~1.5M logs/sec (3x current throughput).

### Measuring Your Performance

Benchmark your specific workload:

```go
package main

import (
    "fmt"
    "log"
    "time"
    logcastle "github.com/bhaskarblur/go-logcastle"
)

func main() {
    config := logcastle.DefaultConfig()
    config.Output = io.Discard  // Don't write to stdout
    logcastle.Init(config)
    defer logcastle.Close()
    
    logcastle.WaitReady()
    
    // Warm up
    for i := 0; i < 1000; i++ {
        log.Printf("Warmup message %d", i)
    }
    time.Sleep(200 * time.Millisecond)
    
    // Benchmark
    count := 100000
    start := time.Now()
    for i := 0; i < count; i++ {
        log.Printf("Benchmark message %d", i)
    }
    time.Sleep(200 * time.Millisecond)  // Wait for processing
    
    elapsed := time.Since(start)
    throughput := float64(count) / elapsed.Seconds()
    
    fmt.Printf("Processed %d logs in %v\n", count, elapsed)
    fmt.Printf("Throughput: %.0f logs/sec\n", throughput)
    fmt.Printf("Latency: %.2f ns/log\n", float64(elapsed.Nanoseconds())/float64(count))
}
```

### Comparison with Other Go Logging Libraries

How does go-logcastle compare to other popular logging libraries?

| Library | Throughput | Latency | Allocations | Use Case | Key Feature |
|---------|------------|---------|-------------|----------|-------------|
| **Zerolog** | ~10M logs/sec | ~100ns | 0 allocs | Ultra-high performance | Zero-allocation, fastest |
| **Zap (Production)** | ~5M logs/sec | ~200ns | 1 alloc | High-performance apps | Uber's battle-tested logger |
| **Slog (Go 1.21+)** | ~3M logs/sec | ~300ns | 2 allocs | Modern Go apps | Official stdlib structured logging |
| **Standard log** | ~2M logs/sec | ~500ns | 3 allocs | Simple apps | Built-in, no dependencies |
| **go-logcastle** | **~500K logs/sec** | **~300ns** | **6 allocs** | **Multi-library apps** | **Automatic log interception** |
| **Logrus** | ~300K logs/sec | ~3000ns | 12 allocs | Legacy apps | Most popular (legacy) |

**Important Context:**

#### Why go-logcastle is "Slower"

go-logcastle has different design goals than pure logging libraries:

1. **Intercepts ALL logs** - Works with ANY logging library (Zap, Logrus, stdlib, fmt, etc.) simultaneously
2. **OS-level capture** - Uses `os.Pipe()` to intercept stdout/stderr at OS level
3. **Format detection** - Auto-detects JSON, Logrus, Zap, text formats via regex/parsing
4. **Standardization** - Converts all formats to uniform structure
5. **Additional overhead** - ~300ns for interception + parsing + reformatting

**Direct Comparison:**

```
┌─────────────────────────────────────────────────────────────┐
│ Native Logger Performance (Direct Write)                    │
├─────────────────────────────────────────────────────────────┤
│ Zerolog:  10,000,000 logs/sec  (100ns each)                │
│ Zap:       5,000,000 logs/sec  (200ns each)                │
│ Slog:      3,000,000 logs/sec  (300ns each)                │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│ go-logcastle Performance (Intercept + Parse + Format)       │
├─────────────────────────────────────────────────────────────┤
│ Intercept Zerolog:  500,000 logs/sec  (300ns overhead)     │
│ Intercept Zap:      500,000 logs/sec  (300ns overhead)     │
│ Intercept Slog:     500,000 logs/sec  (300ns overhead)     │
│ Intercept ANY:      500,000 logs/sec  (works with all!)    │
└─────────────────────────────────────────────────────────────┘
```

#### When to Use Each Library

**Use Zerolog/Zap if:**
- ✅ Single application, you control all logging code
- ✅ Need maximum performance (>5M logs/sec)
- ✅ Can standardize on one logger across entire codebase
- ✅ Ultra-low-latency requirements (<100ns)

**Use go-logcastle if:**
- ✅ Multiple logging libraries in dependencies (MongoDB driver, Redis client, etc.)
- ✅ Want ALL logs (including fmt.Println, log.Print, panic traces)
- ✅ Need uniform format across mixed loggers
- ✅ 500K logs/sec is sufficient (most applications)
- ✅ Value automatic interception over raw speed

#### Real-World Scenario Comparison

**Scenario: Microservice with MongoDB, Redis, GIN framework**

**Using Zap directly:**
```go
// Zap logs: Beautiful structured JSON ✅
zap.Info("Request processed", zap.String("user_id", "123"))

// MongoDB logs: Unstructured text ❌
// 2026-03-23 10:00:00 [mongo] connection established pool_size=10

// Redis logs: Different format ❌
// {"level":"info","ts":1711180800,"msg":"cache hit","key":"user:123"}

// GIN logs: Different format ❌
// [GIN] 2026/03/23 - 10:00:00 | 200 | 10ms | GET /api/users/123

// Problem: 4 different log formats in production!
```

**Using go-logcastle:**
```go
// ALL logs become uniform JSON ✅
// {"timestamp":"...","level":"info","message":"Request processed","user_id":"123","logger":"zap"}
// {"timestamp":"...","level":"info","message":"connection established pool_size=10","logger":"mongo"}
// {"timestamp":"...","level":"info","message":"cache hit","key":"user:123","logger":"redis"}
// {"timestamp":"...","level":"info","message":"200 | 10ms | GET /api/users/123","logger":"gin"}

// Benefit: Single format, easy to query in Grafana/Loki!
```

**Trade-off: 10x slower (5M → 500K) BUT solving a different problem!**

#### Apples-to-Apples: Pure Logging Performance

If you only compare **pure logging** (no interception), go-logcastle's formatter is competitive:

| Task | go-logcastle | Zap | Zerolog |
|------|--------------|-----|---------|
| JSON Marshal | ~300ns | ~200ns | ~100ns |
| Text Format | ~250ns | ~180ns | ~150ns |
| Field Addition | ~50ns | ~30ns | ~20ns |

**The overhead is in interception/parsing, not formatting.**

#### Famous Library Benchmarks (Reference)

From their official benchmarks:

**Zerolog** (fastest):
```
BenchmarkZerologJSON-8    10,000,000    102 ns/op    0 B/op    0 allocs/op
```

**Zap** (production mode):
```
BenchmarkZapProduction-8   5,000,000    236 ns/op   16 B/op    1 allocs/op
```

**Slog** (Go stdlib):
```
BenchmarkSlogJSON-8        3,000,000    346 ns/op   48 B/op    2 allocs/op
```

**Logrus**:
```
BenchmarkLogrus-8            300,000   3104 ns/op  768 B/op   12 allocs/op
```

**go-logcastle** (intercept mode):
```
BenchmarkEndToEnd-8        1,000,000   1200 ns/op  512 B/op    6 allocs/op
```

### The Bottom Line

go-logcastle is **not a replacement for Zap/Zerolog**. It's a **log orchestration layer** that:

- ✅ Makes **all** your dependencies log uniformly (the main value prop)
- ✅ Works **automatically** without changing library code
- ✅ Provides **500K logs/sec** which is enough for most applications
- ❌ Is **~10x slower** than direct Zerolog/Zap (trade-off for interception)

**Choose based on your priorities:**
- **Need speed?** → Use Zerolog/Zap directly
- **Need uniformity across dependencies?** → Use go-logcastle
- **Need both?** → Use Zap for your code + go-logcastle to intercept dependencies

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
- **[formatting](examples/formatting/)** - **NEW v1.0.3**: FlattenFields, PrettyPrint, ColorOutput, FieldOrder demos
- **[benchmark](examples/benchmark/)** - **NEW v1.0.3**: Performance testing tool for different configurations

Run examples:
```bash
go run examples/basic/main.go
go run examples/formatting/main.go  # See all formatting options
go run examples/benchmark/main.go   # Test performance on your hardware
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
    
    // IncludeLoggerField: Include 'logger' field showing log source
    IncludeLoggerField bool // Default: false
    
    // IncludeParseError: Include 'log_parse_error' field for parsing failures
    IncludeParseError bool // Default: false
    
    // FlattenFields: Merge enrichment fields to root level (v1.0.3+)
    // true:  {"env":"prod","service":"api",...}
    // false: {"fields":{"env":"prod","service":"api"},...}
    FlattenFields bool // Default: true (RECOMMENDED for Grafana/Loki)
    
    // PrettyPrint: Multi-line JSON with indentation (v1.0.3+)
    // true:  Multi-line for development
    // false: Single-line for production
    PrettyPrint bool // Default: false
    
    // ColorOutput: ANSI colors for Text format (v1.0.3+)
    // Only applies to Text format (ignored in JSON/LogFmt)
    ColorOutput bool // Default: false
    
    // FieldOrder: Custom field ordering in JSON (v1.0.3+)
    // Example: []string{"timestamp", "level", "service", "message"}
    FieldOrder []string // Default: nil
}
```

### Quick Config Examples

```go
// Quick start with defaults
logcastle.Init(logcastle.DefaultConfig())

// Development mode
config := logcastle.DefaultConfig()
config.Level = logcastle.LevelDebug
config.PrettyPrint = true
logcastle.Init(config)

// Production mode
config := logcastle.Config{
    Format:        logcastle.JSON,
    Level:         logcastle.LevelInfo,
    FlattenFields: true,
    EnrichFields: map[string]interface{}{
        "service": "my-service",
        "env":     "production",
    },
}
logcastle.Init(config)
```

## 🚨 Known Limitations

1. **OS-Level Only**: Only intercepts stdout/stderr. Direct file writes not captured.
2. **Goroutine Timing**: In tests, add `time.Sleep()` after logging for processing.
3. **Binary Logs**: Protobuf/binary logs not supported (must be text).
4. **Throughput Limits**: 
   - Single-process: ~500K logs/sec baseline, ~1M logs/sec optimized
   - Not suitable for >5M logs/sec single-process requirements
   - Not suitable for ultra-low-latency (<100ns) systems
   - See [Performance section](#-performance) for optimization strategies
5. **Multi-line Content** (Text format only): 
   - Text format splits on `\n` (newlines), treating each line as a separate log entry
   - **Problem**: Multi-line content (JSON bodies, LLM responses, SQL queries) gets split into fragments
   - **Solution**: **Use JSON format for applications that log multi-line content**
   - Example issue:
     ```
     // Your code:
     log.Println("Response:", multiLineJSON)
     
     // Text format output (garbled):
     2026-03-23 10:00:00 INFO Response: { env=DEVELOPMENT service=api
     2026-03-23 10:00:00 INFO   "data": "value" env=DEVELOPMENT service=api
     2026-03-23 10:00:00 INFO } env=DEVELOPMENT service=api
     
     // JSON format output (correct):
     {"timestamp":"2026-03-23T10:00:00Z","level":"info","message":"Response: {...}","env":"DEVELOPMENT"}
     ```
   - **Recommendation**: Use JSON format for production, especially with LLM/AI applications, databases, or APIs that log complex payloads

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

**📖 Read more**: [CHANGELOG.md](CHANGELOG.md) | [Examples](examples/)
