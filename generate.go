package main

import (
	//"fmt"
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
	metrics := []string{}
	metric := ""

	port := 1
	tcpCount := 0
	udpCount := 0

	for _, row := range tableData {
		if action == "nat-totals" {
			if strings.Contains(row[0], "IP Family") {
				tcpCount, udpCount = processNatTotals(tableData)
				metric = modelActionMetric + "." + "tcp.connections" + "=" + strconv.Itoa(tcpCount) + dotZero
				metrics = append(metrics, metric)
				metric = modelActionMetric + "." + "udp.connections" + "=" + strconv.Itoa(udpCount) + dotZero
				metrics = append(metrics, metric)
				break
			}

			if row[0] == "Total sessions available" || row[0] == "Select display option" {
				continue
			}
		}

		if row[0] == "" {
			continue
		}

		length := len(row)

		stat := ""
		for i := range length {
			stat = row[0]
			if action == "nat-totals" {
				if strings.Contains(stat, "sessions in use") {
					stat = strings.Replace(stat, "Total sessions in use", "connetions", 1)
				}
			}

			stat = strings.Replace(stat, " ", ".", 1)
			stat = strings.Replace(stat, " (Mbps)", "", 1)

			if i != 0 {
				if length > 2 {
					portNumber := strconv.Itoa(port)
					summary = strings.ToLower(summary)
					summary = strings.Replace(summary, " ", ".", 1)
					metric = modelActionMetric + "." + summary + "." + "port" + portNumber + "." + stat + "="
				} else {
					if action == "nat-totals" {
						metric = modelActionMetric + "." + stat + "="
					} else {
						metric = modelActionMetric + "." + summary + "." + stat + "="
					}
				}

				value := row[i]

				if _, err := strconv.Atoi(value); err == nil {
					value = value + dotZero
				}

				metric = metric + value
				metrics = append(metrics, metric)

				if length > 2 {
					port = port + 1
					if port == length {
						port = 1
					}
				}
			}
		}
	}

	return metrics
}
