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
	for _, row := range tableData {
		line := strings.Join(row, ": ")

		if !strings.Contains(line, "Legal Disclaimer") {
			fmt.Println(line)
		}
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

func processNatTotals(tableData [][]string) (int, int) {
	// Initialize counters
	var tcpCount, udpCount int

	for _, row := range tableData {
		if row[1] == "tcp" {
			tcpCount++
		}
		if row[1] == "udp" {
			udpCount++
		}
	}

	return tcpCount, udpCount
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
		tcpCount, udpCount := processNatTotals(tableData)
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
		sortedDestinationsIPs := CountIPsByColumn(tableData, 7)
		fmt.Println("Destinations IP addresses:")

		for _, row := range sortedDestinationsIPs {
			fmt.Printf("%d %s\n", row.Count, row.IP)
		}
	}
}

func processNatSources(class string, tableData [][]string) {
	if class == "grid table100" {
		sortedSourcesIPs := CountIPsByColumn(tableData, 5)
		fmt.Println("Source IP addresses:")

		for _, row := range sortedSourcesIPs {
			fmt.Printf("%d %s\n", row.Count, row.IP)
		}
	}
}
