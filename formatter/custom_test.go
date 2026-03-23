package formatter

import (
	"testing"
	"time"
)

func TestValidateTemplate(t *testing.T) {
	tests := []struct {
		name      string
		template  string
		shouldErr bool
	}{
		{
			name:      "valid template with message",
			template:  "{{timestamp}} [{{level}}] {{message}}",
			shouldErr: false,
		},
		{
			name:      "missing message placeholder",
			template:  "{{timestamp}} [{{level}}]",
			shouldErr: true,
		},
		{
			name:      "unbalanced braces",
			template:  "{{timestamp} [{{level}}] {{message}}",
			shouldErr: true,
		},
		{
			name:      "unknown placeholder",
			template:  "{{timestamp}} {{unknown}} {{message}}",
			shouldErr: true,
		},
		{
			name:      "empty template",
			template:  "",
			shouldErr: true,
		},
		{
			name:      "all placeholders",
			template:  "{{timestamp}} [{{level}}] {{logger}} {{message}} {{caller}} {{fields}} {{trace_id}}",
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTemplate(tt.template)
			if tt.shouldErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestCompileTemplate(t *testing.T) {
	template := "{{timestamp}} [{{level}}] {{message}}"

	compiled, err := CompileTemplate(template)
	if err != nil {
		t.Fatalf("Failed to compile valid template: %v", err)
	}

	if compiled == nil {
		t.Error("Compiled template is nil")
	}

	if len(compiled.parts) == 0 {
		t.Error("Compiled template has no parts")
	}
}

func TestCustomTemplate_Execute(t *testing.T) {
	template := "[{{level}}] {{message}} ({{logger}})"

	compiled, err := CompileTemplate(template)
	if err != nil {
		t.Fatalf("Failed to compile template: %v", err)
	}

	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     LevelInfo,
		Logger:    "test-logger",
		Message:   "test message",
	}

	result := compiled.Execute(entry)

	expected := "[info] test message (test-logger)"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestCustomFormatterIntegration(t *testing.T) {
	template := "{{timestamp}} [{{level}}] {{message}}"

	formatter, err := NewCustomFormatter(template)
	if err != nil {
		t.Fatalf("Failed to create custom formatter: %v", err)
	}

	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     LevelInfo,
		Message:   "test message",
	}

	output := formatter.Format(entry)
	if len(output) == 0 {
		t.Error("Formatter produced no output")
	}

	// Should contain the message
	if !containsString(output, "test message") {
		t.Error("Output doesn't contain expected message")
	}

	// Should contain the level
	if !containsString(output, "info") {
		t.Error("Output doesn't contain expected level")
	}
}

func TestExtractPlaceholders(t *testing.T) {
	template := "{{timestamp}} [{{level}}] {{message}}"

	placeholders := extractPlaceholders(template)

	expected := 3
	if len(placeholders) != expected {
		t.Errorf("Expected %d placeholders, got %d", expected, len(placeholders))
	}

	if placeholders[0] != "{{timestamp}}" {
		t.Errorf("Expected {{timestamp}}, got %s", placeholders[0])
	}
}

func containsString(data []byte, s string) bool {
	str := string(data)
	return len(str) > 0 && len(s) > 0 && (str == s || len(str) > len(s) && (str[:len(s)] == s || str[len(str)-len(s):] == s || len(str) > len(s)+1 && containsSubstring(str, s)))
}

func containsSubstring(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func BenchmarkCustomTemplateExecute(b *testing.B) {
	template := "{{timestamp}} [{{level}}] {{logger}}: {{message}}"
	compiled, _ := CompileTemplate(template)

	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     LevelInfo,
		Logger:    "bench",
		Message:   "benchmark message",
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		compiled.Execute(entry)
	}
}
