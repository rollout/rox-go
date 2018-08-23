package model

import "time"

type BUID interface {
	GetValue() string
	GetQueryStringParts() map[string]string
}

type DeviceProperties interface {
	GetAllProperties() map[string]string
	RolloutEnvironment() string
	LibVersion() string
	DistinctId() string
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
	ApiKey() string
	DevModeSecret() string
}

type InternalFlags interface {
	IsEnabled(flagName string) bool
}
