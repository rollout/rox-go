package consts

import "os"

func EnvironmentRoxyInternalPath() string {
	return "device/request_configuration"
}

func EnvironmentCDNPath() string {
	rolloutMode := os.Getenv("ROLLOUT_MODE")

	switch rolloutMode {
	case "QA":
		return "https://s3.amazonaws.com/qa-rox-conf.rollout.io/v1/qa"
	case "LOCAL":
		return "https://s3.amazonaws.com/qa-rox-conf.rollout.io/v1/development"
	}
	return "https://s3.amazonaws.com/rox-conf.rollout.io/v1/production"
}

func EnvironmentAPIPath() string {
	rolloutMode := os.Getenv("ROLLOUT_MODE")

	switch rolloutMode {
	case "QA":
		return "https://qax.rollout.io/device/request_configuration"
	case "LOCAL":
		return "http://127.0.0.1:8557/device/request_configuration"
	}
	return "https://x-api.rollout.io/device/request_configuration"
}

func EnvironmentAnalyticsPath() string {
	rolloutMode := os.Getenv("ROLLOUT_MODE")

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
