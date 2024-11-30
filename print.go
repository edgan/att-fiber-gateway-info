package main

import (
	"fmt"
	"strings"
)

// Formats a special two-column row with special formatting
func formatSpecialTwoColumnRow(row []string, columnWidths []int) {
	key := row[zero]
	value := empty

	if len(row) >= two {
		value = row[one]
	}

	modifiedKey := key
	if key != empty {
		modifiedKey = key + colon
	}

	format0 := fmt.Sprintf(columnFormat, columnWidths[zero]+two)
	fmt.Printf(format0, modifiedKey)

	format1 := fmt.Sprintf(columnFormat, columnWidths[one]+two)
	fmt.Printf(format1, value)
}

// Formats a general row for tables with more than two columns or when special formatting is not needed
func formatGeneralRow(row []string, columnWidths []int, numCols int) {
	for i := zero; i < numCols; i++ {
		var cell string
		if i < len(row) {
			cell = row[i]
		} else {
			cell = empty
		}
		format := fmt.Sprintf(columnFormat, columnWidths[i]+two)
		fmt.Printf(format, cell)
	}
}

// Function to print a row with padding and handle special cases
func printRowWithPadding(row []string, columnWidths []int, numCols int, specialFormatting bool) {
	// Skip rows containing "Legal Disclaimer"
	if len(row) > zero && strings.Contains(row[zero], "Legal Disclaimer") {
		return
	}

	if specialFormatting && numCols == two && len(row) >= one {
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

		if i == zero && numCols == two && specialFormatting {
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
	if len(tableData) == zero {
		return
	}

	// Determine the maximum number of columns
	numCols := zero

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
	if action == "nat-totals" && class == gridTable100 {
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
	if len(tableData) > zero {
		// Process and print table data
		printTableData(action, class, flags, currentHeader, tableData)
	}
}

func printTableData(action string, class string, flags *flags, header string, tableData [][]string) {
	// Output the section header if it's available
	if header != empty {
		fmt.Printf("\n%s\n", header)
		fmt.Println(strings.Repeat(dash, len(header)))
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
