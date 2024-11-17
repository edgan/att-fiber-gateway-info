package main

import (
	"fmt"
	"log"
	"strings"
)

func returnActions() []string {
	actions := []string{
		"broadband-status", "device-list", "fiber-status", "home-network-status",
		"ip-allocation", "nat-check", "nat-connections", "nat-destinations",
		"nat-sources", "nat-totals", "reset-connection", "reset-device",
		"reset-firewall", "reset-ip", "reset-wifi", "restart-gateway",
		"system-information",
	}

	return actions
}

// actionsHelp returns the actions available for actions
func actionsHelp() string {
	actions := returnActions()
	actionsHelp := []string{}

	for _, action := range actions {
		actionsHelp = append(actionsHelp, action)
	}

	return fmt.Sprintf("Action to perform (%s)", strings.Join(actionsHelp, ", "))
}

// returnActionPages returns action to page mappings
func returnActionPages() map[string]string {
	actions := returnActions()

	actionPages := map[string]string{
		"broadband-status":    "broadbandstatistics",
		"device-list":         "devices",
		"fiber-status":        "fiberstat",
		"home-network-status": "lanstatistics",
		"ip-allocation":       "ipalloc",
		"restart-gateway":     "reset",
		"system-information":  "sysinfo",
	}

	natActionPrefix := "nat-"
	resetActionPrefix := "reset-"

	for _, action := range actions {
		if strings.HasPrefix(action, natActionPrefix) {
			actionPages[action] = "nattable"
		}
		if strings.HasPrefix(action, resetActionPrefix) {
			actionPages[action] = "reset"
		}
	}

	return actionPages
}

// getActionPage gets the page for a specific action
func getActionPage(actionPages map[string]string, action string) string {
	page, exists := actionPages[action]
	if !exists {
		log.Fatalf("Unknown action: %s", action)
	}
	return page
}
