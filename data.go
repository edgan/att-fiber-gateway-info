package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func extractHeadersAndTableData(action string, doc *goquery.Document, filter string, pretty bool) {
	// Track current section header
	var currentHeader string

	// Process h1, h2, and tables
	doc.Find("h1, h2, table").Each(func(i int, s *goquery.Selection) {
		if !strings.Contains(action, "nat-") {
			// Check if this is an h1 or h2 element
			if s.Is("h1") || s.Is("h2") {
				headerText := strings.TrimSpace(s.Text())
				if headerText != "" {
					currentHeader = headerText
				}
				return
			}
		}

		// Process tables, specifically looking for "grid table100"
		if s.Is("table") {
			validClasses := []string{"table60", "table75", "table100"}
			if class, exists := s.Attr("class"); exists {
				for _, validClass := range validClasses {
					if strings.Contains(class, validClass) {
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

						printData(action, class, currentHeader, pretty, tableData)
					}
				}
			}
		}
	})
}

func extractData(action string, content string, filter string, natActionPrefix string, pretty bool) error {
	// Load the HTML content into goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))

	if err != nil {
		return fmt.Errorf("failed to parse content: %v", err)
	}

	extractHeadersAndTableData(action, doc, filter, pretty)

	return nil
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
