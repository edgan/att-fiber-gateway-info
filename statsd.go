package main

import (
	"log"

	"github.com/DataDog/datadog-go/statsd"
)

func giveMetricsToDatadogStatsd(configs configs, metrics []string, model string) {
	client, err := statsd.New(configs.StatsdIPPort) // udp

	if err != nil {
		logFatalf("Error creating StatsD client: %v", err)
	}

	defer client.Close()

	floatMetrics := processDatadogMetrics(metrics)

	for key, value := range floatMetrics {
		err = client.Gauge(key, value, []string{"gateway:" + model}, 1)
		if err != nil {
			log.Printf("Error sending %s to statsd: %v", key, err)
		}
	}
}
