package influxdb

import (
	"net/url"

	"github.com/pkg/errors"
)

const (
	defaultAddress                = "https://influxdb.alameda.svc.cluster.local:8086"
	defaultUsername               = "datahub"
	defaultPassword               = "datahub"
	defaultRetentionDuration      = "15d"
	defaultRetentionShardDuration = "1d"
)

type Config struct {
	Address                string `mapstructure:"address"`
	Username               string `mapstructure:"username"`
	Password               string `mapstructure:"password"`
	RetentionDuration      string `mapstructure:"retentionDuration"`
	RetentionShardDuration string `mapstructure:"retentionShardDuration"`
}

func NewDefaultConfig() Config {
	var config = Config{
		Address:                defaultAddress,
		Username:               defaultUsername,
		Password:               defaultPassword,
		RetentionDuration:      defaultRetentionDuration,
		RetentionShardDuration: defaultRetentionShardDuration,
	}
	return config
}

func (c *Config) Validate() error {
	_, err := url.Parse(c.Address)
	if err != nil {
		return errors.New("InfluxDB config validate failed: " + err.Error())
	}

	return nil
}
