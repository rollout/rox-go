package consts

import "os"

type EnvironmentAPI int

const (
	// API providers for the SDK endpoints
	PLATFORM_API EnvironmentAPI = iota // SDK will utilitze Platform API endpoints
	ROLLOUT_API                        // SDK will utilitze Rollout API endpoints
)

func EnvironmentRoxyInternalPath() string {
	return "device/request_configuration"
}

func EnvironmentCDNPath(envApi EnvironmentAPI) string {
	rolloutMode := os.Getenv("ROLLOUT_MODE")

	switch rolloutMode {
	case "QA":
		return "https://qa-conf.rollout.io"
	case "LOCAL":
		return "https://development-conf.rollout.io"
	}

	if envApi == PLATFORM_API {
		return "https://rox-conf.cloudbees.io"
	}
	return "https://conf.rollout.io"
}

func EnvironmentAPIPath(envApi EnvironmentAPI) string {
	rolloutMode := os.Getenv("ROLLOUT_MODE")

	switch rolloutMode {
	case "QA":
		return "https://qa-api.rollout.io/device/get_configuration"
	case "LOCAL":
		return "http://127.0.0.1:8557/device/get_configuration"
	}

	if envApi == PLATFORM_API {
		return "https://api.cloudbees.io/device/get_configuration"
	}
	return "https://x-api.rollout.io/device/get_configuration"
}

func EnvironmentStateCDNPath(envApi EnvironmentAPI) string {
	rolloutMode := os.Getenv("ROLLOUT_MODE")

	switch rolloutMode {
	case "QA":
		return "https://qa-statestore.rollout.io"
	case "LOCAL":
		return "https://development-statestore.rollout.io"
	}

	if envApi == PLATFORM_API {
		return "https://rox-state.cloudbees.io"
	}
	return "https://statestore.rollout.io"
}

func EnvironmentStateAPIPath(envApi EnvironmentAPI) string {
	rolloutMode := os.Getenv("ROLLOUT_MODE")

	switch rolloutMode {
	case "QA":
		return "https://qa-api.rollout.io/device/update_state_store"
	case "LOCAL":
		return "http://127.0.0.1:8557/device/update_state_store"
	}

	if envApi == PLATFORM_API {
		return "https://api.cloudbees.io/device/update_state_store"
	}
	return "https://api.cloudbees.io/device/update_state_store"
}

// EnvironmentAnalyticsPath returns the URL for the analytics endpoint.
func EnvironmentAnalyticsPath(envApi EnvironmentAPI) string {
	rolloutMode := os.Getenv("ROLLOUT_MODE")

	switch rolloutMode {
	case "QA":
		return "https://qaanalytic.rollout.io"
	case "LOCAL":
		return "http://127.0.0.1:8787"
	}

	if envApi == PLATFORM_API {
		return "https://fm-analytics.cloudbees.io/impression"
	}
	return "https://analytic.rollout.io"
}

func EnvironmentNotificationsPath(envApi EnvironmentAPI) string {
	rolloutMode := os.Getenv("ROLLOUT_MODE")

	switch rolloutMode {
	case "QA":
		return "https://qax-push.rollout.io/sse"
	case "LOCAL":
		return "http://127.0.0.1:8887/sse"
	}

	if envApi == PLATFORM_API {
		return "https://sdk-notification-service.cloudbees.io/sse"
	}
	return "https://push.rollout.io/sse"
}
