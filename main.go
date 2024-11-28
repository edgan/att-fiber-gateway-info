package main

import (
	"fmt"
	"os"
)

func main() {
	// Initial setup
	colorMode := checkColorTerminal()
	config := loadAppConfig(determineConfigFile())
	actionPrefixes := returnActionPrefixes()
	actionPages := returnActionPages(actionPrefixes)
	actionsDescription := actionsHelp()
	filtersDescription := filtersHelp()
	cookiePath := determineCookiePath()
	action, flags, version := returnFlags(actionsDescription, colorMode, cookiePath, filtersDescription)

	if *version {
		fmt.Println(returnApplicationNameVersion())
		os.Exit(0)
	}

	// Validate flags
	configs, flags := validateFlags(*action, actionPages, config, flags)

	// Create client
	client, err := createGatewayClient(configs, colorMode, flags)
	if err != nil {
		logFatalf("Failed to create router client: %v", err)
	}

	// Retrieve model information
	model, err := client.retrieveAction("system-information", actionPages, configs, flags, "", "model")
	if err != nil {
		logFatalf("Failed to get system-information: %v", err)
	}

	// Handle actions based on flags
	if *flags.AllMetrics {
		executeAllMetrics(actionPages, client, configs, flags, model, *flags.Continuous)
	} else if *flags.Metrics {
		executeRetrieveAction(client, *action, actionPages, configs, flags, model, *flags.Continuous)
	} else {
		executeRetrieveAction(client, *action, actionPages, configs, flags, model, *flags.Continuous)
	}
}
