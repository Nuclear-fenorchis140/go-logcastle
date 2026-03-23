package utils

import (
	"fmt"
	"strconv"
	"time"
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

// FormatTimestamp formats a timestamp according to the specified format
func FormatTimestamp(t time.Time, format TimestampFormat, customFormat string) string {
	switch format {
	case TimestampFormatRFC3339Nano:
		return t.Format(time.RFC3339Nano)
	case TimestampFormatRFC3339:
		return t.Format(time.RFC3339)
	case TimestampFormatRFC3339Millis:
		return t.Format("2006-01-02T15:04:05.000Z07:00")
	case TimestampFormatUnix:
		return strconv.FormatInt(t.Unix(), 10)
	case TimestampFormatUnixMilli:
		return strconv.FormatInt(t.UnixMilli(), 10)
	case TimestampFormatUnixNano:
		return strconv.FormatInt(t.UnixNano(), 10)
	case TimestampFormatDateTime:
		return t.Format("2006-01-02 15:04:05")
	case TimestampFormatCustom:
		if customFormat != "" {
			return t.Format(customFormat)
		}
		// Fallback to RFC3339Nano if custom format not provided
		return t.Format(time.RFC3339Nano)
	default:
		// Default to RFC3339Nano for unknown formats
		return t.Format(time.RFC3339Nano)
	}
}

// GetGoTimeFormat returns the Go time layout string for a timestamp format
func GetGoTimeFormat(format TimestampFormat, customFormat string) string {
	switch format {
	case TimestampFormatRFC3339Nano:
		return time.RFC3339Nano
	case TimestampFormatRFC3339:
		return time.RFC3339
	case TimestampFormatRFC3339Millis:
		return "2006-01-02T15:04:05.000Z07:00"
	case TimestampFormatUnix, TimestampFormatUnixMilli, TimestampFormatUnixNano:
		return "" // Numeric formats don't have a layout
	case TimestampFormatDateTime:
		return "2006-01-02 15:04:05"
	case TimestampFormatCustom:
		if customFormat != "" {
			return customFormat
		}
		return time.RFC3339Nano
	default:
		return time.RFC3339Nano
	}
}

// ValidateTimestampFormat validates a timestamp format configuration
func ValidateTimestampFormat(format TimestampFormat, customFormat string) error {
	if format == TimestampFormatCustom && customFormat == "" {
		return fmt.Errorf("custom timestamp format string must be provided when using TimestampFormatCustom")
	}
	return nil
}
