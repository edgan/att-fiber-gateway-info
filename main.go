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

	// Return action pages
	actionPages := returnActionPages()

	// Return a actions description
	actionsDescription := actionsHelp()

	// Return a filters description
	filtersDescription := filtersHelp()

	// Return a cookie path
	cookiePath := determineCookiePath()

	// Flags
	action, answerNo, answerYes, baseURLFlag, cookieFile, debug, filter, freshCookies,
		passwordFlag, pretty := returnFlags(actionsDescription, colorMode, cookiePath, filtersDescription)

	// Validate flags
	baseURL, loginRequired, page, password := validateFlags(action, actionPages, baseURLFlag, config, filter, passwordFlag)

	client, err := createGatewayClient(baseURL, colorMode, *cookieFile, *debug, *freshCookies)
	if err != nil {
		log.Fatalf("Failed to create router client: %v", err)
	}

	if err := client.getPage(*action, *answerNo, *answerYes, *filter, loginRequired, "nat-", page, password, *pretty); err != nil {
		log.Fatalf("Failed to get %s: %v", action, err)
	}

	os.Exit(0)
}
