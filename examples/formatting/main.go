package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bhaskarblur/go-logcastle"
)

func main() {
	fmt.Println("=== Example 1: Terminal Development (Pretty + Colors) ===")
	terminalConfig()

	fmt.Println("\n\n=== Example 2: Production Grafana/Loki (Flattened, Single-line) ===")
	grafanaConfig()

	fmt.Println("\n\n=== Example 3: Custom Field Ordering (ELK/Logstash format) ===")
	elkConfig()

	fmt.Println("\n\n=== Example 4: Text Format with Colors ===")
	textColorConfig()
}

// terminalConfig: Development mode - pretty printed, colored output
func terminalConfig() {
	logcastle.Reset()

	config := logcastle.Config{
		Format:        logcastle.JSON,
		Level:         logcastle.LevelInfo,
		Output:        os.Stdout,
		FlattenFields: true,  // Flatten enrichment fields
		PrettyPrint:   true,  // Multi-line JSON with indentation
		ColorOutput:   false, // JSON doesn't use colors
		EnrichFields: map[string]interface{}{
			"env":     "development",
			"service": "api-server",
			"version": "1.0.0",
		},
	}

	logcastle.Init(config)
	defer logcastle.Close()
	logcastle.WaitReady()

	log.Println("Server starting on port 8080")
	fmt.Println("User authentication successful")

	time.Sleep(50 * time.Millisecond)
}

// grafanaConfig: Production mode - single-line, flattened for Grafana/Loki
func grafanaConfig() {
	logcastle.Reset()

	config := logcastle.Config{
		Format:        logcastle.JSON,
		Level:         logcastle.LevelInfo,
		Output:        os.Stdout,
		FlattenFields: true,  // CRITICAL: Flatten for Grafana label extraction
		PrettyPrint:   false, // Single-line for log aggregation
		ColorOutput:   false, // No ANSI codes in production
		EnrichFields: map[string]interface{}{
			"env":       "production",
			"service":   "payment-service",
			"region":    "us-east-1",
			"pod":       "payment-7f8c9d-xyz",
			"namespace": "prod",
		},
	}

	logcastle.Init(config)
	defer logcastle.Close()
	logcastle.WaitReady()

	log.Println("Processing payment transaction")
	fmt.Printf("Transaction ID: %s completed\n", "txn_12345")

	time.Sleep(50 * time.Millisecond)
}

// elkConfig: ELK/Logstash format with custom field ordering
func elkConfig() {
	logcastle.Reset()

	config := logcastle.Config{
		Format:        logcastle.JSON,
		Level:         logcastle.LevelInfo,
		Output:        os.Stdout,
		FlattenFields: true,
		PrettyPrint:   true, // For visibility in this example
		FieldOrder: []string{
			"timestamp", // Show timestamp first
			"level",     // Then log level
			"service",   // Then service name (from EnrichFields)
			"env",       // Then environment
			"message",   // Then the actual log message
		},
		EnrichFields: map[string]interface{}{
			"env":      "staging",
			"service":  "user-service",
			"cluster":  "k8s-staging",
			"hostname": "node-42",
		},
	}

	logcastle.Init(config)
	defer logcastle.Close()
	logcastle.WaitReady()

	log.Println("User login successful")
	fmt.Println("Database connection established")

	time.Sleep(50 * time.Millisecond)
}

// textColorConfig: Text format with ANSI colors for terminal
func textColorConfig() {
	logcastle.Reset()

	config := logcastle.Config{
		Format:             logcastle.Text, // Text format instead of JSON
		Level:              logcastle.LevelDebug,
		Output:             os.Stdout,
		ColorOutput:        true, // Enable ANSI colors
		IncludeLoggerField: true, // Show logger source
		EnrichFields: map[string]interface{}{
			"service": "auth-api",
		},
	}

	logcastle.Init(config)
	defer logcastle.Close()
	logcastle.WaitReady()

	// Different log levels will show in different colors
	log.Println("Debug: Verbose debugging info")
	log.Println("INFO: Application started")
	log.Println("WARN: High memory usage detected")
	fmt.Println("ERROR: Failed to connect to database")

	time.Sleep(50 * time.Millisecond)
}