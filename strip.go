package main

import (
	"regexp"
)

// Helper function to remove ANSI escape codes (if any)
func stripAnsi(str string) string {
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return ansiRegex.ReplaceAllString(str, "")
}
