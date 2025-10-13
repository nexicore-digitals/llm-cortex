package utils

import "log"

// HandleError checks if an error is not nil and panics if it is.
// This is a simple utility for handling critical errors that should stop execution.
func HandleError(err error) {
	if err != nil {
		log.Panicf("A critical error occurred: %v", err)
	}
}