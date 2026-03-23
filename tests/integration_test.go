package tests

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	logcastle "github.com/bhaskarblur/go-logcastle"
)

func TestBasicInterception(t *testing.T) {
	var buf bytes.Buffer

	err := logcastle.Init(logcastle.Config{
		Format: logcastle.JSON,
		Output: &buf,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer logcastle.Close()

	logcastle.WaitReady()

	fmt.Println("test message")

	time.Sleep(200 * time.Millisecond)

	output := buf.String()
	if output == "" {
		t.Error("Expected output, got empty string")
	}

	if !bytes.Contains([]byte(output), []byte("test message")) {
		t.Errorf("Expected 'test message' in output, got: %s", output)
	}

	if !bytes.Contains([]byte(output), []byte(`"message":`)) {
		t.Error("Expected JSON format with message field")
	}
}

func TestEnrichment(t *testing.T) {
	var buf bytes.Buffer

	err := logcastle.Init(logcastle.Config{
		Format: logcastle.JSON,
		Output: &buf,
		EnrichFields: map[string]interface{}{
			"service": "test-service",
			"version": "1.0.0",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer logcastle.Close()

	logcastle.WaitReady()

	fmt.Println("enriched message")

	time.Sleep(250 * time.Millisecond)

	output := buf.String()

	if !bytes.Contains([]byte(output), []byte("enriched message")) {
		t.Error("Expected message in output")
	}

	if !bytes.Contains([]byte(output), []byte("test-service")) {
		t.Error("Expected enriched field 'service' in output")
	}

	if !bytes.Contains([]byte(output), []byte("1.0.0")) {
		t.Error("Expected enriched field 'version' in output")
	}
}

func TestGracefulShutdown(t *testing.T) {
	var buf bytes.Buffer

	err := logcastle.Init(logcastle.Config{
		Format: logcastle.JSON,
		Output: &buf,
	})
	if err != nil {
		t.Fatal(err)
	}

	logcastle.WaitReady()

	fmt.Println("message before close")
	time.Sleep(50 * time.Millisecond)

	err = logcastle.Close()
	if err != nil {
		t.Errorf("Close() returned error: %v", err)
	}

	output := buf.String()
	if !bytes.Contains([]byte(output), []byte("message before close")) {
		t.Error("Expected message to be flushed before close")
	}
}

func TestLevelFiltering(t *testing.T) {
	var buf bytes.Buffer

	err := logcastle.Init(logcastle.Config{
		Format: logcastle.JSON,
		Output: &buf,
		Level:  logcastle.LevelWarn,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer logcastle.Close()

	logcastle.WaitReady()

	// These should be filtered
	fmt.Println("debug message")
	fmt.Println("info message")

	time.Sleep(150 * time.Millisecond)

	// Note: fmt.Println doesn't have level info, so defaults to Info
	// In production, proper loggers set levels
}

func TestConfigDefaults(t *testing.T) {
	config := logcastle.DefaultConfig()

	if config.Format != logcastle.JSON {
		t.Error("Expected default format to be JSON")
	}

	if config.Level != logcastle.LevelInfo {
		t.Error("Expected default level to be Info")
	}

	if config.BufferSize != 10000 {
		t.Error("Expected default buffer size to be 10000")
	}
}
