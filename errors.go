package jerrors

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"
)

/*
callerDepth is how many function calls to step back before getting
Caller information. This should be enough to get us back to the
initiating function.
*/
const callerDepth = 2

// callersShow sets how many calling functions to show.
const callersShow = 2

// JError holds our Level and Message data map.
type JError struct {
	Level    Level     `json:"level,omitempty"`
	Time     time.Time `json:"time,omitempty"`
	Message  string    `json:"message"`
	Metadata map[string]string
}

func newError() *JError {
	m := make(map[string]string)
	return &JError{Metadata: m}
}

// New creates a new JError object and returns it.
// args should be in the for of keyString1, valueString1,...
func New(l Level, msg string, args ...interface{}) *JError {
	e := newError()
	if logTime {
		e.Time = time.Now()
	}
	if logCaller {
		e.Metadata["caller"] = getCaller()
	}

	// Convert args to key value pairs
	for i, arg := range args {
		e.Metadata[fmt.Sprint(arg)] = fmt.Sprint(args[i+1])
		if i+2 >= len(args) {
			break
		}
	}

	return e
}

func (e *JError) Error() string {
	if !logLevel {
		e.Level = 0
	}
	j, _ := json.Marshal(e)
	return string(j)
}

// IsError returns true for anything above WARN
func (e *JError) IsError() bool {
	return e.Level.IsError()
}

// IsFatal returns true for anything above ERROR
func (e *JError) IsFatal() bool {
	return e.Level.IsFatal()
}

// SetLevel the Level
func (e *JError) SetLevel(level Level) {
	e.Level = level
}

// Log logs the error with the appropriate logger type.
func (e *JError) Log() {
	if loggingLevel == e.Level {
		log.Println(e)
	}
}

// Fatal logs and exits as a fatal error.
func (e *JError) Fatal() {
	if len(e.Message) > 0 {
		e.Level = FATAL
		log.Fatalln(e)
	}
}

func getCaller() string {
	callers := make([]uintptr, 0, callersShow)
	for i := callerDepth; i <= callerDepth+callersShow; i++ {
		c, _, _, _ := runtime.Caller(i)
		callers = append(callers, c)
	}
	frames := runtime.CallersFrames(callers)
	s := make([]string, 0, callersShow)
	for i := 0; i < callersShow; i++ {
		f, _ := frames.Next()
		s = append([]string{fmt.Sprintf("%s{%d}", f.Function, f.Line)}, s...)
	}
	return strings.Join(s, "->")
}
