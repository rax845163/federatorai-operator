package log

type Config struct {
	OutputLevel string   `mapstructure:"outputLevel"`
	OutputPaths []string `mapstructure:"outputPaths"`
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
