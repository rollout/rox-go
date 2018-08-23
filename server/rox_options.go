package server

import (
	"github.com/rollout/rox-go/core/model"
	"time"
)

type RoxOptionsBuilder struct {
	Version                     string
	DevModeKey                  string
	FetchInterval               time.Duration
	ImpressionHandler           model.ImpressionHandler
	ConfigurationFetchedHandler model.ConfigurationFetchedHandler
	RoxyUrl                     string
}

type roxOptions struct {
	version                     string
	devModeKey                  string
	fetchInterval               time.Duration
	impressionHandler           model.ImpressionHandler
	configurationFetchedHandler model.ConfigurationFetchedHandler
	roxyUrl                     string
}

func NewRoxOptions(builder RoxOptionsBuilder) model.RoxOptions {
	devModeKey := builder.DevModeKey
	if devModeKey == "" {
		devModeKey = "stam"
	}

	version := builder.Version
	if version == "" {
		version = "0.0"
	}

	fetchInterval := builder.FetchInterval
	if fetchInterval > 0 {
		if fetchInterval < 30*time.Second {
			fetchInterval = 30 * time.Second
		}
	} else {
		fetchInterval = 60 * time.Second
	}

	// TODO logger

	return &roxOptions{
		version:                     version,
		devModeKey:                  devModeKey,
		fetchInterval:               fetchInterval,
		impressionHandler:           builder.ImpressionHandler,
		configurationFetchedHandler: builder.ConfigurationFetchedHandler,
		roxyUrl:                     builder.RoxyUrl,
	}
}

func (ro *roxOptions) DevModeKey() string {
	return ro.devModeKey
}

func (ro *roxOptions) Version() string {
	return ro.version
}

func (ro *roxOptions) FetchInterval() time.Duration {
	return ro.fetchInterval
}

func (ro *roxOptions) ImpressionHandler() model.ImpressionHandler {
	return ro.impressionHandler
}

func (ro *roxOptions) ConfigurationFetchedHandler() model.ConfigurationFetchedHandler {
	return ro.configurationFetchedHandler
}

func (ro *roxOptions) RoxyURL() string {
	return ro.roxyUrl
}
