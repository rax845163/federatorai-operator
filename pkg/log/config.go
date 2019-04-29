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

func (c *Config) AppendOutput(outputPath string) {

	for _, path := range c.OutputPaths {
		if path == outputPath {
			return
		}
	}
	c.OutputPaths = append(c.OutputPaths, outputPath)
}
