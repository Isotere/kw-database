package logger

func (l *Logger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}
