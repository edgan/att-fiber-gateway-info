package main

func outputMetrics(action string, datadog bool, header string, model string, summary string, statsdIPPort string, tableData [][]string) {
	actionMetric := returnActionMetric(action)
	modelActionMetric := model + "." + actionMetric

	dotZero := ".0"

	metrics := []string{}

	if action == "fiber-status" {
		metrics = generateFiberMetric(dotZero, header, modelActionMetric, tableData)
	} else {
		metrics = generateNonFiberMetric(dotZero, modelActionMetric, summary, tableData)
	}

	if datadog {
		giveMetricsToDatadogStatsd(metrics, model, statsdIPPort)
	} else {
		printMetrics(metrics)
	}
}
