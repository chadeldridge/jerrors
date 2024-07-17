package jerrors

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrors(t *testing.T) {
	SetConfig(DefaultConfig())

	errs := New()
	errs.NewError(WARN, testMessage, mdTypeKey, mdTypeVal, mdUserKey, mdUserVal)
	require.Equal(t, 1, len(errs.Errors))
	require.Equal(t, WARN, errs.Level)

	// errs.Level should remain WARN
	errs.Add(debugErr)
	require.Equal(t, 2, len(errs.Errors))
	require.Equal(t, WARN, errs.Level)

	// errs.Level should change to ERROR
	errs.Add(errorErr)
	require.Equal(t, 3, len(errs.Errors))
	require.Equal(t, ERROR, errs.Level)
}

func TestErrorsRemove(t *testing.T) {
	SetConfig(DefaultConfig())

	errs := New()
	errs.Add(debugErr)
	errs.Add(infoErr)
	errs.Add(warnErr)
	errs.Add(errorErr)
	require.Equal(t, 4, len(errs.Errors))
	require.Equal(t, ERROR, errs.Level)

	errs.Remove(errorErr)
	require.Equal(t, 3, len(errs.Errors))
	require.Equal(t, WARN, errs.Level)

	e := NewError(ERROR, testMessage, mdTypeKey, mdTypeVal, mdUserKey, mdUserVal)
	errs.Remove(e)
	require.Equal(t, 3, len(errs.Errors))
	require.Equal(t, WARN, errs.Level)
}

func TestErrorsUpdateLevel(t *testing.T) {
	SetConfig(DefaultConfig())

	errs := New()
	errs.Add(debugErr)
	errs.Add(infoErr)
	errs.Add(warnErr)
	errs.Add(errorErr)
	require.Equal(t, ERROR, errs.UpdateLevel())

	// errs.Level should remain unchanged.
	errs.Remove(infoErr)
	require.Equal(t, ERROR, errs.UpdateLevel())

	// errs.Level should change to WARN.
	errs.Remove(errorErr)
	require.Equal(t, WARN, errs.UpdateLevel())
}

func TestErrorsCheck(t *testing.T) {
	SetConfig(DefaultConfig())

	errs := New()
	n, hasErrors := errs.Check()
	require.Equal(t, 0, n)
	require.False(t, hasErrors)

	errs.Add(debugErr)
	errCount, hasErrors := errs.Check()
	require.Equal(t, 1, errCount)
	require.True(t, hasErrors)
}

func TestErrorsFirst(t *testing.T) {
	SetConfig(DefaultConfig())

	errs := New()
	require.Equal(t, ErrNoErrorFound, errs.First())

	errs.Add(debugErr)
	errs.Add(infoErr)
	errs.Add(warnErr)
	errs.Add(errorErr)
	require.Equal(t, debugErr, errs.First())
}

func TestErrorsLast(t *testing.T) {
	SetConfig(DefaultConfig())

	errs := New()
	require.Equal(t, ErrNoErrorFound, errs.Last())

	errs.Add(debugErr)
	errs.Add(infoErr)
	errs.Add(warnErr)
	errs.Add(errorErr)
	require.Equal(t, errorErr, errs.Last())
}

func TestErrorsIsEmpty(t *testing.T) {
	SetConfig(DefaultConfig())

	errs := New()
	require.True(t, errs.IsEmpty())

	errs.Add(debugErr)
	require.False(t, errs.IsEmpty())
}

func TestErrorsIsError(t *testing.T) {
	SetConfig(DefaultConfig())

	errs := New()
	errs.Add(debugErr)
	require.False(t, errs.IsError())

	errs.Add(infoErr)
	require.False(t, errs.IsError())

	errs.Add(warnErr)
	require.False(t, errs.IsError())

	errs.Add(errorErr)
	require.True(t, errs.IsError())

	errs.Add(fatalErr)
	require.True(t, errs.IsError())
}

func TestErrorsIsFatal(t *testing.T) {
	SetConfig(DefaultConfig())

	errs := New()
	errs.Add(debugErr)
	require.False(t, errs.IsFatal())

	errs.Add(infoErr)
	require.False(t, errs.IsFatal())

	errs.Add(warnErr)
	require.False(t, errs.IsFatal())

	errs.Add(errorErr)
	require.False(t, errs.IsFatal())

	errs.Add(fatalErr)
	require.True(t, errs.IsFatal())
}

func TestErrorsSetLevel(t *testing.T) {
	SetConfig(DefaultConfig())

	errs := New()
	require.Equal(t, Level(0), errs.Level)

	errs.SetLevel(WARN)
	require.Equal(t, WARN, errs.Level)
}

func TestErrorsString(t *testing.T) {
	SetConfig(DefaultConfig())

	errs := New()
	errs.Add(debugErr)
	errs.Add(errorErr)
	s := errs.String()
	require.NotEmpty(t, s)
}

func TestErrorsStack(t *testing.T) {
	SetConfig(DefaultConfig())

	errs := New()
	errs2 := New()

	errs.Stack(errs2)
	require.Equal(t, 0, len(errs.Errors))

	errs2.Add(warnErr)
	errs2.Add(errorErr)
	errs.Stack(errs2)
	require.Equal(t, 2, len(errs.Errors))
	require.Equal(t, WARN, errs.Errors[0].Level)
	require.Equal(t, ERROR, errs.Errors[1].Level)

	errs = New()
	errs.Add(debugErr)
	errs.Add(infoErr)

	errs.Stack(errs2)
	require.Equal(t, 4, len(errs.Errors))
	require.Equal(t, WARN, errs.Errors[0].Level)
	require.Equal(t, ERROR, errs.Errors[1].Level)
	require.Equal(t, DEBUG, errs.Errors[2].Level)
	require.Equal(t, INFO, errs.Errors[3].Level)
}

func TestErrorsAppend(t *testing.T) {
	SetConfig(DefaultConfig())

	errs := New()
	errs2 := New()
	errs.Append(errs2)
	require.Equal(t, 0, len(errs.Errors))

	errs2.Add(warnErr)
	errs2.Add(errorErr)
	errs.Append(errs2)
	require.Equal(t, 2, len(errs.Errors))
	require.Equal(t, WARN, errs.Errors[0].Level)
	require.Equal(t, ERROR, errs.Errors[1].Level)

	errs = New()
	errs.Add(debugErr)
	errs.Add(infoErr)
	errs.Append(errs2)
	require.Equal(t, DEBUG, errs.Errors[0].Level)
	require.Equal(t, INFO, errs.Errors[1].Level)
	require.Equal(t, WARN, errs.Errors[2].Level)
	require.Equal(t, ERROR, errs.Errors[3].Level)
}

func TestErrorsToArray(t *testing.T) {
	c := NewConfig()
	c.LoggingLevel = DEBUG
	SetConfig(c)

	errs := New()
	a := errs.ToArray()
	require.Equal(t, 0, len(a))

	errs.Add(debugErr)
	errs.Add(infoErr)
	errs.Add(warnErr)
	errs.Add(errorErr)
	errs.Add(fatalErr)

	a = errs.ToArray()
	require.Equal(t, 5, len(a))
	require.True(t, strings.Contains(a[0], "debug"))
	require.True(t, strings.Contains(a[4], "fatal"))
}

func TestErrorsError(t *testing.T) {
	SetConfig(DefaultConfig())

	errs := New()
	errs.Add(debugErr)
	errs.Add(infoErr)
	errs.Add(warnErr)
	errs.Add(errorErr)
	errs.Add(fatalErr)

	s := errs.Error()
	require.NotEmpty(t, s)
	require.NotContains(t, s, "debug")
	require.Contains(t, s, "info")
	require.Contains(t, s, "fatal")
}

func TestErrorsToLogArray(t *testing.T) {
	SetConfig(DefaultConfig())

	errs := New()
	a := errs.ToLogArray()
	require.Equal(t, 0, len(a))

	errs.Add(debugErr)
	errs.Add(infoErr)
	errs.Add(warnErr)
	errs.Add(errorErr)
	errs.Add(fatalErr)

	a = errs.ToLogArray()
	require.Equal(t, 5, len(a))
	require.Contains(t, a[0], "debug")
	require.Contains(t, a[4], "fatal")
}

func TestErrorsLog(t *testing.T) {
	SetConfig(DefaultConfig())
	var buf bytes.Buffer
	log.SetOutput(&buf)

	errs := New()
	errs.Log()
	require.Empty(t, buf)

	errs.Add(debugErr)
	errs.Add(infoErr)

	errs.Log()
	require.NotEmpty(t, buf)
	require.Contains(t, buf.String(), "debug")
	require.Contains(t, buf.String(), "info")
}

func TestErrorsClear(t *testing.T) {
	SetConfig(DefaultConfig())

	var errs Errors
	errs.Add(errorErr)
	errs.Add(fatalErr)
	require.Equal(t, 2, len(errs.Errors))

	errs.Clear()
	require.Equal(t, 0, len(errs.Errors))
}
