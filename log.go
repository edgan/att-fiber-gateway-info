package main

import (
	"fmt"
	"os"
)

func logFatal(msg interface{}) {
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

func logFatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
