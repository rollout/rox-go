package model

type Environment interface {
	EnvironmentRoxyInternalPath() string
	EnvironmentCDNPath() string
	EnvironmentAPIPath() string
	EnvironmentStateCDNPath() string
	EnvironmentStateAPIPath() string
	EnvironmentAnalyticsPath() string
	EnvironmentNotificationsPath() string
	IsSelfManaged() bool
}
