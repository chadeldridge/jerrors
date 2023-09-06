package jerrors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func testErrorMarshalJSON(t *testing.T, level Level) {
	expect := []byte(fmt.Sprintf("\"%s\"", level.String()))
	b, err := level.MarshalJSON()
	require.Nil(t, err, fmt.Sprintf("got error for MarshalJSON(): %s", err))
	require.NotEmpty(t, b)
	require.Equal(t, expect, b)
}

func testErrorUnmarshalJSON(t *testing.T, expect Level, j string) {
	var l Level
	err := l.UnmarshalJSON([]byte(j))
	require.Nil(t, err, fmt.Sprintf("got error for UnmarshalJSON(): %s", err))
	require.NotEqual(t, 0, l)
	require.Equal(t, expect, l)
}

func TestLevelsDebug(t *testing.T) {
	require.False(t, DEBUG.NotDebug())
	require.True(t, INFO.NotDebug())
	require.True(t, WARN.NotDebug())
	require.True(t, ERROR.NotDebug())
	require.True(t, FATAL.NotDebug())
}

func TestLevelsIsError(t *testing.T) {
	require.False(t, DEBUG.IsError())
	require.False(t, INFO.IsError())
	require.False(t, WARN.IsError())
	require.True(t, ERROR.IsError())
	require.True(t, FATAL.IsError())
}

func TestLevelsIsFatal(t *testing.T) {
	require.False(t, DEBUG.IsFatal())
	require.False(t, INFO.IsFatal())
	require.False(t, WARN.IsFatal())
	require.False(t, ERROR.IsFatal())
	require.True(t, FATAL.IsFatal())
}

// TestLevelsString also tests Error() since String() calls Error()
func TestLevelsString(t *testing.T) {
	require.Equal(t, "debug", DEBUG.String())
	require.Equal(t, "info", INFO.String())
	require.Equal(t, "warn", WARN.String())
	require.Equal(t, "error", ERROR.String())
	require.Equal(t, "fatal", FATAL.String())
}

func TestLevelsMarshalJSON(t *testing.T) {
	testErrorMarshalJSON(t, DEBUG)
	testErrorMarshalJSON(t, INFO)
	testErrorMarshalJSON(t, WARN)
	testErrorMarshalJSON(t, ERROR)
	testErrorMarshalJSON(t, FATAL)
}

func TestLevelsUnmarshalJSON(t *testing.T) {
	testErrorUnmarshalJSON(t, DEBUG, `"debug"`)
	testErrorUnmarshalJSON(t, INFO, `"info"`)
	testErrorUnmarshalJSON(t, WARN, `"warn"`)
	testErrorUnmarshalJSON(t, ERROR, `"error"`)
	testErrorUnmarshalJSON(t, FATAL, `"fatal"`)
}
