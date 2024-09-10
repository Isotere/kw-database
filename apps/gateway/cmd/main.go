package main

import "os"

const (
	fail    = 1
	success = 0
)

func main() {
	os.Exit(run())
}

func run() (exitCode int) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			// log.Error(context.Background(), panicErr)
			exitCode = fail
		}
	}()

	return success
}
