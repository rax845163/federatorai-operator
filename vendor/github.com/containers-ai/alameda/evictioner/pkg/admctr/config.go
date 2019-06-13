package admctr

import "fmt"

type Config struct {
	SvcName string `mapstructure:"service-name"`
	SvcPort int32  `mapstructure:"service-port"`
}

func NewConfig() *Config {
	c := Config{}
	c.init()
	return &c
}

func (c *Config) init() {
	c.SvcName = "admission-controller"
	c.SvcPort = 443
}

func (c *Config) Validate() error {
	if c.SvcPort < 1 || c.SvcPort > 65535 {
		return fmt.Errorf("Admission controller service port %v is not valid", c.SvcPort)
	}
	return nil
}
