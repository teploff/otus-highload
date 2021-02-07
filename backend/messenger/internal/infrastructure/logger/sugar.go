package logger

import (
	"github.com/go-kit/kit/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const celleSkip = 2

type zapSugarLogger func(msg string, keysAndValues ...interface{})

// Log if first key = msg than first value will be interpreted as zap logger message
// (this is important for sentry integration).
func (l zapSugarLogger) Log(kv ...interface{}) error {
	if msg0, ok := kv[0].(string); ok && msg0 == "msg" {
		l(kv[1].(string), kv[2:]...)
	} else {
		l("", kv...)
	}

	return nil
}

// NewZapSugarLogger returns a Go kit log.Logger that sends
// log events to a zap.Logger.
func NewZapSugarLogger(logger *zap.Logger, level zapcore.Level) log.Logger {
	var sugar zapSugarLogger

	sugarLogger := logger.WithOptions(zap.AddCallerSkip(celleSkip)).Sugar()

	switch level {
	case zapcore.DebugLevel:
		sugar = sugarLogger.Debugw
	case zapcore.InfoLevel:
		sugar = sugarLogger.Infow
	case zapcore.WarnLevel:
		sugar = sugarLogger.Warnw
	case zapcore.ErrorLevel:
		sugar = sugarLogger.Errorw
	case zapcore.DPanicLevel:
		sugar = sugarLogger.DPanicw
	case zapcore.PanicLevel:
		sugar = sugarLogger.Panicw
	case zapcore.FatalLevel:
		sugar = sugarLogger.Fatalw
	default:
		sugar = sugarLogger.Infow
	}

	return sugar
}
