package main

import (
	"fmt"
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
	action, flags, version := returnFlags(actionsDescription, colorMode, cookiePath, filtersDescription)

	if *version {
		appVersion := returnApplicationNameVersion()
		fmt.Println(appVersion)
		os.Exit(0)
	}

	// Validate flags
	configs, flags := validateFlags(*action, actionPages, config, flags)

	client, err := createGatewayClient(configs, colorMode, flags)
	if err != nil {
		logFatalf("Failed to create router client: %v", err)
	}

	returnFact := "model"

	// "" is the model variable
	model, err := client.retrieveAction("system-information", actionPages, configs, flags, "", actionPrefixes["nat"], returnFact)

	if err != nil {
		logFatalf("Failed to get %s: %v", action, err)
	}

	if *flags.AllMetrics {
		if *flags.Continuous {
			for {
				allMetrics(actionPages, client, configs, flags, model, actionPrefixes["nat"])
			}
		} else {
			allMetrics(actionPages, client, configs, flags, model, actionPrefixes["nat"])
			os.Exit(0)
		}
	}

	returnFact = ""

	if *flags.Metrics {
		if *flags.Continuous {
			for {
				_, err = client.retrieveAction(*action, actionPages, configs, flags, model, actionPrefixes["nat"], returnFact)

				if err != nil {
					logFatalf("Failed to get %s: %v", action, err)
				}
			}
		}
	}

	_, err = client.retrieveAction(*action, actionPages, configs, flags, model, actionPrefixes["nat"], returnFact)

	if err != nil {
		logFatalf("Failed to get %s: %v", action, err)
	}

	os.Exit(0)
}
