package main

import (
	"os"
)

func executeAllMetrics(actionPages map[string]string, client *GatewayClient, configs Configs, flags *Flags, model string, continuous bool) {
	for {
		allMetrics(actionPages, client, configs, flags, model)
		if !continuous {
			break
		}
	}
	os.Exit(0)
}

func executeRetrieveAction(client *GatewayClient, action string, actionPages map[string]string, configs Configs, flags *Flags, model string, continuous bool) {
	for {
		_, err := client.retrieveAction(action, actionPages, configs, flags, model, "")
		if err != nil {
			logFatalf("Failed to get %s: %v", action, err)
		}
		if !continuous {
			break
		}
	}
	os.Exit(0)
}
