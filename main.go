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
	action, allmetrics, answerNo, answerYes, baseURLFlag, cookieFile, debug, filter, freshCookies, metrics,
		passwordFlag, pretty := returnFlags(actionsDescription, colorMode, cookiePath, filtersDescription)

	// Validate flags
	allmetrics, baseURL, metrics, password := validateFlags(action, actionPages, allmetrics, baseURLFlag, config, filter, metrics, passwordFlag)

	client, err := createGatewayClient(baseURL, colorMode, *cookieFile, *debug, *freshCookies)
	if err != nil {
		log.Fatalf("Failed to create router client: %v", err)
	}

	returnFact := "model"

	model, err := client.retrieveAction("system-information", actionPages, *answerNo, *answerYes, *filter, *metrics, "", actionPrefixes["nat"], password, *pretty, returnFact)

	if err != nil {
		log.Fatalf("Failed to get %s: %v", action, err)
	}

	if *allmetrics {
		actions := []string{"broadband-status", "fiber-status", "home-network-status"}

                returnFact = ""

		for _, action := range actions {
			_, err = client.retrieveAction(action, actionPages, *answerNo, *answerYes, *filter, *metrics, model, actionPrefixes["nat"], password, *pretty, returnFact)

			if err != nil {
				log.Fatalf("Failed to get %s: %v", action, err)
			}
		}

		os.Exit(0)
	}

	returnFact = ""

	_, err = client.retrieveAction(*action, actionPages, *answerNo, *answerYes, *filter, *metrics, model, actionPrefixes["nat"], password, *pretty, returnFact)

	if err != nil {
		log.Fatalf("Failed to get %s: %v", action, err)
	}

	os.Exit(0)
}
