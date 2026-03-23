package utils

import (
	"fmt"
	"sort"
)

// FormatFields converts a map of fields to a sorted string representation
func FormatFields(fields map[string]interface{}) string {
	if len(fields) == 0 {
		return ""
	}

	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	result := ""
	for i, k := range keys {
		if i > 0 {
			result += " "
		}
		result += fmt.Sprintf("%s=%v", k, fields[k])
	}
	return result
}

// SanitizeString removes or escapes problematic characters in log output
// Handles newlines, null bytes, control characters, and ANSI escape codes
func SanitizeString(s string) string {
	if s == "" {
		return s
	}

	// Pre-allocate with same capacity
	result := make([]byte, 0, len(s))

	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '\n':
			// Replace newlines with \\n for single-line output
			result = append(result, '\\', 'n')
		case c == '\r':
			// Replace carriage returns with \\r
			result = append(result, '\\', 'r')
		case c == '\t':
			// Replace tabs with \\t
			result = append(result, '\\', 't')
		case c == 0:
			// Skip null bytes
			continue
		case c == 0x1B && i+1 < len(s) && s[i+1] == '[':
			// Skip ANSI escape sequences (e.g., color codes)
			// Format: ESC [ ... m
			j := i + 2
			for j < len(s) && s[j] != 'm' {
				j++
			}
			if j < len(s) {
				i = j // Skip entire sequence
			}
		case c < 32 && c != '\n' && c != '\r' && c != '\t':
			// Replace other control characters with ?
			result = append(result, '?')
		default:
			result = append(result, c)
		}
	}

	return string(result)
}

// MergeMaps merges two maps, with values from the second map taking precedence
func MergeMaps(base, override map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{}, len(base)+len(override))

	for k, v := range base {
		result[k] = v
	}

	for k, v := range override {
		result[k] = v
	}

	return result
}
