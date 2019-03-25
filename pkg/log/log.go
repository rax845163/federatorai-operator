package log

import (
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZaprLogger(cfg Config) (logr.Logger, error) {

	zapLogger, err := newZapLogger(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "new logger failed")
	}

	logger := zapr.NewLogger(zapLogger)

	return logger, nil
}

func newZapLogger(cfg Config) (*zap.Logger, error) {

	var zapLogLevel zapcore.Level
	zapLogLevel.UnmarshalText([]byte(cfg.OutputLevel))
	zapAtomicLevel := zap.NewAtomicLevelAt(zapLogLevel)

	zapCfg := zap.Config{
		Level:         zapAtomicLevel,
		OutputPaths:   cfg.OutputPaths,
		Encoding:      "console",
		DisableCaller: true,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:    "message",
			LevelKey:      "level",
			TimeKey:       "timestamp",
			NameKey:       "name",
			CallerKey:     "caller",
			StacktraceKey: "stacktrace",
			EncodeLevel:   zapcore.CapitalLevelEncoder,
			EncodeTime:    zapcore.ISO8601TimeEncoder,
			EncodeName:    zapcore.FullNameEncoder,
			EncodeCaller:  zapcore.ShortCallerEncoder,
		},
	}

	zapLogger, err := zapCfg.Build()
	if err != nil {
		return nil, errors.Errorf("new zap logger failed: %s", err.Error())
	}
	defer zapLogger.Sync()

	return zapLogger, nil
}
