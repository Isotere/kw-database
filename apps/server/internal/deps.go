package internal

type Log interface {
	Info(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	WithError(msg string, err error)
}
