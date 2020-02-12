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
	Errors []Error
	Level  Level
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

// Stack adds the args' List to the top of the method's List
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

// MarshalArray all errors in List into a single json string.
func (l *List) MarshalArray() []string {
	switch i := len(l.Errors); i {
	case 0:
		msgs := make([]string, 0, 0)
		return msgs
	default:
		msgs := make([]string, 0, i)
		for _, err := range l.Errors {
			if err.Level >= loggingLevel {
				msgs = append(msgs, err.String())
			}
		}
		return msgs
	}
}

// Marshal all errors in List into a single json string.
func (l *List) Marshal() string {
	msgs := l.MarshalArray()
	j, _ := json.Marshal(msgs)
	return string(j)
}

// Log all messages in the List.
func (l *List) Log() {
	if len(l.Errors) == 0 {
		return
	}

	msgs := l.MarshalArray()
	log.Print(strings.Join(msgs, "\n"))
}

// Fatal converts all errors to a single error and runs Fatal to print error and exit(1).
func (l *List) Fatal() {
	if l.IsFatal() {
		msgs := l.MarshalArray()
		log.Fatal(strings.Join(msgs, "\n"))
	}
}
