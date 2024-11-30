//revive:disable:add-constant
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"edgan/att-fiber-gateway-info/internal/logging"
)

func processDeviceList(tableData [][]string) {
	for _, row := range tableData {
		count := len(row)

		// ipv4 address / name
		if row[keyRow] == "IPv4 Address / Name" {
			row[keyRow] = "IPv4 Address"
		}

		substring := " / "

		if count > one {
			if strings.Contains(row[one], substring) {
				row[one] = strings.Replace(row[one], substring, "Name: ", oneOccurance)
			}
		}

		line := strings.Join(row, ": ")

		// connection-type
		if row[keyRow] == "Connection Type" {
			line = strings.Join(row, ": \n  ")
		}

		fmt.Println(line)
	}
}

func processGeneric(tableData [][]string) {
	// Call prettyPrint with stripAnsiCodes = true and specialFormatting = true
	prettyPrint(tableData, true, true)
}

func processHomeNetworkStatus(tableData [][]string) {
	for _, row := range tableData {
		count := len(row)
		if row[keyRow] != empty && count > 1 {
			row[keyRow] = row[keyRow] + colon
		}
	}

	prettyPrint(tableData, false, false)
}

func processIPAllocation(tableData [][]string) {
	for _, row := range tableData {
		row[columnNumnerOfAction] = empty
	}

	prettyPrint(tableData, false, false)
}

func processNatTotals(tableData [][]string) (icmpCount int, tcpCount int, udpCount int) {
	for _, row := range tableData {
		protocol := row[one]

		if protocol == "icmp" {
			icmpCount++
		}
		if protocol == "tcp" {
			tcpCount++
		}
		if protocol == "udp" {
			udpCount++
		}
	}

	return icmpCount, tcpCount, udpCount
}

func processNatCheck(value string) {
	fmt.Printf("%s.0\n", value)

	connections, err := strconv.Atoi(value)

	if err != nil {
		panic(err)
	}

	if connections >= maxConnections {
		fmt.Printf("\nError: Too many connections\n")
		os.Exit(1)
	}
}

func processNatCheckTotals(action string, class string, tableData [][]string) {
	if class == "table60" {
		for _, row := range tableData {
			if len(row) == 0 || len(row) < 2 {
				continue
			}

			if row[keyRow] == "Total sessions in use" {
				if action == "nat-check" {
					processNatCheck(row[valueRow])
				} else if action == "nat-totals" {
					printNatTotals(row[valueRow])
				}
			}
		}
	}
}

func processNatConnectionsNonPretty(class string, tableData [][]string) {
	if class == gridTable100 {
		for _, row := range tableData {
			line := strings.Join(row, commaSpace)
			fmt.Println(line)
		}
	}
}

func processNatConnectionsPretty(class string, tableData [][]string) {
	if class == gridTable100 {
		prettyPrint(tableData, false, false)
	}
}

func processNatDestinations(class string, tableData [][]string) {
	if class == gridTable100 {
		sortedDestinationsIPs := countIPsByColumn(tableData, columnNumberOfDestinationIPs)
		fmt.Println("Destinations IP addresses:")

		for _, row := range sortedDestinationsIPs {
			fmt.Printf("%d %s\n", row.Count, row.IP)
		}
	}
}

func processNatSources(class string, tableData [][]string) {
	if class == gridTable100 {
		sortedSourcesIPs := countIPsByColumn(tableData, columnNumberOfSourceIPs)
		fmt.Println("Source IP addresses:")

		for _, row := range sortedSourcesIPs {
			fmt.Printf("%d %s\n", row.Count, row.IP)
		}
	}
}

// Processes the "IP Family" case
func processIPFamilyCase(tableData [][]string, modelActionMetric, dotZero string) []string {
	var metrics []string

	icmpCount, tcpCount, udpCount := processNatTotals(tableData)

	metrics = append(metrics,
		fmt.Sprintf("%s.icmp.connections=%d%s", modelActionMetric, icmpCount, dotZero),
		fmt.Sprintf("%s.tcp.connections=%d%s", modelActionMetric, tcpCount, dotZero),
		fmt.Sprintf("%s.udp.connections=%d%s", modelActionMetric, udpCount, dotZero),
	)

	return metrics
}

func processTotalConnections(key, value, modelActionMetric, dotZero string, flags *flags) string {
	stat := processStat(key)

	if strings.Contains(stat, "sessions.in.use") {
		stat = strings.Replace(stat, "Total.sessions.in.use", "connections", oneOccurance)
	}

	returnedValue := processValue(value, flags.Noconvert, dotZero)
	metric := fmt.Sprintf("%s.%s=%s", modelActionMetric, stat, returnedValue)
	return metric
}

// Helper function to process "nat-totals" action
func processNatTotalsAction(tableData [][]string, modelActionMetric, dotZero string, flags *flags) []string {
	var metrics []string

	for _, row := range tableData {
		if len(row) == 0 || len(row) < 2 {
			continue
		}

		if row[keyRow] == empty {
			continue
		}

		switch {
		case strings.Contains(row[keyRow], "IP Family"):
			metrics = append(metrics, processIPFamilyCase(tableData, modelActionMetric, dotZero)...)
			return metrics

		case row[keyRow] == "Total sessions available", row[keyRow] == "Select display option":
			continue

		default:
			metric := processTotalConnections(row[keyRow], row[valueRow], modelActionMetric, dotZero, flags)
			metrics = append(metrics, metric)
		}
	}

	return metrics
}

// Helper function to process general actions without nested for loops
func processGeneralAction(
	tableData [][]string, modelActionMetric, lowerSummary, dotZero string, flags *flags,
) []string {
	var metrics []string
	var rowsToProcess []struct {
		Stat           string
		IsPortSpecific bool
		Cells          []string
	}

	// First phase: Collect necessary data from tableData
	for _, row := range tableData {
		if len(row) == 0 || row[0] == empty {
			continue
		}

		key := row[1]
		stat := processStat(key)
		isPortSpecific := len(row) > 2
		cells := row[1:]

		rowsToProcess = append(rowsToProcess, struct {
			Stat           string
			IsPortSpecific bool
			Cells          []string
		}{
			Stat:           stat,
			IsPortSpecific: isPortSpecific,
			Cells:          cells,
		})
	}

	// Second phase: Process the collected data to generate metrics
	for _, item := range rowsToProcess {
		indices := make([]int, len(item.Cells))
		for i := range item.Cells {
			indices[i] = i
		}
		metrics = append(metrics, processCells(item, modelActionMetric, lowerSummary, dotZero, flags, indices)...)
	}

	return metrics
}

// Helper function to process cells for each row
func processCells(item struct {
	Stat           string
	IsPortSpecific bool
	Cells          []string
}, modelActionMetric, lowerSummary, dotZero string, flags *flags, indices []int) []string {
	var metrics []string

	for _, i := range indices {
		cell := item.Cells[i]
		returnedValue := processValue(cell, flags.Noconvert, dotZero)
		var metric string
		if item.IsPortSpecific {
			// Port-specific metrics
			metric = fmt.Sprintf(
				"%s.%s.port%d.%s=%s", modelActionMetric, lowerSummary, i+1, item.Stat, returnedValue,
			)
		} else {
			metric = fmt.Sprintf("%s.%s.%s=%s", modelActionMetric, lowerSummary, item.Stat, returnedValue)
		}
		metrics = append(metrics, metric)
	}

	return metrics
}

// Helper function to process stat strings
func processStat(stat string) string {
	stat = strings.Replace(stat, space, period, oneOccurance)
	stat = strings.Replace(stat, " (Mbps)", empty, oneOccurance)
	return stat
}

// Helper function to process value strings
func processValue(value string, noconvert *bool, dotZero string) string {
	value = strings.ToLower(value)
	if !*noconvert {
		switch value {
		case "down":
			value = "0"
		case "half", "up":
			value = "1"
		case "full":
			value = "2"
		}
	}
	if _, err := strconv.Atoi(value); err == nil {
		value += dotZero
	}
	return value
}

func processDatadogMetrics(metrics []string) (floatMetrics map[string]float64) {
	for _, metric := range metrics {
		metric = strings.ToLower(strings.TrimSpace(metric))
		splitMetric := strings.Split(metric, equals)

		if len(splitMetric) != 2 {
			logging.LogFatalf("Invalid metric format:", metric)
		}

		if strings.HasSuffix(splitMetric[valueRow], ".0") {
			valueF, err := strconv.ParseFloat(splitMetric[valueRow], floatPrecision)

			if err != nil {
				logging.LogFatalf("Error:", err)
			}

			floatMetrics[splitMetric[keyRow]] = valueF
		}
	}

	return floatMetrics
}
