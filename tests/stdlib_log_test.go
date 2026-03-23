package tests

import (
	"bytes"
	"log"
	"testing"
	"time"

	"github.com/bhaskarblur/go-logcastle"
)

// TestStdlibLogCapture verifies that stdlib log.Println() is properly captured
func TestStdlibLogCapture(t *testing.T) {
	logcastle.Reset() // Reset for test isolation
	var buf bytes.Buffer

	config := logcastle.Config{
		Format:        logcastle.JSON,
		Level:         logcastle.LevelInfo,
		Output:        &buf,
		BufferSize:    100,
		FlushInterval: 10 * time.Millisecond,
	}

	err := logcastle.Init(config)
	if err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}
	defer logcastle.Close()

	logcastle.WaitReady()

	// Test stdlib log package - this should be captured now!
	log.Println("Test message from stdlib log")
	log.Printf("Formatted message: %s", "test")

	// Wait for flush
	time.Sleep(50 * time.Millisecond)

	output := buf.String()

	if !bytes.Contains([]byte(output), []byte("Test message from stdlib log")) {
		t.Errorf("Expected stdlib log.Println to be captured, got: %s", output)
	}

	if !bytes.Contains([]byte(output), []byte("Formatted message: test")) {
		t.Errorf("Expected stdlib log.Printf to be captured, got: %s", output)
	}

	// Verify JSON format
	if !bytes.Contains([]byte(output), []byte(`"level":`)) {
		t.Error("Expected JSON format with level field")
	}
}
