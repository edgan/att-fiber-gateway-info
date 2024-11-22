package main

func allMetrics(actionPages map[string]string, client *GatewayClient, configs Configs, flags *Flags, model string, natActionPrefix string) {
	metricActions := returnMeticsActions()

	returnFact := ""

	for _, action := range metricActions {
		_, err := client.retrieveAction(action, actionPages, configs, flags, model, natActionPrefix, returnFact)

		if err != nil {
			logFatalf("Failed to get %s: %v", action, err)
		}
	}
}
