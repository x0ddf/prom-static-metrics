package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/x0ddf/prom-static-metrics/pkg/config"
	"log"
	"net/http"
	"os"
)

const (
	CONFIG_FILE_VAR = "CONFIG_FILE"
	PORT_VAR        = "PORT"
)

var (
	DefaultConfig = "./config.yaml"
	DefaultPort   = "8080"
)

func main() {

	configFile := os.Getenv(CONFIG_FILE_VAR)
	if configFile == "" {
		configFile = DefaultConfig
	}
	suggestedPort := os.Getenv(PORT_VAR)
	if suggestedPort == "" {
		suggestedPort = DefaultPort
	}
	manager := config.NewMetricsManager(configFile)
	// Expose /metrics HTTP endpoint using the created custom registry.
	http.Handle(
		"/metrics", promhttp.HandlerFor(
			manager.Registry,
			promhttp.HandlerOpts{
				EnableOpenMetrics: true,
			}),
	)
	// To test: curl -H 'Accept: application/openmetrics-text' localhost:8080/metrics
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%s", suggestedPort), nil))

}
