package jerrors

import (
	"encoding/json"
	"log"
	"strings"
)

/*
List is a slice of Errors and a Level indicating the most
critical error level added so far.
*/
type List struct {
	Errors []Error `json:"errors"`
	Level  Level   `json:"level"`
}

// Add an error to the method's List.
func (l *List) Add(err Error) {
	if err.Level != 0 && err.Message != "" {
		if err.Level > l.Level {
			l.Level = err.Level
		}
		//l.Errors = append([]Error{err}, l.Errors...)
		l.Errors = append(l.Errors, err)
	}
}

// Check return true if errors exist and the number of errors present.
func (l *List) Check() (bool, int) {
	if l := len(l.Errors); l > 0 {
		return true, l
	}
	return false, 0
}

// IsError returns true for anything above WARN
func (l *List) IsError() bool {
	return l.Level.IsError()
}

// IsFatal returns true for anything above ERROR
func (l *List) IsFatal() bool {
	return l.Level.IsFatal()
}

// SetLevel the Level
func (l *List) SetLevel(level Level) {
	l.Level = level
}

// Stack adds the args' List to top of the method's List
func (l *List) Stack(errs List) {
	if len(errs.Errors) > 0 {
		// If l is currently empty then overwrite it.
		if len(l.Errors) == 0 {
			l.Level = errs.Level
			l.Errors = errs.Errors
			return
		}
		if errs.Level > l.Level {
			l.Level = errs.Level
		}
		l.Errors = append(errs.Errors, l.Errors...)
	}
}

// Append the arg List to the method's List.
func (l *List) Append(errs List) {
	if len(errs.Errors) > 0 {
		// If l is currently empty then overwrite it.
		if len(l.Errors) == 0 {
			l.Level = errs.Level
			l.Errors = errs.Errors
			return
		}
		if errs.Level > l.Level {
			l.Level = errs.Level
		}
		l.Errors = append(l.Errors, errs.Errors...)
	}
}

// Clear the List. Level = 0; Errors = nil.
func (l *List) Clear() {
	l.Level = 0
	l.Errors = nil
}

func (l *List) toArray(enforceLogLevel bool) []Error {
	switch i := len(l.Errors); i {
	case 0:
		msgs := make([]Error, 0, 0)
		return msgs
	default:
		msgs := []Error{}
		for _, err := range l.Errors {
			if enforceLogLevel {
				if err.Level >= loggingLevel {
					msgs = append(msgs, err)
				}
			} else {
				msgs = append(msgs, err)
			}
		}
		return msgs
	}
}

// Error returns all errors in List as a single json string. Returns empty string if failed.
func (l *List) Error() string {
	msgs := l.toArray(false)
	j, err := json.Marshal(msgs)
	if err != nil {
		return ""
	}
	return string(j)
}

// String is an alternate method name for List.Error().
func (l *List) String() string {
	return l.Error()
}

/*
// MarshalJSON converts a List to json.
func (l *List) MarshalJSON() ([]byte, error) {
		msgs := l.toArray(false)
		j, err := json.Marshal(msgs)
		if err != nil {
			return nil, err
		}
		return j, nil
}

// UnmarshalJSON converts json to a List.
func (l *List) UnmarshalJSON(b []byte) error {
		var a []Error
		err := json.Unmarshal(b, &a)
		if err != nil {
			return err
		}

		l.Errors = a
		for _, e := range a {
			if e.Level > l.Level {
				l.Level = e.Level
			}
		}

		return nil
}
*/

// ToArray returns an array of all marshalled errors in List.
func (l *List) ToArray() []string {
	var a []string
	errs := l.toArray(false)
	for _, e := range errs {
		a = append(a, e.Error())
	}
	return a
}

// ToLogArray returns an array of marshalled errors in List.
// Omits errors below the current logLevel.
func (l *List) ToLogArray() []string {
	var a []string
	errs := l.toArray(true)
	for _, e := range errs {
		a = append(a, e.Error())
	}
	return a
}

// Log all messages in the List.
func (l *List) Log() {
	if len(l.Errors) == 0 {
		return
	}

	msgs := l.ToLogArray()
	log.Print(strings.Join(msgs, "\n"))
}

// Fatal converts all errors to a single error and runs Fatal to print error and exit(1).
func (l *List) Fatal() {
	msgs := l.ToLogArray()
	log.Fatal(strings.Join(msgs, "\n"))
}
