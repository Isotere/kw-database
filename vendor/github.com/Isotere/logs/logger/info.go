package logger

func (l *Logger) Info(args ...interface{}) {
	l.logger.Info(args...)
}
