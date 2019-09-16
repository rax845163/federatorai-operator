package grpc

// Config defines configuration of GRPC protocol
type Config struct {
	Retry uint `mapstructure:"retry"`
}

// NewDefaultConfig returns default configuration of GRPC protocol
func NewDefaultConfig() Config {
	return Config{
		Retry: 3,
	}
}
