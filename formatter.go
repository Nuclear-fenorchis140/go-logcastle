package logcastle

import (
	"bytes"
	"fmt"
	"sort"
	"time"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// ANSI color codes for terminal output
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorGreen  = "\033[32m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
	colorBold   = "\033[1m"
)

// Formatter formats log entries to the desired output format (JSON, Text, LogFmt).
// Supports custom timestamp formats, field flattening, pretty printing, and color output.
//
// Format Options:
//   - JSON: {"timestamp":"...","level":"info","message":"..."}
//   - Text: 2026-03-23T10:00:00Z INFO server started
//   - LogFmt: time=2026-03-23T10:00:00Z level=info msg="server started"
//
// Advanced Features:
//   - FlattenFields: Merge enrichment fields to root (Grafana/Loki optimization)
//   - PrettyPrint: Multi-line JSON with indentation (development/debugging)
//   - ColorOutput: ANSI colors for Text format (terminal readability)
//   - FieldOrder: Custom field ordering in JSON (ELK/Logstash optimization)
type Formatter struct {
	format                Format
	timestampFormat       TimestampFormat
	customTimestampFormat string
	includeLoggerField    bool     // Include "logger" field showing source (stdout, log, gin, mongo, etc.)
	includeParseError     bool     // Include "log_parse_error" field when log parsing fails
	flattenFields         bool     // Flatten EnrichFields to root level (true) vs nested "fields" object (false)
	prettyPrint           bool     // Multi-line JSON with indentation (true) vs single-line (false)
	colorOutput           bool     // ANSI color codes for Text format levels (ERROR=red, WARN=yellow, INFO=green, DEBUG=gray)
	fieldOrder            []string // Preferred field order in JSON output (e.g., ["timestamp", "level", "service"])
}

// NewFormatter creates a new formatter
func NewFormatter(format Format, tsFormat TimestampFormat, customTsFormat string,
	includeLogger bool, includeParseError bool, flattenFields bool,
	prettyPrint bool, colorOutput bool, fieldOrder []string) *Formatter {
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
		flattenFields:         flattenFields,
		prettyPrint:           prettyPrint,
		colorOutput:           colorOutput,
		fieldOrder:            fieldOrder,
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

	// Build output map with optional flattening
	if f.flattenFields && len(entry.Fields) > 0 {
		// Flatten: merge enrichment fields to root level (Logstash/ELK format)
		for key, value := range entry.Fields {
			output[key] = value
		}
	}

	// Add standard fields
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

	// If not flattening, add fields as nested object
	if !f.flattenFields && len(entry.Fields) > 0 {
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

	// Handle field ordering and pretty printing
	var data []byte
	var err error

	if len(f.fieldOrder) > 0 {
		// Custom field ordering for better readability
		data, err = f.marshalWithOrder(output)
	} else if f.prettyPrint {
		// Pretty print with indentation
		data, err = json.MarshalIndent(output, "", "  ")
	} else {
		// Standard single-line JSON
		data, err = json.Marshal(output)
	}

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

	// Colorize level if enabled
	if f.colorOutput {
		buf.WriteString(f.colorizeLevel(entry.Level))
	} else {
		buf.WriteString(entry.Level.String())
	}
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

// marshalWithOrder marshals a map with custom field ordering
// Fields in fieldOrder appear first, followed by remaining fields alphabetically
func (f *Formatter) marshalWithOrder(data map[string]interface{}) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')

	// Track which fields we've written
	written := make(map[string]bool)
	first := true

	// Write ordered fields first
	for _, field := range f.fieldOrder {
		if value, exists := data[field]; exists {
			if !first {
				buf.WriteByte(',')
			}
			first = false

			// Write key
			keyBytes, err := json.Marshal(field)
			if err != nil {
				return nil, err
			}
			buf.Write(keyBytes)
			buf.WriteByte(':')

			// Write value
			valueBytes, err := json.Marshal(value)
			if err != nil {
				return nil, err
			}
			buf.Write(valueBytes)

			written[field] = true
		}
	}

	// Write remaining fields alphabetically
	remainingKeys := make([]string, 0)
	for key := range data {
		if !written[key] {
			remainingKeys = append(remainingKeys, key)
		}
	}
	sort.Strings(remainingKeys)

	for _, key := range remainingKeys {
		if !first {
			buf.WriteByte(',')
		}
		first = false

		// Write key
		keyBytes, err := json.Marshal(key)
		if err != nil {
			return nil, err
		}
		buf.Write(keyBytes)
		buf.WriteByte(':')

		// Write value
		valueBytes, err := json.Marshal(data[key])
		if err != nil {
			return nil, err
		}
		buf.Write(valueBytes)
	}

	buf.WriteByte('}')
	return buf.Bytes(), nil
}

// colorizeLevel adds ANSI color codes to log levels for terminal output
func (f *Formatter) colorizeLevel(level Level) string {
	switch level {
	case LevelError:
		return colorBold + colorRed + "ERROR" + colorReset
	case LevelWarn:
		return colorYellow + "WARN" + colorReset
	case LevelInfo:
		return colorGreen + "INFO" + colorReset
	case LevelDebug:
		return colorGray + "DEBUG" + colorReset
	default:
		return level.String()
	}
}
