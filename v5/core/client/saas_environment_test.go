package client

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rollout/rox-go/v5/core/consts"
)

func TestSaasEnvironmentEnvironment_PlatformAPI_Path(t *testing.T) {
	testsCases := []struct {
		testName string
		osEnv    string
		platform consts.EnvironmentAPI
		expected string
	}{
		{
			testName: "Platform QA API",
			osEnv:    "QA",
			platform: consts.PLATFORM_API,
			expected: "https://api-staging.saas-dev.beescloud.com/events/flag-impressions",
		},
		{
			testName: "Platform Local API",
			osEnv:    "LOCAL",
			platform: consts.PLATFORM_API,
			expected: "http://127.0.0.1:8097/events/flag-impressions",
		},
		{
			testName: "Platform Production API",
			osEnv:    "",
			platform: consts.PLATFORM_API,
			expected: "https://api.cloudbees.io/events/flag-impressions",
		},
		{
			testName: "Rollout QA API",
			osEnv:    "QA",
			platform: consts.ROLLOUT_API,
			expected: "https://qaanalytic.rollout.io",
		},
		{
			testName: "Rollout Local API",
			osEnv:    "LOCAL",
			platform: consts.ROLLOUT_API,
			expected: "http://127.0.0.1:8787",
		},
		{
			testName: "Rollout Production API",
			osEnv:    "",
			platform: consts.ROLLOUT_API,
			expected: "https://analytic.rollout.io",
		},
	}

	for _, tc := range testsCases {
		t.Run(tc.testName, func(t *testing.T) {
			os.Setenv("ROLLOUT_MODE", tc.osEnv)
			saasEnvironment := NewSaasEnvironment(tc.platform)
			assert.Equal(t, tc.expected, saasEnvironment.EnvironmentAnalyticsPath())
		})

		os.Setenv("ROLLOUT_MODE", "")
	}
}
