package benchmarks

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"time"

	logcastle "github.com/yourusername/go-logcastle"
	"github.com/yourusername/go-logcastle/formatter"
	"github.com/yourusername/go-logcastle/parser"
)

// Benchmark parsing performance
func BenchmarkParseJSON(b *testing.B) {
	p := parser.NewParser()
	line := []byte(`{"timestamp":"2026-03-23T10:00:00Z","level":"info","message":"test","user":"alice","id":123}`)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		p.Parse(line)
	}
}

func BenchmarkParseLogrusText(b *testing.B) {
	p := parser.NewParser()
	line := []byte(`time="2026-03-23T10:00:00Z" level=info msg="test message" user=alice`)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		p.Parse(line)
	}
}

// Benchmark formatting performance
func BenchmarkFormatJSON(b *testing.B) {
	f := formatter.NewFormatter(formatter.JSON)
	entry := logcastle.NewLogEntry()
	entry.Message = "test message"
	entry.Fields = map[string]interface{}{
		"user": "alice",
		"id":   123,
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		f.Format(entry)
	}
}

func BenchmarkFormatText(b *testing.B) {
	f := formatter.NewFormatter(formatter.Text)
	entry := logcastle.NewLogEntry()
	entry.Message = "test message"
	entry.Fields = map[string]interface{}{
		"user": "alice",
		"id":   123,
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		f.Format(entry)
	}
}

func BenchmarkFormatCustom(b *testing.B) {
	template := "{{timestamp}} [{{level}}] {{logger}}: {{message}}"
	f, err := formatter.NewCustomFormatter(template)
	if err != nil {
		b.Fatal(err)
	}

	entry := logcastle.NewLogEntry()
	entry.Message = "test message"
	entry.Logger = "bench"
	entry.Level = logcastle.LevelInfo

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		f.Format(entry)
	}
}

// Benchmark end-to-end throughput
func BenchmarkEndToEnd(b *testing.B) {
	var buf bytes.Buffer

	err := logcastle.Init(logcastle.Config{
		Format:        logcastle.JSON,
		Output:        &buf,
		BufferSize:    100000,
		FlushInterval: 1 * time.Second,
	})
	if err != nil {
		b.Fatal(err)
	}
	defer logcastle.Close()

	logcastle.WaitReady()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		fmt.Printf("benchmark message %d\n", i)
	}

	time.Sleep(100 * time.Millisecond)
}

// Throughput test - measures logs/second
func BenchmarkThroughput1K(b *testing.B) {
	benchmarkThroughput(b, 1000)
}

func BenchmarkThroughput10K(b *testing.B) {
	benchmarkThroughput(b, 10000)
}

func BenchmarkThroughput100K(b *testing.B) {
	benchmarkThroughput(b, 100000)
}

func benchmarkThroughput(b *testing.B, numLogs int) {
	var buf bytes.Buffer

	err := logcastle.Init(logcastle.Config{
		Format:        logcastle.JSON,
		Output:        &buf,
		BufferSize:    numLogs,
		FlushInterval: 1 * time.Second,
	})
	if err != nil {
		b.Fatal(err)
	}
	defer logcastle.Close()

	logcastle.WaitReady()

	b.ResetTimer()
	start := time.Now()

	for i := 0; i < numLogs; i++ {
		fmt.Fprintf(os.Stdout, "throughput test message %d\n", i)
	}

	// Wait for all logs to be processed
	time.Sleep(500 * time.Millisecond)

	elapsed := time.Since(start)
	logsPerSecond := float64(numLogs) / elapsed.Seconds()

	b.ReportMetric(logsPerSecond, "logs/sec")
	b.ReportMetric(float64(buf.Len())/float64(numLogs), "bytes/log")
}

// Memory allocation benchmarks
func BenchmarkParseAllocations(b *testing.B) {
	p := parser.NewParser()
	line := []byte(`{"timestamp":"2026-03-23T10:00:00Z","level":"info","message":"test","user":"alice"}`)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		p.Parse(line)
	}
}

func BenchmarkFormatAllocations(b *testing.B) {
	f := formatter.NewFormatter(formatter.JSON)
	entry := logcastle.NewLogEntry()
	entry.Message = "test"

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		f.Format(entry)
	}
}
