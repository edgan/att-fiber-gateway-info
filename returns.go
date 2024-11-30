package main

import (
	"embed"
	"strings"

	"edgan/att-fiber-gateway-info/internal/logging"
)

func returnApplicationName() string {
	applicationName := "att-fiber-gateway-info"

	return applicationName
}

func returnApplicationNameVersion() string {
	version := returnVersion()
	applicationName := returnApplicationName()
	applicationNameVersion := applicationName + space + version

	return applicationNameVersion
}

//go:embed .version
var versionFile embed.FS

func returnVersion() string {
	versionBytes, _ := versionFile.ReadFile(".version")
	version := strings.TrimSpace(string(versionBytes))

	return version
}

func returnActionMetric(action string) string {
	actionMetric := strings.Replace(action, dash, period, twoOccurance)
	actionMetric = strings.Replace(actionMetric, space, period, oneOccurance)

	return actionMetric
}

// returnActionPage returns the page for a specific action
func returnActionPage(action string, actionPages map[string]string) string {
	page, exists := actionPages[action]
	if !exists {
		logging.LogFatalf("Unknown action: %s", action)
	}
	return page
}

// returnActionPages returns action to page mappings
func returnActionPages() map[string]string {
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
		if strings.HasPrefix(action, natActionPrefix) {
			actionPages[action] = "nattable"
		}
		if strings.HasPrefix(action, resetActionPrefix) {
			actionPages[action] = "reset"
		}
	}

	return actionPages
}

func returnPath(page string) string {
	path := "/cgi-bin/" + page + ".ha"

	return path
}

func returnConfigValue(flagValue string, configValue string) string {
	if flagValue != empty {
		return flagValue
	}
	return configValue
}
