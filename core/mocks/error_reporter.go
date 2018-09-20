package mocks

import "github.com/stretchr/testify/mock"

type ErrorReporter struct {
	mock.Mock
}

func (m *ErrorReporter) Report(message string, err error) {
}
