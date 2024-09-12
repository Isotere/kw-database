package logger

func (l *Logger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}
