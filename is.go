package main

import (
	"regexp"
	"strings"
)

// Helper function to validate table classes
func isValidTableClass(class string, validClasses []string) bool {
	for _, validClass := range validClasses {
		if strings.Contains(class, validClass) {
			return true
		}
	}
	return false
}

// Helper function to validate and process table summaries
func isValidSummary(summary string) (bool, string) {
	if strings.Contains(strings.ToLower(summary), "statistics") ||
		summary == "Summary of nattable connections" ||
		summary == "This table displays a summary of session information." {

		shortSummary := summary
		patterns := []string{
			` [Ss]tatistic.*`,
			`Ethernet `,
			`This table displays `,
		}
		for _, pattern := range patterns {
			re := regexp.MustCompile(pattern)
			shortSummary = re.ReplaceAllString(shortSummary, "")
		}
		return true, shortSummary
	}
	return false, ""
}
