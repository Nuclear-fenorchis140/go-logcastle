package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/bhaskarblur/go-logcastle"
)

func TestFlattenFields(t *testing.T) {
	logcastle.Reset()
	defer logcastle.Close()

	var buf bytes.Buffer

	config := logcastle.Config{
		Output:             &buf,
		Format:             logcastle.JSON,
		Level:              logcastle.LevelInfo,
		FlattenFields:      true,
		IncludeLoggerField: false,
		IncludeParseError:  false,
		EnrichFields: map[string]interface{}{
			"service": "test-service",
			"env":     "testing",
		},
	}

	logcastle.Init(config)
	logcastle.WaitReady()
	time.Sleep(100 * time.Millisecond)

	fmt.Println("Test flatten fields")

	time.Sleep(100 * time.Millisecond)
	logcastle.Close()

	output := buf.String()
	if output == "" {
		t.Fatal("No output captured")
	}

	var logEntry map[string]interface{}
	if err := json.Unmarshal([]byte(output), &logEntry); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if _, exists := logEntry["fields"]; exists {
		t.Error("Fields should be flattened, but 'fields' key exists")
	}

	if service, ok := logEntry["service"].(string); !ok || service != "test-service" {
		t.Error("Expected service at root level")
	}

	if env, ok := logEntry["env"].(string); !ok || env != "testing" {
		t.Error("Expected env at root level")
	}
}

func TestNestedFields(t *testing.T) {
	logcastle.Reset()
	defer logcastle.Close()

	var buf bytes.Buffer

	config := logcastle.Config{
		Output:             &buf,
		Format:             logcastle.JSON,
		Level:              logcastle.LevelInfo,
		FlattenFields:      false,
		IncludeLoggerField: false,
		IncludeParseError:  false,
		EnrichFields: map[string]interface{}{
			"service": "test-service",
			"env":     "testing",
		},
	}

	logcastle.Init(config)
	logcastle.WaitReady()
	time.Sleep(100 * time.Millisecond)

	fmt.Println("Test nested fields")

	time.Sleep(100 * time.Millisecond)
	logcastle.Close()

	output := buf.String()
	if output == "" {
		t.Fatal("No output captured")
	}

	var logEntry map[string]interface{}
	if err := json.Unmarshal([]byte(output), &logEntry); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	fieldsObj, exists := logEntry["fields"]
	if !exists {
		t.Fatal("Expected 'fields' key with FlattenFields=false")
	}

	fields, ok := fieldsObj.(map[string]interface{})
	if !ok {
		t.Fatal("Expected 'fields' to be an object")
	}

	if service, ok := fields["service"].(string); !ok || service != "test-service" {
		t.Error("Expected service in fields object")
	}
}
