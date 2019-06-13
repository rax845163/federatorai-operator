package datahub

import (
	"errors"

	influxdb_repository "github.com/containers-ai/alameda/datahub/pkg/repository/influxdb"
	"github.com/containers-ai/alameda/datahub/pkg/repository/prometheus"
	"github.com/containers-ai/alameda/datahub/pkg/repository/weavescope"
	"github.com/containers-ai/alameda/pkg/utils/log"
)

const (
	defaultBindAddress = ":50050"
)

type Config struct {
	BindAddress string                      `mapstructure:"bind-address"`
	Prometheus  *prometheus.Config          `mapstructure:"prometheus"`
	InfluxDB    *influxdb_repository.Config `mapstructure:"influxdb"`
	WeaveScope  *weavescope.Config          `mapstructure:"weavescope"`
	Log         *log.Config                 `mapstructure:"log"`
}

func NewDefaultConfig() Config {

	var (
		defaultlogConfig        = log.NewDefaultConfig()
		defaultPrometheusConfig = prometheus.NewDefaultConfig()
		defaultInfluxDBConfig   = influxdb_repository.NewDefaultConfig()
		defaultWeaveScopeConfig = weavescope.NewDefaultConfig()
		config                  = Config{
			BindAddress: defaultBindAddress,
			Prometheus:  &defaultPrometheusConfig,
			InfluxDB:    &defaultInfluxDBConfig,
			WeaveScope:  &defaultWeaveScopeConfig,
			Log:         &defaultlogConfig,
		}
	)

	return config
}

func (c *Config) Validate() error {

	var err error

	err = c.Prometheus.Validate()
	if err != nil {
		return errors.New("gRPC config validate failed: " + err.Error())
	}

	return nil
}
