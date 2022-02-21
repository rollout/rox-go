package client

import (
	"github.com/rollout/rox-go/v5/core/consts"
	"github.com/rollout/rox-go/v5/core/model"
)

type SelfManagedEnvironment struct {
	serverURL    string
	analyticsURL string
}

func NewSelfManagedEnvironment(options model.SelfManagedOptions) SelfManagedEnvironment {
	return SelfManagedEnvironment{
		serverURL:    options.ServerURL(),
		analyticsURL: options.AnalyticsURL(),
	}
}

func (e SelfManagedEnvironment) EnvironmentRoxyInternalPath() string {
	return consts.EnvironmentRoxyInternalPath()
}

func (e SelfManagedEnvironment) EnvironmentCDNPath() string {
	return ""
}

func (e SelfManagedEnvironment) EnvironmentAPIPath() string {
	return e.serverURL + "/device/get_configuration"
}

func (e SelfManagedEnvironment) EnvironmentStateCDNPath() string {
	return ""
}

func (e SelfManagedEnvironment) EnvironmentStateAPIPath() string {
	return e.serverURL + "/device/update_state_store"
}

func (e SelfManagedEnvironment) EnvironmentAnalyticsPath() string {
	return e.analyticsURL
}

func (e SelfManagedEnvironment) EnvironmentNotificationsPath() string {
	return ""
}

func (e SelfManagedEnvironment) IsSelfManaged() bool {
	return true
}
