package datahub

import (
	"errors"
	"net/url"

	datahubutils "github.com/containers-ai/alameda/operator/pkg/utils/datahub"
)

type Config struct {
	Address string `mapstructure:"address"`
}

func NewConfig() *Config {

	c := Config{}
	c.init()
	return &c
}

func (c *Config) init() {
	c.Address = datahubutils.GetDatahubAddress()
}

func (c *Config) Validate() error {

	var err error

	_, err = url.Parse(c.Address)
	if err != nil {
		return errors.New("datahub config validate failed: " + err.Error())
	}

	return nil
}
