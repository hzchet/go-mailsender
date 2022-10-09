package logger

import (
	"fmt"

	"github.com/TheZeroSlave/zapsentry"
	"github.com/caarlos0/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger = *zap.Logger


type Config struct {
	IsProd bool `env:"IS_PROD" envDefault:"false"`
	Dsn string `env:"SENTRY_ADDRESS" envDefault:""`
}

func New() (Logger, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("read logger configuration failed: %w", err)
	}

	var zapCfg zap.Config
	if cfg.IsProd {
		zapCfg = zap.NewProductionConfig()
	} else {
		zapCfg = zap.NewDevelopmentConfig()
	}

	l, err := zapCfg.Build()
	if err != nil {
		return nil, fmt.Errorf("create logger failed: %w", err)
	}

	l = initSentry(l, cfg.Dsn, "myenv")

	zap.ReplaceGlobals(l)

	return l, nil
}

func initSentry(log *zap.Logger, sentryAddress, environment string) *zap.Logger {
	if sentryAddress == "" {
		return log
	}

	cfg := zapsentry.Configuration{
		Level: zapcore.ErrorLevel,
		Tags: map[string]string{
			"environment": environment,
			"app":         "demoApp",
		},
	}

	core, err := zapsentry.NewCore(cfg, zapsentry.NewSentryClientFromDSN(sentryAddress))
	if err != nil {
		log.Warn("failed to init zap", zap.Error(err))
	}

	return zapsentry.AttachCoreToLogger(core, log)
}
