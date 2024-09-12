package logger

func (l *Logger) Warning(args ...interface{}) {
	l.logger.Warn(args...)
}
