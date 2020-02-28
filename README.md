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
  - [Levels](#levels)
    - [Level Definition](#level-definition)
    - [Level Functions](#level-functions)
  - [Lists](#lists)
    - [Checking List Levels](#checking-list-levels)
    - [List Manipulation](#list-manipulation)
    - [List Conversions](#list-conversions)
    - [List Logging](#list-logging)
  - [Logs](#logs)
    - [Logging Options](#logging-options)
    - [Logging Level](#logging-level)
    - [Log Output](#log-output)

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
In the below output the DEBUG level errors in the jerrors.List were not logged because we set the Log Level to ERROR.
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
Minimal Error with just a Level and Message set. By default jerrors always adds the "caller" to Metadata.
```go
err := jerrors.New(jerrors.ERROR, "simple error message")
```
```
{"time":"2020-02-26T13:11:40.039443772-05:00","level":"error","message":"simple error message","metadata":{"caller":"runtime.main{203}-\u003emain.main{25}"}}
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
jerrors.Error implements the Golang Error interface and can be used in all the same ways.
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
 See [Levels](#levels) for details on jerrors.Level.

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
if err.IsError() {
    err.Log()
}
```

#### IsFatal
IsFatal returns true only if Error.Level == FATAL.
```go
if err.IsFatal() {
    err.Fatal()
}
```

### Error Logging

#### Error.Log()
An Error can be logged directly by calling Error.Log(). By default Log() will print to os.Stderr.
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

## Levels
The Levels enum is used to implement a standardization on error level hierarchy.

### Level Definition

```go
const (
	DEBUG Level = iota + 1
	INFO
	WARN
	ERROR
	FATAL
)
```

A Level of 0 is the nil value. Levels can be compared with normal mathimatical operators.
```go
l := jerrors.ERROR
l == jerrors.ERROR // true
l >  jerrors.DEBUG // true
l >= jerrors.FATAL // false
```

Levels string representation.
```go
	DEBUG: "debug",
	INFO:  "info",
	WARN:  "warn",
	ERROR: "error",
	FATAL: "fatal",
```

### Level Functions
```go

// StringToLevel returns the matched Level. "debug" = DEBUG
// level arg is NOT case sensitive. No match returns 0.
func StringToLevel(level string) Level

// String converts Level to a lowercase string. DEBUG = "debug", etc.
// Returns empty string if Level is 0.
func (l Level) String() string

// NotDebug returns true for all Levels except DEBUG.
func (l Level) NotDebug() bool

// IsError returns true if Level is ERROR or FATAL.
func (l Level) IsError() bool

// IsFatal returns true if Level is FATAL.
func (l Level) IsFatal() bool

// MarshalJSON converts Level to json.
func (l Level) MarshalJSON() ([]byte, error)

// UnmarshalJSON converts json to a Level.
func (l *Level) UnmarshalJSON(b []byte) error
```

## Lists

```go
type List struct {
	Errors []Error `json:"errors"`
	Level  Level   `json:"level"`
}
```
A List holds an array of Errors and a Level. The Level is the most critical Level of any error added to the List with List.Add.

```go
var l jerrors.List
l.Add(jerrors.New(jerrors.DEBUG, "some helpful message"))
fmt.Println(l.Level) // Prints: debug

l.Add(jerrors.New(jerrors.ERROR, "some error message"))
fmt.Println(l.Level) // Prints: error

l.Add(jerrors.New(jerrors.WARN, "some warning message"))
fmt.Println(l.Level) // Prints: error
```

### Checking List Levels

#### Check Function
List.Check looks to see if any Errors exist and returns a bool and the number of errors in List.Errors.
```go
var l List
l.Add(jerrors.New(jerrors.DEBUG, "some helpful message"))
l.Add(jerrors.New(jerrors.ERROR, "some error message"))
if hasErrors, count := l.Check(); hasErrors {
	fmt.Println(count) // Prints: 2
}
```

#### Comparing List.Error Directly
You can access the List Level directly and compaire just like with Error.
```go
if l.Level >= jerrors.ERROR {
	l.Log()
}
```

#### IsError
IsError returns true if List.Level is >= ERROR.
```go
if l.IsError() {
    l.Log()
}
```

#### IsFatal
IsFatal returns true only if List.Level == FATAL.
```go
if l.IsFatal() {
    l.Fatal()
}
```

### List Manipulation

#### Append Lists
You can append a new List to and existing List. Append works just like array Append adding the argument list to the bottom of the list calling the method.
```go
var l jerrors.List
l.Add(jerrors.New(jerrors.DEBUG, "some helpful message"))
l.Add(jerrors.New(jerrors.ERROR, "some error message"))

var l2 jerrors.List
l2.Add(jerrors.New(jerrors.WARN, "some warning message"))

l.Append(l2)
fmt.Println(l.Errors[0].Message) // Prints: some helpful message
```

#### Stack Lists
You can stack two lists together. Stack puts the aurgument List at the top of the List calling the method.
```go
var l jerrors.List
l.Add(jerrors.New(jerrors.DEBUG, "some helpful message"))
l.Add(jerrors.New(jerrors.ERROR, "some error message"))

var l2 jerrors.List
l2.Add(jerrors.New(jerrors.WARN, "some warning message"))

l.Stack(l2)
fmt.Println(l.Errors[0].Message) // Prints: some warning message
```

#### Clear List
Clear deletes all Errors in the List and sets Level to 0.
```go
var l jerrors.List
l.Add(jerrors.New(jerrors.DEBUG, "some helpful message"))
l.Add(jerrors.New(jerrors.ERROR, "some error message"))
hasErrors, count := l.Check() // true, 2
l.Clear()
hasErrors, count = l.Check() // false, 0
```

### List Conversions

#### Error and String
List implements the error interface and can be used the same as any other error. Error converts List.Errors into a single JSON string.
String is an alias of Error.
```go
var l jerrors.List
l.Add(jerrors.New(jerrors.DEBUG, "some helpful message"))
l.Add(jerrors.New(jerrors.ERROR, "some error message"))
fmt.Printf("List of errors:\n%v\n", l.Error())
```

Output:
```
List of errors:
[{"time":"2020-02-27T15:56:01.914821132-05:00","level":"debug","message":"some helpful message","metadata":{"caller":"runtime.main{203}-\u003emain.main{11}"}},{"time":"2020-02-27T15:56:01.914955631-05:00","level":"error","message":"some error message","metadata":{"caller":"runtime.main{203}-\u003emain.main{15}"}}]
```

#### JSON
You can use json.Marshall and json.Unmarshall to convert between List and JSON.
```go
var l List
l.Add(jerrors.New(jerrors.DEBUG, "some helpful message"))
j, _ := json.Marshall(l)
fmt.Println(string(j))

var l2 List
if err := json.Unmarshall(l, &l2); err != nil {
	fmt.Fatal(err)
}

l2.Log()
```

#### ToArray
Converts the List.Errors to an array of JSON strings. Ignores Log Level.
```go
var l jerrors.List
l.Add(jerrors.New(jerrors.DEBUG, "some helpful message"))
l.Add(jerrors.New(jerrors.ERROR, "some error message"))
a := l.ToArray()
fmt.Println(strings.Join(a, "\n"))
```
Output:
```
{"time":"2020-02-28T13:30:21.668092176-05:00","level":"debug","message":"some helpful message","metadata":{"caller":"runtime.main{203}-\u003emain.main{12}"}}
{"time":"2020-02-28T13:30:21.668119405-05:00","level":"error","message":"some error message","metadata":{"caller":"runtime.main{203}-\u003emain.main{13}"}}
```

#### ToLogArray
Converts the List.Errors to an array of JSON strings. Honors Log Level.
```go
var l jerrors.List
l.Add(jerrors.New(jerrors.DEBUG, "some helpful message"))
l.Add(jerrors.New(jerrors.ERROR, "some error message"))
a := l.ToLogArray()
fmt.Println(strings.Join(a, "\n"))
```
Output:
```
{"time":"2020-02-28T13:31:20.453088284-05:00","level":"error","message":"some error message","metadata":{"caller":"runtime.main{203}-\u003emain.main{13}"}}
```

### List Logging

#### List Log
Call log.Print to log the array of Errors. Any Error with a Level >= current Log Level will be logged. See [Logs](#logs) for more details.
```go
var l jerrors.List
l.Add(jerrors.New(jerrors.DEBUG, "some helpful message"))
l.Add(jerrors.New(jerrors.ERROR, "some error message"))
l.Log()
```
Output:
```
{"time":"2020-02-28T13:31:20.453088284-05:00","level":"error","message":"some error message","metadata":{"caller":"runtime.main{203}-\u003emain.main{13}"}}
```

#### List Fatal
Same as [List.Log()](#list-log) but uses log.Fatal to exit with a status of 1.
```go
var l jerrors.List
l.Add(jerrors.New(jerrors.DEBUG, "some helpful message"))
l.Add(jerrors.New(jerrors.ERROR, "some error message"))
l.Log()
```
Output:
```
{"time":"2020-02-28T13:31:20.453088284-05:00","level":"error","message":"some error message","metadata":{"caller":"runtime.main{203}-\u003emain.main{13}"}}
```

## Logs

### Logging Options
SetLogOptions takes in a map[string]bool of options. Current available options.
"caller" - Determines if "caller" should be included in Metadata. Other Metadata will still be logged if any exists. Defaults to true.
"level" - Determines if the Error Level should be included when logging an Error. Defaults to true.
"time" - Determines if Time should be included when logging an Error.

```go
ops := map[string]bool{"caller": false, "time": false}
jerrors.SetLogOptions(ops)

var l jerrors.List
l.Add(jerrors.New(jerrors.ERROR, "some error message", "type", "test"))
l.Log()
```
Output:
```
{"level":"error","message":"some error message","metadata":{"type":"test"}}
```

### Logging Level
The Logging Level determine what should get logged when calling Log() or Fatal on both Error and List. Log/Fatal will log all Errors with a Level >= the current Log Level.
You can use SetLogLevel to change the Logging Level. Defaults to INFO.

```go
var l jerrors.List
l.Add(jerrors.New(jerrors.DEBUG, "some helpful message"))
l.Add(jerrors.New(jerrors.ERROR, "some error message", "type", "test"))

fmt.Println("Default Log Level: INFO")
l.Log()

jerrors.SetLogLevel(jerrors.DEBUG)
fmt.Println("\nLog Level: DEBUG")
l.Log()
```

Output:
```
Default Log Level: INFO
{"time":"2020-02-28T14:39:47.067940286-05:00","level":"error","message":"some error message","metadata":{"caller":"runtime.main{203}-\u003emain.main{12}","type":"test"}}

Log Level: DEBUG
{"time":"2020-02-28T14:39:47.067920642-05:00","level":"debug","message":"some helpful message","metadata":{"caller":"runtime.main{203}-\u003emain.main{11}"}}
{"time":"2020-02-28T14:39:47.067940286-05:00","level":"error","message":"some error message","metadata":{"caller":"runtime.main{203}-\u003emain.main{12}","type":"test"}}
```

### Log Output
Use SetLogOutput to repoint 'log' to output to a different io.Writer. This can be os.Stdout, a buffer, or any other io.Writer.
```go
var l jerrors.List
l.Add(jerrors.New(jerrors.ERROR, "some error message", "type", "test"))
buf := new(bytes.Buffer)
jerrors.SetLogOutput(buf)
l.Log()

line := buf.String()
fmt.Print(line)
```
