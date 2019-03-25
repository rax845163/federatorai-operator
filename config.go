package operator

import (
	"github.com/containers-ai/federatorai-operator/pkg/log"
)

// Config encapsultes configuration of federatorai operator
type Config struct {
	WatchNamespace string        `mapstructure:"watch-namespace"`
	Metrics        MetricsConfig `mapstructure:"metrics"`
	Log            log.Config    `mapstructure:"log"`
}

// NewDefaultConfig creates operator default configuration
func NewDefaultConfig() Config {
	return Config{
		WatchNamespace: "",
		Metrics:        NewDefaultMetricsConfig(),
		Log:            log.NewDefaultConfig(),
	}
}

// MetricsConfig encapsultes configuration of federatorai operator metrics server
type MetricsConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// NewDefaultMetricsConfig creates operator's metrics server default configuration
func NewDefaultMetricsConfig() MetricsConfig {
	return MetricsConfig{
		Host: "0.0.0.0",
		Port: 8383,
	}
}
