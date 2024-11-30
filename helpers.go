package main

import (
	"sort"
)

// Helper function to get map keys as a sorted slice
func getMapKeys(m map[string]string) []string {
	keys := make([]string, zero, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys) // Sort keys alphabetically
	return keys
}

// Helper function to check if a map contains a key
func containsMapKey(m map[string]string, key string) bool {
	_, exists := m[key]
	return exists
}

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
