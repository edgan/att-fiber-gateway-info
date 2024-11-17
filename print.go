package main

import (
	"fmt"
	"os"
	"strconv"
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

func printData(action string, class string, currentHeader string, pretty bool, tableData [][]string) {
	// If the table has data, process it
	if len(tableData) > 0 {
		// Output the section header if it's available
		if currentHeader != "" {
			fmt.Printf("\n%s\n", currentHeader)
			fmt.Println(strings.Repeat("-", len(currentHeader)))
		}
		// Process and print table data
		printTableData(action, class, currentHeader, pretty, tableData)
	}
}

func printTableData(action string, class string, header string, pretty bool, tableData [][]string) {
	if action == "device-list" {
		for _, row := range tableData {
			count := len(row)

			// ipv4 address / name
			if row[0] == "IPv4 Address / Name" {
				row[0] = "IPv4 Address"
			}

			substring := " / "

			if count > 1 {
				if strings.Contains(row[1], substring) {
					row[1] = strings.Replace(row[1], substring, "Name: ", 1)
				}
			}

			line := strings.Join(row, ": ")

			// connection-type
			if row[0] == "Connection Type" {
				line = strings.Join(row, ": \n  ")
			}

			fmt.Println(line)
		}
	} else if action == "home-network-status" {
		for _, row := range tableData {
			count := len(row)
			if row[0] != "" && count > 1 {
				row[0] = row[0] + ":"
			}
		}

		prettyPrint(tableData)
	} else if action == "ip-allocation" {
		for _, row := range tableData {
			row[4] = ""
		}

		prettyPrint(tableData)
	} else if action == "nat-check" || action == "nat-totals" {
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

			fmt.Printf("Total number of tcp connections: %d\n", tcpCount)
			fmt.Printf("Total number of udp connections: %d\n", udpCount)
		}
	} else if action == "nat-connections" && !pretty {
		if class == "grid table100" {
			for _, row := range tableData {
				line := strings.Join(row, ", ")
				fmt.Println(line)
			}
		}
	} else if action == "nat-connections" && pretty {
		if class == "grid table100" {
			prettyPrint(tableData)
		}
	} else if action == "nat-destinations" {
		if class == "grid table100" {
			sortedDestinationsIPs := CountIPsByColumn(tableData, 7)
			fmt.Println("Destinations IP addresses:")

			for _, row := range sortedDestinationsIPs {
				fmt.Printf("%d %s\n", row.Count, row.IP)
			}
		}
	} else if action == "nat-sources" {
		if class == "grid table100" {
			sortedSourcesIPs := CountIPsByColumn(tableData, 5)
			fmt.Println("Source IP addresses:")

			for _, row := range sortedSourcesIPs {
				fmt.Printf("%d %s\n", row.Count, row.IP)
			}
		}
	} else {
		for _, row := range tableData {
			line := strings.Join(row, ": ")

			if !strings.Contains(line, "Legal Disclaimer") {
				fmt.Println(line)
			}
		}
	}
}
