package main

import (
	"strings"
)

func generateFiberMetric(dotZero string, header string, modelActionMetric string) []string {
	metrics := []string{}
	fiber := "Currently"

	if strings.Contains(header, fiber) {
		keyValue := empty
		keyValue = strings.Replace(header, fiber, empty, one)
		keyValue = strings.Replace(keyValue, "\u00A0\u00A0 ", equals, one)
		keyValue = strings.Replace(keyValue, space, period, one)
		metric := modelActionMetric + period + keyValue + dotZero
		metrics = append(metrics, metric)
	}

	return metrics
}

func generateNonFiberMetric(
	action, dotZero string, flags *flags, modelActionMetric, summary string, tableData [][]string,
) []string {
	var metrics []string
	lowerSummary := strings.ToLower(strings.Replace(summary, space, period, one))

	if action == "nat-totals" {
		metrics = processNatTotalsAction(tableData, modelActionMetric, dotZero, flags)
	} else {
		metrics = processGeneralAction(tableData, modelActionMetric, lowerSummary, dotZero, flags)
	}

	return metrics
}
