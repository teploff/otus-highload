// Package logger provides logger functionality (using zap logger instead).
package logger

import (
	"os"

	"gateway/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Option func(logger *zap.Logger) *zap.Logger

//func WithSentry(cfg *config.ErrorTrackerConfig) Option {
//	return func(logger *zap.Logger) *zap.Logger {
//		return wrapZapWithSentry(logger, cfg)
//	}
//}

func New(cfg *config.LoggerConfig, opts ...Option) *zap.Logger {
	var options []zap.Option

	prodConfig := zap.NewProductionEncoderConfig()
	prodConfig.TimeKey = "T"
	prodConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	options = append(options, zap.AddStacktrace(zap.ErrorLevel))
	options = append(options, zap.Development())

	core := zapcore.NewCore(
		encoder,
		os.Stdout,
		getLogLevel(cfg.Level),
	)

	logger := zap.New(core, options...)
	for _, opt := range opts {
		logger = opt(logger)
	}

	return logger
}

// getLogLevel unmarshal text to a zap level notation.
//
// level - text logging notation.
func getLogLevel(level string) zapcore.Level {
	lvl := zap.DebugLevel
	_ = lvl.UnmarshalText([]byte(level))

	return lvl
}
