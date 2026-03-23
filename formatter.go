package logcastle

import (
	"bytes"
	"fmt"
	"sort"
	"time"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Formatter formats log entries to the desired output format (JSON, Text, LogFmt).
// Supports custom timestamp formats and preserves parse error information.
type Formatter struct {
	format                Format
	timestampFormat       TimestampFormat
	customTimestampFormat string
	includeLoggerField    bool // Whether to include logger field in output
	includeParseError     bool // Whether to include parse error field in output
}

// NewFormatter creates a new formatter
func NewFormatter(format Format, tsFormat TimestampFormat, customTsFormat string, includeLogger bool, includeParseError bool) *Formatter {
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
		includeLoggerField:    includeLogger,
		includeParseError:     includeParseError,
	}
}

// Format converts a LogEntry into formatted bytes ready for output.
// Returns newline-terminated bytes in the configured format (JSON/Text/LogFmt).
func (f *Formatter) Format(entry *LogEntry) []byte {
	switch f.format {
	case JSON:
		return f.formatJSON(entry)
	case Text:
		return f.formatText(entry)
	case LogFmt:
		return f.formatLogFmt(entry)
	default:
		return f.formatJSON(entry)
	}
}

func (f *Formatter) formatJSON(entry *LogEntry) []byte {
	output := make(map[string]interface{})

	output["timestamp"] = f.formatTimestamp(entry.Timestamp)
	output["level"] = entry.Level.String()
	output["message"] = entry.Message

	// Optionally include logger field
	if f.includeLoggerField && entry.Logger != "" {
		output["logger"] = entry.Logger
	}

	if entry.Caller != "" {
		output["caller"] = entry.Caller
	}

	if len(entry.Fields) > 0 {
		output["fields"] = entry.Fields
	}

	if entry.TraceID != "" {
		output["trace_id"] = entry.TraceID
	}

	if entry.SpanID != "" {
		output["span_id"] = entry.SpanID
	}

	// Optionally include parse error
	if f.includeParseError && entry.ParseError != "" {
		output["log_parse_error"] = entry.ParseError
	}

	data, err := json.Marshal(output)
	if err != nil {
		return nil
	}

	return append(data, '\n')
}

func (f *Formatter) formatText(entry *LogEntry) []byte {
	var buf bytes.Buffer

	// timestamp level [logger] message
	buf.WriteString(f.formatTimestamp(entry.Timestamp))
	buf.WriteByte(' ')
	buf.WriteString(entry.Level.String())
	buf.WriteByte(' ')

	// Optionally include logger field
	if f.includeLoggerField && entry.Logger != "" {
		buf.WriteByte('[')
		buf.WriteString(entry.Logger)
		buf.WriteString("] ")
	}

	buf.WriteString(entry.Message)

	// Optionally include parse error
	if f.includeParseError && entry.ParseError != "" {
		buf.WriteString(" [parse_error: ")
		buf.WriteString(entry.ParseError)
		buf.WriteByte(']')
	}

	// Add fields
	if len(entry.Fields) > 0 {
		buf.WriteByte(' ')
		f.writeFields(&buf, entry.Fields)
	}

	buf.WriteByte('\n')
	return buf.Bytes()
}

func (f *Formatter) formatLogFmt(entry *LogEntry) []byte {
	var buf bytes.Buffer

	// time=... level=... msg=...
	fmt.Fprintf(&buf, "time=%s ", entry.Timestamp.Format(time.RFC3339Nano))
	fmt.Fprintf(&buf, "level=%s ", entry.Level.String())

	// Optionally include logger field
	if f.includeLoggerField && entry.Logger != "" {
		fmt.Fprintf(&buf, "logger=%s ", entry.Logger)
	}

	fmt.Fprintf(&buf, "msg=%q", entry.Message)

	// Optionally include parse error
	if f.includeParseError && entry.ParseError != "" {
		fmt.Fprintf(&buf, " parse_error=%q", entry.ParseError)
	}

	// Add fields
	if len(entry.Fields) > 0 {
		buf.WriteByte(' ')
		f.writeFields(&buf, entry.Fields)
	}

	buf.WriteByte('\n')
	return buf.Bytes()
}

func (f *Formatter) writeFields(buf *bytes.Buffer, fields map[string]interface{}) {
	// Sort keys for consistent output
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	first := true
	for _, k := range keys {
		if !first {
			buf.WriteByte(' ')
		}
		first = false

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
