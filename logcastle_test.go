package logcastle

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestLogCastle_BasicInterception(t *testing.T) {
	var buf bytes.Buffer

	err := Init(Config{
		Format: JSON,
		Output: &buf,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer Close()

	// Wait for interception to be ready
	WaitReady()

	// Write some logs - use fmt.Println which writes to os.Stdout
	fmt.Println("test message")

	// Wait for processing and flush
	time.Sleep(200 * time.Millisecond)

	output := buf.String()
	if output == "" {
		t.Error("Expected output, got empty string")
	}

	if !bytes.Contains([]byte(output), []byte("test message")) {
		t.Errorf("Expected 'test message' in output, got: %s", output)
	}

	// Verify it's valid JSON
	if !bytes.Contains([]byte(output), []byte(`"message":`)) {
		t.Error("Expected JSON format with message field")
	}
}

func TestParser_JSON(t *testing.T) {
	parser := NewParser()

	input := []byte(`{"timestamp":"2026-03-23T10:00:00Z","level":"info","message":"test"}`)
	entry := parser.Parse(input)

	if entry.Message != "test" {
		t.Errorf("Expected message 'test', got '%s'", entry.Message)
	}

	if entry.Level != LevelInfo {
		t.Errorf("Expected level Info, got %v", entry.Level)
	}
}

func TestFormatter_JSON(t *testing.T) {
	formatter := NewFormatter(JSON, TimestampFormatRFC3339Nano, "")
	entry := NewLogEntry()
	entry.Message = "test message"
	entry.Level = LevelInfo
	entry.Logger = "test"

	output := formatter.Format(entry)
	if output == nil {
		t.Fatal("Expected non-nil output")
	}

	if !bytes.Contains(output, []byte("test message")) {
		t.Errorf("Expected 'test message' in output, got: %s", output)
	}
}

func BenchmarkParse(b *testing.B) {
	parser := NewParser()
	line := []byte(`{"timestamp":"2026-03-23T10:00:00Z","level":"info","message":"test","user":"alice"}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser.Parse(line)
	}
}

func BenchmarkFormat(b *testing.B) {
	formatter := NewFormatter(JSON, TimestampFormatRFC3339Nano, "")
	entry := NewLogEntry()
	entry.Message = "test message"
	entry.Fields = map[string]interface{}{
		"user": "alice",
		"id":   123,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		formatter.Format(entry)
	}
}
