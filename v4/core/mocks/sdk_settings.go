package mocks

import "github.com/stretchr/testify/mock"

type SdkSettings struct {
	mock.Mock
}

func (m *SdkSettings) APIKey() string {
	args := m.Called()
	return args.String(0)
}

func (m *SdkSettings) DevModeSecret() string {
	args := m.Called()
	return args.String(0)
}
