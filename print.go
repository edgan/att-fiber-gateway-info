package main

import (
	"fmt"
	"strings"
)

// Function to print a row with padding
func printRowWithPadding(row []string, columnWidths []int) {
	for i, cell := range row {
		fmt.Print(cell + strings.Repeat(" ", columnWidths[i]-len(cell)+2))
	}

	fmt.Println()
}

func prettyPrint(tableData [][]string) {
	for _, row := range tableData {
		// Find the maximum width of each column
		columnWidths := make([]int, len(tableData[0]))

		for _, row := range tableData {
			for i, cell := range row {
				if len(cell) > columnWidths[i] {
					columnWidths[i] = len(cell)
				}
			}
		}

		printRowWithPadding(row, columnWidths)
	}
}

func printMetrics(metrics []string) {
	for _, m := range metrics {
		fmt.Println(strings.ToLower(m))
	}
}

func printData(action string, class string, currentHeader string, flags *Flags, model string, tableData [][]string) {
	// If the table has data, process it
	if len(tableData) > 0 {
		// Process and print table data
		printTableData(action, class, flags, currentHeader, model, tableData)
	}
}

func printTableData(action string, class string, flags *Flags, header string, model string, tableData [][]string) {
	// Output the section header if it's available
	if header != "" {
		fmt.Printf("\n%s\n", header)
		fmt.Println(strings.Repeat("-", len(header)))
	}

	if action == "device-list" {
		processDeviceList(tableData)
	} else if action == "home-network-status" {
		processHomeNetworkStatus(tableData)
	} else if action == "ip-allocation" {
		processIPAllocation(tableData)
	} else if action == "nat-check" || action == "nat-totals" {
		processNatCheckTotals(action, class, tableData)
	} else if action == "nat-connections" && !*flags.Pretty {
		processNatConnectionsNonPretty(class, tableData)
	} else if action == "nat-connections" && *flags.Pretty {
		processNatConnectionsPretty(class, tableData)
	} else if action == "nat-destinations" {
		processNatDestinations(class, tableData)
	} else if action == "nat-sources" {
		processNatSources(class, tableData)
	} else {
		processGeneric(tableData)
	}
}
