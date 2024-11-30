// att-fiber-gateway-info
// A golang command line tool to pull values from the pages of a AT&T Fiber gateways.
package main

import (
	"fmt"
	"os"

	"edgan/att-fiber-gateway-info/internal/color"
	"edgan/att-fiber-gateway-info/internal/logging"
)

func main() {
	colorMode := color.CheckColorTerminal()
	config := loadAppConfig(determineConfigFile())
	actionPages := returnActionPages()
	actionsDescription := actionsHelp()
	filtersDescription := filtersHelp()
	cookiePath := determineCookiePath()
	action, flags, version := returnFlags(actionsDescription, colorMode, cookiePath, filtersDescription)

	if *version {
		fmt.Println(returnApplicationNameVersion())
		os.Exit(zero)
	}

	// Validate flags
	configs, flags := validateFlags(*action, actionPages, config, flags)

	// Create client
	client, err := createGatewayClient(configs, colorMode, flags)
	if err != nil {
		logging.LogFatalf("Failed to create router client: %v", err)
	}

	// Retrieve model information
	model, err := client.retrieveAction("system-information", actionPages, configs, flags, empty, "model")
	if err != nil {
		logging.LogFatalf("Failed to get system-information: %v", err)
	}

	// Handle actions based on flags
	if *flags.AllMetrics {
		executeAllMetrics(actionPages, client, configs, flags, model, *flags.Continuous)
	} else {
		executeRetrieveAction(client, *action, actionPages, configs, flags, model, *flags.Continuous)
	}
}
