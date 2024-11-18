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

func printMetrics(action string, header string, model string, statsHeaderSuffix string, tableData [][]string) {
	shortHeader := ""
	fiber := "Currently"

	if action == "fiber-status" {
		if strings.Contains(header, fiber) {
			shortHeader = strings.Replace(header, fiber, "", 1)
			shortHeader = strings.Replace(shortHeader, "\u00A0\u00A0 ", "=", 1)
			shortHeader = strings.Replace(shortHeader, " ", ".", 1)
		}
	}

	actionMetric := strings.Replace(action, "-", ".", 2)
	actionMetric = strings.Replace(actionMetric, " ", ".", 1)
	actionMetric = actionMetric + "." + shortHeader

	if action == "fiber-status" {
		if strings.Contains(header, fiber) {
			metric := model + "." + actionMetric
			fmt.Println(strings.ToLower(metric))
		}
	}

	count := 1

	if action != "fiber-status" {
		if strings.HasSuffix(header, statsHeaderSuffix) {
			metric := ""
			for _, row := range tableData {
				if row[0] == "" {
					continue
				}

				length := len(row)

				stat := ""
				for i := range length {
					stat = row[0]
					stat = strings.Replace(stat, " ", ".", 1)
					if i != 0 {
						if length > 2 {
							portNumber := strconv.Itoa(count)
							metric = model + "." + actionMetric + "port" + portNumber + "." + stat + "="
						} else {
							metric = model + "." + actionMetric + stat + "="
						}
						if _, err := strconv.Atoi(row[i]); err == nil {
							row[i] = row[i] + ".0"
						}
						row[i] = metric + row[i] + "\n"
						if length > 2 {
							count = count + 1
							if count == length {
								count = 1
							}
						}
					}
				}

				line := ""

				for i, column := range row {
					if i != 0 {
						line = line + column
					}
				}

				fmt.Printf(strings.ToLower(line))
			}

		}
	}
}

func printData(action string, class string, currentHeader string, metrics bool, model string, pretty bool, tableData [][]string) {
	// If the table has data, process it
	if len(tableData) > 0 {
		// Process and print table data
		printTableData(action, class, currentHeader, metrics, model, pretty, tableData)
	}
}

func printTableData(action string, class string, header string, metrics bool, model string, pretty bool, tableData [][]string) {
	if !metrics {
		// Output the section header if it's available
		if header != "" {
			fmt.Printf("\n%s\n", header)
			fmt.Println(strings.Repeat("-", len(header)))
		}
	}

	statsHeaderSuffix := " Statistics"

	if metrics {
		printMetrics(action, header, model, statsHeaderSuffix, tableData)
	} else if action == "device-list" {
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
