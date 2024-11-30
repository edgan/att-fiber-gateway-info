package main

import (
	"flag"
	"fmt"
	"reflect"
	"strings"

	"edgan/att-fiber-gateway-info/internal/logging"

	"github.com/fatih/color"
)

type configs struct {
	BaseURL      string
	Password     string
	StatsdIPPort string
}

type flags struct {
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

func returnFlags(
	actionDescription string, colorMode bool, cookiePath string, filterDescription string,
) (*string, *flags, *bool) {
	// action is a special case where there can be more than one action per
	// run, and hence it doesn't work as part of the flags struct.
	action := flag.String("action", empty, actionDescription)

	flags := &flags{
		AllMetrics: flag.Bool("allmetrics", false, "Return all metrics"),
		AnswerNo:   flag.Bool("no", false, "Answer no to any questions"),
		AnswerYes:  flag.Bool("yes", false, "Answer yes to any questions"),
		BaseURL:    flag.String("url", empty, "Gateway base URL"),
		Continuous: flag.Bool("continuous", false, "Continuously repeat metrics"),
		CookieFile: flag.String("cookiefile", cookiePath, "File to save session cookies"),
		Datadog:    flag.Bool("datadog", false, "Send metrics to datadog"),
		Debug:      flag.Bool("debug", false, "Enable debug mode"),
		Filter:     flag.String("filter", empty, filterDescription),

		FreshCookies: flag.Bool(
			"fresh", false,
			"Do not use existing cookies (Warning: If always used the gateway will run out of sessions.)",
		),

		Metrics:      flag.Bool("metrics", false, "Return metrics based on the data instead the data"),
		Noconvert:    flag.Bool("noconvert", false, "Do not convert text metrics to float values"),
		Password:     flag.String("password", empty, "Gateway password"),
		Pretty:       flag.Bool("pretty", false, "Enable pretty mode for nat-connections"),
		StatsdIPPort: flag.String("statsdipport", empty, "Statsd ip port, like 127.0.0.1:8125"),
	}

	version := flag.Bool("version", false, "Show version")

	flag.Usage = func() {
		usage(colorMode)
	}

	flag.Parse()

	return action, flags, version
}

func validateMetricsFlags(flags *flags, returnFlags *flags) {
	if *flags.Continuous && !(*flags.AllMetrics || *flags.Metrics) {
		logging.LogFatal("-continuous must not be set without -allmetrics or -metrics.")
	}

	if *flags.AllMetrics || *flags.Datadog {
		*returnFlags.Metrics = true
	}
}

func validateMetricActionsFlags(action string, flags *flags) {
	if *flags.Metrics && !*flags.AllMetrics {
		if !contains(metricActions, action) {
			logging.LogFatal(
				fmt.Sprintf(
					"Action must be one of these (%s) when -metrics is enabled.",
					strings.Join(metricActions, commaSpace),
				),
			)
		}
	}
}

func validateActionsAndFilterFlags(action string, actionPages map[string]string, flags *flags) {
	// Action validation
	if !*flags.AllMetrics && !containsMapKey(actionPages, action) {
		actionsHelp := getMapKeys(actionPages)
		logging.LogFatal(fmt.Sprintf("Action must be one of these (%s)", strings.Join(actionsHelp, commaSpace)))
	}

	// Filter validation
	if *flags.Filter != empty {
		if !contains(filters, *flags.Filter) {
			logging.LogFatal(fmt.Sprintf("Filter must be one of these (%s)", strings.Join(filters, commaSpace)))
		}
	}
}

func validateFlags(
	action string, actionPages map[string]string, config *config, flags *flags) (configs configs, returnFlags *flags,
) {
	returnFlags = flags

	configs.BaseURL = returnConfigValue(*flags.BaseURL, config.BaseURL)
	configs.Password = returnConfigValue(*flags.Password, config.Password)
	configs.StatsdIPPort = returnConfigValue(*flags.StatsdIPPort, config.StatsdIPPort)

	validateMetricsFlags(flags, returnFlags)

	validateMetricActionsFlags(action, flags)

	validateActionsAndFilterFlags(action, actionPages, flags)

	return configs, returnFlags
}

func usage(colorMode bool) {
	// Define Sprintf functions based on colorMode
	var blueSprintf, boldGreenSprintf, cyanSprintf, greenSprintf func(format string, a ...any) string

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
		if f.DefValue != empty {
			s += blueSprintf(" (default: %v)", f.DefValue)
		}

		// Add the usage description
		if f.Usage != empty {
			s += "\n    \t" + cyanSprintf(f.Usage)
		}

		fmt.Println(s)
	})
}
