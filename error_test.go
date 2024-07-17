package jerrors

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testMessage = "test error"
	mdTypeKey   = "type"
	mdTypeVal   = "test"
	mdUserKey   = "user"
	mdUserVal   = "test1"
	debugErr    = NewError(DEBUG, testMessage, mdTypeKey, mdTypeVal, mdUserKey, mdUserVal)
	infoErr     = NewError(INFO, testMessage, mdTypeKey, mdTypeVal, mdUserKey, mdUserVal)
	warnErr     = NewError(WARN, testMessage, mdTypeKey, mdTypeVal, mdUserKey, mdUserVal)
	errorErr    = NewError(ERROR, testMessage, mdTypeKey, mdTypeVal, mdUserKey, mdUserVal)
	fatalErr    = NewError(FATAL, testMessage, mdTypeKey, mdTypeVal, mdUserKey, mdUserVal)
)

func TestError(t *testing.T) {
	SetConfig(DefaultConfig())

	err := NewError(ERROR, testMessage, mdTypeKey, mdTypeVal, mdUserKey, mdUserVal)
	require.Equal(t, err.Level, ERROR)
	require.Equal(t, err.Message, testMessage)
	_, ok := err.Metadata["caller"]
	require.False(t, ok)
	require.Equal(t, mdTypeVal, err.Metadata[mdTypeKey])
	require.Equal(t, mdUserVal, err.Metadata[mdUserKey])
	require.False(t, err.Time.IsZero())
}

func TestErrorWithCaller(t *testing.T) {
	c := DefaultConfig()
	c.LogCaller = true
	SetConfig(c)

	err := NewError(ERROR, testMessage, mdTypeKey, mdTypeVal, mdUserKey, mdUserVal)
	require.Equal(t, err.Level, ERROR)
	require.Equal(t, err.Message, testMessage)
	_, ok := err.Metadata["caller"]
	require.True(t, ok)
	require.NotEqual(t, err.Metadata["caller"], "")
	require.Equal(t, err.Metadata[mdTypeKey], mdTypeVal)
	require.Equal(t, err.Metadata[mdUserKey], mdUserVal)
	require.False(t, err.Time.IsZero())
}

func TestErrorNoLevel(t *testing.T) {
	c := DefaultConfig()
	c.LogLevel = false
	SetConfig(c)

	err := NewError(ERROR, testMessage, mdTypeKey, mdTypeVal, mdUserKey, mdUserVal)
	require.Equal(t, err.Message, testMessage)
	_, ok := err.Metadata["caller"]
	require.False(t, ok)
	require.Equal(t, err.Metadata[mdTypeKey], mdTypeVal)
	require.Equal(t, err.Metadata[mdUserKey], mdUserVal)
	require.False(t, err.Time.IsZero())

	s := err.String()
	require.NotContains(t, s, "level")
}

func TestErrorNoTime(t *testing.T) {
	c := DefaultConfig()
	c.LogTime = false
	SetConfig(c)

	err := NewError(ERROR, testMessage, mdTypeKey, mdTypeVal, mdUserKey, mdUserVal)
	require.Equal(t, err.Level, ERROR)
	require.Equal(t, err.Message, testMessage)
	_, ok := err.Metadata["caller"]
	require.False(t, ok)
	require.Equal(t, err.Metadata[mdTypeKey], mdTypeVal)
	require.Equal(t, err.Metadata[mdUserKey], mdUserVal)
	require.Empty(t, err.Time)
}

func TestErrorChangeFields(t *testing.T) {
	SetConfig(DefaultConfig())

	err := NewError(ERROR, testMessage, mdTypeKey, mdTypeVal, mdUserKey, mdUserVal)
	// Verify our defaults are correct
	require.Equal(t, err.Level, ERROR)
	require.Equal(t, err.Message, testMessage)
	_, ok := err.Metadata["caller"]
	require.False(t, ok)
	require.Equal(t, err.Metadata[mdTypeKey], mdTypeVal)
	require.Equal(t, err.Metadata[mdUserKey], mdUserVal)
	require.False(t, err.Time.IsZero())

	newMessage := "change message"
	newKey := "myKey"
	newVal := "my bigger value"

	err.Level = FATAL
	err.Message = newMessage
	err.Metadata[newKey] = newVal
	require.Equal(t, err.Level, FATAL)
	require.Equal(t, err.Message, newMessage)
	require.Equal(t, err.Metadata[newKey], newVal)
}

func TestErrorAddMetadata(t *testing.T) {
	SetConfig(DefaultConfig())

	err := NewError(ERROR, testMessage)
	err.AddMetadata("key1", "val1", "key2", "val2")
	require.Equal(t, err.Metadata["key1"], "val1")
	require.Equal(t, err.Metadata["key2"], "val2")

	// Odd number of args.
	err = NewError(ERROR, testMessage)
	err.AddMetadata("key1", "val1", "key2")
	require.Equal(t, err.Metadata["key1"], "val1")
	require.NotContains(t, err.Metadata, "key2")
}

func TestErrorEqual(t *testing.T) {
	SetConfig(DefaultConfig())

	// True
	err1 := NewError(ERROR, testMessage, mdTypeKey, mdTypeVal, mdUserKey, mdUserVal)
	err2 := err1
	require.True(t, err1.Equal(err2))

	// Different level.
	err2.Level = INFO
	require.False(t, err1.Equal(err2))

	// Different Message.
	err2 = err1
	err2.Message = "new message"
	require.False(t, err1.Equal(err2))

	// Different time.
	err3 := NewError(ERROR, "different message", mdTypeKey, mdTypeVal, mdUserKey, mdUserVal)
	require.False(t, err1.Equal(err3))

	// Different Metadata.
	err3.Metadata["newKey"] = "newValue"
	require.False(t, err1.Equal(err3))
}

func TestErrorString(t *testing.T) {
	SetConfig(DefaultConfig())

	err := NewError(ERROR, testMessage, mdTypeKey, mdTypeVal, mdUserKey, mdUserVal)
	s := err.String()
	require.Contains(t, s, `"time":"`)
	require.Contains(t, s, `"level":"error","message":"test error","metadata":{"type":"test","user":"test1"}}`)

	// LogLevel = false
	config.LogLevel = false
	s = err.String()
	require.Contains(t, s, `"time":"`)
	require.NotContains(t, s, `"level":"error"`)
	require.Contains(t, s, `"message":"test error","metadata":{"type":"test","user":"test1"}}`)
}

func TestErrorUnmarshalJSON(t *testing.T) {
	SetConfig(DefaultConfig())
	//`{"time":"2024-07-17T12:41:40.064205262-04:00","level":"error","message":"test error","metadata":{"type":"test"}}`,
	j := []byte(
		`{"time":"2024-07-17T12:41:40.064205262-04:00","level":"error","message":"test error","metadata":{"type":"test"}}`,
	)
	/*
		e := NewError(ERROR, testMessage, mdTypeKey, mdTypeVal, mdUserKey, mdUserVal)
		j, err := json.Marshal(e)
		if err != nil {
			log.Println(err)
		}
	*/

	jerr := &Error{}
	// jerr.UnmarshalJSON(j)
	err := json.Unmarshal(j, jerr)
	require.Nil(t, err)

	require.Equal(t, ERROR, jerr.Level)
	require.Equal(t, testMessage, jerr.Message)
	require.Equal(t, mdTypeVal, jerr.Metadata[mdTypeKey])
	require.False(t, jerr.Time.IsZero())
}

func TestErrorIsError(t *testing.T) {
	SetConfig(DefaultConfig())

	err := debugErr
	require.False(t, err.IsError())

	err.Level = INFO
	require.False(t, err.IsError())

	err.Level = WARN
	require.False(t, err.IsError())

	err.Level = ERROR
	require.True(t, err.IsError())

	err.Level = FATAL
	require.True(t, err.IsError())
}

func TestErrorIsFatal(t *testing.T) {
	SetConfig(DefaultConfig())

	err := debugErr
	require.False(t, err.IsFatal())

	err.Level = INFO
	require.False(t, err.IsFatal())

	err.Level = WARN
	require.False(t, err.IsFatal())

	err.Level = ERROR
	require.False(t, err.IsFatal())

	err.Level = FATAL
	require.True(t, err.IsFatal())
}

func TestErrorLog(t *testing.T) {
	SetConfig(DefaultConfig())

	// Set log to buffer.
	var buf bytes.Buffer
	log.SetOutput(&buf)

	// Grab a predefined error and log it.
	err := errorErr
	err.Log()

	require.NotEmpty(t, buf)
	require.Contains(t, buf.String(), "test error")

	// Set log back to stderr just in case.
	log.SetOutput(os.Stderr)
}
