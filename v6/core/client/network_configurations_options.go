package client

import "github.com/rollout/rox-go/v6/core/model"

type NetworkConfigurationsBuilder struct {
	GetConfigApiEndpoint     string
	GetConfigCloudEndpoint   string
	SendStateApiEndpoint     string
	SendStateCloudEndpoint   string
	AnalyticsEndpoint        string
	PushNotificationEndpoint string
}

type NetworkConfigurations struct {
	getConfigApiEndpoint     string
	getConfigCloudEndpoint   string
	sendStateApiEndpoint     string
	sendStateCloudEndpoint   string
	analyticsEndpoint        string
	pushNotificationEndpoint string
}

func NewNetworkConfigurationsOptions(builder NetworkConfigurationsBuilder) model.NetworkConfigurationsOptions {
	return NetworkConfigurations{
		getConfigApiEndpoint:     builder.GetConfigApiEndpoint,
		getConfigCloudEndpoint:   builder.GetConfigCloudEndpoint,
		sendStateApiEndpoint:     builder.SendStateApiEndpoint,
		sendStateCloudEndpoint:   builder.SendStateCloudEndpoint,
		analyticsEndpoint:        builder.AnalyticsEndpoint,
		pushNotificationEndpoint: builder.PushNotificationEndpoint,
	}
}

func (s NetworkConfigurations) GetConfigApiEndpoint() string {
	return s.getConfigApiEndpoint
}

func (s NetworkConfigurations) GetConfigCloudEndpoint() string {
	return s.getConfigCloudEndpoint
}

func (s NetworkConfigurations) SendStateApiEndpoint() string {
	return s.sendStateApiEndpoint
}

func (s NetworkConfigurations) SendStateCloudEndpoint() string {
	return s.sendStateCloudEndpoint
}

func (s NetworkConfigurations) AnalyticsEndpoint() string {
	return s.analyticsEndpoint
}

func (s NetworkConfigurations) PushNotificationEndpoint() string {
	return s.pushNotificationEndpoint
}
