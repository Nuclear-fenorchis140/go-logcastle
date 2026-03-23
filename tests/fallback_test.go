package tests

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	logcastle "github.com/bhaskarblur/go-logcastle"
)

func TestFallbackParsing(t *testing.T) {
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

	// Write unparseable log
	fmt.Println("This is completely unstructured text")

	time.Sleep(200 * time.Millisecond)

	output := buf.String()
	if output == "" {
		t.Error("Expected output even for unparseable log")
	}

	// Should contain the message
	if !bytes.Contains([]byte(output), []byte("unstructured text")) {
		t.Error("Expected message to be captured")
	}

	// Should contain parse error indicator
	if !bytes.Contains([]byte(output), []byte("log_parse_error")) {
		t.Error("Expected log_parse_error field for unparseable log")
	}
}

func TestMalformedJSON(t *testing.T) {
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

	// Write malformed JSON
	fmt.Println(`{"level":"info","message":"Missing closing brace"`)

	time.Sleep(200 * time.Millisecond)

	output := buf.String()
	if output == "" {
		t.Error("Expected output even for malformed JSON")
	}

	// Should still capture the content
	if !bytes.Contains([]byte(output), []byte("Missing closing brace")) {
		t.Error("Expected message to be captured from malformed JSON")
	}

	// Should contain parse error
	if !bytes.Contains([]byte(output), []byte("log_parse_error")) {
		t.Error("Expected log_parse_error field for malformed JSON")
	}
}

func TestTimestampFormats(t *testing.T) {
	formats := []struct {
		format   logcastle.TimestampFormat
		name     string
		contains string // Expected substring in output
	}{
		{logcastle.TimestampFormatRFC3339Nano, "RFC3339Nano", "T"},
		{logcastle.TimestampFormatRFC3339, "RFC3339", "T"},
		{logcastle.TimestampFormatUnix, "Unix", "17"},        // Unix timestamp starts with 17...
		{logcastle.TimestampFormatDateTime, "DateTime", " "}, // Contains space
	}

	for _, tc := range formats {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer

			err := logcastle.Init(logcastle.Config{
				Format:          logcastle.JSON,
				Output:          &buf,
				TimestampFormat: tc.format,
			})
			if err != nil {
				t.Fatal(err)
			}
			defer logcastle.Close()

			logcastle.WaitReady()

			fmt.Println("test message")
			time.Sleep(150 * time.Millisecond)

			output := buf.String()
			if output == "" {
				t.Error("Expected output")
			}

			// Check for timestamp field
			if !bytes.Contains([]byte(output), []byte("timestamp")) {
				t.Error("Expected timestamp field in output")
			}
		})
	}
}

func TestCustomTimestampFormat(t *testing.T) {
	var buf bytes.Buffer

	customFormat := "2006-01-02 15:04:05"
	err := logcastle.Init(logcastle.Config{
		Format:                logcastle.JSON,
		Output:                &buf,
		TimestampFormat:       logcastle.TimestampFormatCustom,
		CustomTimestampFormat: customFormat,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer logcastle.Close()

	logcastle.WaitReady()

	fmt.Println("custom timestamp test")
	time.Sleep(150 * time.Millisecond)

	output := buf.String()
	if output == "" {
		t.Error("Expected output")
	}

	// Should have timestamp in custom format (YYYY-MM-DD HH:MM:SS)
	if !bytes.Contains([]byte(output), []byte("timestamp")) {
		t.Error("Expected timestamp field")
	}
}
