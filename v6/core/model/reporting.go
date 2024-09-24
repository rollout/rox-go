package model

type ErrorReporter interface {
	Report(message string, err error)
}
