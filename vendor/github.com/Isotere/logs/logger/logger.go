package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// https://betterstack.com/community/guides/logging/go/zap/#logging-errors-with-zap

var (
	logEncoding      = "console"
	logErrorFileName = []string{"stderr", "logs/app_error.log"}
	logFileName      = []string{"stdout", "logs/app.log"}
)

type Logger struct {
	logger *zap.SugaredLogger
}

func New(logLevel string) (*Logger, error) {
	logger := Logger{}

	var encoderCfg zapcore.EncoderConfig

	if logLevel == LogLevelDevel {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
	}

	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderCfg.EncodeName = zapcore.FullNameEncoder

	cfg := &zap.Config{
		Encoding:          logEncoding,
		Level:             zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths:       logFileName,
		ErrorOutputPaths:  logErrorFileName,
		DisableCaller:     false,
		DisableStacktrace: false,
		EncoderConfig:     encoderCfg,

		InitialFields: nil,
	}

	lg := zap.Must(cfg.Build())
	logger.logger = lg.Sugar()

	return &logger, nil
}

func (l *Logger) child(fields map[string]interface{}) *zap.SugaredLogger {
	fieldsToAdd := make([]interface{}, 0)

	for k, v := range fields {
		fieldsToAdd = append(fieldsToAdd, zap.Any(k, v))
	}

	return l.logger.With(fieldsToAdd...)
}
