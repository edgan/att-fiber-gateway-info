package main

import (
	"fmt"
	"strconv"
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

func generateNonFiberMetric(action string, dotZero string, flags *Flags, modelActionMetric string, summary string, tableData [][]string) []string {
	metrics := make([]string, 0, len(tableData)*2) // Preallocate for efficiency
	port := 1
	tcpCount, udpCount := 0, 0

	lowerSummary := strings.ToLower(strings.Replace(summary, " ", ".", 1)) // Precompute summary once

	for _, row := range tableData {
		if len(row) == 0 || row[0] == "" {
			continue
		}

		if action == "nat-totals" {
			if strings.Contains(row[0], "IP Family") {
				tcpCount, udpCount = processNatTotals(tableData)
				metrics = append(metrics,
					fmt.Sprintf("%s.tcp.connections=%d%s", modelActionMetric, tcpCount, dotZero),
					fmt.Sprintf("%s.udp.connections=%d%s", modelActionMetric, udpCount, dotZero),
				)
				break
			}

			if row[0] == "Total sessions available" || row[0] == "Select display option" {
				continue
			}
		}

		stat := strings.Replace(row[0], " ", ".", 1)
		stat = strings.Replace(stat, " (Mbps)", "", 1)

		if action == "nat-totals" && strings.Contains(stat, "sessions in use") {
			stat = strings.Replace(stat, "Total sessions in use", "connections", 1)
		}

		for i := 1; i < len(row); i++ { // Start at 1 since stat is from row[0]
			value := row[i]
			if _, err := strconv.Atoi(value); err == nil {
				value += dotZero
			}

			var metric string
			if len(row) > 2 { // Handle port-specific metrics
				metric = fmt.Sprintf("%s.%s.port%d.%s=%s", modelActionMetric, lowerSummary, port, stat, value)
				port++
				if port == len(row) {
					port = 1
				}
			} else if action == "nat-totals" {
				metric = fmt.Sprintf("%s.%s=%s", modelActionMetric, stat, value)
			} else {
				metric = fmt.Sprintf("%s.%s.%s=%s", modelActionMetric, lowerSummary, stat, value)
			}

			metrics = append(metrics, metric)
		}
	}

	return metrics
}
