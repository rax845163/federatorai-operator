package evictioner

import (
	"github.com/containers-ai/alameda/evictioner/pkg/admctr"
	"github.com/containers-ai/alameda/evictioner/pkg/datahub"
	"github.com/containers-ai/alameda/evictioner/pkg/eviction"
	"github.com/containers-ai/alameda/pkg/utils/log"
)

// Config is evict configuration
type Config struct {
	Log      *log.Config      `mapstructure:"log"`
	Eviction *eviction.Config `mapstructure:"eviction"`
	Datahub  *datahub.Config  `mapstructure:"datahub"`
	AdmCtr   *admctr.Config   `mapstructure:"admission-controller"`
}

// NewDefaultConfig returns Config instance
func NewDefaultConfig() Config {

	var (
		defaultlogConfig      = log.NewDefaultConfig()
		defaultDatahubConfig  = datahub.NewConfig()
		defaultEvictionConfig = eviction.NewDefaultConfig()
		defaultAdmCtlConfig   = admctr.NewConfig()
		config                = Config{
			Log:      &defaultlogConfig,
			Datahub:  defaultDatahubConfig,
			Eviction: &defaultEvictionConfig,
			AdmCtr:   defaultAdmCtlConfig,
		}
	)

	return config
}

func (c *Config) Validate() error {
	return nil
}
