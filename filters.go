package main

import (
	"fmt"
	"strings"
)

func returnFilters() []string {
	filters := []string{"icmp", "ipv4", "ipv6", "tcp", "udp"}

	return filters
}

// filtersHelp returns the filters available for actions
func filtersHelp() string {
	filters := returnFilters()
	filtersHelp := []string{}

	for _, filter := range filters {
		filtersHelp = append(filtersHelp, filter)
	}

	return fmt.Sprintf("Filter to perform (%s)", strings.Join(filtersHelp, ", "))
}
