package main

import (
	"log"

	"edgan/att-fiber-gateway-info/internal/logging"

	"github.com/DataDog/datadog-go/statsd"
)

func giveMetricsToDatadogStatsd(configs configs, metrics []string, model string) {
	client, err := statsd.New(configs.StatsdIPPort) // udp

	if err != nil {
		logging.LogFatalf("Error creating StatsD client: %v", err)
	}

	defer client.Close()

	floatMetrics := processDatadogMetrics(metrics)

	for key, value := range floatMetrics {
		err = client.Gauge(key, value, []string{"gateway:" + model}, one)
		if err != nil {
			log.Printf("Error sending %s to statsd: %v", key, err)
		}
	}
}
