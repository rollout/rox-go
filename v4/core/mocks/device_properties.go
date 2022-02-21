package mocks

import "github.com/stretchr/testify/mock"

type DeviceProperties struct {
	mock.Mock
}

func (m *DeviceProperties) GetAllProperties() map[string]string {
	args := m.Called()
	return args.Get(0).(map[string]string)
}

func (m *DeviceProperties) RolloutEnvironment() string {
	args := m.Called()
	return args.String(0)
}

func (m *DeviceProperties) LibVersion() string {
	args := m.Called()
	return args.String(0)
}

func (m *DeviceProperties) DistinctID() string {
	args := m.Called()
	return args.String(0)
}

func (m *DeviceProperties) RolloutKey() string {
	args := m.Called()
	return args.String(0)
}
