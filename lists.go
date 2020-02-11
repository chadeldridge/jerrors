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
func (e *List) Check() (bool, int) {
	if l := len(e.Errors); l > 0 {
		return true, l
	}
	return false, 0
}

// IsError returns true for anything above WARN
func (e *List) IsError() bool {
	return e.Level.IsError()
}

// IsFatal returns true for anything above ERROR
func (e *List) IsFatal() bool {
	return e.Level.IsFatal()
}

// SetLevel the Level
func (e *List) SetLevel(level Level) {
	e.Level = level
}

// Add the error to the method's List.
func (e *List) Add(err Error) {
	if err.Level != 0 && err.Message != nil {
		if err.Level > e.Level {
			e.Level = err.Level
		}
		//e.Errors = append([]Error{err}, e.Errors...)
		e.Errors = append(e.Errors, err)
	}
}

// Stack adds the args List to the top of the method's List
func (e *List) Stack(errs List) {
	if len(errs.Errors) > 0 {
		// If e is currently empty then overwrite it.
		if len(e.Errors) == 0 {
			e.Level = errs.Level
			e.Errors = errs.Errors
			return
		}
		if errs.Level > e.Level {
			e.Level = errs.Level
		}
		e.Errors = append(errs.Errors, e.Errors...)
	}
}

// Append the arg List to the method's List.
func (e *List) Append(errs List) {
	if len(errs.Errors) > 0 {
		// If e is currently empty then overwrite it.
		if len(e.Errors) == 0 {
			e.Level = errs.Level
			e.Errors = errs.Errors
			return
		}
		if errs.Level > e.Level {
			e.Level = errs.Level
		}
		e.Errors = append(e.Errors, errs.Errors...)
	}
}

// Clear the List. Level = 0; Errors = nil.
func (e *List) Clear() {
	e.Level = 0
	e.Errors = nil
}

// MarshalArray all errors in List into a single json string.
func (e *List) MarshalArray() []string {
	switch l := len(e.Errors); l {
	case 0:
		msgs := make([]string, 0, 0)
		return msgs
	default:
		msgs := make([]string, 0, l)
		for _, err := range e.Errors {
			if err.Level >= loggingLevel {
				msgs = append(msgs, err.String())
			}
		}
		return msgs
	}
}

// Marshal all errors in List into a single json string.
func (e *List) Marshal() string {
	msgs := e.MarshalArray()
	j, _ := json.Marshal(msgs)
	return string(j)
}

// Log each message then clear the List.
func (e *List) Log() {
	if len(e.Errors) == 0 {
		return
	}

	msgs := e.MarshalArray()
	for _, err := range msgs {
		log.Println(err)
	}

	e.Clear()
}

// Fatal converts all errors to a single error and runs Fatal to print error and exit(1).
func (e *List) Fatal() {
	if e.IsFatal() {
		msgs := e.MarshalArray()
		log.Fatal(strings.Join(msgs, "\n"))
	}
}
