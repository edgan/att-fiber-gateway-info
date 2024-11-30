package logging

import (
	"fmt"
)

// debugLog prints debug messages when debug mode is enabled
func DebugLog(enabled bool, message string) {
	if enabled {
		fmt.Printf("Debug: %s\n", message)
	}
}
