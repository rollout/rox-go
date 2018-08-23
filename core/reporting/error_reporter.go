package reporting

type ErrorReporter interface {
	Report(message string, err error)
}

type errorReporter struct {
}

func NewErrorReporter() ErrorReporter {
	// TODO
	return &errorReporter{}
}

func (errorReporter) Report(message string, err error) {
	// TODO
}
