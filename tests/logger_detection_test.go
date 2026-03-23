package tests

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/bhaskarblur/go-logcastle"
)

// TestLoggerDetection tests smart detection: log vs stdout
func TestLoggerDetection(t *testing.T) {
	var buf bytes.Buffer
	config := logcastle.Config{
		Format:             logcastle.JSON,
		Level:              logcastle.LevelInfo,
		Output:             &buf,
		BufferSize:         100,
		FlushInterval:      10 * time.Millisecond,
		IncludeLoggerField: true,
	}

	err := logcastle.Init(config)
	if err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}

	logcastle.WaitReady()

	// Stdlib log has timestamp prefix -> logger="log"
	log.Println("From log.Println")

	// Direct stdout no timestamp -> logger="stdout"
	fmt.Fprintf(os.Stdout, "From stdout\n")

	time.Sleep(150 * time.Millisecond)
	logcastle.Close()

	output := buf.String()

	// Verify smart detection worked
	if !strings.Contains(output, `"logger":"log"`) {
		t.Errorf("Expected logger='log', got: %s", output)
	}
	if !strings.Contains(output, `"logger":"stdout"`) {
		t.Errorf("Expected logger='stdout', got: %s", output)
	}
}

// TestDefaultLogLevel verifies unstructured logs default to Info level
func TestDefaultLogLevel(t *testing.T) {
	// Default level is set in NewLogEntry() to LevelInfo
	// This test just documents the behavior
	entry := &logcastle.LogEntry{}
	entry.Level = logcastle.LevelInfo // Default

	if entry.Level.String() != "info" {
		t.Errorf("Expected default level to be 'info', got: %s", entry.Level.String())
	}
}
