package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func extractHeadersAndTableData(action string, configs Configs, doc *goquery.Document, flags *Flags, model string, returnFact string) string {
	fact := ""
	currentHeader := ""
	actionPrefixes := returnActionPrefixes()
	validClasses := []string{"table60", "table75", "table100"}

	doc.Find("h1, h2, table").Each(func(i int, s *goquery.Selection) {
		// Update current header if applicable
		if !strings.Contains(action, actionPrefixes["nat"]) && (s.Is("h1") || s.Is("h2")) {
			headerText := strings.TrimSpace(s.Text())
			if headerText != "" {
				currentHeader = headerText
			}
			return
		}

		if s.Is("table") {
			// Validate table class
			class, hasClass := s.Attr("class")
			if !hasClass || !isValidTableClass(class, validClasses) {
				return
			}

			shortSummary := ""
			// Validate table summary if metrics are to be returned
			if action != "fiber-status" && returnFact == "" && *flags.Metrics {
				summary, hasSummary := s.Attr("summary")
				if !hasSummary {
					return
				}

				debugLog(*flags.Debug, summary)
				var isValid bool
				isValid, shortSummary = isValidSummary(summary)
				if !isValid {
					return
				}
			}

			// Extract table data
			tableData := extractTableData(s)

			// Handle return facts or output data
			if returnFact == "model" {
				fact = strings.Replace(tableData[1][1], "-", "", 1)
			} else if *flags.Metrics {
				debugLog(*flags.Debug, "outputMetrics")
				outputMetrics(action, configs, flags, currentHeader, model, shortSummary, tableData)
			} else {
				printData(action, class, currentHeader, flags, model, tableData)
			}
		}
	})

	return fact
}

// Helper function to validate table classes
func isValidTableClass(class string, validClasses []string) bool {
	for _, validClass := range validClasses {
		if strings.Contains(class, validClass) {
			return true
		}
	}
	return false
}

// Helper function to validate and process table summaries
func isValidSummary(summary string) (bool, string) {
	if strings.Contains(strings.ToLower(summary), "statistics") ||
		summary == "Summary of nattable connections" ||
		summary == "This table displays a summary of session information." {

		shortSummary := summary
		patterns := []string{
			` [Ss]tatistic.*`,
			`Ethernet `,
			`This table displays `,
		}
		for _, pattern := range patterns {
			re := regexp.MustCompile(pattern)
			shortSummary = re.ReplaceAllString(shortSummary, "")
		}
		return true, shortSummary
	}
	return false, ""
}

// Helper function to extract table data
func extractTableData(s *goquery.Selection) [][]string {
	var tableData [][]string
	s.Find("tr").Each(func(j int, row *goquery.Selection) {
		var rowData []string
		row.Find("th, td").Each(func(k int, cell *goquery.Selection) {
			cellText := extractCellText(cell)
			rowData = append(rowData, cellText)
		})
		if len(rowData) > 0 {
			tableData = append(tableData, rowData)
		}
	})
	return tableData
}

// Helper function to extract and process cell text
func extractCellText(cell *goquery.Selection) string {
	pre := cell.Find("pre")
	if pre.Length() > 0 {
		htmlContent, err := pre.Html()
		if err != nil {
			return ""
		}
		cellText := htmlContent
		replacements := map[string]string{
			"Wi-Fi<br/>": "Wi-Fi: ",
			"<br/>":      "\n  ",
			"<br>":       "\n",
		}
		for old, new := range replacements {
			cellText = strings.ReplaceAll(cellText, old, new)
		}
		return cellText
	}
	return strings.TrimSpace(cell.Text())
}

func extractData(action string, configs Configs, content string, flags *Flags, model string, returnFact string) (string, error) {
	fact := ""

	// Load the HTML content into goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))

	if err != nil {
		debugLog(*flags.Debug, "Failed to goquery.NewDocumentFromReader in extractData")
		return fact, fmt.Errorf("failed to parse content: %v", err)
	}

	fact = extractHeadersAndTableData(action, configs, doc, flags, model, returnFact)

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
