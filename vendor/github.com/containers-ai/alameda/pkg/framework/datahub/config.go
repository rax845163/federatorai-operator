package datahub

// Config datahub service configuration
type Config struct {
	Address string `mapstructure:"address"`
}

func NewDefaultConfig() Config {
	return Config{
		Address: "datahub.alameda.svc:50050",
	}
}
