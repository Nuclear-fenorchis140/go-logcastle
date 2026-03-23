package tests

import (
	"bytes"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/bhaskarblur/go-logcastle"
)

// TestCleanOutputByDefault tests that logs are clean by default
func TestCleanOutputByDefault(t *testing.T) {
	logcastle.Reset() // Reset for test isolation
	var buf bytes.Buffer
	config := logcastle.DefaultConfig()
	config.Output = &buf
	config.BufferSize = 100
	config.FlushInterval = 10 * time.Millisecond

	err := logcastle.Init(config)
	if err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}

	logcastle.WaitReady()
	log.Println("Clean log message")
	time.Sleep(100 * time.Millisecond)
	logcastle.Close()

	output := buf.String()

	// Should NOT include logger or parse_error fields
	if strings.Contains(output, `"logger"`) {
		t.Errorf("Expected logger field omitted by default, got: %s", output)
	}
	if strings.Contains(output, `"log_parse_error"`) {
		t.Errorf("Expected log_parse_error field omitted by default, got: %s", output)
	}
	if !strings.Contains(output, "Clean log message") {
		t.Errorf("Expected message content, got: %s", output)
	}
}
