// Package logging is for logging functions
package logging

import (
	"fmt"
	"os"
)

// LogFatal is print + exit
func LogFatal(msg any) {
	var output string
	switch v := msg.(type) {
	case string:
		output = v
	case error:
		output = v.Error()
	default:
		output = "Unknown error"
	}

	fmt.Fprintf(os.Stderr, "Error: %s\n", output)
	os.Exit(1)
}

// LogFatalf is printf + exit
func LogFatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
