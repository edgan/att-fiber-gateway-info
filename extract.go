//revive:disable:add-constant
package main

import (
	"fmt"
	"strings"

	"edgan/att-fiber-gateway-info/internal/logging"

	"github.com/PuerkitoBio/goquery"
)

// DataContext holds the context data for processing
type DataContext struct {
	action        string
	configs       configs
	flags         *flags
	model         string
	returnFact    string
	currentHeader string
	shortSummary  string
	class         string
	fact          *string
}

// Handles return facts or outputs data based on the flags and action
func handleReturnFactsOrOutputData(ctx *DataContext, tableData [][]string) {
	if ctx.returnFact == "model" {
		if len(tableData) > 1 && len(tableData[1]) > 1 {
			*ctx.fact = strings.Replace(tableData[1][1], dash, empty, oneOccurance)
		}
	} else if *ctx.flags.Metrics {
		logging.DebugLog(*ctx.flags.Debug, "outputMetrics")
		outputMetrics(ctx.action, ctx.configs, ctx.flags, ctx.currentHeader, ctx.model, ctx.shortSummary, tableData)
	} else {
		printData(ctx.action, ctx.class, ctx.currentHeader, ctx.flags, tableData)
	}
}

// Extracts headers and table data from the document
func extractHeadersAndTableData(
	action string, configs configs, doc *goquery.Document, flags *flags, model string, returnFact string,
) string {
	fact := empty
	currentHeader := empty
	validClasses := []string{"table60", "table75", "table100"}

	// Initialize the context
	ctx := DataContext{
		action:     action,
		configs:    configs,
		flags:      flags,
		model:      model,
		returnFact: returnFact,
		fact:       &fact,
	}

	doc.Find("h1, h2, table").Each(func(_ int, s *goquery.Selection) {
		// Update current header if applicable
		if updateCurrentHeader(action, s, &currentHeader) {
			return
		}

		if s.Is("table") {
			// Validate table class
			class, hasClass := s.Attr("class")
			if !hasClass || !isValidTableClass(class, validClasses) {
				return
			}

			shortSummary := empty

			// Validate table summary if metrics are to be returned
			if !validateTableSummary(action, returnFact, flags, s, &shortSummary) {
				return
			}

			// Extract table data
			tableData := extractTableData(s)

			// Update context with current header, short summary, and class
			ctx.currentHeader = currentHeader
			ctx.shortSummary = shortSummary
			ctx.class = class

			// Handle return facts or output data
			handleReturnFactsOrOutputData(&ctx, tableData)
		}
	})

	return fact
}

// Updates the current header if applicable
func updateCurrentHeader(action string, s *goquery.Selection, currentHeader *string) bool {
	if !strings.Contains(action, natActionPrefix) && (s.Is("h1") || s.Is("h2")) {
		headerText := strings.TrimSpace(s.Text())
		if headerText != empty {
			*currentHeader = headerText
		}
		return true
	}
	return false
}

// Validates the table summary if metrics are to be returned
func validateTableSummary(
	action string, returnFact string, flags *flags, s *goquery.Selection, shortSummary *string,
) bool {
	if action != "fiber-status" && returnFact == empty && *flags.Metrics {
		summary, hasSummary := s.Attr("summary")
		if !hasSummary {
			return false
		}

		logging.DebugLog(*flags.Debug, summary)
		var isValid bool
		isValid, *shortSummary = isValidSummary(summary)
		if !isValid {
			return false
		}
	}
	return true
}

// Modify the extractTableData function to include table headers
func extractTableData(s *goquery.Selection) [][]string {
	var tableData [][]string
	s.Find("tr").Each(func(_ int, row *goquery.Selection) {
		var rowData []string

		// Extract "th" and "td" content
		row.Find("th, td").Each(func(_ int, cell *goquery.Selection) {
			cellText := strings.TrimSpace(cell.Text())
			rowData = append(rowData, cellText)
		})

		// Add row data if not empty
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
			return empty
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

func extractData(
	action string, configs configs, content string, flags *flags, model string, returnFact string,
) (string, error) {
	fact := empty

	// Load the HTML content into goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))

	if err != nil {
		logging.DebugLog(*flags.Debug, "Failed to goquery.NewDocumentFromReader in extractData")
		return fact, fmt.Errorf("failed to parse content: %v", err)
	}

	fact = extractHeadersAndTableData(action, configs, doc, flags, model, returnFact)

	return fact, err
}

func extractContentSub(htmlStr string) (string, error) {
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))

	if err != nil {
		return empty, err
	}

	// Find the content-sub div
	contentSub := doc.Find("#content-sub").First()

	// Get the HTML of the content-sub div
	htmlContent, err := contentSub.Html()

	if err != nil {
		return empty, err
	}

	return htmlContent, nil
}
