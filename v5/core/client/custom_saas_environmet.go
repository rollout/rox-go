package client

import (
	"strings"

	"github.com/rollout/rox-go/v5/core/model"
)

type CustomSaasEnvironment struct{
	getConfigApiURL string
	getConfigCloudURL string
	sendStateApiURL string
	sendStateCloudURL string
	analyticsURL string
	pushNotificationURL string
}

func TrimSuffix(s, suffix string) string {
    if strings.HasSuffix(s, suffix) {
        s = s[:len(s)-len(suffix)]
    }
    return s
}

func NewCustomSaasEnvironment(options model.NetworkConfigurationsOptions) CustomSaasEnvironment {
	return CustomSaasEnvironment{
		getConfigApiURL:    TrimSuffix(options.GetConfigApiEndpoint(), "/"),
		getConfigCloudURL: 	TrimSuffix(options.GetConfigCloudEndpoint(), "/"),
		sendStateApiURL: 	TrimSuffix(options.SendStateApiEndpoint(), "/"),
		sendStateCloudURL: 	TrimSuffix(options.SendStateCloudEndpoint(), "/"),
		analyticsURL: 		TrimSuffix(options.AnalyticsEndpoint(), "/"),
		pushNotificationURL:TrimSuffix(options.PushNotificationEndpoint(), "/"),
	}
}

func (e CustomSaasEnvironment) EnvironmentRoxyInternalPath() string {
	return "device/request_configuration"
}

func (e CustomSaasEnvironment) EnvironmentCDNPath() string {
	return e.getConfigCloudURL
}

func (e CustomSaasEnvironment) EnvironmentAPIPath() string {
	return e.getConfigApiURL
}

func (e CustomSaasEnvironment) EnvironmentStateCDNPath() string {
	return e.sendStateCloudURL
}

func (e CustomSaasEnvironment) EnvironmentStateAPIPath() string {
	return e.sendStateApiURL
}

func (e CustomSaasEnvironment) EnvironmentAnalyticsPath() string {
	return e.analyticsURL
}

func (e CustomSaasEnvironment) EnvironmentNotificationsPath() string {
	return e.pushNotificationURL
}

func (e CustomSaasEnvironment) IsSelfManaged() bool {
	return false
}
