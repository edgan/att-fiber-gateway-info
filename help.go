package main

import (
	"fmt"
	"strings"
)

// actionsHelp returns the actions available for actions
func actionsHelp() string {
	actions := returnActions()
	actionsHelp := []string{}

	actionsHelp = append(actionsHelp, actions...)

	return fmt.Sprintf("Action to perform (%s)", strings.Join(actionsHelp, ", "))
}

// filtersHelp returns the filters available for actions
func filtersHelp() string {
	filters := returnFilters()
	filtersHelp := []string{}

	filtersHelp = append(filtersHelp, filters...)

	return fmt.Sprintf("Filter to perform (%s)", strings.Join(filtersHelp, ", "))
}
