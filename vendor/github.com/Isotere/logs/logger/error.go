package logger

func (l *Logger) Error(args ...interface{}) {
	l.logger.Error(args...)
}
