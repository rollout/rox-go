package client

import "github.com/rollout/rox-go/core/model"

type selfManagedOptions struct {
	serverURL    string
	analyticsURL string
}

func NewSelfManagedOptions(serverURL string, analyticsURL string) model.SelfManagedOptions {
	return selfManagedOptions{
		serverURL:    serverURL,
		analyticsURL: analyticsURL,
	}
}

func (s selfManagedOptions) AnalyticsURL() string {
	return s.analyticsURL
}

func (s selfManagedOptions) ServerURL() string {
	return s.serverURL
}
