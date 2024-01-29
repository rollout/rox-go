package client

import (
	"github.com/rollout/rox-go/v5/core/consts"
)

type SaasEnvironment struct {
	EnvironmentAPI consts.EnvironmentAPI
}

func NewSaasEnvironment(envAPI consts.EnvironmentAPI) SaasEnvironment {
	return SaasEnvironment{
		EnvironmentAPI: envAPI,
	}
}

func (e SaasEnvironment) EnvironmentRoxyInternalPath() string {
	return consts.EnvironmentRoxyInternalPath()
}

func (e SaasEnvironment) EnvironmentCDNPath() string {
	return consts.EnvironmentCDNPath()
}

func (e SaasEnvironment) EnvironmentAPIPath() string {
	return consts.EnvironmentAPIPath()
}

func (e SaasEnvironment) EnvironmentStateCDNPath() string {
	return consts.EnvironmentStateCDNPath()
}

func (e SaasEnvironment) EnvironmentStateAPIPath() string {
	return consts.EnvironmentStateAPIPath()
}

func (e SaasEnvironment) EnvironmentAnalyticsPath() string {
	return consts.EnvironmentAnalyticsPath(e.EnvironmentAPI)
}

func (e SaasEnvironment) EnvironmentNotificationsPath() string {
	return consts.EnvironmentNotificationsPath()
}

func (e SaasEnvironment) IsSelfManaged() bool {
	return false
}
