package formatter

import (
	"fmt"
	"strings"
	"time"
)

// CustomTemplate represents a compiled custom log format template
type CustomTemplate struct {
	template string
	parts    []templatePart
}

type templatePart struct {
	isPlaceholder bool
	value         string
}

// Supported placeholders
const (
	PlaceholderTimestamp = "{{timestamp}}"
	PlaceholderLevel     = "{{level}}"
	PlaceholderLogger    = "{{logger}}"
	PlaceholderMessage   = "{{message}}"
	PlaceholderCaller    = "{{caller}}"
	PlaceholderFields    = "{{fields}}"
	PlaceholderSource    = "{{source}}"
	PlaceholderTraceID   = "{{trace_id}}"
	PlaceholderSpanID    = "{{span_id}}"
	// Runtime fields
	PlaceholderHostname  = "{{hostname}}"
	PlaceholderPID       = "{{pid}}"
	PlaceholderGoroutine = "{{goroutine}}"
	PlaceholderEnv       = "{{env}}"
)

// ValidPlaceholders lists all valid placeholders
var ValidPlaceholders = []string{
	PlaceholderTimestamp,
	PlaceholderLevel,
	PlaceholderLogger,
	PlaceholderMessage,
	PlaceholderCaller,
	PlaceholderFields,
	PlaceholderSource,
	PlaceholderTraceID,
	PlaceholderSpanID,
	PlaceholderHostname,
	PlaceholderPID,
	PlaceholderGoroutine,
	PlaceholderEnv,
}

// CompileTemplate compiles a custom template string into an executable template
// Returns an error if the template contains invalid placeholders or syntax
func CompileTemplate(template string) (*CustomTemplate, error) {
	if template == "" {
		return nil, fmt.Errorf("template cannot be empty")
	}

	// Validate template
	if err := ValidateTemplate(template); err != nil {
		return nil, err
	}

	// Parse template into parts
	parts := parseTemplate(template)

	return &CustomTemplate{
		template: template,
		parts:    parts,
	}, nil
}

// ValidateTemplate validates a template string
// Checks for:
// - Required placeholder presence (at least timestamp and message)
// - Invalid placeholder syntax
// - Unknown placeholders
func ValidateTemplate(template string) error {
	// Must contain at least message placeholder
	if !strings.Contains(template, PlaceholderMessage) {
		return fmt.Errorf("template must contain %s placeholder", PlaceholderMessage)
	}

	// Check for incomplete placeholders
	if strings.Count(template, "{{") != strings.Count(template, "}}") {
		return fmt.Errorf("template has unbalanced braces: {{ and }} must match")
	}

	// Find all placeholders
	placeholders := extractPlaceholders(template)

	// Validate each placeholder
	for _, ph := range placeholders {
		if !isValidPlaceholder(ph) {
			return fmt.Errorf("unknown placeholder: %s (valid: %v)", ph, ValidPlaceholders)
		}
	}

	return nil
}

// extractPlaceholders finds all {{...}} patterns in the template
func extractPlaceholders(template string) []string {
	var placeholders []string
	start := 0

	for {
		// Find next {{
		openIdx := strings.Index(template[start:], "{{")
		if openIdx == -1 {
			break
		}
		openIdx += start

		// Find matching }}
		closeIdx := strings.Index(template[openIdx:], "}}")
		if closeIdx == -1 {
			break
		}
		closeIdx += openIdx + 2

		placeholder := template[openIdx:closeIdx]
		placeholders = append(placeholders, placeholder)

		start = closeIdx
	}

	return placeholders
}

// isValidPlaceholder checks if a placeholder is in the valid list
func isValidPlaceholder(placeholder string) bool {
	for _, valid := range ValidPlaceholders {
		if placeholder == valid {
			return true
		}
	}
	return false
}

// parseTemplate parses a template into static and placeholder parts
func parseTemplate(template string) []templatePart {
	var parts []templatePart
	remaining := template

	for len(remaining) > 0 {
		// Find next placeholder
		idx := strings.Index(remaining, "{{")
		if idx == -1 {
			// No more placeholders, add remaining as static
			if remaining != "" {
				parts = append(parts, templatePart{
					isPlaceholder: false,
					value:         remaining,
				})
			}
			break
		}

		// Add static part before placeholder
		if idx > 0 {
			parts = append(parts, templatePart{
				isPlaceholder: false,
				value:         remaining[:idx],
			})
		}

		// Find end of placeholder
		endIdx := strings.Index(remaining[idx:], "}}")
		if endIdx == -1 {
			// Malformed, treat rest as static
			parts = append(parts, templatePart{
				isPlaceholder: false,
				value:         remaining[idx:],
			})
			break
		}

		endIdx += idx + 2
		placeholder := remaining[idx:endIdx]

		parts = append(parts, templatePart{
			isPlaceholder: true,
			value:         placeholder,
		})

		remaining = remaining[endIdx:]
	}

	return parts
}

// Execute executes the template with the given log entry
func (ct *CustomTemplate) Execute(entry *LogEntry) string {
	var result strings.Builder

	for _, part := range ct.parts {
		if part.isPlaceholder {
			result.WriteString(ct.resolvePlaceholder(part.value, entry))
		} else {
			result.WriteString(part.value)
		}
	}

	return result.String()
}

// resolvePlaceholder resolves a placeholder to its actual value
func (ct *CustomTemplate) resolvePlaceholder(placeholder string, entry *LogEntry) string {
	switch placeholder {
	case PlaceholderTimestamp:
		// Note: Custom templates use RFC3339Nano by default
		// For custom timestamp formatting, use NewFormatterWithTimestamp
		return entry.Timestamp.Format(time.RFC3339Nano)
	case PlaceholderLevel:
		return entry.Level.String()
	case PlaceholderLogger:
		return entry.Logger
	case PlaceholderMessage:
		return entry.Message
	case PlaceholderCaller:
		return entry.Caller
	case PlaceholderSource:
		return entry.Source
	case PlaceholderTraceID:
		return entry.TraceID
	case PlaceholderSpanID:
		return entry.SpanID
	case PlaceholderFields:
		return formatFieldsCustom(entry.Fields)
	default:
		return placeholder // Return as-is if unknown
	}
}

// formatFieldsCustom formats fields for custom template
func formatFieldsCustom(fields map[string]interface{}) string {
	if len(fields) == 0 {
		return ""
	}

	var result strings.Builder
	first := true

	for k, v := range fields {
		if !first {
			result.WriteString(", ")
		}
		first = false
		result.WriteString(fmt.Sprintf("%s=%v", k, v))
	}

	return result.String()
}
