package tcp

import (
	"fmt"
)

type Logger interface {
	Info(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	WithError(msg string, err error)
}

type FakeLogger struct{}

func NewFakeLogger() *FakeLogger {
	return &FakeLogger{}
}

func (f *FakeLogger) Info(args ...interface{}) {
	fmt.Println(args...)
}

func (f *FakeLogger) Warning(args ...interface{}) {
	fmt.Println(args...)
}

func (f *FakeLogger) Error(args ...interface{}) {
	fmt.Println(args...)
}

func (f *FakeLogger) Fatal(args ...interface{}) {
	fmt.Println(args...)
}

func (f *FakeLogger) WithError(msg string, err error) {
	fmt.Println(msg, err)
}
