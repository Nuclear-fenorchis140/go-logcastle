package formatter

import (
	"bytes"
	"fmt"
	"sort"
	"time"

	jsoniter "github.com/json-iterator/go"
	logcastle "github.com/bhaskarblur/go-logcastle"
)

var jsonAPI = jsoniter.ConfigCompatibleWithStandardLibrary

// Format represents the output format for logs
type Format string

const (
	JSON   Format = "json"
	Text   Format = "text"
	LogFmt Format = "logfmt"
	Custom Format = "custom"
)

// Type aliases for convenience
type LogEntry = logcastle.LogEntry
type Level = logcastle.Level

const (
	LevelDebug = logcastle.LevelDebug
	LevelInfo  = logcastle.LevelInfo
	LevelWarn  = logcastle.LevelWarn
	LevelError = logcastle.LevelError
	LevelFatal = logcastle.LevelFatal
)

// TimestampFormat represents timestamp formatting options
type TimestampFormat string

const (
	TimestampFormatRFC3339Nano   TimestampFormat = "rfc3339nano"
	TimestampFormatRFC3339       TimestampFormat = "rfc3339"
	TimestampFormatRFC3339Millis TimestampFormat = "rfc3339milli"
	TimestampFormatUnix          TimestampFormat = "unix"
	TimestampFormatUnixMilli     TimestampFormat = "unixmilli"
	TimestampFormatUnixNano      TimestampFormat = "unixnano"
	TimestampFormatDateTime      TimestampFormat = "datetime"
	TimestampFormatCustom        TimestampFormat = "custom"
)

// Formatter formats log entries to the desired output format
type Formatter struct {
	format                Format
	customTemplate        string
	customCompiled        *CustomTemplate
	timestampFormat       TimestampFormat
	customTimestampFormat string
}

// NewFormatter creates a new formatter with the specified format
func NewFormatter(format Format) *Formatter {
	if format == "" {
		format = JSON
	}
	return &Formatter{
		format:          format,
		timestampFormat: TimestampFormatRFC3339Nano, // Default
	}
}

// NewFormatterWithTimestamp creates a formatter with custom timestamp format
func NewFormatterWithTimestamp(format Format, tsFormat TimestampFormat, customTsFormat string) *Formatter {
	if format == "" {
		format = JSON
	}
	if tsFormat == "" {
		tsFormat = TimestampFormatRFC3339Nano
	}
	return &Formatter{
		format:                format,
		timestampFormat:       tsFormat,
		customTimestampFormat: customTsFormat,
	}
}

// NewCustomFormatter creates a formatter with a custom template
// The template is validated at creation time
func NewCustomFormatter(template string) (*Formatter, error) {
	compiled, err := CompileTemplate(template)
	if err != nil {
		return nil, fmt.Errorf("invalid custom template: %w", err)
	}

	return &Formatter{
		format:          Custom,
		customTemplate:  template,
		customCompiled:  compiled,
		timestampFormat: TimestampFormatRFC3339Nano,
	}, nil
}

// Format formats a log entry according to the configured format
func (f *Formatter) Format(entry *LogEntry) []byte {
	switch f.format {
	case JSON:
		return f.formatJSON(entry)
	case Text:
		return f.formatText(entry)
	case LogFmt:
		return f.formatLogFmt(entry)
	case Custom:
		return f.formatCustom(entry)
	default:
		return f.formatJSON(entry)
	}
}

// formatJSON formats as JSON (default format)
func (f *Formatter) formatJSON(entry *LogEntry) []byte {
	output := make(map[string]interface{})

	output["timestamp"] = f.formatTimestamp(entry.Timestamp)
	output["level"] = entry.Level.String()
	output["message"] = entry.Message

	if entry.Logger != "" {
		output["logger"] = entry.Logger
	}

	if entry.Caller != "" {
		output["caller"] = entry.Caller
	}

	if entry.Source != "" {
		output["source"] = entry.Source
	}

	if entry.TraceID != "" {
		output["trace_id"] = entry.TraceID
	}

	if entry.SpanID != "" {
		output["span_id"] = entry.SpanID
	}

	// Add parse error if present
	if entry.ParseError != "" {
		output["log_parse_error"] = entry.ParseError
	}

	// Add fields
	if len(entry.Fields) > 0 {
		for k, v := range entry.Fields {
			output[k] = v
		}
	}

	data, err := jsonAPI.Marshal(output)
	if err != nil {
		return []byte(fmt.Sprintf(`{"error":"failed to marshal log: %v"}`, err))
	}

	return append(data, '\n')
}

// formatText formats as human-readable text
func (f *Formatter) formatText(entry *LogEntry) []byte {
	var buf bytes.Buffer

	buf.WriteString(f.formatTimestamp(entry.Timestamp))
	buf.WriteString(" [")
	buf.WriteString(entry.Level.String())
	buf.WriteString("] ")

	if entry.Logger != "" {
		buf.WriteString("[")
		buf.WriteString(entry.Logger)
		buf.WriteString("] ")
	}

	buf.WriteString(entry.Message)

	// Add parse error if present
	if entry.ParseError != "" {
		buf.WriteString(" [parse_error: ")
		buf.WriteString(entry.ParseError)
		buf.WriteString("]")
	}

	if len(entry.Fields) > 0 {
		buf.WriteString(" ")
		f.writeFields(&buf, entry.Fields)
	}

	buf.WriteString("\n")
	return buf.Bytes()
}

// formatLogFmt formats as logfmt (key=value pairs)
func (f *Formatter) formatLogFmt(entry *LogEntry) []byte {
	var buf bytes.Buffer

	buf.WriteString("timestamp=")
	buf.WriteString(f.formatTimestamp(entry.Timestamp))
	buf.WriteString(" level=")
	buf.WriteString(entry.Level.String())

	if entry.Logger != "" {
		buf.WriteString(" logger=")
		buf.WriteString(entry.Logger)
	}

	buf.WriteString(" message=\"")
	buf.WriteString(entry.Message)
	buf.WriteString("\"")

	// Add parse error if present
	if entry.ParseError != "" {
		buf.WriteString(" log_parse_error=\"")
		buf.WriteString(entry.ParseError)
		buf.WriteString("\"")
	}

	if len(entry.Fields) > 0 {
		buf.WriteString(" ")
		f.writeFields(&buf, entry.Fields)
	}

	buf.WriteString("\n")
	return buf.Bytes()
}

// formatCustom formats using the custom template
func (f *Formatter) formatCustom(entry *LogEntry) []byte {
	if f.customCompiled == nil {
		return []byte("ERROR: custom template not compiled\n")
	}

	result := f.customCompiled.Execute(entry)
	return append([]byte(result), '\n')
}

// writeFields writes fields to a buffer
func (f *Formatter) writeFields(buf *bytes.Buffer, fields map[string]interface{}) {
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i, k := range keys {
		if i > 0 {
			buf.WriteString(" ")
		}
		fmt.Fprintf(buf, "%s=%v", k, fields[k])
	}
}

// formatTimestamp formats a timestamp according to the configured format
func (f *Formatter) formatTimestamp(t time.Time) string {
	switch f.timestampFormat {
	case TimestampFormatRFC3339Nano:
		return t.Format(time.RFC3339Nano)
	case TimestampFormatRFC3339:
		return t.Format(time.RFC3339)
	case TimestampFormatRFC3339Millis:
		return t.Format("2006-01-02T15:04:05.000Z07:00")
	case TimestampFormatUnix:
		return fmt.Sprintf("%d", t.Unix())
	case TimestampFormatUnixMilli:
		return fmt.Sprintf("%d", t.UnixMilli())
	case TimestampFormatUnixNano:
		return fmt.Sprintf("%d", t.UnixNano())
	case TimestampFormatDateTime:
		return t.Format("2006-01-02 15:04:05")
	case TimestampFormatCustom:
		if f.customTimestampFormat != "" {
			return t.Format(f.customTimestampFormat)
		}
		return t.Format(time.RFC3339Nano)
	default:
		return t.Format(time.RFC3339Nano)
	}
}
