package client

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rollout/rox-go/v5/core/consts"
)

type Paths struct {
	CDN           string
	API           string
	StateCDN      string
	StateAPI      string
	Analytics     string
	Notifications string
}

func TestSaasEnvironmentEnvironment_PlatformAPI_Path(t *testing.T) {
	testsCases := []struct {
		testName string
		osEnv    string
		platform consts.EnvironmentAPI
		expected *Paths
	}{
		{
			testName: "Platform QA API",
			osEnv:    "QA",
			platform: consts.PLATFORM_API,
			expected: &Paths{
				CDN:           "https://qa-conf.rollout.io",
				API:           "https://qa-api.rollout.io/device/get_configuration",
				StateCDN:      "https://qa-statestore.rollout.io",
				StateAPI:      "https://qa-api.rollout.io/device/update_state_store",
				Analytics:     "https://qaanalytic.rollout.io",
				Notifications: "https://qax-push.rollout.io/sse",
			},
		},
		{
			testName: "Platform Local API",
			osEnv:    "LOCAL",
			platform: consts.PLATFORM_API,
			expected: &Paths{
				CDN:           "https://development-conf.rollout.io",
				API:           "http://127.0.0.1:8557/device/get_configuration",
				StateCDN:      "https://development-statestore.rollout.io",
				StateAPI:      "http://127.0.0.1:8557/device/update_state_store",
				Analytics:     "http://127.0.0.1:8787",
				Notifications: "http://127.0.0.1:8887/sse",
			},
		},
		{
			testName: "Platform Production API",
			osEnv:    "",
			platform: consts.PLATFORM_API,
			expected: &Paths{
				CDN:           "https://rox-conf.cloudbees.io",
				API:           "https://api.cloudbees.io/device/get_configuration",
				StateCDN:      "https://rox-state.cloudbees.io",
				StateAPI:      "https://api.cloudbees.io/device/update_state_store",
				Analytics:     "https://fm-analytics.cloudbees.io/impression",
				Notifications: "https://sdk-notification-service.cloudbees.io/sse",
			},
		},
		{
			testName: "Rollout QA API",
			osEnv:    "QA",
			platform: consts.ROLLOUT_API,
			expected: &Paths{
				CDN:           "https://qa-conf.rollout.io",
				API:           "https://qa-api.rollout.io/device/get_configuration",
				StateCDN:      "https://qa-statestore.rollout.io",
				StateAPI:      "https://qa-api.rollout.io/device/update_state_store",
				Analytics:     "https://qaanalytic.rollout.io",
				Notifications: "https://qax-push.rollout.io/sse",
			},
		},
		{
			testName: "Rollout Local API",
			osEnv:    "LOCAL",
			platform: consts.ROLLOUT_API,
			expected: &Paths{
				CDN:           "https://development-conf.rollout.io",
				API:           "http://127.0.0.1:8557/device/get_configuration",
				StateCDN:      "https://development-statestore.rollout.io",
				StateAPI:      "http://127.0.0.1:8557/device/update_state_store",
				Analytics:     "http://127.0.0.1:8787",
				Notifications: "http://127.0.0.1:8887/sse",
			},
		},
		{
			testName: "Rollout Production API",
			osEnv:    "",
			platform: consts.ROLLOUT_API,
			expected: &Paths{
				CDN:           "https://conf.rollout.io",
				API:           "https://x-api.rollout.io/device/get_configuration",
				StateCDN:      "https://statestore.rollout.io",
				StateAPI:      "https://api.cloudbees.io/device/update_state_store",
				Analytics:     "https://analytic.rollout.io",
				Notifications: "https://push.rollout.io/sse",
			},
		},
	}

	for _, tc := range testsCases {
		t.Run(tc.testName, func(t *testing.T) {
			os.Setenv("ROLLOUT_MODE", tc.osEnv)
			saasEnvironment := NewSaasEnvironment(tc.platform)
			assert.Equal(t, tc.expected.CDN, saasEnvironment.EnvironmentCDNPath())
			assert.Equal(t, tc.expected.API, saasEnvironment.EnvironmentAPIPath())
			assert.Equal(t, tc.expected.StateCDN, saasEnvironment.EnvironmentStateCDNPath())
			assert.Equal(t, tc.expected.StateAPI, saasEnvironment.EnvironmentStateAPIPath())
			assert.Equal(t, tc.expected.Analytics, saasEnvironment.EnvironmentAnalyticsPath())
			assert.Equal(t, tc.expected.Notifications, saasEnvironment.EnvironmentNotificationsPath())
		})

		os.Setenv("ROLLOUT_MODE", "")
	}
}
