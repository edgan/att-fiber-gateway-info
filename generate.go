package main

import (
	"strings"
)

func generateFiberMetric(dotZero string, header string, modelActionMetric string, tableData [][]string) []string {
	metrics := []string{}
	fiber := "Currently"

	if strings.Contains(header, fiber) {
		keyValue := ""
		keyValue = strings.Replace(header, fiber, "", 1)
		keyValue = strings.Replace(keyValue, "\u00A0\u00A0 ", "=", 1)
		keyValue = strings.Replace(keyValue, " ", ".", 1)
		metric := modelActionMetric + "." + keyValue + dotZero
		metrics = append(metrics, metric)
	}

	return metrics
}

func generateNonFiberMetric(action, dotZero string, flags *Flags, modelActionMetric, summary string, tableData [][]string) []string {
	var metrics []string
	lowerSummary := strings.ToLower(strings.Replace(summary, " ", ".", 1))

	if action == "nat-totals" {
		metrics = processNatTotalsAction(tableData, modelActionMetric, dotZero, flags)
	} else {
		metrics = processGeneralAction(tableData, modelActionMetric, lowerSummary, dotZero, flags)
	}

	return metrics
}
