package main

import (
	"strings"
)

func generateFiberMetric(dotZero string, header string, modelActionMetric string) []string {
	metrics := []string{}
	fiber := "Currently"

	if strings.Contains(header, fiber) {
		keyPlusValue := empty
		keyPlusValue = strings.Replace(header, fiber, empty, oneOccurance)
		keyPlusValue = strings.Replace(keyPlusValue, "\u00A0\u00A0 ", equals, oneOccurance)
		keyPlusValue = strings.Replace(keyPlusValue, space, period, oneOccurance)
		metric := modelActionMetric + period + keyPlusValue + dotZero
		metrics = append(metrics, metric)
	}

	return metrics
}

func generateNonFiberMetric(
	action, dotZero string, flags *flags, modelActionMetric, summary string, tableData [][]string,
) []string {
	var metrics []string
	lowerSummary := strings.ToLower(strings.Replace(summary, space, period, oneOccurance))

	if action == "nat-totals" {
		metrics = processNatTotalsAction(tableData, modelActionMetric, dotZero, flags)
	} else {
		metrics = processGeneralAction(tableData, modelActionMetric, lowerSummary, dotZero, flags)
	}

	return metrics
}
