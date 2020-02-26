# jerrors
A simple golang Error module for creating errors and lists of errors in a JSON format.

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
