package jerrors

import (
	"testing"
)

func errorIfListIsError(t *testing.T, errs List) {
	if b := errs.IsError(); b {
		t.Errorf("list IsError returned (%v), level is (%v)\n", b, errs.Level)
	}
}

func errorIfNotListIsError(t *testing.T, errs List) {
	if b := errs.IsError(); !b {
		t.Errorf("list IsError returned (%v), level is (%v)\n", b, errs.Level)
	}
}

func errorIfListIsFatal(t *testing.T, errs List) {
	if b := errs.IsFatal(); b {
		t.Errorf("list IsFatal returned (%v), level is (%v)\n", b, errs.Level)
	}
}

func TestListBasic(t *testing.T) {
	var errs List
	errs.Add(New(WARN, "test message 1", "test", "1", "user", "test1"))
	if len(errs.Errors) == 0 {
		t.Error("list is empty\n")
	}
	if errs.Level != WARN {
		t.Errorf("list level (%v), expected (%v)\n", errs.Level, WARN)
	}

	// errs.Level should remain WARN
	errs.Add(New(DEBUG, "test message 2", "test", "2", "user", "test1"))
	if len(errs.Errors) < 2 {
		t.Errorf("list error not added: length (%v)\n", len(errs.Errors))
	}
	if errs.Level != WARN {
		t.Errorf("list level (%v), expected (%v)\n", errs.Level, WARN)
	}

	// errs.Level should change to ERROR
	errs.Add(New(ERROR, "test message 3", "test", "3", "user", "test1"))
	if len(errs.Errors) < 3 {
		t.Errorf("list error not added: length (%v)\n", len(errs.Errors))
	}
	if errs.Level != ERROR {
		t.Errorf("list level (%v), expected (%v)\n", errs.Level, ERROR)
	}

	errs = List{}
	if hasErrors, count := errs.Check(); !hasErrors {
		t.Errorf("list Check returned (%v) with (%v) errors, expected (true) and (3)\n", hasErrors, count)
	}
}

func TestListIsError(t *testing.T) {
	var errs List
	errs.Add(New(DEBUG, "test message 1", "test", "1", "user", "test1"))
	errorIfListIsError(t, errs)
	errs.SetLevel(INFO)
	errorIfListIsError(t, errs)
	errs.SetLevel(WARN)
	errorIfListIsError(t, errs)
	errs.SetLevel(ERROR)
	errorIfNotListIsError(t, errs)
	errs.SetLevel(FATAL)
	errorIfNotListIsError(t, errs)
}

func TestListIsFatal(t *testing.T) {
	var errs List
	errs.Add(New(DEBUG, "test message 1", "test", "1", "user", "test1"))
	errorIfListIsFatal(t, errs)
	errs.SetLevel(INFO)
	errorIfListIsFatal(t, errs)
	errs.SetLevel(WARN)
	errorIfListIsFatal(t, errs)
	errs.SetLevel(ERROR)
	errorIfListIsFatal(t, errs)
	errs.SetLevel(FATAL)
	if b := errs.IsFatal(); !b {
		t.Errorf("list IsFatal returned (%v), level is (%v)\n", b, errs.Level)
	}
}

func TestListStack(t *testing.T) {
	var errs List
	errs.Add(New(DEBUG, "test message 1", "test", "1", "user", "test1"))
	errs.Add(New(DEBUG, "test message 2", "test", "2", "user", "test1"))
	var errs2 List
	errs2.Add(New(DEBUG, "test message 3", "test", "3", "user", "test1"))
	errs2.Add(New(DEBUG, "test message 4", "test", "4", "user", "test1"))
	errs.Stack(errs2)
	if len(errs.Errors) < 4 {
		t.Errorf("list Stack error count (%v), expected (4)\n", len(errs.Errors))
	}
	if errs.Errors[0].Metadata["test"] != "3" {
		t.Errorf("list Stack first error is (%v), expected (3)\n", errs.Errors[0].Metadata["test"])
	}
	if errs.Errors[3].Metadata["test"] != "2" {
		t.Errorf("list Stack last error is (%v), expected (4)\n", errs.Errors[3].Metadata["test"])
	}
}

func TestListAppend(t *testing.T) {
	var errs List
	errs.Add(New(DEBUG, "test message 1", "test", "1", "user", "test1"))
	errs.Add(New(DEBUG, "test message 2", "test", "2", "user", "test1"))
	var errs2 List
	errs2.Add(New(DEBUG, "test message 3", "test", "3", "user", "test1"))
	errs2.Add(New(DEBUG, "test message 4", "test", "4", "user", "test1"))
	errs.Append(errs2)
	if len(errs.Errors) < 4 {
		t.Errorf("list Append error count (%v), expected (4)\n", len(errs.Errors))
	}
	if errs.Errors[0].Metadata["test"] != "1" {
		t.Errorf("list Append first error is (%v), expected (1)\n", errs.Errors[0].Metadata["test"])
	}
	if errs.Errors[3].Metadata["test"] != "4" {
		t.Errorf("list Append last error is (%v), expected (4)\n", errs.Errors[3].Metadata["test"])
	}
}

func TestListClear(t *testing.T) {
	var errs List
	errs.Add(New(DEBUG, "test message 1", "test", "1", "user", "test1"))
	errs.Add(New(DEBUG, "test message 2", "test", "2", "user", "test1"))
	errs.Clear()
	if len(errs.Errors) != 0 {
		t.Errorf("list Clear error count (%v), expected (0)\n", len(errs.Errors))
	}
	if errs.Level != 0 {
		t.Errorf("list Clear level (%v), expected (0)\n", errs.Level)
	}
}
