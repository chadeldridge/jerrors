package jerrors

import (
	"encoding/json"
	"fmt"
	"log"
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

func newError() Error {
	m := make(map[string]string)
	return Error{Metadata: m}
}

// NewError creates a new Error object and returns it.
// args should be in the for of keyString1, valueString1,...
func NewError(level Level, msg string, args ...interface{}) Error {
	e := newError()
	e.Level = level
	e.Message = msg

	if config.LogTime {
		t := time.Now()
		e.Time = &t
	}

	if config.LogCaller {
		e.Metadata["caller"] = getCaller()
	}

	// Convert args to key value pairs
	for i, arg := range args {
		if i%2 != 0 {
			continue
		}

		if i+1 > len(args) {
			break
		}

		e.Metadata[fmt.Sprint(arg)] = fmt.Sprint(args[i+1])
	}

	return e
}

func (e *Error) AddMetadata(args ...interface{}) {
	l := len(args)

	for i, arg := range args {
		if i%2 != 0 {
			continue
		}
		if l < i+1 {
			return
		}

		e.Metadata[fmt.Sprint(arg)] = fmt.Sprint(args[i+1])
	}
}

func (e *Error) String() string {
	if !config.LogLevel {
		e.Level = 0
	}

	j, _ := json.Marshal(e)
	return string(j)
}

func (e *Error) Error() string {
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
func (e *Error) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &e)
	return err
}

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
