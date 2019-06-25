package eviction

type triggerThreshold struct {
	CPU    float64 `mapstructure:"cpu"`
	Memory float64 `mapstructure:"memory"`
}

// Config is eviction configuration
type Config struct {
	CheckCycle              int64 `mapstructure:"check-cycle"`
	Enable                  bool  `mapstructure:"enable"`
	PurgeContainerCPUMemory bool  `mapstructure:"purge-container-cpu-memory"`
}

// NewDefaultConfig returns Config instance
func NewDefaultConfig() Config {
	return Config{
		CheckCycle:              3,
		Enable:                  false,
		PurgeContainerCPUMemory: false,
	}
}

func (c *Config) Validate() error {
	return nil
}
