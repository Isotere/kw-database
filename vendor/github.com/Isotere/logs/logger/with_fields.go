package logger

func (l *Logger) WithFields(msg string, fields map[string]interface{}) {
	child := l.child(fields)
	child.Error(msg)
}
