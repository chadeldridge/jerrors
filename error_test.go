package jerrors

import (
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
	require.Equal(t, err.Metadata[mdTypeKey], mdTypeVal)
	require.Equal(t, err.Metadata[mdUserKey], mdUserVal)
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
