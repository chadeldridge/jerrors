package jerrors

import (
	"bytes"
	"regexp"
	"strconv"
	"testing"
	"time"
)

func ops(logCaller, logLevel, logTime bool, level Level) {
	SetOptions(map[string]bool{
		"caller": logCaller,
		"level":  logLevel,
		"time":   logTime,
	})
	SetLogLevel(level)
}

func opsDefault() {
	ops(true, true, true, INFO)
}

func setTestLog() *bytes.Buffer {
	buf := new(bytes.Buffer)
	SetLogOutput(buf)
	return buf
}

func getBufferString(buf *bytes.Buffer) string {
	line := buf.String()
	if len(line) == 0 {
		return ""
	}

	line = line[0 : len(line)-1]
	return line
}

func errorIfMatchString(t *testing.T, pattern string, s string) {
	if m, _ := regexp.MatchString(pattern, s); m {
		t.Errorf("pattern (%s) found:\n%s", pattern, s)
	}
}

func errorIfNotMatchString(t *testing.T, pattern string, s string) {
	if m, _ := regexp.MatchString(pattern, s); !m {
		t.Errorf("pattern (%s) not found:\n%s", pattern, s)
	}
}

func errorIfIsError(t *testing.T, err *Error) {
	if e := err.IsError(); e {
		t.Errorf("error is (%v) but IsError returned (%v):\n", err.Level, e)
	}
}

func errorIfNotIsError(t *testing.T, err *Error) {
	if e := err.IsError(); !e {
		t.Errorf("error is (%v) but IsError returned (%v):\n", err.Level, e)
	}
}

func errorIfIsFatal(t *testing.T, err *Error) {
	if e := err.IsFatal(); e {
		t.Errorf("error is (%v) but IsFatal returned (%v):\n", err.Level, e)
	}
}

func errorIfNotIsFatal(t *testing.T, err *Error) {
	if e := err.IsFatal(); !e {
		t.Errorf("error is (%v) but IsFatal returned (%v):\n", err.Level, e)
	}
}

func TestNewError(t *testing.T) {
	opsDefault()
	err := New(ERROR, "test error", "type", "test", "user", "test1")
	if err.Level != ERROR {
		t.Errorf("error level (%v), expected (%v)", err.Level, ERROR)
	}
	if err.Message != "test error" {
		t.Errorf("error message (%v), expected (test error)", err.Message)
	}
	if err.Metadata["caller"] == "" {
		t.Error("error metadata 'caller' is blank")
	}
	if err.Metadata["type"] != "test" {
		t.Errorf("error metadata 'type' (%v), expected (test)", err.Message)
	}
	if err.Metadata["user"] != "test1" {
		t.Errorf("error metadata 'user' (%v), expected (test1)", err.Message)
	}
	if err.Time.IsZero() {
		t.Errorf("error time (%v), expected non zero time", err.Time)
	}
}

func TestNewErrorNoCaller(t *testing.T) {
	ops(false, true, true, INFO)
	err := New(ERROR, "test error", "type", "test", "user", "test1")
	if err.Level != ERROR {
		t.Errorf("error level (%v), expected (%v)", err.Level, ERROR)
	}
	if err.Message != "test error" {
		t.Errorf("error message (%v), expected (test error)", err.Message)
	}
	if err.Metadata["caller"] != "" {
		t.Error("error metadata 'caller' exists")
	}
	if err.Metadata["type"] != "test" {
		t.Errorf("error metadata 'type' (%v), expected (test)", err.Message)
	}
	if err.Metadata["user"] != "test1" {
		t.Errorf("error metadata 'user' (%v), expected (test1)", err.Message)
	}
	if err.Time.IsZero() {
		t.Errorf("error time (%v), expected non zero time", err.Time)
	}
}

func TestNewErrorNoLevel(t *testing.T) {
	ops(true, false, true, INFO)
	err := New(ERROR, "test error", "type", "test", "user", "test1")
	if err.Level != ERROR {
		t.Errorf("error level (%v), expected (%v)", err.Level, ERROR)
	}
	if err.Message != "test error" {
		t.Errorf("error message (%v), expected (test error)", err.Message)
	}
	if err.Metadata["caller"] == "" {
		t.Error("error metadata 'caller' is blank")
	}
	if err.Metadata["type"] != "test" {
		t.Errorf("error metadata 'type' (%v), expected (test)", err.Message)
	}
	if err.Metadata["user"] != "test1" {
		t.Errorf("error metadata 'user' (%v), expected (test1)", err.Message)
	}
	if err.Time.IsZero() {
		t.Errorf("error time (%v), expected non zero time", err.Time)
	}
	s := err.Error()
	errorIfMatchString(t, `"level"`, s)
}

func TestNewErrorNoTime(t *testing.T) {
	ops(true, true, false, INFO)
	err := New(ERROR, "test error", "type", "test", "user", "test1")
	if err.Level != ERROR {
		t.Errorf("error level (%v), expected (%v)", err.Level, ERROR)
	}
	if err.Message != "test error" {
		t.Errorf("error message (%v), expected (test error)", err.Message)
	}
	if err.Metadata["caller"] == "" {
		t.Error("error metadata 'caller' is blank")
	}
	if err.Metadata["type"] != "test" {
		t.Errorf("error metadata 'type' (%v), expected (test)", err.Message)
	}
	if err.Metadata["user"] != "test1" {
		t.Errorf("error metadata 'user' (%v), expected (test1)", err.Message)
	}
	if !err.Time.IsZero() {
		t.Errorf("error time (%v), expected zero time", err.Time)
	}
}

func TestErrorLog(t *testing.T) {
	opsDefault()
	buf := setTestLog()
	err := New(ERROR, "test error", "type", "test", "user", "test1")
	err.Log()
	s := getBufferString(buf)
	errorIfNotMatchString(t, `"level":"`+err.Level.String()+`"`, s)
	errorIfNotMatchString(t, `"message":"test error"`, s)
	errorIfNotMatchString(t, `"caller":"[a-zA-Z]+`, s)
	errorIfNotMatchString(t, `"type":"test"`, s)
	errorIfNotMatchString(t, `"user":"test1"`, s)
	errorIfNotMatchString(t, `"time":"`+strconv.Itoa(time.Now().Year()), s)
}

func TestErrorLogNoLevel(t *testing.T) {
	ops(true, false, true, INFO)
	buf := setTestLog()
	err := New(ERROR, "test error", "type", "test", "user", "test1")
	err.Log()
	s := getBufferString(buf)
	errorIfMatchString(t, `"level":"`+err.Level.String()+`"`, s)
	errorIfNotMatchString(t, `"message":"test error"`, s)
	errorIfNotMatchString(t, `"caller":"[a-zA-Z]+`, s)
	errorIfNotMatchString(t, `"type":"test"`, s)
	errorIfNotMatchString(t, `"user":"test1"`, s)
	errorIfNotMatchString(t, `"time":"`+strconv.Itoa(time.Now().Year()), s)
}

func TestErrorLogCorrectLevel(t *testing.T) {
	ops(true, true, true, INFO)
	buf := setTestLog()
	err := New(DEBUG, "test error", "type", "test", "user", "test1")
	err.Log()
	s := getBufferString(buf)
	if s != "" {
		t.Errorf("error logged at level(%v), logLevel set at (%v):\n%s", ERROR, INFO, s)
	}
}

func TestErrorSetLevel(t *testing.T) {
	opsDefault()
	err := New(DEBUG, "test error", "type", "test", "user", "test1")
	if err.Level != DEBUG {
		t.Errorf("error level (%v), expected (%v)", err.Level, DEBUG)
	}
	err.SetLevel(ERROR)
	if err.Level != ERROR {
		t.Errorf("error level (%v), expected (%v)", err.Level, ERROR)
	}
}

func TestErrorIsError(t *testing.T) {
	opsDefault()
	err := New(DEBUG, "test error", "type", "test", "user", "test1")
	errorIfIsError(t, err)

	err.SetLevel(INFO)
	errorIfIsError(t, err)

	err.SetLevel(WARN)
	errorIfIsError(t, err)

	err.SetLevel(ERROR)
	errorIfNotIsError(t, err)

	err.SetLevel(FATAL)
	errorIfNotIsError(t, err)
}

func TestErrorIsFatal(t *testing.T) {
	opsDefault()
	err := New(DEBUG, "test error", "type", "test", "user", "test1")
	errorIfIsFatal(t, err)

	err.SetLevel(INFO)
	errorIfIsFatal(t, err)

	err.SetLevel(WARN)
	errorIfIsFatal(t, err)

	err.SetLevel(ERROR)
	errorIfIsFatal(t, err)

	err.SetLevel(FATAL)
	errorIfNotIsFatal(t, err)
}
