package client

import "github.com/rollout/rox-go/core/model"

type SelfManagedOptionsBuilder struct {
	ServerURL    string
	AnalyticsURL string
}

type selfManagedOptions struct {
	serverURL    string
	analyticsURL string
}

func NewSelfManagedOptions(builder SelfManagedOptionsBuilder) model.SelfManagedOptions {
	return selfManagedOptions{
		serverURL:    builder.ServerURL,
		analyticsURL: builder.AnalyticsURL,
	}
}

func (s selfManagedOptions) AnalyticsURL() string {
	return s.analyticsURL
}

func (s selfManagedOptions) ServerURL() string {
	return s.serverURL
}
