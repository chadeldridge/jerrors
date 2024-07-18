package jerrors

import (
	"encoding/json"
	"log"
	"strings"
)

var ErrNoErrorFound = NewError(0, "No error found")

// Errors is a slice of Errors and with Level showing the highest error level added.
type Errors struct {
	Errors []Error `json:"errors"`
	Level  Level   `json:"level"`
}

func New() Errors {
	return Errors{
		Errors: []Error{},
		Level:  0,
	}
}

// New creates a new Error and adds it to the List.
func (e *Errors) NewError(level Level, msg string, args ...interface{}) {
	e.Add(NewError(level, msg, args...))
}

// Add an error to the method's List.
func (e *Errors) Add(err Error) {
	if err.Level > e.Level {
		e.Level = err.Level
	}

	e.Errors = append(e.Errors, err)
}

func (e *Errors) Remove(error Error) bool {
	for i, err := range e.Errors {
		if err.Equal(error) {
			e.Errors = append(e.Errors[:i], e.Errors[i+1:]...)
			e.UpdateLevel()
			return true
		}
	}

	return false
}

// UpdateLevel sets Errors.Level the the highest one in the Errors list and returns Errors.Level.
func (e *Errors) UpdateLevel() Level {
	var l Level
	for _, err := range e.Errors {
		if err.Level > l {
			l = err.Level
		}
	}

	e.Level = l
	return e.Level
}

// Check if Errors is not empty and return the number of Errors.
func (e *Errors) Check() (int, bool) {
	if l := len(e.Errors); l > 0 {
		return l, true
	}

	return 0, false
}

// First returns the first (e.Errors[0]) Error in the Errors List. Reutnrs ErrNoErrorFound if
// Errors.Errors is empt. ErrNoErrorFound = NewError(0, "No errors found")
func (e *Errors) First() Error {
	if len(e.Errors) == 0 {
		return ErrNoErrorFound
	}

	return e.Errors[0]
}

// Last returns the last Error in the Errors List. Returns ErrNoErrorFound if Errors.Errors is empty.
func (e *Errors) Last() Error {
	if len(e.Errors) == 0 {
		return ErrNoErrorFound
	}

	return e.Errors[len(e.Errors)-1]
}

// IsEmpty checks to see if Errors.Errors is empty.
func (e *Errors) IsEmpty() bool { return len(e.Errors) == 0 }

// IsError returns true for anything above WARN
func (e *Errors) IsError() bool { return e.Level.IsError() }

// IsFatal returns true for anything above ERROR
func (e *Errors) IsFatal() bool { return e.Level.IsFatal() }

// SetLevel overrides the Level of the List
func (e *Errors) SetLevel(level Level) { e.Level = level }

// String is an alternate method name for List.Error()
func (e *Errors) String() string { return e.Error() }

// Clear the List and Level
func (e *Errors) Clear() { *e = New() }

// Stack adds the arg List to top of the method's List
func (e *Errors) Stack(errs Errors) {
	if len(errs.Errors) == 0 {
		return
	}

	// If l is currently empty then overwrite it
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

// Append the arg List to the method's List.
func (e *Errors) Append(errs Errors) {
	if len(errs.Errors) == 0 {
		return
	}

	// If l is currently empty then overwrite it.
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

func (e *Errors) toArray(enforceLogLevel bool) []Error {
	switch i := len(e.Errors); i {
	case 0:
		return []Error{}
	default:
		msgs := []Error{}
		for _, err := range e.Errors {
			if !enforceLogLevel && err.Level < config.LoggingLevel {
				continue
			}

			msgs = append(msgs, err)
		}
		return msgs
	}
}

// Error returns all errors in List as a single json string. Returns empty string if failed.
func (e *Errors) Error() string {
	msgs := e.toArray(false)
	j, err := json.Marshal(msgs)
	if err != nil {
		return ""
	}

	return string(j)
}

/*
// MarshalJSON converts a List to json.
func (e *Errors) MarshalJSON() ([]byte, error) {
	msgs := e.toArray(false)
	j, err := json.Marshal(msgs)
	if err != nil {
		return nil, err
	}

	return j, nil
}

// UnmarshalJSON converts json to a List.
func (e *Errors) UnmarshalJSON(b []byte) error {
	var a []Error

	r := json.Unmarshal(b, &a)
	if r != nil {
		return r
	}

	e.Errors = a
	for _, err := range a {
		if err.Level > e.Level {
			e.Level = err.Level
		}
	}

	return nil
}
*/

func (e *Errors) Pretty() string {
	// msgs := e.toArray(false)
	j, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return ""
	}

	return string(j)
}

// ToArray returns an array of all marshalled errors in List.
func (e *Errors) ToArray() []string {
	var a []string

	errs := e.toArray(false)
	for _, e := range errs {
		a = append(a, e.Error())
	}

	return a
}

// ToLogArray returns an array of marshalled errors in List.
// Omits errors below the current logLevee.
func (e *Errors) ToLogArray() []string {
	var a []string

	errs := e.toArray(true)
	for _, e := range errs {
		a = append(a, e.Error())
	}

	return a
}

// Log all messages in the List.
func (e *Errors) Log() {
	if len(e.Errors) == 0 {
		return
	}

	msgs := e.ToLogArray()
	log.Print(strings.Join(msgs, "\n"))
}

// Fatal converts all errors to a single error and runs Fatal to print error and exit(1).
func (e *Errors) Fatal(msg string) {
	msgs := e.ToLogArray()
	log.Fatal(msg + strings.Join(msgs, "\n"))
}
