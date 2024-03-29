package jerrors

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	want := Config{
		LogLevel:      true,
		LogTime:       true,
		LoggingLevel:  INFO,
		LogCaller:     false,
		CallerDepth:   2,
		CallersToShow: 2,
	}

	got := NewConfig()
	require.Equal(t, want, got)
}

func TestNewConfig(t *testing.T) {
	want := DefaultConfig()
	got := NewConfig()

	require.Equal(t, want, got)
}

func TestSetConfig(t *testing.T) {
	want := DefaultConfig()
	want.LogTime = false
	want.LoggingLevel = ERROR
	SetConfig(want)

	got := GetConfig()

	require.Equal(t, want.LogLevel, got.LogLevel)
	require.Equal(t, want.LogTime, got.LogTime)
	require.Equal(t, want.LoggingLevel, got.LoggingLevel)
}
