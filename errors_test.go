package jerrors

import (
	"bytes"
	"regexp"
	"strconv"
	"testing"
	"time"
)

func ops(logCaller, logLevel, logTime bool, level Level) {
	o := map[string]bool{
		"caller": logCaller,
		"level":  logLevel,
		"time":   logTime,
	}
	SetOptions(o)
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
	line = line[0 : len(line)-1]
	return line
}

func testMatchString(t *testing.T, pattern string, s string) {
	if m, _ := regexp.MatchString(pattern, s); !m {
		t.Errorf("pattern (%s) not found:\n%s", pattern, s)
	}
}

func TestNewError(t *testing.T) {
	opsDefault()
	err := New(ERROR, "test error", "type", "test", "user", "test1")
	if err.Level != ERROR {
		t.Errorf("error level (%v), expected (Error)", err.Level)
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
		t.Errorf("error level (%v), expected (Error)", err.Level)
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
		t.Errorf("error level (%v), expected (Error)", err.Level)
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
	e := err.Error()
	if matched, _ := regexp.MatchString(`"level"`, e); matched {
		t.Errorf("error level found, expected missing level:\n%s", e)
	}
}

func TestNewErrorNoTime(t *testing.T) {
	ops(true, true, false, INFO)
	err := New(ERROR, "test error", "type", "test", "user", "test1")
	if err.Level != ERROR {
		t.Errorf("error level (%v), expected (Error)", err.Level)
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
	testMatchString(t, `"level":"error"`, s)
	testMatchString(t, `"message":"test error"`, s)
	testMatchString(t, `"caller":"[a-zA-Z]+`, s)
	testMatchString(t, `"type":"test"`, s)
	testMatchString(t, `"user":"test1"`, s)
	testMatchString(t, `"time":"`+strconv.Itoa(time.Now().Year()), s)
}
