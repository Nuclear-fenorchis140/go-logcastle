package tests

import (
	"bytes"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/bhaskarblur/go-logcastle"
)

func TestColorOutput(t *testing.T) {
	logcastle.Reset()
	defer logcastle.Close()

	var buf bytes.Buffer

	config := logcastle.Config{
		Output:             &buf,
		Format:             logcastle.Text,
		Level:              logcastle.LevelDebug,
		ColorOutput:        true,
		IncludeLoggerField: false,
		IncludeParseError:  false,
	}

	logcastle.Init(config)
	logcastle.WaitReady()
	time.Sleep(100 * time.Millisecond)

	log.Println("Test message with colors")

	time.Sleep(100 * time.Millisecond)
	logcastle.Close()

	output := buf.String()
	if output == "" {
		t.Fatal("No output captured")
	}

	if !strings.Contains(output, "\033[") {
		t.Error("Expected ANSI color codes in output with ColorOutput=true, but none found")
	}

	hasReset := strings.Contains(output, "\033[0m")
	if !hasReset {
		t.Error("Expected color reset code in colored output")
	}
}

func TestNoColorOutput(t *testing.T) {
	logcastle.Reset()
	defer logcastle.Close()

	var buf bytes.Buffer

	config := logcastle.Config{
		Output:             &buf,
		Format:             logcastle.Text,
		Level:              logcastle.LevelInfo,
		ColorOutput:        false,
		IncludeLoggerField: false,
		IncludeParseError:  false,
	}

	logcastle.Init(config)
	logcastle.WaitReady()
	time.Sleep(100 * time.Millisecond)

	log.Println("Test message without colors")

	time.Sleep(100 * time.Millisecond)
	logcastle.Close()

	output := buf.String()
	if output == "" {
		t.Fatal("No output captured")
	}

	if strings.Contains(output, "\033[") {
		t.Error("Expected no ANSI color codes with ColorOutput=false, but found some")
	}
}
