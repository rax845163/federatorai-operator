package admission_controller

import (
	"crypto/tls"

	"github.com/containers-ai/alameda/admission-controller/pkg/service"
	"github.com/containers-ai/alameda/pkg/framework/datahub"
	"github.com/containers-ai/alameda/pkg/grpc"
	"github.com/containers-ai/alameda/pkg/utils/log"
	"github.com/pkg/errors"
)

type JsonPatchValidationFuncName = string

const (
	JsonPatchValidationFuncOpenshift3_9 = "openshift3.9"
)

// Config contains the server (the webhook) cert and key.
type Config struct {
	CACertFile              string                      `mapstructure:"caCertFile"`
	CertFile                string                      `mapstructure:"tlsCertFile"`
	KeyFile                 string                      `mapstructure:"tlsPrivateKeyFile"`
	Enable                  bool                        `mapstructure:"enable"`
	JsonPatchValidationFunc JsonPatchValidationFuncName `mapstructure:"jsonPatchValidationFunc"`
	DeployedNamespace       string                      `mapstructure:"deployedNamespace"`
	Log                     *log.Config                 `mapstructure:"log"`
	Datahub                 *datahub.Config             `mapstructure:"datahub"`
	Port                    int32                       `mapstructure:"port"`
	Service                 *service.Config             `mapstructure:"service"`
	GRPC                    *grpc.Config                `mapstructure:"gRPC"`
}

func NewDefaultConfig() Config {

	defaultDatahubConfig := datahub.NewDefaultConfig()
	defaultLogConfig := log.NewDefaultConfig()
	defaultSvcConfig := service.NewDefaultConfig()
	return Config{
		CACertFile:              "",
		CertFile:                "",
		KeyFile:                 "",
		Enable:                  false,
		JsonPatchValidationFunc: "",
		DeployedNamespace:       "alameda",
		Log:                     &defaultLogConfig,
		Datahub:                 &defaultDatahubConfig,
		Port:                    8000,
		Service:                 defaultSvcConfig,
	}
}

func (c Config) ConfigTLS() (*tls.Config, error) {
	sCert, err := tls.LoadX509KeyPair(c.CertFile, c.KeyFile)
	if err != nil {
		return nil, errors.Errorf("get tls config failed: %s", err.Error())
	}
	return &tls.Config{
		Certificates: []tls.Certificate{sCert},
		// TODO: uses mutual tls after we agree on what cert the apiserver should use.
		// ClientAuth:   tls.RequireAndVerifyClientCert,
	}, nil
}
