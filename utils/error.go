package utils

import (
	"log"
	"os"
)

// HandleError handles errors with optional context and fatal flag.
// If context is empty, the error's string is used as context.
// If fatal is false (default), the program continues after logging.
func HandleError(err error, args ...interface{}) {
	if err == nil {
		return
	}

	// Default values
	context := err.Error()
	fatal := false

	// Parse optional arguments
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			context = v
		case bool:
			fatal = v
		}
	}

	// Log the error
	log.Printf("[ERROR] %s: %v", context, err)

	// Exit if fatal
	if fatal {
		os.Exit(1)
	}
}