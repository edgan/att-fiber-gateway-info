package main

import (
	"fmt"
	"strings"
)

// actionsHelp returns the actions available for actions
func actionsHelp() string {
	actions := returnActions()
	actionsHelp := []string{}

	for _, action := range actions {
		actionsHelp = append(actionsHelp, action)
	}

	return fmt.Sprintf("Action to perform (%s)", strings.Join(actionsHelp, ", "))
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
