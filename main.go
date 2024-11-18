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
	action, answerNo, answerYes, baseURLFlag, cookieFile, debug, filter, freshCookies, metrics,
		passwordFlag, pretty := returnFlags(actionsDescription, colorMode, cookiePath, filtersDescription)

	// Validate flags
	baseURL, loginRequired, page, password := validateFlags(action, actionPages, baseURLFlag, config, filter, passwordFlag)

	client, err := createGatewayClient(baseURL, colorMode, *cookieFile, *debug, *freshCookies)
	if err != nil {
		log.Fatalf("Failed to create router client: %v", err)
	}

	returnFact := "model"
	model, err := client.getPage("system-information", *answerNo, *answerYes, *filter, false, *metrics, "", actionPrefixes["nat"], "sysinfo", password, *pretty, returnFact)
	if err != nil {
		log.Fatalf("Failed to get %s: %v", action, err)
	}

	returnFact = ""
	_, err = client.getPage(*action, *answerNo, *answerYes, *filter, loginRequired, *metrics, model, actionPrefixes["nat"], page, password, *pretty, returnFact)
	if err != nil {
		log.Fatalf("Failed to get %s: %v", action, err)
	}

	os.Exit(0)
}
