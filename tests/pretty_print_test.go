package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/bhaskarblur/go-logcastle"
)

func TestPrettyPrint(t *testing.T) {
	logcastle.Reset()
	defer logcastle.Close()

	var buf bytes.Buffer

	config := logcastle.Config{
		Output:             &buf,
		Format:             logcastle.JSON,
		Level:              logcastle.LevelInfo,
		PrettyPrint:        true,
		FlattenFields:      false,
		IncludeLoggerField: false,
		IncludeParseError:  false,
		EnrichFields: map[string]interface{}{
			"service": "test-service",
		},
	}

	logcastle.Init(config)
	logcastle.WaitReady()
	time.Sleep(100 * time.Millisecond)

	fmt.Println("Pretty print test")

	time.Sleep(100 * time.Millisecond)
	logcastle.Close()

	output := buf.String()
	if output == "" {
		t.Fatal("No output captured")
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) <= 1 {
		t.Error("Expected multi-line output with PrettyPrint=true")
	}

	hasIndentation := false
	for _, line := range lines {
		if strings.HasPrefix(line, "  ") {
			hasIndentation = true
			break
		}
	}
	if !hasIndentation {
		t.Error("Expected indentation in pretty-printed JSON")
	}

	var logEntry map[string]interface{}
	if err := json.Unmarshal([]byte(output), &logEntry); err != nil {
		t.Fatalf("Pretty-printed output is not valid JSON: %v", err)
	}
}

func TestSingleLineJSON(t *testing.T) {
	logcastle.Reset()
	defer logcastle.Close()

	var buf bytes.Buffer

	config := logcastle.Config{
		Output:             &buf,
		Format:             logcastle.JSON,
		Level:              logcastle.LevelInfo,
		PrettyPrint:        false,
		FlattenFields:      false,
		IncludeLoggerField: false,
		IncludeParseError:  false,
	}

	logcastle.Init(config)
	logcastle.WaitReady()
	time.Sleep(100 * time.Millisecond)

	fmt.Println("Single line test")

	time.Sleep(100 * time.Millisecond)
	logcastle.Close()

	output := buf.String()
	if output == "" {
		t.Fatal("No output captured")
	}

	trimmed := strings.TrimSpace(output)
	lines := strings.Split(trimmed, "\n")
	if len(lines) != 1 {
		t.Error("Expected single-line output with PrettyPrint=false")
	}

	var logEntry map[string]interface{}
	if err := json.Unmarshal([]byte(output), &logEntry); err != nil {
		t.Fatalf("Output is not valid JSON: %v", err)
	}
}
