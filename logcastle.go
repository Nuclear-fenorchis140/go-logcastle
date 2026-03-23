// Package logcastle provides high-performance centralized log orchestration.
// It intercepts logs from any library, standardizes the format, and writes to configured outputs.
//
// Basic Usage:
//
//	import "github.com/bhaskarblur/go-logcastle"
//
//	func main() {
//	    logcastle.Init(logcastle.Config{
//	        Format: logcastle.JSON, // or Text, LogFmt
//	        Level:  logcastle.LevelInfo,
//	    })
//	    defer logcastle.Close()
//
//	    // All logs now intercepted and standardized
//	    log.Println("Server started")
//	    fmt.Println("Processing request")
//	}
//
// Formatting Examples:
//
// 1. Development Mode (Pretty + Readable):
//
//	logcastle.Config{
//	    Format:        logcastle.JSON,
//	    PrettyPrint:   true,  // Multi-line JSON
//	    FlattenFields: true,  // Clean structure
//	}
//
// 2. Production Grafana/Loki (Optimized for Label Extraction):
//
//	logcastle.Config{
//	    Format:        logcastle.JSON,
//	    FlattenFields: true,  // CRITICAL: Fields at root for Loki labels
//	    PrettyPrint:   false, // Single-line for log aggregation
//	    FieldOrder:    []string{"timestamp", "level", "service", "env"},
//	}
//
// 3. Terminal Output with Colors:
//
//	logcastle.Config{
//	    Format:      logcastle.Text, // Human-readable
//	    ColorOutput: true,           // ERROR=red, WARN=yellow, INFO=green
//	}
//
// 4. ELK/Logstash (Custom Field Ordering):
//
//	logcastle.Config{
//	    Format:        logcastle.JSON,
//	    FlattenFields: true,
//	    FieldOrder:    []string{"timestamp", "level", "service", "message"},
//	}
//
// Format Types:
//   - JSON (default): {"timestamp":"2026-03-23T10:00:00Z","level":"info","message":"server started"}
//   - Text: 2026-03-23T10:00:00Z INFO server started
//   - LogFmt: time=2026-03-23T10:00:00Z level=info msg="server started"
//
// Formatting Options:
//   - FlattenFields (default: true): Merge enrichment fields to root level
//   - PrettyPrint (default: false): Multi-line JSON with indentation
//   - ColorOutput (default: false): ANSI colors for Text format
//   - FieldOrder (default: nil): Custom field ordering in JSON
package logcastle

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

var (
	defaultCastle *Castle
	once          sync.Once
)

// Format represents the output format for logs.
// Three formats are available:
//
// JSON (default): Structured JSON format for production and log aggregators.
//
//	Example: {"timestamp":"2026-03-23T10:00:00Z","level":"info","message":"server started"}
//
// Text: Human-readable format for terminal viewing and development.
//
//	Example: 2026-03-23T10:00:00Z INFO server started
//
// LogFmt: Key=value pairs format compatible with Logstash and other parsers.
//
//	Example: time=2026-03-23T10:00:00Z level=info msg="server started"
type Format string

const (
	// JSON is the default format - structured, machine-readable, ideal for Grafana/Loki/ELK
	JSON Format = "json"
	// Text is human-readable format with optional ANSI colors for terminal output
	Text Format = "text"
	// LogFmt is key=value format compatible with Heroku/Splunk/Logstash
	LogFmt Format = "logfmt"
)

// TimestampFormat represents the format for timestamps in log entries.
// Choose from RFC3339, Unix epoch, or custom formats for flexibility.
type TimestampFormat string

const (
	// TimestampFormatRFC3339Nano is the default format with nanosecond precision
	TimestampFormatRFC3339Nano TimestampFormat = "rfc3339nano"
	// TimestampFormatRFC3339 is standard RFC3339
	TimestampFormatRFC3339 TimestampFormat = "rfc3339"
	// TimestampFormatRFC3339Millis includes millisecond precision
	TimestampFormatRFC3339Millis TimestampFormat = "rfc3339milli"
	// TimestampFormatUnix is Unix timestamp (seconds)
	TimestampFormatUnix TimestampFormat = "unix"
	// TimestampFormatUnixMilli is Unix timestamp in milliseconds
	TimestampFormatUnixMilli TimestampFormat = "unixmilli"
	// TimestampFormatUnixNano is Unix timestamp in nanoseconds
	TimestampFormatUnixNano TimestampFormat = "unixnano"
	// TimestampFormatDateTime is human-readable format
	TimestampFormatDateTime TimestampFormat = "datetime"
	// TimestampFormatCustom allows custom format string
	TimestampFormatCustom TimestampFormat = "custom"
)

// Config configures the log castle behavior and output.
// Set Format, Level, Output destination, and other options to customize logging.
//
// Quick Start - Use defaults:
//
//	logcastle.Init(logcastle.DefaultConfig())
//
// Custom Configuration:
//
//	logcastle.Init(logcastle.Config{
//	    Format: logcastle.JSON,  // or Text, LogFmt
//	    Level:  logcastle.LevelInfo,
//	    EnrichFields: map[string]interface{}{
//	        "service": "api-server",
//	        "env":     "production",
//	    },
//	})
type Config struct {
	// Format specifies the output format: JSON (default, structured), Text (readable), or LogFmt (key=value).
	//
	// JSON example: {\"timestamp\":\"2026-03-23T10:00:00Z\",\"level\":\"info\",\"message\":\"server started\"}
	// Text example: 2026-03-23T10:00:00Z INFO server started
	// LogFmt example: time=2026-03-23T10:00:00Z level=info msg=\"server started\"
	//
	// Default: JSON
	Format Format

	// Level sets the minimum log level to capture (Debug, Info, Warn, Error, Fatal).
	// Logs below this level are filtered out.
	Level Level

	// Output is where to write standardized logs (default: os.Stdout)
	Output io.Writer

	// BufferSize controls internal buffering (default: 10000).
	// Increase for high throughput, decrease for low latency.
	BufferSize int

	// FlushInterval is how often to flush buffered logs
	FlushInterval time.Duration

	// EnrichFields adds additional fields to all logs
	EnrichFields map[string]interface{}

	// TimestampFormat specifies how timestamps are formatted in output.
	// Default: RFC3339Nano. See TimestampFormat constants for options.
	TimestampFormat TimestampFormat

	// CustomTimestampFormat is used when TimestampFormat is TimestampFormatCustom
	// Uses Go time format layout (e.g., "2006-01-02 15:04:05")
	CustomTimestampFormat string

	// IncludeLoggerField controls whether to include the 'logger' field in output.
	// When false, the logger field is omitted from formatted logs (default: false)
	IncludeLoggerField bool

	// IncludeParseError controls whether to include 'log_parse_error' field in output.
	// When false, parse error messages are omitted from formatted logs (default: false)
	IncludeParseError bool

	// FlattenFields merges EnrichFields to root level instead of nested "fields" object.
	// This is CRITICAL for Grafana/Loki label extraction and readability.
	//
	// When true (default):
	//   {"timestamp":"...","level":"info","env":"prod","service":"api"}
	//
	// When false:
	//   {"timestamp":"...","level":"info","fields":{"env":"prod","service":"api"}}
	//
	// Default: true (recommended for production observability)
	FlattenFields bool

	// PrettyPrint formats JSON with indentation for terminal viewing and debugging.
	// Use true for development/debugging, false for production log aggregation.
	//
	// When true:
	//   {
	//     "timestamp": "2026-03-23T10:00:00Z",
	//     "level": "info",
	//     "message": "server started"
	//   }
	//
	// When false (default):
	//   {"timestamp":"2026-03-23T10:00:00Z","level":"info","message":"server started"}
	//
	// Default: false (single-line for production)
	PrettyPrint bool

	// ColorOutput adds ANSI color codes to terminal output (Text format only, ignored in JSON).
	// Improves visual distinction between log levels in terminal/console.
	//
	// When true (Text format):
	//   ERROR appears in bold red
	//   WARN appears in yellow
	//   INFO appears in green
	//   DEBUG appears in gray
	//
	// Example output:
	//   2026-03-23T10:00:00Z \033[32mINFO\033[0m server started
	//
	// Default: false (no colors)
	ColorOutput bool

	// FieldOrder specifies which fields appear first in JSON output.
	// Useful for optimizing log readability in ELK/Logstash/Splunk where field order matters.
	//
	// Example for ELK/Logstash:
	//   FieldOrder: []string{"timestamp", "level", "service", "env", "message"}
	//
	// Output:
	//   {"timestamp":"...","level":"info","service":"api","env":"prod","message":"...","caller":"..."}
	//
	// Fields not in FieldOrder appear after, in alphabetical order.
	// Default: nil (standard order: timestamp, level, message, then others)
	FieldOrder []string
}

// DefaultConfig returns a Config with sensible defaults for most applications.
//
// Defaults:
//   - Format: JSON (structured logs)
//   - Level: LevelInfo (filters out debug logs)
//   - Output: os.Stdout (standard output)
//   - BufferSize: 10000 (high throughput)
//   - FlushInterval: 100ms (balance between latency and performance)
//   - FlattenFields: true (Grafana/Loki optimized)
//   - PrettyPrint: false (single-line for production)
//   - ColorOutput: false (no ANSI codes)
//   - IncludeLoggerField: false (clean output)
//   - IncludeParseError: false (clean output)
//
// Example usage:
//
//	config := logcastle.DefaultConfig()
//	config.Level = logcastle.LevelDebug
//	config.EnrichFields = map[string]interface{}{
//	    "service": "api-server",
//	    "env":     "production",
//	}
//	logcastle.Init(config)
func DefaultConfig() Config {
	return Config{
		Format:             JSON,
		Level:              LevelInfo,
		Output:             os.Stdout,
		BufferSize:         10000,
		FlushInterval:      100 * time.Millisecond,
		EnrichFields:       make(map[string]interface{}),
		TimestampFormat:    TimestampFormatRFC3339Nano,
		FlattenFields:      true,  // Flatten for better readability
		PrettyPrint:        false, // Single-line for production
		ColorOutput:        false, // No colors by default
		IncludeLoggerField: false, // Clean output
		IncludeParseError:  false, // Clean output
	}
}

// Init initializes the log castle with the given config.
// Call this once at startup. All fmt.Print*, log.Print*, etc. will be intercepted.
// Returns error if initialization fails (e.g., invalid config).
func Init(config Config) error {
	var err error
	once.Do(func() {
		// Apply defaults
		if config.Output == nil {
			config.Output = os.Stdout
		}
		if config.BufferSize == 0 {
			config.BufferSize = 10000
		}
		if config.FlushInterval == 0 {
			config.FlushInterval = 100 * time.Millisecond
		}
		if config.EnrichFields == nil {
			config.EnrichFields = make(map[string]interface{})
		}

		defaultCastle, err = newCastle(config)
		if err != nil {
			return
		}

		err = defaultCastle.start()
	})
	return err
}

// WaitReady blocks until log interception is fully active.
// Use in tests or when you need immediate log capture guarantee.
func WaitReady() {
	if defaultCastle != nil {
		<-defaultCastle.ready
	}
}

// Close gracefully shuts down the log castle and flushes buffered logs.
// Always call this before application exit  (use defer).
func Close() error {
	if defaultCastle != nil {
		return defaultCastle.stop()
	}
	return nil
}

// Reset forcefully resets the logcastle for testing purposes.
// WARNING: Only use this in tests! Not safe for production use.
func Reset() {
	if defaultCastle != nil {
		defaultCastle.stop()
	}

	// Reset the sync.Once to allow re-initialization
	once = sync.Once{}
	defaultCastle = nil
}

// Castle is the main log orchestrator
type Castle struct {
	config Config

	originalStdout *os.File
	originalStderr *os.File

	stdoutReader *os.File
	stdoutWriter *os.File
	stderrReader *os.File
	stderrWriter *os.File

	parser    *Parser
	formatter *Formatter
	writer    *BufferedWriter

	done  chan struct{}
	ready chan struct{} // Signals when interception is ready
	wg    sync.WaitGroup
	mu    sync.Mutex // Protects initialization state
}

func newCastle(config Config) (*Castle, error) {
	parser := NewParser()
	formatter := NewFormatter(config.Format, config.TimestampFormat, config.CustomTimestampFormat,
		config.IncludeLoggerField, config.IncludeParseError, config.FlattenFields,
		config.PrettyPrint, config.ColorOutput, config.FieldOrder)
	writer := NewBufferedWriter(config.Output, config.BufferSize, config.FlushInterval)

	return &Castle{
		config:         config,
		originalStdout: os.Stdout,
		originalStderr: os.Stderr,
		parser:         parser,
		formatter:      formatter,
		writer:         writer,
		done:           make(chan struct{}),
		ready:          make(chan struct{}),
	}, nil
}

func (c *Castle) start() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Create pipes for stdout and stderr using os.Pipe
	var err error
	c.stdoutReader, c.stdoutWriter, err = os.Pipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	c.stderrReader, c.stderrWriter, err = os.Pipe()
	if err != nil {
		c.stdoutReader.Close()
		c.stdoutWriter.Close()
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the buffered writer first
	c.writer.Start()

	// Start processing goroutines with ready synchronization
	var startWg sync.WaitGroup
	startWg.Add(2)

	c.wg.Add(2)
	go func() {
		startWg.Done()
		c.processStream(c.stdoutReader, "stdout")
	}()
	go func() {
		startWg.Done()
		c.processStream(c.stderrReader, "stderr")
	}()

	// Wait for goroutines to start
	startWg.Wait()

	// Replace os.Stdout and os.Stderr BEFORE signaling ready
	os.Stdout = c.stdoutWriter
	os.Stderr = c.stderrWriter

	// CRITICAL: Reconfigure stdlib log package to use new stderr
	// The log package stores os.Stderr at import time, so we must update it
	log.SetOutput(c.stderrWriter)
	log.SetFlags(log.LstdFlags) // Keep default flags

	// Small delay to ensure pipes are fully established
	time.Sleep(10 * time.Millisecond)

	// Signal that interception is ready
	close(c.ready)

	return nil
}

func (c *Castle) stop() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Restore original stdout/stderr first to prevent deadlock
	if c.originalStdout != nil {
		os.Stdout = c.originalStdout
	}
	if c.originalStderr != nil {
		os.Stderr = c.originalStderr
		// Restore stdlib log package to use original stderr
		log.SetOutput(c.originalStderr)
	}

	// Close the done channel
	select {
	case <-c.done:
		// Already closed
	default:
		close(c.done)
	}

	// Close pipe writers to unblock readers
	if c.stdoutWriter != nil {
		c.stdoutWriter.Close()
	}
	if c.stderrWriter != nil {
		c.stderrWriter.Close()
	}

	// Wait for processing goroutines to finish
	c.wg.Wait()

	// Close pipe readers
	if c.stdoutReader != nil {
		c.stdoutReader.Close()
	}
	if c.stderrReader != nil {
		c.stderrReader.Close()
	}

	// Stop the writer and flush remaining data
	if err := c.writer.Stop(); err != nil {
		return fmt.Errorf("failed to flush logs: %w", err)
	}

	return nil
}

func (c *Castle) processStream(r io.Reader, source string) {
	defer c.wg.Done()

	scanner := NewScanner(r)
	for {
		select {
		case <-c.done:
			return
		default:
			if !scanner.Scan() {
				// Check for errors
				if err := scanner.Err(); err != nil && err != io.EOF {
					// Log error to original stderr (best effort)
					fmt.Fprintf(c.originalStderr, "logcastle: scanner error: %v\n", err)
				}
				return
			}

			line := scanner.Bytes()
			if len(line) == 0 {
				continue
			}

			// Parse the log line
			entry := c.parser.Parse(line)

			// Add source if not already set
			if entry.Source == "" {
				entry.Source = source
			}

			// Filter by level
			if entry.Level < c.config.Level {
				continue
			}

			// Enrich with configured fields
			for k, v := range c.config.EnrichFields {
				if entry.Fields == nil {
					entry.Fields = make(map[string]interface{})
				}
				if _, exists := entry.Fields[k]; !exists {
					entry.Fields[k] = v
				}
			}

			// Format to standard output
			formatted := c.formatter.Format(entry)
			if formatted != nil {
				// Write to buffered writer
				if err := c.writer.Write(formatted); err != nil {
					// Best effort logging to original stderr
					fmt.Fprintf(c.originalStderr, "logcastle: write error: %v\n", err)
				}
			}
		}
	}
}
