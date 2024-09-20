package client

import (
	"github.com/rollout/rox-go/v6/core/consts"
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
	return consts.EnvironmentCDNPath(e.EnvironmentAPI)
}

func (e SaasEnvironment) EnvironmentAPIPath() string {
	return consts.EnvironmentAPIPath(e.EnvironmentAPI)
}

func (e SaasEnvironment) EnvironmentStateCDNPath() string {
	return consts.EnvironmentStateCDNPath(e.EnvironmentAPI)
}

func (e SaasEnvironment) EnvironmentStateAPIPath() string {
	return consts.EnvironmentStateAPIPath(e.EnvironmentAPI)
}

func (e SaasEnvironment) EnvironmentAnalyticsPath() string {
	return consts.EnvironmentAnalyticsPath(e.EnvironmentAPI)
}

func (e SaasEnvironment) EnvironmentNotificationsPath() string {
	return consts.EnvironmentNotificationsPath(e.EnvironmentAPI)
}

func (e SaasEnvironment) IsSelfManaged() bool {
	return false
}
