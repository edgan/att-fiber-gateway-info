package main

import (
	"fmt"
)

// debugLog prints debug messages when debug mode is enabled
func debugLog(enabled bool, message string) {
	if enabled {
		fmt.Printf("Debug: %s\n", message)
	}
}
