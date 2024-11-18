package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/fatih/color"
)

func returnFlags(actionDescription string, colorMode bool, cookiePath string, filterDescription string) (*string, *bool, *bool, *bool, *string, *string, *bool, *string, *bool, *bool, *string, *bool) {
	action := flag.String("action", "", actionDescription)
	allmetrics := flag.Bool("allmetrics", false, "Return all metrics")
	answerNo := flag.Bool("no", false, "Answer no to any questions")
	answerYes := flag.Bool("yes", false, "Answer yes to any questions")
	baseURLFlag := flag.String("url", "", "Gateway base URL")
	cookieFile := flag.String("cookiefile", cookiePath, "File to save session cookies")
	debug := flag.Bool("debug", false, "Enable debug mode")
	filter := flag.String("filter", "", filterDescription)

	freshCookies := flag.Bool(
		"fresh", false,
		"Do not use existing cookies (Warning: If always used the gateway will run out of sessions.)",
	)

	metrics := flag.Bool("metrics", false, "Return metrics instead from table data")
	passwordFlag := flag.String("password", "", "Gateway password")
	pretty := flag.Bool("pretty", false, "Enable pretty mode for nat-connections")

	if colorMode {
		// Replace the default Usage with our colored version
		flag.Usage = ColoredUsage
	}

	flag.Parse()

	return action, allmetrics, answerNo, answerYes, baseURLFlag, cookieFile, debug, filter, freshCookies, metrics, passwordFlag, pretty
}

func validateFlags(actionFlag *string, actionPages map[string]string, allmetrics *bool, baseURLFlag *string, config *Config, filterFlag *string, metrics *bool, passwordFlag *string) (*bool, string, *bool, string) {
	var baseURL string

	if *baseURLFlag == "" {
		baseURL = config.BaseURL
	} else {
		baseURL = *baseURLFlag
	}

	var password string

	if *passwordFlag == "" {
		password = config.Password
	} else {
		password = *passwordFlag
	}

	if *allmetrics {
		*metrics = true
	}

	// Action validation
	isValidAction := false
	actionsHelp := []string{}
	for action := range actionPages {
		actionsHelp = append(actionsHelp, action)
	}

	for _, a := range actionsHelp {
		if *actionFlag == a {
			isValidAction = true
			break
		}
	}

	if !isValidAction && !*allmetrics {
		actionError := fmt.Sprintf("Action must be one of these (%s)", strings.Join(actionsHelp, ", "))
		log.Fatal(actionError)
	}

	// Filter validation
	isValidFilter := false
	isValidFilter = *filterFlag == "" // Default to valid if filter is empty (optional)
	filters := returnFilters()

	if *filterFlag != "" {
		for _, f := range filters {
			if *filterFlag == f {
				isValidFilter = true
				break
			}
		}
	}

	if !isValidFilter {
		filterError := fmt.Sprintf("Filter must be one of these (%s)", strings.Join(filters, ", "))
		log.Fatal(filterError)
	}

	return allmetrics, baseURL, metrics, password
}

func ColoredUsage() {
	// Create color functions
	blue := color.New(color.FgBlue)
	boldGreen := color.New(color.FgGreen, color.Bold)
	cyan := color.New(color.FgCyan)
	green := color.New(color.FgGreen)

	// Print usage header
	fmt.Printf("Usage of %s:\n", green.Sprintf(os.Args[0]))

	flag.VisitAll(func(f *flag.Flag) {
		// Format flag name with color
		s := fmt.Sprintf("  ")
		s += boldGreen.Sprintf("-%s", f.Name)

		// Use reflection to get the type of the flag's value and clean it
		flagType := reflect.TypeOf(f.Value).Elem().Name()
		if strings.HasSuffix(flagType, "Value") {
			flagType = strings.TrimSuffix(flagType, "Value")
		}
		s += blue.Sprintf(" %s", flagType)

		// Add default value if it exists and isn't empty
		if f.DefValue != "" {
			s += blue.Sprintf(" (default: %v)", f.DefValue)
		}

		// Add the usage description in cyan
		if f.Usage != "" {
			s += "\n    \t" + cyan.Sprintf(f.Usage)
		}

		fmt.Println(s)
	})
}
