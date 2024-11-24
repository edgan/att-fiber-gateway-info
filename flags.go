package main

import (
	"flag"
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/color"
)

type Configs struct {
	BaseURL      string
	Password     string
	StatsdIPPort string
}

type Flags struct {
	AllMetrics   *bool
	AnswerNo     *bool
	AnswerYes    *bool
	BaseURL      *string
	Continuous   *bool
	CookieFile   *string
	Datadog      *bool
	Debug        *bool
	Filter       *string
	FreshCookies *bool
	Metrics      *bool
	Password     *string
	Pretty       *bool
	StatsdIPPort *string
}

func returnFlags(actionDescription string, colorMode bool, cookiePath string, filterDescription string) (*string, *Flags, *bool) {
	// action is a special case where there can be more than one action per run, and hence it doesn't work as part of
	// the flags struct.
	action := flag.String("action", "", actionDescription)

	flags := &Flags{
		AllMetrics: flag.Bool("allmetrics", false, "Return all metrics"),
		AnswerNo:   flag.Bool("no", false, "Answer no to any questions"),
		AnswerYes:  flag.Bool("yes", false, "Answer yes to any questions"),
		BaseURL:    flag.String("url", "", "Gateway base URL"),
		Continuous: flag.Bool("continuous", false, "Continuously repeat metrics"),
		CookieFile: flag.String("cookiefile", cookiePath, "File to save session cookies"),
		Datadog:    flag.Bool("datadog", false, "Send metrics to datadog"),
		Debug:      flag.Bool("debug", false, "Enable debug mode"),
		Filter:     flag.String("filter", "", filterDescription),

		FreshCookies: flag.Bool(
			"fresh", false,
			"Do not use existing cookies (Warning: If always used the gateway will run out of sessions.)",
		),

		Metrics:      flag.Bool("metrics", false, "Return metrics based on the data instead the data"),
		Password:     flag.String("password", "", "Gateway password"),
		Pretty:       flag.Bool("pretty", false, "Enable pretty mode for nat-connections"),
		StatsdIPPort: flag.String("statsdipport", "", "Statsd ip port, like 127.0.0.1:8125"),
	}

	version := flag.Bool("version", false, "Show version")

	flag.Usage = func() {
		Usage(colorMode)
	}

	flag.Parse()

	return action, flags, version
}

func validateFlags(action string, actionPages map[string]string, config *Config, flags *Flags) (Configs, *Flags) {
	var configs Configs

	if *flags.BaseURL == "" {
		configs.BaseURL = config.BaseURL
	} else {
		configs.BaseURL = *flags.BaseURL
	}

	if *flags.Password == "" {
		configs.Password = config.Password
	} else {
		configs.Password = *flags.Password
	}

	if *flags.Continuous && (!*flags.AllMetrics && !*flags.Metrics) {
		logFatal("-continuous must not be set without -allmetrics or -metrics.")
	}

	if *flags.StatsdIPPort == "" {
		configs.StatsdIPPort = config.StatsdIPPort
	} else {
		configs.StatsdIPPort = *flags.StatsdIPPort
	}

	if *flags.AllMetrics {
		*flags.Metrics = true
	}

	if *flags.Metrics && !*flags.AllMetrics {
		metricActions := returnMeticsActions()

		inSlice := false

		for _, metricAction := range metricActions {
			if action == metricAction {
				inSlice = true
			}
		}

		if !inSlice {
			metricActionError := fmt.Sprintf("Action must be one of these (%s) when -metrics is enabled.", strings.Join(metricActions, ", "))
			logFatal(metricActionError)
		}
	}

	if *flags.Datadog && !*flags.Metrics {
		datadogError := "Metrics must be enabled when enabling datadog"
		logFatal(datadogError)
	}

	// Action validation
	isValidAction := false
	actionsHelp := []string{}
	for action := range actionPages {
		actionsHelp = append(actionsHelp, action)
	}

	for _, a := range actionsHelp {
		if action == a {
			isValidAction = true
			break
		}
	}

	if !isValidAction && !*flags.AllMetrics {
		actionError := fmt.Sprintf("Action must be one of these (%s)", strings.Join(actionsHelp, ", "))
		logFatal(actionError)
	}

	// Filter validation
	isValidFilter := false
	isValidFilter = *flags.Filter == "" // Default to valid if filter is empty (optional)
	filters := returnFilters()

	if *flags.Filter != "" {
		for _, f := range filters {
			if *flags.Filter == f {
				isValidFilter = true
				break
			}
		}
	}

	if !isValidFilter {
		filterError := fmt.Sprintf("Filter must be one of these (%s)", strings.Join(filters, ", "))
		logFatal(filterError)
	}

	return configs, flags
}

func Usage(colorMode bool) {
	// Create color functions
	blue := color.New(color.FgBlue)
	boldGreen := color.New(color.FgGreen, color.Bold)
	cyan := color.New(color.FgCyan)
	green := color.New(color.FgGreen)

	applicationNameVersion := returnApplicationNameVersion()

	if colorMode {
		applicationNameVersion = green.Sprintf(applicationNameVersion)
	}

	fmt.Println(applicationNameVersion)

	usage := "\nUsage:\n"

	if colorMode {
		usage = green.Sprintf(usage)
	}

	fmt.Print(usage)

	flag.VisitAll(func(f *flag.Flag) {
		// Format flag name with color
		s := "  "

		if colorMode {
			s += boldGreen.Sprintf("-%s", f.Name)
		} else {
			s += fmt.Sprintf("-%s", f.Name)
		}

		// Use reflection to get the type of the flag's value and clean it
		flagType := reflect.TypeOf(f.Value).Elem().Name()
		flagType = strings.TrimSuffix(flagType, "Value")

		if colorMode {
			s += blue.Sprintf(" %s", flagType)
		} else {
			s += fmt.Sprintf(" %s", flagType)
		}

		// Add default value if it exists and isn't empty
		if f.DefValue != "" {
			if colorMode {
				s += blue.Sprintf(" (default: %v)", f.DefValue)
			} else {
				s += fmt.Sprintf(" (default: %v)", f.DefValue)
			}
		}

		// Add the usage description in cyan
		if f.Usage != "" {
			if colorMode {
				s += "\n    \t" + cyan.Sprintf(f.Usage)
			} else {
				s += "\n    \t" + fmt.Sprintf(f.Usage)
			}
		}

		fmt.Println(s)
	})
}
