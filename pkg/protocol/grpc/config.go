package grpc

import (
	"github.com/containers-ai/federatorai-operator/pkg/log"
)

// Config defines configuration of GRPC protocol
type Config struct {
	Log log.Config `mapstructure:"log"`
}

// NewDefaultConfig returns default configuration of GRPC protocol
func NewDefaultConfig() Config {
	return Config{
		Log: log.NewDefaultConfig(),
	}
}
