# jerrors Go JSON Error Framework
A simple golang Error module for creating errors and lists of errors in a JSON format.

## Contents
- [jerrors Go JSON Error Framework](#jerrors-go-json-error-framework)
	- [Contents](#contents)
	- [Installation](#installation)
	- [Quick Start](#quick-start)
	- [Error](#error)
		- [Creating A Error](#creating-a-error)
		- [Accessing Metadata](#accessing-metadata)
		- [Checking Error Levels](#checking-error-levels)
			- [Direct Comparison](#direct-comparison)
			- [IsError](#iserror)
			- [IsFatal](#isfatal)
		- [Error Logging](#error-logging)
			- [Error.Log()](#errorlog)
			- [Error.Fatal()](#errorfatal)
	- [Levels](#levels)
		- [Level Definition](#level-definition)
	- [Errors](#errors)
		- [Creating New Errors Errors](#creating-new-errors-errors)
		- [Adding An Error to Errors](#adding-an-error-to-errors)
		- [Checking Errors Levels](#checking-errors-levels)
			- [Check Function](#check-function)
			- [Comparing Errors.Error Directly](#comparing-errorserror-directly)
			- [IsError](#iserror-1)
			- [IsFatal](#isfatal-1)
		- [Errors Manipulation](#errors-manipulation)
			- [Append Errors](#append-errors)
			- [Stack Errorss](#stack-errorss)
			- [Clear Errors](#clear-errors)
		- [Errors Conversions](#errors-conversions)
			- [Error and String](#error-and-string)
			- [JSON](#json)
			- [ToArray](#toarray)
			- [ToLogArray](#tologarray)
		- [Errors Logging](#errors-logging)
			- [Errors Log](#errors-log)
			- [Errors Fatal](#errors-fatal)
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
func main() {
	// Create a new Errors object to store our errors in.
	errs := jerrors.New()

	// Create some errors.
	err1 := jerrors.NewError(jerrors.ERROR, "simple error message")
	err2 := jerrors.NewError(jerrors.WARN, "error message with metadata", "key1", "value1", "key2", "value2")

	// Add the errors to the Errors object.
	errs.Add(err1)
	errs.Add(err2)

	// The Level of the Errors object will be the highest level error in the list.
	fmt.Printf("errs.Level: %s\n", errs.Level)
	// Output: errs.Level: error

	fmt.Println(errs.Pretty())
	/* Output:
	{
	  "errors": [
	    {
	      "time": "2024-07-18T13:09:25.355507403-04:00",
	      "level": "error",
	      "message": "simple error message"
	    },
	    {
	      "time": "2024-07-18T13:09:25.355507652-04:00",
	      "level": "warn",
	      "message": "error message with metadata",
	      "metadata": {
	        "key1": "value1",
	        "key2": "value2"
	      }
	    },
	  ],
	  "level": "fatal"
	}
	*/

	// Create a second Errors object with different errors.
	errs2 := jerrors.New()
	err3 := jerrors.NewError(jerrors.FATAL, "fatal error message", "key1", "value1")
	err4 := jerrors.NewError(jerrors.DEBUG, "debug error message", "key1", "value1", "key2", "value2", "key3", "value3")
	errs2.Add(err3)
	errs2.Add(err4)

	fmt.Printf("errs.Level: %s\n", errs2.Level)
	// Output: errs.Level: fatal

	// Combine the two error lists.
	errs.Append(errs2)

	// Does errs contain any error with a Level of ERROR or higher?
	if errs.IsError() {
		fmt.Println("errs contains at least one error with a level of error or higher\n")
	}

	// Append() updates the level to the highest error level in both lists.
	fmt.Printf("errs.Level: %s\n", errs.Level)
	// Output: errs.Level: fatal

	// Normal marshal functionality.
	j, err := json.Marshal(errs)
	if err != nil {
		log.Fatalf("error marshalling errors: %s", err)
	}
	SendToClient(j)

	// Set log to write to stdout instead of stderr. You can pass any io.Writer like bytes.Buffer as well.
	jerrors.SetLogOutput(os.Stdout)

	// Log to stdout.
	errs.Log()
	/* Output:
	{"time":"2024-07-18T13:09:25.355507403-04:00","level":"error","message":"simple error message"}
	{"time":"2024-07-18T13:09:25.355507652-04:00","level":"warn","message":"error message with metadata","metadata":{"key1":"value1","key2":"value2"}}
	{"time":"2024-07-18T13:09:25.355552851-04:00","level":"fatal","message":"fatal error message","metadata":{"key1":"value1"}}
	{"time":"2024-07-18T13:09:25.355554016-04:00","level":"debug","message":"debug error message","metadata":{"key1":"value1","key2":"value2","key3":"value3"}}
	*/
}
```

## Error
```go
type Error struct {
    Time     time.Time         `json:"time,omitempty"`
    Level    Level             `json:"level,omitempty"`
    Message  string            `json:"message"`
    Metadata map[string]string `json:"metadata,omitempty"`
}
```

### Creating A Error
You can create an Error with just a Level and Message or you can add as many key value pairs of metadata as you want.
```go
err := jerrors.NewError(jerrors.ERROR, "simple error message")
err.Log()

errWithMeta := jerrors.NewError(
	jerrors.ERROR,
	"this error has metadata key pairs",
	"my",
	"meta",
	"key",
	"pairs",
	"sev",
	4,
)
errWithMeta.Log()
```

Output:
```
{"time":"2024-04-06T12:50:32.674656319-04:00","level":"error","message":"simple error message"}
{"time":"2024-04-06T12:50:32.674788908-04:00","level":"error","message":"this error has metadata key pairs","metadata":{"key":"pairs","my":"meta","sev":"4"}}
```

### Accessing Metadata
The Metadata of an Error can be accessed directly as a map[string]string.
```go
err := jerrors.NewError(jerrors.ERROR, "simple error message", "user", "bob")
if err.Metadata["user"] == "bob" {
    fmt.Printf("Why %s?! %v", err.Metadata["user"], err)
}
```

### Checking Error Levels
 See [Levels](#levels) for details on jerrors.Level.

#### Direct Comparison
```go
err := jerrors.NewError(jerrors.ERROR, "simple error message")
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
err := jerrors.NewError(jerrors.ERROR, "simple error message")
err.Log()
```
Output
```
{"time":"2020-02-26T13:11:40.038906297-05:00","level":"error","message":"simple error message","metadata":{"caller":"runtime.main{203}-\u003emain.main{11}"}}
```

#### Error.Fatal()
Fatal sets the Error.Level to FATAL before logging the error and exiting with a status of 1.
```go
err := jerrors.NewError(jerrors.ERROR, "simple error message")
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

## Errors
```go
type Errors struct {
	Errors []Error `json:"errors"`
	Level  Level   `json:"level"`
}
```
Errors contains a list of Error and the highest Level of any Error added to the list using Errors.Add().

### Creating New Errors Errors
```go
errs := jerrors.New()
```

### Adding An Error to Errors
```go
errs.Add(jerrors.NewError(jerrors.DEBUG, "some debug message"))
fmt.Println(errs.Level) // Prints: debug

errs.Add(jerrors.NewError(jerrors.ERROR, "some error message", "with", "metadata"))
fmt.Println(errs.Level) // Prints: error

err := jerrors.NewError(jerrors.WARN, "some warning message", "with", "metadata")
errs.Add(err)
fmt.Println(errs.Level) // Prints: error

errs.Log()
```

Output:
```
debug
error
error
{"time":"2024-04-06T11:47:37.206520985-04:00","level":"debug","message":"some debug message"}
{"time":"2024-04-06T11:47:37.206584203-04:00","level":"error","message":"some error message","metadata":{"with":"metadata"}}
{"time":"2024-04-06T11:47:37.206589403-04:00","level":"warn","message":"some warning message","metadata":{"with":"metadata"}}
```

### Checking Errors Levels

#### Check Function
Errors.Check() checks to see if Errors.Errors is not empty. Returns the number of errors and bool.
Errors.IsEmpty() returns true if there are no errors in Errors.Error.
```go
errs := jerrors.New()
errs.Add(jerrors.NewError(jerrors.DEBUG, "some debug message"))
errs.Add(jerrors.NewError(jerrors.ERROR, "some error message"))
if hasErrors, count :=errs.Check(); hasErrors {
	fmt.Println(count) // Prints: 2
}

if errs.IsEmpty() {
    fmt.Println("no errors found")
}
```

#### Comparing Errors.Error Directly
You can access the Errors Level directly and compaire just like with Error.
```go
iferrs.Level >= jerrors.ERROR {
	l.Log()
}
```

#### IsError
IsError returns true if Errors.Level is >= ERROR.
```go
iferrs.IsError() {
   errs.Log()
}
```

#### IsFatal
IsFatal returns true only if Errors.Level == FATAL.
```go
iferrs.IsFatal() {
   errs.Fatal()
}
```

### Errors Manipulation

#### Append Errors
You can append a new Errors to and existing List. Append works just like array Append adding the argument list to the bottom of the list calling the method.
```go
var l jerrors.Errors
errs.Add(jerrors.NewError(jerrors.DEBUG, "some helpful message"))
errs.Add(jerrors.NewError(jerrors.ERROR, "some error message"))

var l2 jerrors.Errors
l2.Add(jerrors.NewError(jerrors.WARN, "some warning message"))

errs.Append(l2)
fmt.Println(errs.Errors[0].Message) // Prints: some helpful message
```

#### Stack Errorss
You can stack two lists together. Stack puts the aurgument Errors at the top of the List calling the method.
```go
var l jerrors.Errors
errs.Add(jerrors.NewError(jerrors.DEBUG, "some helpful message"))
errs.Add(jerrors.NewError(jerrors.ERROR, "some error message"))

var l2 jerrors.Errors
l2.Add(jerrors.NewError(jerrors.WARN, "some warning message"))

errs.Stack(l2)
fmt.Println(errs.Errors[0].Message) // Prints: some warning message
```

#### Clear Errors
Clear deletes all Errors in the Errors and sets Level to 0.
```go
var l jerrors.Errors
errs.Add(jerrors.NewError(jerrors.DEBUG, "some helpful message"))
errs.Add(jerrors.NewError(jerrors.ERROR, "some error message"))
hasErrors, count :=errs.Check() // true, 2
errs.Clear()
hasErrors, count =errs.Check() // false, 0
```

### Errors Conversions

#### Error and String
Errors implements the error interface and can be used the same as any other error. Error converts List.Errors into a single JSON string.
String is an alias of Error.
```go
var l jerrors.Errors
errs.Add(jerrors.NewError(jerrors.DEBUG, "some helpful message"))
errs.Add(jerrors.NewError(jerrors.ERROR, "some error message"))
fmt.Printf("Errors of errors:\n%v\n",errs.Error())
```

Output:
```
Errors of errors:
[{"time":"2020-02-27T15:56:01.914821132-05:00","level":"debug","message":"some helpful message","metadata":{"caller":"runtime.main{203}-\u003emain.main{11}"}},{"time":"2020-02-27T15:56:01.914955631-05:00","level":"error","message":"some error message","metadata":{"caller":"runtime.main{203}-\u003emain.main{15}"}}]
```

#### JSON
You can use json.Marshall and json.Unmarshall to convert between Errors and JSON.
```go
var l Errors
errs.Add(jerrors.NewError(jerrors.DEBUG, "some helpful message"))
j, _ := json.Marshall(l)
fmt.Println(string(j))

var l2 Errors
if err := json.Unmarshall(l, &l2); err != nil {
	fmt.Fatal(err)
}

l2.Log()
```

#### ToArray
Converts the Errors.Errors to an array of JSON strings. Ignores Log Level.
```go
var l jerrors.Errors
errs.Add(jerrors.NewError(jerrors.DEBUG, "some helpful message"))
errs.Add(jerrors.NewError(jerrors.ERROR, "some error message"))
a :=errs.ToArray()
fmt.Println(strings.Join(a, "\n"))
```
Output:
```
{"time":"2020-02-28T13:30:21.668092176-05:00","level":"debug","message":"some helpful message","metadata":{"caller":"runtime.main{203}-\u003emain.main{12}"}}
{"time":"2020-02-28T13:30:21.668119405-05:00","level":"error","message":"some error message","metadata":{"caller":"runtime.main{203}-\u003emain.main{13}"}}
```

#### ToLogArray
Converts the Errors.Errors to an array of JSON strings. Honors Log Level.
```go
var l jerrors.Errors
errs.Add(jerrors.NewError(jerrors.DEBUG, "some helpful message"))
errs.Add(jerrors.NewError(jerrors.ERROR, "some error message"))
a :=errs.ToLogArray()
fmt.Println(strings.Join(a, "\n"))
```
Output:
```
{"time":"2020-02-28T13:31:20.453088284-05:00","level":"error","message":"some error message","metadata":{"caller":"runtime.main{203}-\u003emain.main{13}"}}
```

### Errors Logging

#### Errors Log
Call log.Print to log the array of Errors. Any Error with a Level >= current Log Level will be logged. See [Logs](#logs) for more details.
```go
var l jerrors.Errors
errs.Add(jerrors.NewError(jerrors.DEBUG, "some helpful message"))
errs.Add(jerrors.NewError(jerrors.ERROR, "some error message"))
errs.Log()
```
Output:
```
{"time":"2020-02-28T13:31:20.453088284-05:00","level":"error","message":"some error message","metadata":{"caller":"runtime.main{203}-\u003emain.main{13}"}}
```

#### Errors Fatal
Same as [Errors.Log()](#list-log) but uses log.Fatal to exit with a status of 1.
```go
var l jerrors.Errors
errs.Add(jerrors.NewError(jerrors.DEBUG, "some helpful message"))
errs.Add(jerrors.NewError(jerrors.ERROR, "some error message"))
errs.Log()
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

var l jerrors.Errors
errs.Add(jerrors.NewError(jerrors.ERROR, "some error message", "type", "test"))
errs.Log()
```
Output:
```
{"level":"error","message":"some error message","metadata":{"type":"test"}}
```

### Logging Level
The Logging Level determine what should get logged when calling Log() or Fatal on both Error and Errors. Log/Fatal will log all Errors with a Level >= the current Log Level.
You can use SetLogLevel to change the Logging Level. Defaults to INFO.

```go
var l jerrors.Errors
errs.Add(jerrors.NewError(jerrors.DEBUG, "some helpful message"))
errs.Add(jerrors.NewError(jerrors.ERROR, "some error message", "type", "test"))

fmt.Println("Default Log Level: INFO")
errs.Log()

jerrors.SetLogLevel(jerrors.DEBUG)
fmt.Println("\nLog Level: DEBUG")
errs.Log()
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
var l jerrors.Errors
errs.Add(jerrors.NewError(jerrors.ERROR, "some error message", "type", "test"))
buf := new(bytes.Buffer)
jerrors.SetLogOutput(buf)
errs.Log()

line := buf.String()
fmt.Print(line)
```
