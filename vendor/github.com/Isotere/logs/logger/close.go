package logger

func (l *Logger) Close() {
	_ = l.logger.Sync()
}
