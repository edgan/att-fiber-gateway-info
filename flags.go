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
	Noconvert    *bool
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
		Noconvert:    flag.Bool("noconvert", false, "Do not convert text metrics to float values"),
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

	// Helper function to get the flag value or default to config value
	getConfigValue := func(flagValue string, configValue string) string {
		if flagValue != "" {
			return flagValue
		}
		return configValue
	}

	// Assign config values using helper function
	configs.BaseURL = getConfigValue(*flags.BaseURL, config.BaseURL)
	configs.Password = getConfigValue(*flags.Password, config.Password)
	configs.StatsdIPPort = getConfigValue(*flags.StatsdIPPort, config.StatsdIPPort)

	if *flags.Continuous && !(*flags.AllMetrics || *flags.Metrics) {
		logFatal("-continuous must not be set without -allmetrics or -metrics.")
	}

	if *flags.AllMetrics {
		*flags.Metrics = true
	}

	if *flags.Metrics && !*flags.AllMetrics {
		metricActions := returnMeticsActions()
		if !contains(metricActions, action) {
			logFatal(fmt.Sprintf("Action must be one of these (%s) when -metrics is enabled.", strings.Join(metricActions, ", ")))
		}
	}

	if *flags.Datadog && !*flags.Metrics {
		logFatal("Metrics must be enabled when enabling datadog")
	}

	// Action validation
	if !*flags.AllMetrics && !containsMapKey(actionPages, action) {
		actionsHelp := getMapKeys(actionPages)
		logFatal(fmt.Sprintf("Action must be one of these (%s)", strings.Join(actionsHelp, ", ")))
	}

	// Filter validation
	if *flags.Filter != "" {
		filters := returnFilters()
		if !contains(filters, *flags.Filter) {
			logFatal(fmt.Sprintf("Filter must be one of these (%s)", strings.Join(filters, ", ")))
		}
	}

	return configs, flags
}

func Usage(colorMode bool) {
	// Define Sprintf functions based on colorMode
	var blueSprintf, boldGreenSprintf, cyanSprintf, greenSprintf func(format string, a ...interface{}) string

	if colorMode {
		blueSprintf = color.New(color.FgBlue).Sprintf
		boldGreenSprintf = color.New(color.FgGreen, color.Bold).Sprintf
		cyanSprintf = color.New(color.FgCyan).Sprintf
		greenSprintf = color.New(color.FgGreen).Sprintf
	} else {
		blueSprintf = fmt.Sprintf
		boldGreenSprintf = fmt.Sprintf
		cyanSprintf = fmt.Sprintf
		greenSprintf = fmt.Sprintf
	}

	applicationNameVersion := returnApplicationNameVersion()
	fmt.Println(greenSprintf(applicationNameVersion))

	fmt.Print(greenSprintf("\nUsage:\n"))

	flag.VisitAll(func(f *flag.Flag) {
		// Format flag name with color
		s := "  " + boldGreenSprintf("-%s", f.Name)

		// Get the type of the flag's value
		flagType := reflect.TypeOf(f.Value).Elem().Name()
		flagType = strings.TrimSuffix(flagType, "Value")
		s += blueSprintf(" %s", flagType)

		// Add default value if it exists
		if f.DefValue != "" {
			s += blueSprintf(" (default: %v)", f.DefValue)
		}

		// Add the usage description
		if f.Usage != "" {
			s += "\n    \t" + cyanSprintf(f.Usage)
		}

		fmt.Println(s)
	})
}
