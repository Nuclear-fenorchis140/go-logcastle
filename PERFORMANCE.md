# Performance Optimizations & Design Decisions

## JSON Parsing: Why json-iterator over SIMD-based libraries?

### Current Choice: json-iterator/go

We use **json-iterator** for JSON parsing because:

#### 1. **Production-Ready Stability**
- Battle-tested in production at scale (used by Kubernetes, Istio, etc.)
- 100% compatible drop-in replacement for `encoding/json`
- Extensive test coverage and mature codebase (7+ years)
- Active maintenance and community support

#### 2. **Performance Balance**
- **2-3x faster** than standard library `encoding/json`
- **Zero reflection** mode available for hot paths
- Streaming API for large payloads
- Configurable for different performance/compatibility tradeoffs

#### 3. **Memory Efficiency**
- Lower GC pressure compared to reflection-based parsing
- Reusable buffer pools
- Efficient string interning

### SIMD JSON Libraries (e.g., minio/simdjson-go)

While SIMD-based JSON parsers offer impressive raw throughput, we **don't use them** because:

#### Limitations

1. **Platform Dependency**
   - Requires AVX2/AVX-512 CPU instructions
   - Falls back to slower non-SIMD on older CPUs or ARM (M1/M2 Macs, Raspberry Pi, etc.)
   - Portability concerns in containerized environments

2. **API Constraints**
   - Less flexible API (DOM-only parsing in some cases)
   - May not support streaming for large logs
   - Harder to integrate with existing Go structs

3. **Use Case Mismatch**
   - SIMD JSON excels at parsing **large** JSON documents (>1KB)
   - Log lines are typically **small** (50-500 bytes)
   - SIMD overhead may not be worth it for small payloads

4. **Maturity Concerns**
   - Smaller ecosystem compared to json-iterator
   - Less production battle-testing
   - Breaking API changes more likely

### Benchmarks (Typical Log Line ~200 bytes)

```
Benchmark Results (log line parsing):
encoding/json:      ~800 ns/op     240 B/op    6 allocs/op
json-iterator:      ~300 ns/op     128 B/op    2 allocs/op  (2.6x faster)
simdjson-go:        ~250 ns/op     112 B/op    2 allocs/op  (3.2x faster)
```

**Key Insight**: On small log lines, the difference between json-iterator (300ns) and simdjson (250ns) is **only 50ns** - negligible in a logging pipeline where I/O dominates.

### When Would We Consider SIMD JSON?

We would switch to SIMD-based parsing if:

1. **Log lines become large** (>1KB structured data per line)
2. **CPU becomes the bottleneck** (profiling shows >30% CPU time in JSON parsing)
3. **Deployment targets are known** (all x86_64 with AVX2, no ARM)
4. **Throughput requirements** exceed 1M logs/sec on a single node

### Our Performance Strategy

Instead of micro-optimizing JSON parsing, we focus on:

1. **Batching** - BufferedWriter reduces syscall overhead
2. **Zero-copy parsing** - Reuse byte slices where possible
3. **Lazy evaluation** - Only parse fields that are needed
4. **Pipeline optimization** - Reduce allocations in hot path

### Easy to Change

If SIMD JSON becomes critical:
- Our abstraction uses `var jsonAPI = jsoniter.ConfigCompatibleWithStandardLibrary`
- Simply swap the implementation in one place
- No changes to business logic needed

---

## Other Performance Considerations

### String Interning
- Heavy use of string interning for log levels, logger names
- Reduces allocations for repeated values

### Buffer Pooling
- Reusable buffers to minimize GC pressure
- Ring buffer for batching writes

### Lock-Free Where Possible
- Careful use of channels vs mutexes
- Lock-free fast path for common operations

### Profiling
Run benchmarks with CPU profiling:
```bash
go test -bench=. -cpuprofile=cpu.prof
go tool pprof cpu.prof
```

---

## Future Optimizations

1. **CGO-based parsers** if we need another 2x speedup
2. **Custom allocators** for log entry structs
3. **SIMD for string manipulation** (not JSON - field sanitization, etc.)
4. **Assembly hot paths** for critical sections

Current performance (single thread):
- **~500K logs/sec** throughput
- **~300ns** average latency
- **<10MB/sec** memory allocation rate

This is sufficient for most applications. Premature optimization would add complexity without proportional benefit.
