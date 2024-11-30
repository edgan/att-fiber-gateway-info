package main

import (
	"fmt"
	"strings"
)

func actionsHelp() string {
	actionsHelp := []string{}

	actionsHelp = append(actionsHelp, actions...)

	return fmt.Sprintf("Action to perform (%s)", strings.Join(actionsHelp, commaSpace))
}

func filtersHelp() string {
	filtersHelp := []string{}

	filtersHelp = append(filtersHelp, filters...)

	return fmt.Sprintf("Filter to perform (%s)", strings.Join(filtersHelp, commaSpace))
}
