package main

import (
	"log"
	"os"
)

func main() {
	// Color mode detection
	colorMode := checkColorTerminal()

	// Load config file
	configFile := determineConfigFile()
	config := loadAppConfig(configFile)

	// Return action prefixes
	actionPrefixes := returnActionPrefixes()

	// Return action pages
	actionPages := returnActionPages(actionPrefixes)

	// Return a actions description
	actionsDescription := actionsHelp()

	// Return a filters description
	filtersDescription := filtersHelp()

	// Return a cookie path
	cookiePath := determineCookiePath()

	// Flags
	action, flags := returnFlags(actionsDescription, colorMode, cookiePath, filtersDescription)

	// Validate flags
	configs, flags := validateFlags(*action, actionPages, config, flags)

	client, err := createGatewayClient(configs, colorMode, flags)
	if err != nil {
		log.Fatalf("Failed to create router client: %v", err)
	}

	returnFact := "model"

	// "" is the model variable
	model, err := client.retrieveAction("system-information", actionPages, configs, flags, "", actionPrefixes["nat"], returnFact)

	if err != nil {
		log.Fatalf("Failed to get %s: %v", action, err)
	}

	if *flags.AllMetrics {
		metricActions := returnMeticsActions()

		returnFact = ""

		for _, action := range metricActions {
			_, err = client.retrieveAction(action, actionPages, configs, flags, model, actionPrefixes["nat"], returnFact)

			if err != nil {
				log.Fatalf("Failed to get %s: %v", action, err)
			}
		}

		os.Exit(0)
	}

	returnFact = ""

	_, err = client.retrieveAction(*action, actionPages, configs, flags, model, actionPrefixes["nat"], returnFact)

	if err != nil {
		log.Fatalf("Failed to get %s: %v", action, err)
	}

	os.Exit(0)
}
