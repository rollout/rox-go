package mocks

import "github.com/stretchr/testify/mock"

// Debug(message string, err interface{})
// Warn(message string, err interface{})
// Error(message string, err interface{})
type Logger struct {
	mock.Mock
}

// Debug provides a mock function with given fields: message, err
func (m *Logger) Debug(message string, err interface{}) {
	m.Called(message, err)
}

// Error provides a mock function with given fields: message, err
func (m *Logger) Error(message string, err interface{}) {
	m.Called(message, err)
}

// Warn provides a mock function with given fields: message, err
func (m *Logger) Warn(message string, err interface{}) {
	m.Called(message, err)
}
