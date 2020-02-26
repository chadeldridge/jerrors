# jerrors JSON Error Framework
A simple golang Error module for creating errors and lists of errors in a JSON format.


## Contents
- [jerrors JSON Error Framework](#jerrors-json-error-framework)
  - [Contents](#contents)
  - [Quick Start](#quick-start)
  - [Errors](#errors)
  - [Creating New Error](#creating-new-errors)
    - [Error To String](#error-to-string)
    - [Accessing Metadata](#accessing-metadata)
    - [Checking Error Levels](#checking-error-levels)
    - [Error Logging](#error-logging)

## Installation
To install jerrors you must first has [Go](https://golang.org/) installed and setup.
1. Install jerrors module.
```ssh
$ go get -u github.com/chadeldridge/jerrors
```
2. Import jerrors in your code:
```go
import "github.com/chadeldridge/jerrors"
```

## Quick Start
```go
package main

import (
	"fmt"

	"github.com/chadeldridge/jerrors"
)

func main() {
	// Simple error message.
	err := jerrors.New(jerrors.DEBUG, "a Debug level error message", "type", "test",
		"user", "testuser")

	// Print DEBUG or higher level errors when Log is called. Default is INFO.
	jerrors.SetLogLevel(jerrors.DEBUG)
	err.Log()

	jerrors.SetLogLevel(jerrors.ERROR)

	// Error List
	var l jerrors.List
	// l.Level is set to ERROR
	l.Add(err)
	// l.Level remains ERROR
	l.Add(jerrors.New(jerrors.ERROR, "an Error level error message", "type", "test",
		"app", "testapp1", "user", "testuser"))
	// l.Level is set to FATAL
	l.Add(jerrors.New(jerrors.FATAL, "a Fatal error message", "type", "test", "app",
		"testapp1", "user", "testuser"))

	if has, count := l.Check(); has {
		l.Add(jerrors.New(jerrors.DEBUG,
			fmt.Sprintf("error list contained %v errors", count)))

		if l.Level == jerrors.FATAL {
			l.Fatal()
		} else {
			l.Log()
		}
	}
}
```
In the below output code the DEBUG level errors in the jerrors.List were not logged because we set the Log Level to ERROR.
```
{"time":"2020-02-26T13:11:40.038906297-05:00","level":"debug","message":"a Debug level error message","metadata":{"caller":"runtime.main{203}-\u003emain.main{11}","type":"test","user":"testuser"}}
{"time":"2020-02-26T13:11:40.039443772-05:00","level":"error","message":"an Error level error message","metadata":{"app":"testapp1","caller":"runtime.main{203}-\u003emain.main{25}","type":"test","user":"testuser"}}
{"time":"2020-02-26T13:11:40.039478791-05:00","level":"fatal","message":"a Fatal error message","metadata":{"app":"testapp1","caller":"runtime.main{203}-\u003emain.main{28}","type":"test","user":"testuser"}}
```

## Errors
```go
type Error struct {
	Time     time.Time         `json:"time,omitempty"`
	Level    Level             `json:"level,omitempty"`
	Message  string            `json:"message"`
	Metadata map[string]string `json:"metadata,omitempty"`
}
```

### Creating New Errors
Error with no added Metadata. By default jerrors always adds the caller func to Metadata.
```go
err := jerrors.New(jerrors.ERROR, "simple error message")
```

Error with extra Metadata. Metadata is converted to key:value pairs.
```go
err := jerrors.New(jerrors.ERROR, "simple error message", "add", "metadata",
	"in", "key", "value", "pairs")
```

Formatted error messages.
```go
e := funcReturnsGoError()
err := jerrors.New(jerrors.ERROR, fmt.Sprintf("error: %v", e))
```

### Error To string
jerrors.Error implements the Error interface and can be used in all the same ways.
```go
err := jerrors.New(jerrors.ERROR, "simple error message")
jsonString := err.Error()
fmt.Printf("Error returned: %v", err)
```

### Accessing Metadata
The Metadata of an Error can be accessed directly as a map[string]string.
```go
err := jerrors.New(jerrors.ERROR, "simple error message", "user", "bob")
if err.Metadata["user"] == "bob" {
	fmt.Printf("Why %s?! %v", err.Metadata["user"], err)
}
```

### Checking Error Levels

#### Direct Comparison
```go
err := jerrors.New(jerrors.ERROR, "simple error message")
if err.Level == jerrors.ERROR {
	err.Log()
}
```

#### IsError
IsError returns true if Error.Level is >= ERROR.
```go
if err.IsError {
	err.Log()
}
```

#### IsFatal
IsFatal returns true only if Error.Level == FATAL.
```go
if err.IsFatal {
	err.Fatal()
}
```

### Error Logging

#### Error.Log()
An Error can be logged directly by calling Error.Log(). By default Log() will print to STDERR.
```go
err := jerrors.New(jerrors.ERROR, "simple error message")
err.Log()
```
Output
```
{"time":"2020-02-26T13:11:40.038906297-05:00","level":"error","message":"simple error message","metadata":{"caller":"runtime.main{203}-\u003emain.main{11}"}}
```

#### Error.Fatal()
Fatal sets the Error.Level to FATAL before logging the error and exiting with a status of 1.
```go
err := jerrors.New(jerrors.ERROR, "simple error message")
err.Fatal()
```
Output
```
{"time":"2020-02-26T13:11:40.038906297-05:00","level":"fatal","message":"simple error message","metadata":{"caller":"runtime.main{203}-\u003emain.main{11}"}}
```
