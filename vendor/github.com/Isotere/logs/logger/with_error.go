package logger

import "go.uber.org/zap"

func (l *Logger) WithError(msg string, err error) {
	l.logger.Errorw(msg, zap.Error(err))
}
