package jerrors

import (
	"io"
	"log"
	"os"
)

func init() {
	// Setup default log options
	log.SetFlags(0)
}

// SetLogOutput sets the logging destination.
// Example:
// var buf bytes.Buffer
// log.SetOutput(&buf)
func SetLogOutput(w io.Writer) {
	if w == nil {
		log.SetOutput(os.Stderr)
		return
	}
	log.SetOutput(w)
}
