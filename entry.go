package logcastle

import "time"

// LogEntry represents a parsed and standardized log entry from any logging library.
// All logs are normalized into this structure for consistent processing.
type LogEntry struct {
	// Timestamp is when the log was created.
	Timestamp  time.Time              `json:"timestamp"`
	Level      Level                  `json:"level"`
	Logger     string                 `json:"logger"`
	Message    string                 `json:"message"`
	Caller     string                 `json:"caller,omitempty"`
	Fields     map[string]interface{} `json:"fields,omitempty"`
	TraceID    string                 `json:"trace_id,omitempty"`
	SpanID     string                 `json:"span_id,omitempty"`
	Source     string                 `json:"-"`                         // Internal use only
	// ParseError indicates if log parsing failed (e.g., malformed JSON, unstructured text).
	ParseError string                 `json:"log_parse_error,omitempty"`
}

// NewLogEntry creates a new log entry with default values.
// Used internally by parsers to initialize log entries.
func NewLogEntry() *LogEntry {
	return &LogEntry{
		Timestamp: time.Now().UTC(),
		Level:     LevelInfo,
		Fields:    make(map[string]interface{}),
	}
}
