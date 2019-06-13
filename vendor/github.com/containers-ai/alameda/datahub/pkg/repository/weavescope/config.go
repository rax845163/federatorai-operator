package weavescope

// Config Configuration of Prometheus datasource
type Config struct {
	URL string `mapstructure:"url"`
}

const (
	defaultURL = "https://weavescope:4041"
)

// NewDefaultConfig Provide default configuration
func NewDefaultConfig() Config {
	var config = Config{
		URL: defaultURL,
	}
	return config
}
