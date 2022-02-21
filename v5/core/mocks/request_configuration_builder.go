package mocks

import (
	"github.com/rollout/rox-go/v5/core/model"
	"github.com/stretchr/testify/mock"
)

type RequestConfigurationBuilder struct {
	mock.Mock
}

func (m *RequestConfigurationBuilder) BuildForRoxy() model.RequestData {
	args := m.Called()
	return args.Get(0).(model.RequestData)
}

func (m *RequestConfigurationBuilder) BuildForCDN() model.RequestData {
	args := m.Called()
	return args.Get(0).(model.RequestData)
}

func (m *RequestConfigurationBuilder) BuildForAPI() model.RequestData {
	args := m.Called()
	return args.Get(0).(model.RequestData)
}
