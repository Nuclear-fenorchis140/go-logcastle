package formatter

import (
	"bytes"
	"encoding/json"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// RuntimeFields provides access to runtime and environment information for log enrichment.
// Collects hostname, PID, environment name, and custom fields once at startup.
type RuntimeFields struct {
	hostname  string
	pid       int
	env       string
	once      sync.Once
	envFields map[string]string
}

var runtimeFields = &RuntimeFields{}

// InitRuntimeFields initializes global runtime context for all logs.
// Call once at app startup. environment: "prod"/"dev"/etc, customEnvFields: datacenter, region, etc.
func InitRuntimeFields(environment string, customEnvFields map[string]string) {
	runtimeFields.once.Do(func() {
		runtimeFields.hostname, _ = os.Hostname()
		runtimeFields.pid = os.Getpid()
		runtimeFields.env = environment
		runtimeFields.envFields = customEnvFields
	})
}

// GetRuntimeFields returns current runtime information as a map.
// Includes hostname, PID, environment, goroutine count, and custom env fields.
func GetRuntimeFields() map[string]interface{} {
	fields := make(map[string]interface{})

	if runtimeFields.hostname != "" {
		fields["hostname"] = runtimeFields.hostname
	}
	if runtimeFields.pid > 0 {
		fields["pid"] = runtimeFields.pid
	}
	if runtimeFields.env != "" {
		fields["environment"] = runtimeFields.env
	}

	// Add custom environment fields
	for k, v := range runtimeFields.envFields {
		fields[k] = v
	}

	// Add goroutine count (dynamic)
	fields["goroutines"] = runtime.NumGoroutine()

	return fields
}

// JSONFormatter provides customizable JSON formatting with global fields and field ordering.
// Thread-safe. Supports runtime fields, custom field order, and pretty printing.
type JSONFormatter struct {
	// FieldOrder specifies the order of fields in JSON output for readability.
	// Empty = default alphabetical. Example: ["timestamp", "level", "message"]
	FieldOrder []string

	// GlobalFields are automatically added to every log entry.
	// Use for service name, version, region, etc. Thread-safe updates via SetGlobalField.
	GlobalFields map[string]interface{}

	// IncludeRuntimeFields adds hostname, PID, etc. to logs
	IncludeRuntimeFields bool

	// TimestampFormat for formatting timestamps
	TimestampFormat string

	// PrettyPrint enables indented JSON output
	PrettyPrint bool

	mu sync.RWMutex
}

// NewJSONFormatter creates a JSON formatter with sensible defaults.
// Default field order: timestamp, level, message. Call SetGlobalField to add metadata.
func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{
		FieldOrder:           []string{"timestamp", "level", "message"},
		GlobalFields:         make(map[string]interface{}),
		IncludeRuntimeFields: false,
		TimestampFormat:      time.RFC3339Nano,
		PrettyPrint:          false,
	}
}

// SetGlobalField adds or updates a global field that appears in all logs.
// Thread-safe. Example: f.SetGlobalField("service", "api-gateway")
func (f *JSONFormatter) SetGlobalField(key string, value interface{}) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.GlobalFields[key] = value
}

// SetGlobalFields sets multiple global fields at once
func (f *JSONFormatter) SetGlobalFields(fields map[string]interface{}) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for k, v := range fields {
		f.GlobalFields[k] = v
	}
}

// RemoveGlobalField removes a global field
func (f *JSONFormatter) RemoveGlobalField(key string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.GlobalFields, key)
}

// CustomLogEntry represents a log entry for JSON custom formatter
type CustomLogEntry struct {
	Timestamp time.Time
	Level     string
	Logger    string
	Message   string
	Caller    string
	Fields    map[string]interface{}
	TraceID   string
	SpanID    string
	Source    string
}

// Format formats a log entry as customizable JSON with global and runtime fields.
// Returns newline-terminated JSON bytes. Honors FieldOrder if set.
func (f *JSONFormatter) Format(entry *CustomLogEntry) []byte {
	f.mu.RLock()
	defer f.mu.RUnlock()

	// Build output map
	output := make(map[string]interface{})

	// Add global fields first (can be overridden)
	for k, v := range f.GlobalFields {
		output[k] = v
	}

	// Add runtime fields if enabled
	if f.IncludeRuntimeFields {
		for k, v := range GetRuntimeFields() {
			output[k] = v
		}
	}

	// Add log entry fields
	output["timestamp"] = entry.Timestamp.Format(f.TimestampFormat)
	output["level"] = entry.Level
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

	// Merge entry fields
	for k, v := range entry.Fields {
		output[k] = v
	}

	// Format based on configuration
	if f.PrettyPrint {
		return f.formatPretty(output)
	}

	if len(f.FieldOrder) > 0 {
		return f.formatOrdered(output)
	}

	// Default JSON encoding
	data, _ := json.Marshal(output)
	return append(data, '\n')
}

// formatOrdered returns JSON with fields in specified order
func (f *JSONFormatter) formatOrdered(output map[string]interface{}) []byte {
	var buf bytes.Buffer
	buf.WriteByte('{')

	first := true

	// Output ordered fields first
	for _, key := range f.FieldOrder {
		if val, exists := output[key]; exists {
			if !first {
				buf.WriteByte(',')
			}
			first = false

			// Write key
			buf.WriteByte('"')
			buf.WriteString(key)
			buf.WriteString(`":`) // Use raw string

			// Write value
			f.writeJSONValue(&buf, val)

			// Remove from map so we don't duplicate
			delete(output, key)
		}
	}

	// Output remaining fields (alphabetically)
	if len(output) > 0 {
		keys := make([]string, 0, len(output))
		for k := range output {
			keys = append(keys, k)
		}
		// Sort for consistency
		for i := 0; i < len(keys); i++ {
			for j := i + 1; j < len(keys); j++ {
				if keys[i] > keys[j] {
					keys[i], keys[j] = keys[j], keys[i]
				}
			}
		}

		for _, key := range keys {
			if !first {
				buf.WriteByte(',')
			}
			first = false

			buf.WriteByte('"')
			buf.WriteString(key)
			buf.WriteString(`":`) // Use raw string
			f.writeJSONValue(&buf, output[key])
		}
	}

	buf.WriteByte('}')
	buf.WriteByte('\n')
	return buf.Bytes()
}

// formatPretty returns indented JSON
func (f *JSONFormatter) formatPretty(output map[string]interface{}) []byte {
	data, _ := json.MarshalIndent(output, "", "  ")
	return append(data, '\n')
}

// writeJSONValue writes a value in JSON format
func (f *JSONFormatter) writeJSONValue(buf *bytes.Buffer, val interface{}) {
	switch v := val.(type) {
	case string:
		buf.WriteByte('"')
		// Escape special characters
		for _, c := range v {
			switch c {
			case '"':
				buf.WriteString(`\"`)
			case '\\':
				buf.WriteString(`\\`)
			case '\n':
				buf.WriteString(`\n`)
			case '\r':
				buf.WriteString(`\r`)
			case '\t':
				buf.WriteString(`\t`)
			default:
				buf.WriteRune(c)
			}
		}
		buf.WriteByte('"')
	case int:
		buf.WriteString(strconv.Itoa(v))
	case int64:
		buf.WriteString(strconv.FormatInt(v, 10))
	case float64:
		buf.WriteString(strconv.FormatFloat(v, 'f', -1, 64))
	case bool:
		buf.WriteString(strconv.FormatBool(v))
	case nil:
		buf.WriteString("null")
	default:
		// Fallback to json.Marshal for complex types
		data, _ := json.Marshal(v)
		buf.Write(data)
	}
}

// Example usage:
//
//	formatter := NewJSONFormatter()
//	formatter.SetGlobalField("service", "my-api")
//	formatter.SetGlobalField("version", "1.0.0")
//	formatter.IncludeRuntimeFields = true
//	formatter.FieldOrder = []string{"timestamp", "level", "service", "message"}
//
//	InitRuntimeFields("production", map[string]string{
//	    "region": "us-east-1",
//	    "cluster": "prod-cluster-1",
//	})
