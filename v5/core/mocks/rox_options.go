package mocks

import (
	"time"

	"github.com/rollout/rox-go/v5/core/model"
	"github.com/stretchr/testify/mock"
)

type RoxOptions struct {
	mock.Mock
}

func (m *RoxOptions) DevModeKey() string {
	args := m.Called()
	return args.String(0)
}

func (m *RoxOptions) Version() string {
	args := m.Called()
	return args.String(0)
}

func (m *RoxOptions) FetchInterval() time.Duration {
	args := m.Called()
	return args.Get(0).(time.Duration)
}

func (m *RoxOptions) ImpressionHandler() model.ImpressionHandler {
	args := m.Called()
	result := args.Get(0)
	if result == nil {
		return nil
	}
	return result.(model.ImpressionHandler)
}

func (m *RoxOptions) ConfigurationFetchedHandler() model.ConfigurationFetchedHandler {
	args := m.Called()
	result := args.Get(0)
	if result == nil {
		return nil
	}
	return result.(model.ConfigurationFetchedHandler)
}

func (m *RoxOptions) RoxyURL() string {
	args := m.Called()
	return args.String(0)

}

func (m *RoxOptions) SelfManagedOptions() model.SelfManagedOptions {
	args := m.Called()
	result := args.Get(0)
	if result == nil {
		return nil
	}
	return result.(model.SelfManagedOptions)
}

func (m *RoxOptions) DynamicPropertyRuleHandler() model.DynamicPropertyRuleHandler {
	args := m.Called()
	result := args.Get(0)
	if result == nil {
		return nil
	}
	return result.(model.DynamicPropertyRuleHandler)
}

func (m *RoxOptions) NetworkConfigurationsOptions() model.NetworkConfigurationsOptions {
	return nil
}