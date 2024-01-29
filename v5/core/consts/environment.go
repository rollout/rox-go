package consts

import "os"

type EnvironmentAPI int

const (
	PLATFORM_API EnvironmentAPI = iota // SDK will utilitze Platform API endpoints
	ROLLOUT_API                        // SDK will utilitze Rollout API endpoints
)

func EnvironmentRoxyInternalPath() string {
	return "device/request_configuration"
}

func EnvironmentCDNPath() string {
	rolloutMode := os.Getenv("ROLLOUT_MODE")

	switch rolloutMode {
	case "QA":
		return "https://qa-conf.rollout.io"
	case "LOCAL":
		return "https://development-conf.rollout.io"
	}
	return "https://conf.rollout.io"
}

func EnvironmentAPIPath() string {
	rolloutMode := os.Getenv("ROLLOUT_MODE")

	switch rolloutMode {
	case "QA":
		return "https://qa-api.rollout.io/device/get_configuration"
	case "LOCAL":
		return "http://127.0.0.1:8557/device/get_configuration"
	}
	return "https://x-api.rollout.io/device/get_configuration"
}

func EnvironmentStateCDNPath() string {
	rolloutMode := os.Getenv("ROLLOUT_MODE")

	switch rolloutMode {
	case "QA":
		return "https://qa-statestore.rollout.io"
	case "LOCAL":
		return "https://development-statestore.rollout.io"
	}
	return "https://statestore.rollout.io"
}

func EnvironmentStateAPIPath() string {
	rolloutMode := os.Getenv("ROLLOUT_MODE")

	switch rolloutMode {
	case "QA":
		return "https://qa-api.rollout.io/device/update_state_store"
	case "LOCAL":
		return "http://127.0.0.1:8557/device/update_state_store"
	}
	return "https://x-api.rollout.io/device/update_state_store"
}

// EnvironmentAnalyticsPath returns the URL for the analytics endpoint.
// envApi: PLATFORM_API or ROLLOUT_API (default) identifies if the SDK should use the Platform API or Rollout API endpoints.
func EnvironmentAnalyticsPath(envApi EnvironmentAPI) string {
	rolloutMode := os.Getenv("ROLLOUT_MODE")

	if envApi == PLATFORM_API {
		switch rolloutMode {
		case "QA":
			return "https://api-staging.saas-dev.beescloud.com/events/flag-impressions"
		case "LOCAL":
			return "http://127.0.0.1:8097/events/flag-impressions"
		}
		return "https://api.cloudbees.io/events/flag-impressions"
	}

	switch rolloutMode {
	case "QA":
		return "https://qaanalytic.rollout.io"
	case "LOCAL":
		return "http://127.0.0.1:8787"
	}
	return "https://analytic.rollout.io"
}

func EnvironmentNotificationsPath() string {
	rolloutMode := os.Getenv("ROLLOUT_MODE")

	switch rolloutMode {
	case "QA":
		return "https://qax-push.rollout.io/sse"
	case "LOCAL":
		return "http://127.0.0.1:8887/sse"
	}
	return "https://push.rollout.io/sse"
}
