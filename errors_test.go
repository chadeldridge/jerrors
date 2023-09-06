package jerrors

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrors(t *testing.T) {
	c := DefaultConfig()
	c.SetConfig()

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

func TestErrorsCheck(t *testing.T) {
	c := DefaultConfig()
	c.SetConfig()

	errs := New()
	errs.Add(debugErr)
	errCount, hasErrors := errs.Check()
	require.Equal(t, 1, errCount)
	require.True(t, hasErrors)
}

func TestErrorsIsError(t *testing.T) {
	c := DefaultConfig()
	c.SetConfig()

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
	c := DefaultConfig()
	c.SetConfig()

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

func TestErrorsString(t *testing.T) {
	c := DefaultConfig()
	c.SetConfig()

	errs := New()
	errs.Add(debugErr)
	errs.Add(errorErr)
	s := errs.String()
	require.NotEmpty(t, s)
}

func TestErrorsStack(t *testing.T) {
	c := DefaultConfig()
	c.SetConfig()

	errs := New()
	errs.Add(debugErr)
	errs.Add(infoErr)

	errs2 := New()
	errs2.Add(warnErr)
	errs2.Add(errorErr)

	errs.Stack(errs2)
	require.Equal(t, 4, len(errs.Errors))
	require.Equal(t, WARN, errs.Errors[0].Level)
	require.Equal(t, ERROR, errs.Errors[1].Level)
	require.Equal(t, DEBUG, errs.Errors[2].Level)
	require.Equal(t, INFO, errs.Errors[3].Level)
}

func TestErrorsAppend(t *testing.T) {
	c := DefaultConfig()
	c.SetConfig()

	errs := New()
	errs.Add(debugErr)
	errs.Add(infoErr)

	errs2 := New()
	errs2.Add(warnErr)
	errs2.Add(errorErr)

	errs.Append(errs2)
	require.Equal(t, DEBUG, errs.Errors[0].Level)
	require.Equal(t, INFO, errs.Errors[1].Level)
	require.Equal(t, WARN, errs.Errors[2].Level)
	require.Equal(t, ERROR, errs.Errors[3].Level)
}

func TestErrorsToArray(t *testing.T) {
	c := NewConfig()
	c.LoggingLevel = DEBUG
	c.SetConfig()

	errs := New()
	errs.Add(debugErr)
	errs.Add(infoErr)
	errs.Add(warnErr)
	errs.Add(errorErr)
	errs.Add(fatalErr)

	a := errs.ToArray()
	require.Equal(t, 5, len(a))
	require.True(t, strings.Contains(a[0], "debug"))
	require.True(t, strings.Contains(a[4], "fatal"))
}

func TestErrorsClear(t *testing.T) {
	c := DefaultConfig()
	c.SetConfig()

	var errs Errors
	errs.Add(errorErr)
	errs.Add(fatalErr)
	require.Equal(t, 2, len(errs.Errors))

	errs.Clear()
	require.Equal(t, 0, len(errs.Errors))
}
