package main

import (
	"fmt"
	"strings"
)

// Formats a special two-column row with special formatting
func formatSpecialTwoColumnRow(row []string, columnWidths []int) {
	key := row[0]
	value := ""

	if len(row) >= 2 {
		value = row[1]
	}

	modifiedKey := key
	if key != "" {
		modifiedKey = key + ":"
	}

	format0 := fmt.Sprintf("%%-%ds", columnWidths[0]+2)
	fmt.Printf(format0, modifiedKey)

	format1 := fmt.Sprintf("%%-%ds", columnWidths[1]+2)
	fmt.Printf(format1, value)
}

// Formats a general row for tables with more than two columns or when special formatting is not needed
func formatGeneralRow(row []string, columnWidths []int, numCols int) {
	for i := 0; i < numCols; i++ {
		var cell string
		if i < len(row) {
			cell = row[i]
		} else {
			cell = ""
		}
		format := fmt.Sprintf("%%-%ds", columnWidths[i]+2)
		fmt.Printf(format, cell)
	}
}

// Function to print a row with padding and handle special cases
func printRowWithPadding(row []string, columnWidths []int, numCols int, specialFormatting bool) {
	// Skip rows containing "Legal Disclaimer"
	if len(row) > 0 && strings.Contains(row[0], "Legal Disclaimer") {
		return
	}

	if specialFormatting && numCols == 2 && len(row) >= 1 {
		formatSpecialTwoColumnRow(row, columnWidths)
	} else {
		formatGeneralRow(row, columnWidths, numCols)
	}

	fmt.Println()
}

// Updates the column widths based on a single row
func updateColumnWidths(row []string, columnWidths []int, numCols int, stripAnsiCodes bool, specialFormatting bool) {
	for i, cell := range row {
		cellContent := cell

		if stripAnsiCodes {
			cellContent = stripAnsi(cell)
		}

		cellLen := len(cellContent)

		if i == 0 && numCols == 2 && specialFormatting {
			cellLen++ // Account for the added colon
		}

		if cellLen > columnWidths[i] {
			columnWidths[i] = cellLen
		}
	}
}

// Calculates the maximum width for each column
func calculateColumnWidths(tableData [][]string, numCols int, stripAnsiCodes bool, specialFormatting bool) []int {
	// Initialize slice to hold max width of each column
	columnWidths := make([]int, numCols)

	// Calculate maximum width for each column
	for _, row := range tableData {
		updateColumnWidths(row, columnWidths, numCols, stripAnsiCodes, specialFormatting)
	}
	return columnWidths
}

func prettyPrint(tableData [][]string, stripAnsiCodes bool, specialFormatting bool) {
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

	// Calculate column widths using the new function
	columnWidths := calculateColumnWidths(tableData, numCols, stripAnsiCodes, specialFormatting)

	// Now, iterate over tableData and print each row with padding
	for _, row := range tableData {
		printRowWithPadding(row, columnWidths, numCols, specialFormatting)
	}
}

func printNatConnectionTotals(action string, class string, tableData [][]string) {
	if action == "nat-totals" && class == "grid table100" {
		icmpCount, tcpCount, udpCount := processNatTotals(tableData)
		fmt.Printf("Total number of icmp connections: %d\n", icmpCount)
		fmt.Printf("Total number of tcp connections: %d\n", tcpCount)
		fmt.Printf("Total number of udp connections: %d\n", udpCount)
	}
}

// Processes the "nat-totals" action
func printNatTotals(value string) {
	fmt.Printf("%s: %s\n", "Total number of connections", value)
}

func printMetrics(metrics []string) {
	for _, m := range metrics {
		fmt.Println(strings.ToLower(m))
	}
}

func printData(action string, class string, currentHeader string, flags *flags, tableData [][]string) {
	// If the table has data, process it
	if len(tableData) > 0 {
		// Process and print table data
		printTableData(action, class, flags, currentHeader, tableData)
	}
}

func printTableData(action string, class string, flags *flags, header string, tableData [][]string) {
	// Output the section header if it's available
	if header != "" {
		fmt.Printf("\n%s\n", header)
		fmt.Println(strings.Repeat("-", len(header)))
	}

	switch action {
	case "device-list":
		processDeviceList(tableData)
	case "home-network-status":
		processHomeNetworkStatus(tableData)
	case "ip-allocation":
		processIPAllocation(tableData)
	case "nat-check":
		processNatCheckTotals(action, class, tableData)
	case "nat-totals":
		processNatCheckTotals(action, class, tableData)
		printNatConnectionTotals(action, class, tableData)
	case "nat-connections":
		if !*flags.Pretty {
			processNatConnectionsNonPretty(class, tableData)
		} else {
			processNatConnectionsPretty(class, tableData)
		}
	case "nat-destinations":
		processNatDestinations(class, tableData)
	case "nat-sources":
		processNatSources(class, tableData)
	default:
		processGeneric(tableData)
	}
}
