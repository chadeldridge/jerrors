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

var levelStrings = map[Level]string{
	DEBUG: "debug",
	INFO:  "info",
	WARN:  "warn",
	ERROR: "error",
	FATAL: "fatal",
}

// StringToLevel returns the match Level. "debug" = DEBUG
// level arg is NOT case sensitive. No match returns 0.
func StringToLevel(level string) Level {
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

// String converts Level to a lowercase string. DEBUG = "debug", etc.
// Returns empty string if Level is 0.
func (l Level) String() string {
	if val, ok := levelStrings[l]; ok {
		return val
	}
	return ""
}

// NotDebug returns true for all Levels except DEBUG.
func (l Level) NotDebug() bool {
	if l > 1 {
		return true
	}
	return false
}

// IsError returns true if Level is ERROR or FATAL.
func (l Level) IsError() bool {
	if l >= 4 {
		return true
	}
	return false
}

// IsFatal returns true if Level is FATAL.
func (l Level) IsFatal() bool {
	if l == FATAL {
		return true
	}
	return false
}

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

	*l = StringToLevel(s)
	return nil
}
