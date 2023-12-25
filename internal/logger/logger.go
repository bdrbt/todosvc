// Package logger is just engine-agnostic wrapper/helper to
// formalize all log requests
package logger

import (
	"log"

	"github.com/bdrbt/todo/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(mode int) *zap.Logger {
	lcfg := zap.Config{
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	switch mode {
	case config.EnvDevelopment:
		{
			lcfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
			lcfg.Encoding = "console"
			lcfg.Development = true
		}
	case config.EnvTesting:
		{
			lcfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
			lcfg.Encoding = "console"
			lcfg.Development = false
		}
	default:
		{
			lcfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
			lcfg.Encoding = "json"
			lcfg.Development = false
		}
	}

	logger, err := lcfg.Build()
	if err != nil {
		log.Fatalf("error creating logger:%v", err)
	}

	return logger
}
