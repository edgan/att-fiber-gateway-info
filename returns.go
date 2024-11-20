package main

import (
	"embed"
	"log"
	"strings"
)

func returnApplicationName() string {
	applicationName := "att-fiber-gateway-info"

	return applicationName
}

func returnApplicationNameVersion() string {
	version := returnVersion()
	applicationName := returnApplicationName()
	applicationNameVersion := applicationName + " " + version

	return applicationNameVersion
}

//go:embed .version
var versionFile embed.FS

func returnVersion() string {
	versionBytes, _ := versionFile.ReadFile(".version")
	version := strings.TrimSpace(string(versionBytes))

	return version
}

func returnActionMetric(action string, flags *Flags) string {
	actionMetric := strings.Replace(action, "-", ".", 2)
	actionMetric = strings.Replace(actionMetric, " ", ".", 1)

	return actionMetric
}

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

// returnActionPage returns the page for a specific action
func returnActionPage(action string, actionPages map[string]string) string {
	page, exists := actionPages[action]
	if !exists {
		log.Fatalf("Unknown action: %s", action)
	}
	return page
}

// returnActionPages returns action to page mappings
func returnActionPages(actionPrefixes map[string]string) map[string]string {
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

	for _, action := range actions {
		if strings.HasPrefix(action, actionPrefixes["nat"]) {
			actionPages[action] = "nattable"
		}
		if strings.HasPrefix(action, actionPrefixes["reset"]) {
			actionPages[action] = "reset"
		}
	}

	return actionPages
}

// returnActionPages returns action to page mappings
func returnActionPrefixes() map[string]string {
	actionPrefixes := map[string]string{
		"nat":   "nat-",
		"reset": "reset-",
	}

	return actionPrefixes
}

func returnFilters() []string {
	filters := []string{"icmp", "ipv4", "ipv6", "tcp", "udp"}

	return filters
}

func returnMeticsActions() []string {
	metricActions := []string{"broadband-status", "fiber-status", "home-network-status", "nat-totals"}

	return metricActions
}

func returnPath(page string) string {
	path := "/cgi-bin/" + page + ".ha"

	return path
}
