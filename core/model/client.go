package model

import (
	"fmt"
	"time"
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

type RoxOptions interface {
	DevModeKey() string
	Version() string
	FetchInterval() time.Duration
	ImpressionHandler() ImpressionHandler
	ConfigurationFetchedHandler() ConfigurationFetchedHandler
	RoxyURL() string
}

type SdkSettings interface {
	APIKey() string
	DevModeSecret() string
}

type InternalFlags interface {
	IsEnabled(flagName string) bool
}
