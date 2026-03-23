package logcastle

import (
	"bytes"
	"regexp"
	"time"

	jsoniter "github.com/json-iterator/go"
)

var jsonParser = jsoniter.ConfigCompatibleWithStandardLibrary

// Parser detects and parses logs from multiple formats (JSON, Logrus, Zap, etc.).
// Automatically identifies format and extracts structured data.
type Parser struct {
	// Patterns for detecting different log formats
	logrusTextPattern *regexp.Regexp
	zapConsolePattern *regexp.Regexp
}

// NewParser creates a parser with regex patterns for common log formats.
// Supports JSON, Logrus text, Zap console, and unstructured text.
func NewParser() *Parser {
	return &Parser{
		logrusTextPattern: regexp.MustCompile(`^time="([^"]+)"\s+level=(\w+)\s+msg="([^"]+)"`),
		zapConsolePattern: regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}T[\d:.]+Z?)\s+(\w+)\s+(.+?)(?:\t(.+))?$`),
	}
}

// Parse attempts to parse a log line and extract structured data.
// Always returns a LogEntry (never nil). Sets ParseError field when parsing fails.
// Tries JSON, Logrus, Zap, then falls back to plain text.
func (p *Parser) Parse(line []byte) *LogEntry {
	entry := NewLogEntry()

	// Try JSON first (most common in production)
	if len(line) > 0 && line[0] == '{' {
		if p.parseJSON(line, entry) {
			return entry
		}
		// JSON parse failed
		entry.ParseError = "failed to parse as JSON"
	}

	// Try logrus text format
	if bytes.Contains(line, []byte("level=")) {
		if p.parseLogrusText(line, entry) {
			if entry.ParseError == "" {
				return entry
			}
		}
		// Logrus parse failed
		if entry.ParseError == "" {
			entry.ParseError = "failed to parse as logrus format"
		}
	}

	// Try zap console format
	if p.parseZapConsole(line, entry) {
		if entry.ParseError == "" {
			return entry
		}
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
	if err := jsonParser.Unmarshal(line, &data); err != nil {
		return false
	}

	// Extract standard fields
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

	if level, ok := data["level"].(string); ok {
		entry.Level = ParseLevel(level)
		delete(data, "level")
	} else if level, ok := data["severity"].(string); ok {
		entry.Level = ParseLevel(level)
		delete(data, "severity")
	}

	if msg, ok := data["message"].(string); ok {
		entry.Message = msg
		delete(data, "message")
	} else if msg, ok := data["msg"].(string); ok {
		entry.Message = msg
		delete(data, "msg")
	}

	if logger, ok := data["logger"].(string); ok {
		entry.Logger = logger
		delete(data, "logger")
	}

	if caller, ok := data["caller"].(string); ok {
		entry.Caller = caller
		delete(data, "caller")
	}

	if traceID, ok := data["trace_id"].(string); ok {
		entry.TraceID = traceID
		delete(data, "trace_id")
	}

	if spanID, ok := data["span_id"].(string); ok {
		entry.SpanID = spanID
		delete(data, "span_id")
	}

	// Detect logger type from JSON structure
	if entry.Logger == "" {
		entry.Logger = p.detectLoggerFromJSON(data)
	}

	// Remaining fields go into Fields
	if len(data) > 0 {
		entry.Fields = data
	}

	return true
}

func (p *Parser) detectLoggerFromJSON(data map[string]interface{}) string {
	// Zap specific fields
	if _, ok := data["zapLogger"]; ok {
		return "zap"
	}
	// Zerolog specific
	if _, ok := data["zerolog"]; ok {
		return "zerolog"
	}
	// Logrus specific
	if _, ok := data["logrus"]; ok {
		return "logrus"
	}
	return "json"
}

func (p *Parser) parseLogrusText(line []byte, entry *LogEntry) bool {
	matches := p.logrusTextPattern.FindSubmatch(line)
	if len(matches) < 4 {
		return false
	}

	// Parse timestamp
	if t, err := time.Parse(time.RFC3339, string(matches[1])); err == nil {
		entry.Timestamp = t
	}

	// Parse level
	entry.Level = ParseLevel(string(matches[2]))

	// Parse message
	entry.Message = string(matches[3])

	entry.Logger = "logrus"

	// Try to extract additional fields from the rest of the line
	remaining := line[len(matches[0]):]
	if len(remaining) > 0 {
		p.parseKeyValues(remaining, entry)
	}

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
	}

	// Parse level
	entry.Level = ParseLevel(string(matches[2]))

	// Parse message
	entry.Message = string(matches[3])

	entry.Logger = "zap"

	// Parse fields if present
	if len(matches) > 4 && len(matches[4]) > 0 {
		p.parseKeyValues(matches[4], entry)
	}

	return true
}

func (p *Parser) parseGeneric(line []byte, entry *LogEntry) {
	message := string(bytes.TrimSpace(line))
	entry.Message = message

	// Smart detection: stdlib log vs fmt vs generic stdout
	// Stdlib log format: "2026/03/23 19:37:55 message"
	if len(message) > 19 && message[4] == '/' && message[7] == '/' && message[10] == ' ' {
		// Has stdlib log timestamp prefix
		entry.Logger = "log"
	} else if len(message) > 0 {
		// No timestamp prefix - likely fmt.Println or direct stdout write
		entry.Logger = "stdout"
	} else {
		entry.Logger = "unknown"
	}
}

func (p *Parser) parseKeyValues(data []byte, entry *LogEntry) {
	// Simple key=value parser
	parts := bytes.Split(data, []byte(" "))
	for _, part := range parts {
		if kv := bytes.SplitN(part, []byte("="), 2); len(kv) == 2 {
			key := string(bytes.TrimSpace(kv[0]))
			value := string(bytes.Trim(bytes.TrimSpace(kv[1]), `"`))
			if entry.Fields == nil {
				entry.Fields = make(map[string]interface{})
			}
			entry.Fields[key] = value
		}
	}
}
