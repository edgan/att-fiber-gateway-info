package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func extractHeadersAndTableData(action string, configs Configs, doc *goquery.Document, flags *Flags, model string, natActionPrefix string, returnFact string) string {
	fact := ""
	// Track current section header
	var currentHeader string

	// Process h1, h2, and tables
	doc.Find("h1, h2, table").Each(func(i int, s *goquery.Selection) {
		if !strings.Contains(action, natActionPrefix) {
			// Check if this is an h1 or h2 element
			if s.Is("h1") || s.Is("h2") {
				headerText := strings.TrimSpace(s.Text())
				if headerText != "" {
					currentHeader = headerText
				}
				return
			}
		}

		// Process tables
		if s.Is("table") {
			// Check for valid "class" attribute
			class, hasClass := s.Attr("class")
			validClasses := []string{"table60", "table75", "table100"}
			isValidTable := false

			if hasClass {
				for _, validClass := range validClasses {
					if strings.Contains(class, validClass) {
						isValidTable = true
						break
					}
				}
			}

			// Skip table if it doesn't match any valid class
			if !isValidTable {
				return
			}

			shortSummary := ""

			if action != "fiber-status" && returnFact == "" && *flags.Metrics {
				// Check for a valid "summary" attribute
				summary, hasSummary := s.Attr("summary")
				isValidSummary := false

				if hasSummary {
					debugLog(*flags.Debug, summary)
					// broadband-status / home-network-status, nat-totals
					if strings.Contains(strings.ToLower(summary), "statistics") || summary == "Summary of nattable connections" || summary == "This table displays a summary of session information." {
						re := regexp.MustCompile(` [Ss]tatistic.*`)
						shortSummary = re.ReplaceAllString(summary, "")
						re = regexp.MustCompile(`Ethernet `)
						shortSummary = re.ReplaceAllString(shortSummary, "")
						re = regexp.MustCompile(`This table displays `)
						shortSummary = re.ReplaceAllString(shortSummary, "")

						isValidSummary = true
					}
				}

				// Skip table if it doesn't match the summary filter
				if !isValidSummary {
					return
				}
			}

			// Extract rows from the table
			var tableData [][]string
			s.Find("tr").Each(func(j int, row *goquery.Selection) {
				var rowData []string

				// Extract "th" and "td" content
				row.Find("th, td").Each(func(k int, cell *goquery.Selection) {
					cellText := strings.TrimSpace(cell.Text())
					// Check if there is a <pre> tag inside the <td>
					pre := cell.Find("pre")

					if pre.Length() > 0 {
						// Process <pre> content with <br> tags replaced by newlines
						htmlContent, err := pre.Html() // Get HTML inside <pre>, which returns (string, error)

						if err != nil {
							cellText = "" // or handle error appropriately
						} else {
							cellText = strings.ReplaceAll(htmlContent, "Wi-Fi<br/>", "Wi-Fi: ") // Special case for WiFi
							cellText = strings.ReplaceAll(cellText, "<br/>", "\n  ")            // Replace <br /> with newline
							cellText = strings.ReplaceAll(cellText, "<br>", "\n")               // Handle <br> tag as well
						}
					} else {
						// No <pre> tag; get text normally
						cellText = strings.TrimSpace(cell.Text())
					}
					rowData = append(rowData, cellText)
				})

				// Add row data if not empty
				if len(rowData) > 0 {
					tableData = append(tableData, rowData)
				}
			})

			if returnFact != "" {
				if returnFact == "model" {
					fact = strings.Replace(tableData[1][1], "-", "", 1)
				}
			} else {
				if *flags.Metrics {
					debugLog(*flags.Debug, "outputMetrics")
					outputMetrics(action, configs, flags, currentHeader, model, shortSummary, tableData)
				} else {
					printData(action, class, currentHeader, flags, model, tableData)
				}
			}
		}
	})

	return fact
}

func extractData(action string, configs Configs, content string, flags *Flags, model string, natActionPrefix string, returnFact string) (string, error) {
	fact := ""

	// Load the HTML content into goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))

	if err != nil {
		debugLog(*flags.Debug, "Failed to goquery.NewDocumentFromReader in extractData")
		return fact, fmt.Errorf("failed to parse content: %v", err)
	}

	fact = extractHeadersAndTableData(action, configs, doc, flags, model, natActionPrefix, returnFact)

	return fact, err
}

func extractContentSub(htmlStr string) (string, error) {
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))

	if err != nil {
		return "", err
	}

	// Find the content-sub div
	contentSub := doc.Find("#content-sub").First()

	// Get the HTML of the content-sub div
	htmlContent, err := contentSub.Html()

	if err != nil {
		return "", err
	}

	return htmlContent, nil
}

// Define a function that takes a column index as a parameter and returns sorted IP counts
func CountIPsByColumn(tableData [][]string, column int) []struct {
	IP    string
	Count int
} {
	// Variable to count occurrences of each IP address in the specified column
	ipCount := make(map[string]int)

	for i, row := range tableData {
		if i != 0 {
			ipCount[row[column]]++
		}
	}

	// Convert the map to a slice of structs for sorting
	var sortedIPs []struct {
		IP    string
		Count int
	}
	for ip, count := range ipCount {
		sortedIPs = append(sortedIPs, struct {
			IP    string
			Count int
		}{IP: ip, Count: count})
	}

	// Sort the slice by count in descending order
	sort.Slice(sortedIPs, func(i, j int) bool {
		return sortedIPs[i].Count > sortedIPs[j].Count
	})

	return sortedIPs
}
