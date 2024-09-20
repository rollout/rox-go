package client

import (
	"strings"

	"github.com/rollout/rox-go/v6/core/model"
)

type CustomEnvironment struct {
	getConfigApiURL     string
	getConfigCloudURL   string
	sendStateApiURL     string
	sendStateCloudURL   string
	analyticsURL        string
	pushNotificationURL string
}

func TrimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

func NewCustomEnvironment(options model.NetworkConfigurationsOptions) CustomEnvironment {
	return CustomEnvironment{
		getConfigApiURL:     TrimSuffix(options.GetConfigApiEndpoint(), "/"),
		getConfigCloudURL:   TrimSuffix(options.GetConfigCloudEndpoint(), "/"),
		sendStateApiURL:     TrimSuffix(options.SendStateApiEndpoint(), "/"),
		sendStateCloudURL:   TrimSuffix(options.SendStateCloudEndpoint(), "/"),
		analyticsURL:        TrimSuffix(options.AnalyticsEndpoint(), "/"),
		pushNotificationURL: TrimSuffix(options.PushNotificationEndpoint(), "/"),
	}
}

func (e CustomEnvironment) EnvironmentRoxyInternalPath() string {
	return "device/request_configuration"
}

func (e CustomEnvironment) EnvironmentCDNPath() string {
	return e.getConfigCloudURL
}

func (e CustomEnvironment) EnvironmentAPIPath() string {
	return e.getConfigApiURL
}

func (e CustomEnvironment) EnvironmentStateCDNPath() string {
	return e.sendStateCloudURL
}

func (e CustomEnvironment) EnvironmentStateAPIPath() string {
	return e.sendStateApiURL
}

func (e CustomEnvironment) EnvironmentAnalyticsPath() string {
	return e.analyticsURL
}

func (e CustomEnvironment) EnvironmentNotificationsPath() string {
	return e.pushNotificationURL
}

func (e CustomEnvironment) IsSelfManaged() bool {
	return false
}
