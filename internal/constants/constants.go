package constants

import "time"

// Default configuration values
const (
	DefaultBufferSize    = 10000
	DefaultFlushInterval = 100 * time.Millisecond
)

// Format constants
const (
	JSONFormat   = "json"
	TextFormat   = "text"
	LogFmtFormat = "logfmt"
	CustomFormat = "custom"
)

// Log level strings
const (
	LevelDebugStr = "debug"
	LevelInfoStr  = "info"
	LevelWarnStr  = "warn"
	LevelErrorStr = "error"
	LevelFatalStr = "fatal"
)

// Template placeholders for custom formatter
const (
	PlaceholderTimestamp = "{{timestamp}}"
	PlaceholderLevel     = "{{level}}"
	PlaceholderLogger    = "{{logger}}"
	PlaceholderMessage   = "{{message}}"
	PlaceholderCaller    = "{{caller}}"
	PlaceholderFields    = "{{fields}}"
)

// Time format constants
const (
	RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"
	RFC3339Micro = "2006-01-02T15:04:05.000000Z07:00"
)

// TimestampFormat types for user configuration
type TimestampFormat string

const (
	// RFC3339Nano is the default format with nanosecond precision
	// Example: 2026-03-23T10:15:30.123456789Z
	TimestampFormatRFC3339Nano TimestampFormat = "rfc3339nano"

	// RFC3339 is standard RFC3339 with second precision
	// Example: 2026-03-23T10:15:30Z
	TimestampFormatRFC3339 TimestampFormat = "rfc3339"

	// RFC3339Millis includes millisecond precision
	// Example: 2026-03-23T10:15:30.123Z
	TimestampFormatRFC3339Millis TimestampFormat = "rfc3339milli"

	// Unix is Unix timestamp (seconds since epoch)
	// Example: 1711187730
	TimestampFormatUnix TimestampFormat = "unix"

	// UnixMilli is Unix timestamp in milliseconds
	// Example: 1711187730123
	TimestampFormatUnixMilli TimestampFormat = "unixmilli"

	// UnixNano is Unix timestamp in nanoseconds
	// Example: 1711187730123456789
	TimestampFormatUnixNano TimestampFormat = "unixnano"

	// DateTime is human-readable format
	// Example: 2026-03-23 10:15:30
	TimestampFormatDateTime TimestampFormat = "datetime"

	// Custom allows user to specify their own format string
	TimestampFormatCustom TimestampFormat = "custom"
)

// Parse error field name
const (
	ParseErrorField = "log_parse_error"
)
