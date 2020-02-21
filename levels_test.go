package jerrors

import "testing"

func errorIfNotLevelNotDebug(t *testing.T, l Level) {
	if b := l.NotDebug(); !b {
		t.Errorf("level is (%v) but NotDebug returned (%v):\n", l, b)
	}
}

func errorIfLevelIsError(t *testing.T, l Level) {
	if b := l.IsError(); b {
		t.Errorf("level is (%v) but IsError returned (%v):\n", l, b)
	}
}

func errorIfNotLevelIsError(t *testing.T, l Level) {
	if b := l.IsError(); !b {
		t.Errorf("level is (%v) but IsError returned (%v):\n", l, b)
	}
}

func errorIfLevelIsFatal(t *testing.T, l Level) {
	if b := l.IsFatal(); b {
		t.Errorf("level is (%v) but IsFatal returned (%v):\n", l, b)
	}
}

func TestLevelNotDebug(t *testing.T) {
	if b := DEBUG.NotDebug(); b {
		t.Errorf("level is (%v) but NotDebug returned (%v):\n", DEBUG, b)
	}
	errorIfNotLevelNotDebug(t, INFO)
	errorIfNotLevelNotDebug(t, WARN)
	errorIfNotLevelNotDebug(t, ERROR)
	errorIfNotLevelNotDebug(t, FATAL)
}

func TestLevelIsError(t *testing.T) {
	errorIfLevelIsError(t, DEBUG)
	errorIfLevelIsError(t, INFO)
	errorIfLevelIsError(t, WARN)
	errorIfNotLevelIsError(t, ERROR)
	errorIfNotLevelIsError(t, FATAL)
}

func TestLevelIsFatal(t *testing.T) {
	errorIfLevelIsFatal(t, DEBUG)
	errorIfLevelIsFatal(t, INFO)
	errorIfLevelIsFatal(t, WARN)
	errorIfLevelIsFatal(t, ERROR)
	if b := FATAL.IsFatal(); !b {
		t.Errorf("level is (%v) but NotDebug returned (%v):\n", DEBUG, b)
	}
}

func TestLevelString(t *testing.T) {
	s := ERROR.String()
	if s == "" {
		t.Error("level is (error) but String returned empty\n")
		return
	}
	if s != "error" {
		t.Errorf("level is (error) but String returned (%v):\n", s)
	}
}

func TestLevelMarshalJSON(t *testing.T) {
	b, err := ERROR.MarshalJSON()
	if err != nil {
		t.Errorf("level MarshalJSON returned error: %v\n", err)
		return
	}
	if len(b) == 0 {
		t.Error("level is (\"error\") but MarshalJSON returned empty\n")
		return
	}
	s := string(b)
	if s != `"error"` {
		t.Errorf("level is (\"error\") but String returned (%v):\n", s)
	}
}

func TestLevelUnmarshalJSON(t *testing.T) {
	l := ERROR
	b, err := l.MarshalJSON()
	if err != nil {
		t.Errorf("level MarshalJSON returned error: %v\n", err)
		return
	}
	if len(b) == 0 {
		t.Error("level is (\"error\") but MarshalJSON returned empty\n")
		return
	}

	// Clear level so we know if it gets updated during unmarshal.
	l = 0
	err = l.UnmarshalJSON(b)
	if err != nil {
		t.Errorf("level UnmarshalJSON returned error: %v\n", err)
		return
	}
	if l != ERROR {
		t.Errorf("level UnmarshalJSON is (%v), expected (%v):\n", l, ERROR)
	}
}
