package network_test

import (
	"github.com/rollout/rox-go/core/network"
	"github.com/stretchr/testify/mock"
)

type mockedSdkSettings struct {
	mock.Mock
}

func (m mockedSdkSettings) ApiKey() string {
	args := m.Called()
	return args.String(0)
}

func (m mockedSdkSettings) DevModeSecret() string {
	args := m.Called()
	return args.String(0)
}

type mockedBUID struct {
	mock.Mock
}

func (m mockedBUID) GetValue() string {
	args := m.Called()
	return args.String(0)
}

func (m mockedBUID) GetQueryStringParts() map[string]string {
	args := m.Called()
	return args.Get(0).(map[string]string)
}

type mockedDeviceProperties struct {
	mock.Mock
}

func (m mockedDeviceProperties) GetAllProperties() map[string]string {
	args := m.Called()
	return args.Get(0).(map[string]string)
}

func (m mockedDeviceProperties) RolloutEnvironment() string {
	args := m.Called()
	return args.String(0)
}

func (m mockedDeviceProperties) LibVersion() string {
	args := m.Called()
	return args.String(0)
}

func (m mockedDeviceProperties) DistinctId() string {
	args := m.Called()
	return args.String(0)
}

func (m mockedDeviceProperties) RolloutKey() string {
	args := m.Called()
	return args.String(0)
}

type mockedRequest struct {
	mock.Mock
}

func (m mockedRequest) SendGet(requestData network.RequestData) (response *network.Response, err error) {
	args := m.Called(requestData)
	return args.Get(0).(*network.Response), args.Error(1)
}

func (m mockedRequest) SendPost(uri string, content interface{}) (response *network.Response, err error) {
	args := m.Called(uri, content)
	return args.Get(0).(*network.Response), args.Error(1)
}

type mockedRequestConfigurationBuilder struct {
	mock.Mock
}

func (m mockedRequestConfigurationBuilder) BuildForRoxy() network.RequestData {
	args := m.Called()
	return args.Get(0).(network.RequestData)
}

func (m mockedRequestConfigurationBuilder) BuildForCDN() network.RequestData {
	args := m.Called()
	return args.Get(0).(network.RequestData)
}

func (m mockedRequestConfigurationBuilder) BuildForAPI() network.RequestData {
	args := m.Called()
	return args.Get(0).(network.RequestData)
}
