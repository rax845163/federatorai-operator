package prometheus

import (
	"github.com/pkg/errors"
	"net/url"
)

const (
	defaultURL             = "https://prometheus-k8s.openshift-monitoring:9091"
	defaultBearerTokenFile = "/var/run/secrets/kubernetes.io/serviceaccount/token"
)

// Config Configuration of Prometheus datasource
type Config struct {
	URL             string     `mapstructure:"url"`
	BearerTokenFile string     `mapstructure:"bearer-token-file"`
	TLSConfig       *TLSConfig `mapstructure:"tls-config"`

	bearerToken string
}

// TLSConfig Configuration of tls connnection
type TLSConfig struct {
	InsecureSkipVerify bool `mapstructure:"insecure-skip-verify"`
}

// NewDefaultConfig Provide default configuration
func NewDefaultConfig() Config {

	var config = Config{
		URL:             defaultURL,
		BearerTokenFile: defaultBearerTokenFile,
		TLSConfig: &TLSConfig{
			InsecureSkipVerify: true,
		},
	}
	return config
}

// Validate Confirm the configuration is validate
func (c *Config) Validate() error {

	var err error

	_, err = url.Parse(c.URL)
	if err != nil {
		return errors.Errorf("prometheus config validate failed: %s", err.Error())
	}

	return nil
}
