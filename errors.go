package jerrors

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"
)

var (
	logCaller    bool
	logLevel     bool
	logTime      bool
	loggingLevel Level
)

func init() {
	// Setup default log options
	logCaller = true
	logLevel = true
	logTime = true
	loggingLevel = INFO
}

/*
callerDepth is how many function calls to step back before getting
Caller information. This should be enough to get us back to the
initiating function.
*/
const callerDepth = 2

// callersShow sets how many calling functions to show.
const callersShow = 2

// Error holds our Level and Message data map.
type Error struct {
	Level   Level
	Message map[string]string
}

func newError() Error {
	msg := make(map[string]string)
	return Error{Message: msg}
}

// SetOptions for additional error data.
func SetOptions(options map[string]bool) {
	for k, v := range options {
		switch k {
		case "caller":
			logCaller = v
		case "level":
			logLevel = v
		case "time":
			logTime = v
		}
	}
}

// SetLogLevel sets what level to log. Will log set level and above.
func SetLogLevel(level Level) {
	loggingLevel = level
}

// New creates a new Error object and returns it.
// args should be in the for of keyString1, valueString1,...
func New(l Level, args ...interface{}) Error {
	e := newError()
	if logTime {
		e.Message["time"] = time.Now().String()
	}
	if logCaller {
		e.Message["caller"] = getCaller()
	}

	// Convert args to key value pairs
	for i, arg := range args {
		e.Message[fmt.Sprint(arg)] = fmt.Sprint(args[i+1])
		if i+2 >= len(args) {
			break
		}
	}

	return e
}

func (e *Error) String() string {
	if logLevel {
		e.Message["level"] = e.Level.String()
	}
	j, _ := json.Marshal(e.Message)
	return string(j)
}

// IsError returns true for anything above WARN
func (e *Error) IsError() bool {
	return e.Level.IsError()
}

// IsFatal returns true for anything above ERROR
func (e *Error) IsFatal() bool {
	return e.Level.IsFatal()
}

// SetLevel the Level
func (e *Error) SetLevel(level Level) {
	e.Level = level
}

// Log logs the error with the appropriate logger type.
func (e *Error) Log() {
	if loggingLevel == e.Level {
		log.Println(e.String())
	}
}

// Fatal logs and exits as a fatal error.
func (e *Error) Fatal() {
	if len(e.Message) > 0 {
		e.Level = FATAL
		log.Fatalln(e.String())
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
