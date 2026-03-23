package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/bhaskarblur/go-logcastle"
)

func TestStdlibTimestampStripping(t *testing.T) {
	logcastle.Reset()
	defer logcastle.Close()

	var buf bytes.Buffer

	config := logcastle.Config{
		Output:             &buf,
		Format:             logcastle.JSON,
		Level:              logcastle.LevelInfo,
		IncludeLoggerField: true,
	}

	logcastle.Init(config)
	logcastle.WaitReady()
	time.Sleep(100 * time.Millisecond)

	// Test message - log.Println adds stdlib timestamp prefix automatically
	testMessage := "This is a test message from stdlib log"
	log.Println(testMessage)

	time.Sleep(100 * time.Millisecond)
	logcastle.Close()

	output := buf.String()
	if output == "" {
		t.Fatal("No output captured")
	}

	var logEntry map[string]interface{}
	if err := json.Unmarshal([]byte(output), &logEntry); err != nil {
		t.Fatalf("Failed to parse JSON: %v\nOutput: %s", err, output)
	}

	message, ok := logEntry["message"].(string)
	if !ok {
		t.Fatal("message field not found or not a string")
	}

	// The message should NOT contain the stdlib log timestamp prefix
	// Stdlib format is "2026/03/23 19:37:55 message"
	// We want just "message" without the timestamp
	if strings.Contains(message, "2026/") || strings.Contains(message, "2025/") ||
		strings.Contains(message, "/03/") || strings.Contains(message, "/01/") {
		t.Errorf("Message still contains stdlib timestamp prefix. Full message: %s", message)
	}

	// The message should be just the actual message content
	if !strings.Contains(message, testMessage) {
		t.Errorf("Expected message to contain '%s', got: %s", testMessage, message)
	}

	// Verify logger is correctly identified as "log"
	logger, ok := logEntry["logger"].(string)
	if !ok || logger != "log" {
		t.Errorf("Expected logger='log', got: %v", logEntry["logger"])
	}

	fmt.Printf("✅ Success: Stdlib timestamp stripped. Message: %s\n", message)
}
