package main

func outputMetrics(action string, configs Configs, flags *Flags, header string, model string, summary string, tableData [][]string) {
	actionMetric := returnActionMetric(action, flags)
	modelActionMetric := model + "." + actionMetric

	dotZero := ".0"

	metrics := []string{}

	if action == "fiber-status" {
		metrics = generateFiberMetric(dotZero, header, modelActionMetric, tableData)
	} else {
		metrics = generateNonFiberMetric(action, dotZero, flags, modelActionMetric, summary, tableData)
	}

	if *flags.Datadog {
		giveMetricsToDatadogStatsd(configs, metrics, model)
	} else {
		printMetrics(metrics)
	}
}
