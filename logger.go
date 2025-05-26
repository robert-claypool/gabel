package main

import (
	"fmt"
	"os"
)

var debugMode bool

// InitLogger sets up logging based on debug flag
func InitLogger(debug bool) {
	debugMode = debug
}

// LogDebug prints debug messages to stderr when debug mode is enabled
func LogDebug(format string, args ...interface{}) {
	if debugMode {
		fmt.Fprintf(os.Stderr, "[DEBUG] "+format+"\n", args...)
	}
}

// LogError prints error messages to stderr
func LogError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "[ERROR] "+format+"\n", args...)
}