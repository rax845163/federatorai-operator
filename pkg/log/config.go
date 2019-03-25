package log

type Config struct {
	OutputLevel string   `mapstructure:"output-level"`
	OutputPaths []string `mapstructure:"output-paths"`
}

func NewDefaultConfig() Config {
	return Config{
		OutputLevel: "info",
		OutputPaths: []string{"stdout"},
	}
}
