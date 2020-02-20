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

// LevelStrings is a slice of our ErrorLevel enum
var LevelStrings = map[Level]string{
	DEBUG: "debug",
	INFO:  "info",
	WARN:  "warn",
	ERROR: "error",
	FATAL: "fatal",
}

// StringToLevel returns the match Level. "warn" = WARN
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

// NotDebug returns true for all Levels except DEBUG
func (l Level) NotDebug() bool {
	switch l {
	case DEBUG:
		return false
	default:
		return true
	}
}

// IsError returns true for ERROR or FATAL
func (l Level) IsError() bool {
	switch l {
	case ERROR, FATAL:
		return true
	default:
		return false
	}
}

// IsFatal returns true for anything above ERROR
func (l Level) IsFatal() bool {
	switch l {
	case FATAL:
		return true
	default:
		return false
	}
}

// String converts ErrorLevel value to a string value
func (l Level) String() string {
	return LevelStrings[l]
}

// MarshalJSON converts Level to a json string
func (l Level) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(l.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON converts a json string to a Level
func (l *Level) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	*l = StringToLevel(s)
	return nil
}
