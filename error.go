package jerrors

import (
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"runtime"
	"strings"
	"time"
)

// Error holds our Level and Message data map.
type Error struct {
	Time     *time.Time        `json:"time,omitempty"`
	Level    Level             `json:"level,omitempty"`
	Message  string            `json:"message,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// NewError creates a new Error object and returns it.
// args should be in the for of keyString1, valueString1,...
func NewError(level Level, msg string, args ...interface{}) Error {
	// Create a base error.
	e := Error{Metadata: make(map[string]string)}

	// Set the error level and message.
	e.Level = level
	e.Message = msg

	// Check if we should log the time.
	if config.LogTime {
		t := time.Now()
		e.Time = &t
	}

	// Check if we should log the caller.
	if config.LogCaller {
		e.Metadata["caller"] = getCaller()
	}

	// Convert args to key value pairs
	e.AddMetadata(args...)

	return e
}

// Merge converts args into string pairs and adds them to the Error's Metadata.
func (e *Error) AddMetadata(args ...interface{}) {
	l := len(args)

	for i, arg := range args {
		if i%2 != 0 {
			continue
		}

		if i+1 == l {
			return
		}

		e.Metadata[fmt.Sprint(arg)] = fmt.Sprint(args[i+1])

		// If this was the last key, we're done.
		if i+2 >= l {
			return
		}
	}
}

// Equal returns true if the Error is equal to the given Error. Equal does not compare Time.
func (e Error) Equal(error Error) bool {
	if e.Level != error.Level || e.Message != error.Message {
		return false
	}

	return maps.Equal(e.Metadata, error.Metadata)
}

// String returns the string representation of the Error.
func (e Error) String() string {
	if !config.LogLevel {
		e.Level = 0
	}

	j, _ := json.Marshal(e)
	return string(j)
}

// Error returns the string representation of the Error.
func (e Error) Error() string {
	return e.String()
}

/*
// MarshalJSON converts Error to a json byte array.
func (e *Error) MarshalJSON() ([]byte, error) {
	if !logLevel {
		e.Level = 0
	}
	j, _ := json.Marshal(e)
	return j, nil
}
*/

// UnmarshalJSON converts a json byte array to an Error.
/*
func (e *Error) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, e)
	return err
}
*/

// IsError returns true for anything above WARN
func (e *Error) IsError() bool {
	return e.Level.IsError()
}

// IsFatal returns true for anything above ERROR
func (e *Error) IsFatal() bool {
	return e.Level.IsFatal()
}

// Log logs the error with the appropriate logger type.
func (e *Error) Log() {
	if e.Level >= config.LoggingLevel {
		log.Println(e)
	}
}

// Fatal logs and exits as a fatal error.
func (e *Error) Fatal() {
	if len(e.Message) > 0 {
		e.Level = FATAL
		log.Fatalln(e)
	}
}

func getCaller() string {
	callers := make([]uintptr, 0, config.CallersToShow)
	for i := config.CallerDepth; i <= config.CallerDepth+config.CallersToShow; i++ {
		c, _, _, _ := runtime.Caller(i)
		callers = append(callers, c)
	}
	frames := runtime.CallersFrames(callers)
	s := make([]string, 0, config.CallersToShow)
	for i := 0; i < config.CallersToShow; i++ {
		f, _ := frames.Next()
		s = append([]string{fmt.Sprintf("%s{%d}", f.Function, f.Line)}, s...)
	}
	return strings.Join(s, "->")
}
