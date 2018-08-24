package mocks

import (
	"github.com/rollout/rox-go/core/network"
	"github.com/stretchr/testify/mock"
)

type RequestConfigurationBuilder struct {
	mock.Mock
}

func (m RequestConfigurationBuilder) BuildForRoxy() network.RequestData {
	args := m.Called()
	return args.Get(0).(network.RequestData)
}

func (m RequestConfigurationBuilder) BuildForCDN() network.RequestData {
	args := m.Called()
	return args.Get(0).(network.RequestData)
}

func (m RequestConfigurationBuilder) BuildForAPI() network.RequestData {
	args := m.Called()
	return args.Get(0).(network.RequestData)
}
