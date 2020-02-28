package jerrors

import (
	"io"
	"log"
)

var (
	logCaller    bool
	logLevel     bool
	logTime      bool
	loggingLevel Level
)

func init() {
	// Setup default log options
	logCaller = true
	logLevel = true
	logTime = true
	loggingLevel = INFO
	log.SetFlags(0)
}

// SetLogOptions for additional error data.
func SetLogOptions(options map[string]bool) {
	for k, v := range options {
		switch k {
		case "caller":
			logCaller = v
		case "level":
			logLevel = v
		case "time":
			logTime = v
		}
	}
}

// SetLogLevel sets what level to log. Will log set level and above.
func SetLogLevel(level Level) {
	loggingLevel = level
}

// SetLogOutput sets the logging destination.
func SetLogOutput(w io.Writer) {
	log.SetOutput(w)
}
