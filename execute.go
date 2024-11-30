package main

import (
	"edgan/att-fiber-gateway-info/internal/logging"
)

func executeAllMetrics(
	actionPages map[string]string, client *gatewayClient, configs configs,
	flags *flags, model string, continuous bool,
) {
	for {
		allMetrics(actionPages, client, configs, flags, model)
		if !continuous {
			break
		}
	}
}

func executeRetrieveAction(
	client *gatewayClient, action string, actionPages map[string]string,
	configs configs, flags *flags, model string, continuous bool,
) {
	for {
		_, err := client.retrieveAction(action, actionPages, configs, flags, model, empty)
		if err != nil {
			logging.LogFatalf("Failed to get %s: %v", action, err)
		}
		if !continuous {
			break
		}
	}
}
