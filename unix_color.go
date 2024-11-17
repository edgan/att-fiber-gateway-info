//go:build !windows
// +build !windows

package main

import (
	"os"
	"strings"
)

func isColorTerminal() bool {
	// Check if the TERM environment variable indicates a color terminal
	term := os.Getenv("TERM")

	if term == "" {
		return false
	}

	// Common color-enabled terminals
	colorTerminals := []string{"xterm", "xterm-256color", "screen", "screen-256color", "linux", "tmux"}
	for _, t := range colorTerminals {
		if strings.Contains(term, t) {
			return true
		}
	}

	// For UNIX-like systems (Linux/macOS), check if the output is connected to a terminal
	return isTerminal()
}

func isTerminal() bool {
	// Check if stdout is connected to a terminal
	fileInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}
