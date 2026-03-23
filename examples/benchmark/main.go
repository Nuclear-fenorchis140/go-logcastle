package main

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/bhaskarblur/go-logcastle"
)

func main() {
	fmt.Println("🚀 go-logcastle Performance Benchmark Tool")
	fmt.Println("=" + string(make([]byte, 60)))
	fmt.Println()

	configs := []struct {
		name   string
		config logcastle.Config
	}{
		{
			name: "Maximum Throughput Mode",
			config: logcastle.Config{
				Format:             logcastle.JSON,
				Level:              logcastle.LevelWarn,
				FlattenFields:      true,
				PrettyPrint:        false,
				IncludeLoggerField: false,
				IncludeParseError:  false,
				BufferSize:         50000,
				FlushInterval:      500 * time.Millisecond,
				Output:             io.Discard,
			},
		},
		{
			name: "Balanced Mode (Default)",
			config: logcastle.Config{
				Format:        logcastle.JSON,
				Level:         logcastle.LevelInfo,
				FlattenFields: true,
				PrettyPrint:   false,
				BufferSize:    10000,
				FlushInterval: 100 * time.Millisecond,
				Output:        io.Discard,
			},
		},
		{
			name: "Low-Latency Mode",
			config: logcastle.Config{
				Format:        logcastle.JSON,
				Level:         logcastle.LevelDebug,
				FlattenFields: true,
				BufferSize:    1000,
				FlushInterval: 10 * time.Millisecond,
				Output:        io.Discard,
			},
		},
		{
			name: "Development Mode",
			config: logcastle.Config{
				Format:             logcastle.JSON,
				Level:              logcastle.LevelDebug,
				PrettyPrint:        true,
				IncludeLoggerField: true,
				IncludeParseError:  true,
				BufferSize:         5000,
				FlushInterval:      50 * time.Millisecond,
				Output:             io.Discard,
			},
		},
		{
			name: "Text Format with Colors",
			config: logcastle.Config{
				Format:             logcastle.Text,
				ColorOutput:        true,
				IncludeLoggerField: true,
				Level:              logcastle.LevelDebug,
				Output:             io.Discard,
			},
		},
	}

	for i, cfg := range configs {
		fmt.Printf("[%d/%d] Testing: %s\n", i+1, len(configs), cfg.name)
		throughput := benchmarkConfig(cfg.config)
		fmt.Printf("       Result: %.0f logs/sec\n", throughput)
		fmt.Println()
		time.Sleep(500 * time.Millisecond) // Let CPU cool down
	}

	fmt.Println("=" + string(make([]byte, 60)))
	fmt.Println("✅ Benchmark complete!")
	fmt.Println()
	fmt.Println("💡 Tips:")
	fmt.Println("  - Use Maximum Throughput Mode for high-volume production")
	fmt.Println("  - Use Balanced Mode for most applications")
	fmt.Println("  - Use Development Mode for local debugging")
	fmt.Println("  - Avoid Text+Colors in production (slowest)")
}

func benchmarkConfig(config logcastle.Config) float64 {
	logcastle.Reset()

	if err := logcastle.Init(config); err != nil {
		panic(err)
	}
	defer logcastle.Close()

	logcastle.WaitReady()

	// Warmup
	for i := 0; i < 1000; i++ {
		log.Printf("Warmup message %d", i)
	}
	time.Sleep(200 * time.Millisecond)

	// Benchmark
	count := 100000
	start := time.Now()

	for i := 0; i < count; i++ {
		log.Printf("Benchmark message %d with some data: user_id=%d, action=view", i, i%1000)
	}

	// Wait for processing
	time.Sleep(200 * time.Millisecond)

	elapsed := time.Since(start)
	throughput := float64(count) / elapsed.Seconds()

	return throughput
}
