package parser

import (
	"bytes"
	"regexp"
	"time"

	jsoniter "github.com/json-iterator/go"
)

var jsonAPI = jsoniter.ConfigCompatibleWithStandardLibrary

// LogEntry represents a standardized log entry
type LogEntry struct {
	Timestamp  time.Time
	Level      Level
	Logger     string
	Message    string
	Caller     string
	Fields     map[string]interface{}
	TraceID    string
	SpanID     string
	Source     string
	ParseError string // Set when parsing fails
}

// Level represents log severity
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	case LevelFatal:
		return "fatal"
	default:
		return "unknown"
	}
}

// ParseLevel converts a string to a Level
func ParseLevel(s string) Level {
	switch s {
	case "debug", "DEBUG", "Debug":
		return LevelDebug
	case "info", "INFO", "Info":
		return LevelInfo
	case "warn", "WARN", "Warn", "warning", "WARNING", "Warning":
		return LevelWarn
	case "error", "ERROR", "Error":
		return LevelError
	case "fatal", "FATAL", "Fatal", "panic", "PANIC":
		return LevelFatal
	default:
		return LevelInfo
	}
}

// NewLogEntry creates a new log entry with defaults
func NewLogEntry() *LogEntry {
	return &LogEntry{
		Timestamp: time.Now(),
		Level:     LevelInfo,
		Fields:    make(map[string]interface{}),
	}
}

// Parser parses log lines from various formats
type Parser struct {
	logrusTextPattern *regexp.Regexp
	zapConsolePattern *regexp.Regexp
}

// NewParser creates a new parser
func NewParser() *Parser {
	return &Parser{
		logrusTextPattern: regexp.MustCompile(`^time="([^"]+)"\s+level=(\w+)\s+msg="([^"]+)"`),
		zapConsolePattern: regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}T[\d:.]+Z?)\s+(\w+)\s+(.+?)(?:\t(.+))?$`),
	}
}

// Parse attempts to parse a log line and extract structured data
// Always returns a LogEntry (never nil). On parse failure, sets ParseError field
func (p *Parser) Parse(line []byte) *LogEntry {
	entry := NewLogEntry()

	// Try JSON first (most common in production)
	if len(line) > 0 && line[0] == '{' {
		if p.parseJSON(line, entry) {
			return entry
		}
		entry.ParseError = "failed to parse as JSON"
	}

	// Try logrus text format
	if bytes.Contains(line, []byte("level=")) {
		if p.parseLogrusText(line, entry) {
			if entry.ParseError == "" {
				entry.ParseError = ""
			}
			return entry
		}
		if entry.ParseError == "" {
			entry.ParseError = "failed to parse as logrus format"
		}
	}

	// Try zap console format
	if p.parseZapConsole(line, entry) {
		if entry.ParseError == "" {
			entry.ParseError = ""
		}
		return entry
	}

	// Fallback: treat as generic text
	p.parseGeneric(line, entry)
	if entry.ParseError == "" {
		entry.ParseError = "parsed as unstructured text"
	}
	return entry
}

func (p *Parser) parseJSON(line []byte, entry *LogEntry) bool {
	var data map[string]interface{}
	if err := jsonAPI.Unmarshal(line, &data); err != nil {
		return false
	}

	// Extract timestamp
	if ts, ok := data["timestamp"].(string); ok {
		if t, err := time.Parse(time.RFC3339Nano, ts); err == nil {
			entry.Timestamp = t
		}
		delete(data, "timestamp")
	} else if ts, ok := data["time"].(string); ok {
		if t, err := time.Parse(time.RFC3339Nano, ts); err == nil {
			entry.Timestamp = t
		}
		delete(data, "time")
	} else if ts, ok := data["ts"].(string); ok {
		if t, err := time.Parse(time.RFC3339Nano, ts); err == nil {
			entry.Timestamp = t
		}
		delete(data, "ts")
	}

	// Extract level
	if level, ok := data["level"].(string); ok {
		entry.Level = ParseLevel(level)
		delete(data, "level")
	} else if level, ok := data["severity"].(string); ok {
		entry.Level = ParseLevel(level)
		delete(data, "severity")
	}

	// Extract message
	if msg, ok := data["message"].(string); ok {
		entry.Message = msg
		delete(data, "message")
	} else if msg, ok := data["msg"].(string); ok {
		entry.Message = msg
		delete(data, "msg")
	}

	// Extract logger
	if logger, ok := data["logger"].(string); ok {
		entry.Logger = logger
		delete(data, "logger")
	}

	// Extract caller
	if caller, ok := data["caller"].(string); ok {
		entry.Caller = caller
		delete(data, "caller")
	}

	// Extract trace info
	if traceID, ok := data["trace_id"].(string); ok {
		entry.TraceID = traceID
		delete(data, "trace_id")
	}
	if spanID, ok := data["span_id"].(string); ok {
		entry.SpanID = spanID
		delete(data, "span_id")
	}

	// Detect logger from JSON structure
	p.detectLoggerFromJSON(data, entry)

	// Remaining fields go into Fields map
	entry.Fields = data

	return true
}

func (p *Parser) parseLogrusText(line []byte, entry *LogEntry) bool {
	matches := p.logrusTextPattern.FindSubmatch(line)
	if len(matches) < 4 {
		return false
	}

	// Parse timestamp
	if t, err := time.Parse(time.RFC3339Nano, string(matches[1])); err == nil {
		entry.Timestamp = t
	}

	// Parse level
	entry.Level = ParseLevel(string(matches[2]))

	// Parse message
	entry.Message = string(matches[3])

	entry.Logger = "logrus"
	return true
}

func (p *Parser) parseZapConsole(line []byte, entry *LogEntry) bool {
	matches := p.zapConsolePattern.FindSubmatch(line)
	if len(matches) < 4 {
		return false
	}

	// Parse timestamp
	if t, err := time.Parse(time.RFC3339Nano, string(matches[1])); err == nil {
		entry.Timestamp = t
	} else if t, err := time.Parse("2006-01-02T15:04:05.999Z", string(matches[1])); err == nil {
		entry.Timestamp = t
	}

	// Parse level
	entry.Level = ParseLevel(string(matches[2]))

	// Parse message
	entry.Message = string(matches[3])

	entry.Logger = "zap"
	return true
}

func (p *Parser) parseGeneric(line []byte, entry *LogEntry) {
	entry.Message = string(bytes.TrimSpace(line))
	entry.Logger = "unknown"
}

func (p *Parser) detectLoggerFromJSON(data map[string]interface{}, entry *LogEntry) {
	// Logrus often has these fields
	if _, hasTime := data["time"]; hasTime {
		if _, hasMsg := data["msg"]; hasMsg {
			entry.Logger = "logrus"
			return
		}
	}

	// Zap often has these fields
	if _, hasTs := data["ts"]; hasTs {
		entry.Logger = "zap"
		return
	}

	// Zerolog patterns
	if _, hasLevel := data["level"]; hasLevel {
		if _, hasMessage := data["message"]; hasMessage {
			entry.Logger = "zerolog"
			return
		}
	}
}
