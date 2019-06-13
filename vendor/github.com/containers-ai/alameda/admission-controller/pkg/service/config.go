package service

type Config struct {
	Name string `mapstructure:"name"`
	Port int32  `mapstructure:"port"`
}

func NewDefaultConfig() *Config {

	c := &Config{
		Name: "admission-controller",
		Port: 443,
	}
	return c
}
