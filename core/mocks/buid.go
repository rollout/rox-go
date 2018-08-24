package mocks

import "github.com/stretchr/testify/mock"

type BUID struct {
	mock.Mock
}

func (m BUID) GetValue() string {
	args := m.Called()
	return args.String(0)
}

func (m BUID) GetQueryStringParts() map[string]string {
	args := m.Called()
	return args.Get(0).(map[string]string)
}
