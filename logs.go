package jerrors

import (
	"io"
	"log"
)

func init() {
	// Setup default log options
	log.SetFlags(0)
}

// SetLogOutput sets the logging destination.
func SetLogOutput(w io.Writer) {
	log.SetOutput(w)
}
