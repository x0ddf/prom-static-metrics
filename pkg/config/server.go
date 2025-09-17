package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"log"
	"slices"
	"sync"
)

type MetricsManager struct {
	MetricsConfig     MetricsConfig
	configPath        string
	Registry          *prometheus.Registry
	registeredMetrics map[string]prometheus.Gauge
	mu                sync.Mutex
}

func NewMetricsManager(configPath string) *MetricsManager {
	manager := &MetricsManager{
		configPath:        configPath,
		Registry:          prometheus.NewRegistry(),
		registeredMetrics: make(map[string]prometheus.Gauge),
	}
	manager.Load()
	manager.UpdateRegistry()
	return manager

}

type MetricsConfig struct {
	Metrics []Metric `mapstructure:"metrics" yaml:"metrics" json:"metrics"`
}
type Metric struct {
	Name        string      `mapstructure:"name" yaml:"name" json:"name"`
	Description string      `mapstructure:"description" yaml:"description" json:"description"`
	Value       interface{} `mapstructure:"value" yaml:"value" json:"value"`
}

func (mm *MetricsManager) Load() {
	viper.SetConfigFile(mm.configPath)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	var config MetricsConfig
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal(err)
	}
	mm.MetricsConfig = config
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed:", e.Name)
		var updatedConfig MetricsConfig
		marshalErr := viper.Unmarshal(&updatedConfig)
		if err != nil {
			log.Printf("fail to unmarshal changed config:%v", marshalErr)
		} else {
			mm.MetricsConfig = updatedConfig
			mm.UpdateRegistry()
		}
	})
	viper.WatchConfig()
}

func (mm *MetricsManager) UpdateRegistry() {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	// clean old metrics (if they are deleted) | update existing
	for metricName, metric := range mm.registeredMetrics {
		if idx := slices.IndexFunc(mm.MetricsConfig.Metrics, func(metric Metric) bool {
			return metricName == metric.Name
		}); idx >= 0 {
			mm.registeredMetrics[metricName] = buildMetric(&mm.MetricsConfig.Metrics[idx])
			log.Printf("metric:%v found in the new config, updated", metricName)
		} else {
			log.Printf("metric:%v not found in the new config, purged", metricName)
			mm.Registry.Unregister(metric)
			delete(mm.registeredMetrics, metricName)
		}
	}
	for _, newMetric := range mm.MetricsConfig.Metrics {
		if _, ok := mm.registeredMetrics[newMetric.Name]; !ok {
			md := buildMetric(&newMetric)
			registerErr := mm.Registry.Register(md)
			if registerErr != nil {
				log.Printf("fail to register metric:%v | discarded", registerErr)
			} else {
				mm.registeredMetrics[newMetric.Name] = md
			}
		}
	}
	log.Println("metrics updated")

}

func buildMetric(m *Metric) prometheus.Gauge {
	built := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        m.Name,
		Help:        m.Description,
		ConstLabels: map[string]string{m.Name: fmt.Sprintf("%v", m.Value)},
	})
	built.Set(1)
	return built
}
