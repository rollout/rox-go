package model

import (
	"fmt"
	"time"

	"github.com/rollout/rox-go/v4/core/context"
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

type RoxOptions interface {
	DevModeKey() string
	Version() string
	FetchInterval() time.Duration
	ImpressionHandler() ImpressionHandler
	ConfigurationFetchedHandler() ConfigurationFetchedHandler
	RoxyURL() string
	SelfManagedOptions() SelfManagedOptions
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
}
