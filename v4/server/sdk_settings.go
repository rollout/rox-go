package server

import "github.com/rollout/rox-go/v4/core/model"

type sdkSettings struct {
	apiKey        string
	devModeSecret string
}

func NewSdkSettings(apiKey, devModeSecret string) model.SdkSettings {
	return sdkSettings{
		apiKey:        apiKey,
		devModeSecret: devModeSecret,
	}
}

func (s sdkSettings) APIKey() string {
	return s.apiKey
}

func (s sdkSettings) DevModeSecret() string {
	return s.devModeSecret
}
