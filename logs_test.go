package jerrors

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogsSetLogOutput(t *testing.T) {
	var buf bytes.Buffer
	SetLogOutput(&buf)

	errs := New()
	errs.Add(debugErr)

	errs.Log()
	require.NotEmpty(t, buf)
	require.Contains(t, buf.String(), "debug")

	buf.Reset()
	SetLogOutput(nil)
	errs.Log()
	require.Len(t, buf.String(), 0)
}
