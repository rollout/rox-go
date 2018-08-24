package mocks

import "github.com/stretchr/testify/mock"

type InternalFlags struct {
	mock.Mock
}

func (m *InternalFlags) IsEnabled(flagName string) bool {
	args := m.Called(flagName)
	return args.Bool(0)
}
