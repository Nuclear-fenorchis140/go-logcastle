package main

import (
	"fmt"
	"os"
	"time"

	logcastle "github.com/yourusername/go-logcastle"
)

func main() {
	fmt.Println("=== Fallback Parsing Example ===")
	fmt.Println()
	fmt.Println("This example shows how go-logcastle handles logs that can't be parsed:")
	fmt.Println()

	// Initialize logcastle
	config := logcastle.Config{
		Format:        logcastle.JSON,
		Level:         logcastle.LevelInfo,
		Output:        os.Stdout,
		BufferSize:    1000,
		FlushInterval: 50 * time.Millisecond,
	}

	err := logcastle.Init(config)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer logcastle.Close()

	logcastle.WaitReady()

	fmt.Println("Valid JSON log:")
	fmt.Println(`{"level":"info","message":"This is valid JSON"}`)

	fmt.Println("\nMalformed JSON log (missing closing brace):")
	fmt.Println(`{"level":"info","message":"Missing closing brace"`)

	fmt.Println("\nPlain text log (no structure):")
	fmt.Println("Just some plain text without any structure")

	fmt.Println("\nPartially structured log:")
	fmt.Println("level=info but missing other fields")

	fmt.Println("\n--- Output (notice log_parse_error field) ---")

	time.Sleep(200 * time.Millisecond)

	fmt.Println("\n=== All logs captured, even unparseable ones! ===")
	fmt.Println("Check the 'log_parse_error' field to identify parsing issues.")
}
