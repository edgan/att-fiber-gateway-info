package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func processDeviceList(tableData [][]string) {
	for _, row := range tableData {
		count := len(row)

		// ipv4 address / name
		if row[0] == "IPv4 Address / Name" {
			row[0] = "IPv4 Address"
		}

		substring := " / "

		if count > 1 {
			row[1] = strings.Replace(row[1], substring, "Name: ", 1)
		}

		line := strings.Join(row, ": ")

		// connection-type
		if row[0] == "Connection Type" {
			line = strings.Join(row, ": \n  ")
		}

		fmt.Println(line)
	}
}

func processGeneric(tableData [][]string) {
	if len(tableData) == 0 {
		return
	}

	// Determine the maximum number of columns
	numCols := 0
	for _, row := range tableData {
		if len(row) > numCols {
			numCols = len(row)
		}
	}

	// Initialize slice to hold max width of each column
	colWidths := make([]int, numCols)

	// Calculate maximum width for each column
	for _, row := range tableData {
		for i, cell := range row {
			cellLen := len(stripAnsi(cell))
			if i == 0 && numCols == 2 {
				cellLen += 1 // Account for the ":" after the key
			}
			if cellLen > colWidths[i] {
				colWidths[i] = cellLen
			}
		}
	}

	// Print each row with proper alignment
	for _, row := range tableData {
		// For tables with two columns, add ":" after the first column
		if numCols == 2 && len(row) >= 1 {
			if strings.Contains(row[0], "Legal Disclaimer") {
				continue
			}

			key := row[0] + ":"

			if row[0] == "" {
				key = row[0]
			}

			format0 := fmt.Sprintf("%%-%ds", colWidths[0]+2)
			fmt.Printf(format0, key)

			// Check if the second column exists
			if len(row) >= 2 {
				format1 := fmt.Sprintf("%%-%ds", colWidths[1]+2)
				fmt.Printf(format1, row[1])
			}
		} else {
			// For tables with more than two columns
			for i := 0; i < numCols; i++ {
				var cell string
				if i < len(row) {
					cell = row[i]
				} else {
					cell = ""
				}
				// Left-align the content within the column width
				format := fmt.Sprintf("%%-%ds", colWidths[i]+2) // Add extra space for padding
				fmt.Printf(format, cell)
			}
		}
		fmt.Println()
	}
}

func processHomeNetworkStatus(tableData [][]string) {
	for _, row := range tableData {
		count := len(row)
		if row[0] != "" && count > 1 {
			row[0] = row[0] + ":"
		}
	}

	prettyPrint(tableData)
}

func processIPAllocation(tableData [][]string) {
	for _, row := range tableData {
		row[4] = ""
	}

	prettyPrint(tableData)
}

func processNatTotals(tableData [][]string) (int, int, int) {
	// Initialize counters
	var icmpCount, tcpCount, udpCount int

	for _, row := range tableData {
		if row[1] == "icmp" {
			icmpCount++
		}
		if row[1] == "tcp" {
			tcpCount++
		}
		if row[1] == "udp" {
			udpCount++
		}
	}

	return icmpCount, tcpCount, udpCount
}

func processNatCheckTotals(action string, class string, tableData [][]string) {
	if class == "table60" {
		for _, row := range tableData {
			if len(row) > 0 && row[0] == "Total sessions in use" {
				if action == "nat-check" {
					fmt.Printf("%s.0\n", row[1])

					maxConnections := 8192
					connections, err := strconv.Atoi(row[1])

					if err != nil {
						panic(err)
					}

					if connections >= maxConnections {
						fmt.Printf("\nError: Too many connections\n")
						os.Exit(1)
					}
				}

				if action == "nat-totals" {
					fmt.Printf("%s: %s\n", "Total number of connections", row[1])
				}
			}
		}
	}

	if action == "nat-totals" && class == "grid table100" {
		icmpCount, tcpCount, udpCount := processNatTotals(tableData)
		fmt.Printf("Total number of icmp connections: %d\n", icmpCount)
		fmt.Printf("Total number of tcp connections: %d\n", tcpCount)
		fmt.Printf("Total number of udp connections: %d\n", udpCount)
	}
}

func processNatConnectionsNonPretty(class string, tableData [][]string) {
	if class == "grid table100" {
		for _, row := range tableData {
			line := strings.Join(row, ", ")
			fmt.Println(line)
		}
	}
}

func processNatConnectionsPretty(class string, tableData [][]string) {
	if class == "grid table100" {
		prettyPrint(tableData)
	}
}

func processNatDestinations(class string, tableData [][]string) {
	if class == "grid table100" {
		sortedDestinationsIPs := countIPsByColumn(tableData, 7)
		fmt.Println("Destinations IP addresses:")

		for _, row := range sortedDestinationsIPs {
			fmt.Printf("%d %s\n", row.Count, row.IP)
		}
	}
}

func processNatSources(class string, tableData [][]string) {
	if class == "grid table100" {
		sortedSourcesIPs := countIPsByColumn(tableData, 5)
		fmt.Println("Source IP addresses:")

		for _, row := range sortedSourcesIPs {
			fmt.Printf("%d %s\n", row.Count, row.IP)
		}
	}
}

// Helper function to process "nat-totals" action
func processNatTotalsAction(tableData [][]string, modelActionMetric, dotZero string, flags *Flags) []string {
	var metrics []string
	for _, row := range tableData {
		if len(row) == 0 || row[0] == "" {
			continue
		}

		switch {
		case strings.Contains(row[0], "IP Family"):
			icmpCount, tcpCount, udpCount := processNatTotals(tableData)
			metrics = append(metrics,
				fmt.Sprintf("%s.icmp.connections=%d%s", modelActionMetric, icmpCount, dotZero),
				fmt.Sprintf("%s.tcp.connections=%d%s", modelActionMetric, tcpCount, dotZero),
				fmt.Sprintf("%s.udp.connections=%d%s", modelActionMetric, udpCount, dotZero),
			)
			return metrics

		case row[0] == "Total sessions available", row[0] == "Select display option":
			continue

		default:
			stat := processStat(row[0])
			if strings.Contains(stat, "sessions.in.use") {
				stat = strings.Replace(stat, "Total.sessions.in.use", "connections", 1)
			}
			value := processValue(row[1], flags.Noconvert, dotZero)
			metric := fmt.Sprintf("%s.%s=%s", modelActionMetric, stat, value)
			metrics = append(metrics, metric)
		}
	}
	return metrics
}

// Helper function to process general actions
func processGeneralAction(tableData [][]string, modelActionMetric, lowerSummary, dotZero string, flags *Flags) []string {
	var metrics []string
	for _, row := range tableData {
		if len(row) == 0 || row[0] == "" {
			continue
		}

		stat := processStat(row[0])

		for i, cell := range row[1:] {
			value := processValue(cell, flags.Noconvert, dotZero)
			var metric string
			if len(row) > 2 {
				// Port-specific metrics
				metric = fmt.Sprintf("%s.%s.port%d.%s=%s", modelActionMetric, lowerSummary, i+1, stat, value)
			} else {
				metric = fmt.Sprintf("%s.%s.%s=%s", modelActionMetric, lowerSummary, stat, value)
			}
			metrics = append(metrics, metric)
		}
	}
	return metrics
}

// Helper function to process stat strings
func processStat(stat string) string {
	stat = strings.Replace(stat, " ", ".", 1)
	stat = strings.Replace(stat, " (Mbps)", "", 1)
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
