package tests

import (
	"bytes"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/bhaskarblur/go-logcastle"
)

func TestFieldOrder(t *testing.T) {
	logcastle.Reset()
	defer logcastle.Close()

	var buf bytes.Buffer

	config := logcastle.Config{
		Output:        &buf,
		Format:        logcastle.JSON,
		Level:         logcastle.LevelInfo,
		FlattenFields: true,
		FieldOrder: []string{
			"timestamp",
			"level",
			"message",
			"service",
		},
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

	fmt.Println("Field order test")

	time.Sleep(100 * time.Millisecond)
	logcastle.Close()

	output := buf.String()
	if output == "" {
		t.Fatal("No output captured")
	}

	timestampPos := findFieldPosition(output, "timestamp")
	levelPos := findFieldPosition(output, "level")
	messagePos := findFieldPosition(output, "message")
	servicePos := findFieldPosition(output, "service")

	if timestampPos == -1 || levelPos == -1 || messagePos == -1 || servicePos == -1 {
		t.Fatal("One or more required fields not found")
	}

	if !(timestampPos < levelPos && levelPos < messagePos && messagePos < servicePos) {
		t.Error("Fields are not in the specified order")
	}
}

func findFieldPosition(jsonStr string, fieldName string) int {
	pattern := fmt.Sprintf(`"%s"`, fieldName)
	re := regexp.MustCompile(pattern)
	loc := re.FindStringIndex(jsonStr)
	if loc == nil {
		return -1
	}
	return loc[0]
}
