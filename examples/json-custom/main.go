package main

import (
	"fmt"
	"os"
	"time"

	"github.com/bhaskarblur/go-logcastle/formatter"
)

func main() {
	fmt.Println("=== Custom JSON Formatter Demo ===\n")

	// Initialize runtime fields (do this once at app startup)
	formatter.InitRuntimeFields("production", map[string]string{
		"region":     "us-east-1",
		"cluster":    "prod-cluster-1",
		"datacenter": "dc1",
	})

	// Create custom JSON formatter
	jsonFormatter := formatter.NewJSONFormatter()

	// Set global fields (added to every log)
	jsonFormatter.SetGlobalField("service", "api-gateway")
	jsonFormatter.SetGlobalField("version", "2.1.0")
	jsonFormatter.SetGlobalField("build", "abc123def")

	// Enable runtime fields (hostname, PID, goroutines, etc.)
	jsonFormatter.IncludeRuntimeFields = true

	// Customize field order (most important fields first)
	jsonFormatter.FieldOrder = []string{
		"timestamp",
		"level",
		"service",
		"version",
		"environment",
		"message",
	}

	// Create sample log entry
	entry := &formatter.CustomLogEntry{
		Timestamp: time.Now(),
		Level:     "info",
		Logger:    "http",
		Message:   "Request processed successfully",
		Caller:    "handler.go:42",
		Fields: map[string]interface{}{
			"method":      "POST",
			"path":        "/api/users",
			"status_code": 201,
			"duration_ms": 45.3,
			"user_id":     "usr_123",
		},
		TraceID: "trace-abc-123",
		SpanID:  "span-def-456",
		Source:  "stdout",
	}

	fmt.Println("--- Default Compact JSON ---")
	output := jsonFormatter.Format(entry)
	os.Stdout.Write(output)

	fmt.Println("\n--- Pretty Printed JSON ---")
	jsonFormatter.PrettyPrint = true
	output = jsonFormatter.Format(entry)
	os.Stdout.Write(output)
	jsonFormatter.PrettyPrint = false

	fmt.Println("\n--- Adding More Global Fields at Runtime ---")
	jsonFormatter.SetGlobalField("deployment_id", "deploy-789")
	jsonFormatter.SetGlobalField("replica", 3)

	entry2 := &formatter.CustomLogEntry{
		Timestamp: time.Now(),
		Level:     "error",
		Logger:    "database",
		Message:   "Connection timeout",
		Fields: map[string]interface{}{
			"db_host": "postgres-primary",
			"timeout": "5s",
		},
	}

	output = jsonFormatter.Format(entry2)
	os.Stdout.Write(output)

	fmt.Println("\n--- Removing a Global Field ---")
	jsonFormatter.RemoveGlobalField("build")
	output = jsonFormatter.Format(entry2)
	os.Stdout.Write(output)

	fmt.Println("\n--- Without Runtime Fields ---")
	jsonFormatter.IncludeRuntimeFields = false
	output = jsonFormatter.Format(entry2)
	os.Stdout.Write(output)

	fmt.Println("\n=== Key Features Demonstrated ===")
	fmt.Println("✓ Global fields (service, version) added to all logs")
	fmt.Println("✓ Runtime fields (hostname, PID, goroutines) included")
	fmt.Println("✓ Custom field order for consistent output")
	fmt.Println("✓ Dynamic field management (add/remove at runtime)")
	fmt.Println("✓ Environment context (region, cluster) propagated")
	fmt.Println("✓ Pretty printing for debugging")
}
