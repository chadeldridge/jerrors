package jerrors

import (
	"bytes"
	"encoding/json"
	"strings"
)

// Level is an enum
type Level int

const (
	// DEBUG error level
	DEBUG Level = iota + 1
	// INFO error level
	INFO
	// WARN error level
	WARN
	// ERROR error level
	ERROR
	// FATAL error level
	FATAL
)

// GetLevel returns the matching Level. "debug" = DEBUG. level arg is NOT case sensitive. No match returns 0.
func GetLevel(level string) Level {
	switch l := strings.ToLower(level); l {
	case "debug":
		return DEBUG
	case "info":
		return INFO
	case "warn":
		return WARN
	case "error":
		return ERROR
	case "fatal":
		return FATAL
	default:
		return 0
	}
}

// String returns Level as a lowercase string. DEBUG = "debug"
func (l Level) String() string {
	return []string{
		"",
		"debug",
		"info",
		"warn",
		"error",
		"fatal",
	}[l]
}

// NotDebug returns true if the provided Level is any Level except DEBUG.
func (l Level) NotDebug() bool { return l > DEBUG }

// IsError returns true if Level is ERROR or higher.
func (l Level) IsError() bool { return l >= ERROR }

// IsFatal returns true if Level is FATAL or higher.
func (l Level) IsFatal() bool { return l >= FATAL }

// MarshalJSON converts Level to json.
func (l Level) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(l.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON converts json to a Level.
func (l *Level) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	*l = GetLevel(s)
	return nil
}
