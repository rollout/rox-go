package model

import (
	"fmt"
	"time"

	"github.com/rollout/rox-go/v5/core/context"
)

type BUID interface {
	fmt.Stringer
	GetValue() string
	GetQueryStringParts() map[string]string
}

type DeviceProperties interface {
	GetAllProperties() map[string]string
	RolloutEnvironment() string
	LibVersion() string
	DistinctID() string
	RolloutKey() string
}

type SelfManagedOptions interface {
	ServerURL() string
	AnalyticsURL() string
}

type NetworkConfigurationsOptions interface {
	GetConfigApiEndpoint() string
	GetConfigCloudEndpoint() string
	SendStateApiEndpoint() string
	SendStateCloudEndpoint() string
	AnalyticsEndpoint() string
	PushNotificationEndpoint() string
}

type RoxOptions interface {
	DevModeKey() string
	Version() string
	FetchInterval() time.Duration
	ImpressionHandler() ImpressionHandler
	ConfigurationFetchedHandler() ConfigurationFetchedHandler
	RoxyURL() string
	SelfManagedOptions() SelfManagedOptions
	DynamicPropertyRuleHandler() DynamicPropertyRuleHandler
	NetworkConfigurationsOptions() NetworkConfigurationsOptions
}

type SdkSettings interface {
	APIKey() string
	DevModeSecret() string
}

type InternalFlags interface {
	IsEnabled(flagName string) bool
}

type DynamicAPI interface {
	IsEnabled(name string, defaultValue bool, ctx context.Context) bool
	Value(name string, defaultValue string, options []string, ctx context.Context) string
	GetInt(name string, defaultValue int, options []int, ctx context.Context) int
	GetDouble(name string, defaultValue float64, options []float64, ctx context.Context) float64
}

type DynamicPropertyRuleHandler = func(DynamicPropertyRuleHandlerArgs) interface{}

type DynamicPropertyRuleHandlerArgs struct {
	PropName string
	Context  context.Context
}
