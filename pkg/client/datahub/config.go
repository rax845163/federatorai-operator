package datahub

// Config provides a configuration struct to connect to Alameda-Datahub
type Config struct {
	Address string `mapstructure:"address"`
}

// NewDefaultConfig returns default configuration
func NewDefaultConfig() Config {
	return Config{
		Address: "alameda-datahub.federatorai.svc.cluster.local:50050",
	}
}
