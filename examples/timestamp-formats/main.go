package main

import (
	"fmt"
	"log"
	"os"
	"time"

	logcastle "github.com/bhaskarblur/go-logcastle"
)

func main() {
	fmt.Println("=== Timestamp Format Examples ===")
	fmt.Println()

	// Example 1: Default (RFC3339Nano)
	fmt.Println("1. Default format (RFC3339Nano):")
	runWithFormat(logcastle.TimestampFormatRFC3339Nano, "", "Default ISO format with nanoseconds")

	// Example 2: RFC3339 (second precision)
	fmt.Println("\n2. RFC3339 (second precision):")
	runWithFormat(logcastle.TimestampFormatRFC3339, "", "Standard RFC3339")

	// Example 3: RFC3339 with milliseconds
	fmt.Println("\n3. RFC3339 with milliseconds:")
	runWithFormat(logcastle.TimestampFormatRFC3339Millis, "", "ISO format with milliseconds")

	// Example 4: Unix timestamp
	fmt.Println("\n4. Unix timestamp (seconds since epoch):")
	runWithFormat(logcastle.TimestampFormatUnix, "", "Unix timestamp")

	// Example 5: Unix milliseconds
	fmt.Println("\n5. Unix milliseconds:")
	runWithFormat(logcastle.TimestampFormatUnixMilli, "", "Unix timestamp in milliseconds")

	// Example 6: DateTime format
	fmt.Println("\n6. Human-readable DateTime:")
	runWithFormat(logcastle.TimestampFormatDateTime, "", "Human-friendly format")

	// Example 7: Custom format
	fmt.Println("\n7. Custom format (YYYY-MM-DD HH:MM:SS):")
	runWithFormat(logcastle.TimestampFormatCustom, "2006-01-02 15:04:05", "Custom format")

	fmt.Println("\n=== All formats demonstrated! ===")
}

func runWithFormat(tsFormat logcastle.TimestampFormat, customFormat, description string) {
	config := logcastle.Config{
		Format:                logcastle.JSON,
		Level:                 logcastle.LevelInfo,
		Output:                os.Stdout,
		TimestampFormat:       tsFormat,
		CustomTimestampFormat: customFormat,
		BufferSize:            1000,
		FlushInterval:         50 * time.Millisecond,
	}

	err := logcastle.Init(config)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer logcastle.Close()

	logcastle.WaitReady()

	log.Printf("Example: %s", description)
	time.Sleep(100 * time.Millisecond)
}
