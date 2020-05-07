package server

import (
	"time"

	"github.com/rollout/rox-go/core/logging"
	"github.com/rollout/rox-go/core/model"
)

type RoxOptionsBuilder struct {
	Version                     string
	DevModeKey                  string
	FetchInterval               time.Duration
	Logger                      logging.Logger
	ImpressionHandler           model.ImpressionHandler
	ConfigurationFetchedHandler model.ConfigurationFetchedHandler
	RoxyURL                     string
	SelfManagedOptions          model.SelfManagedOptions
}

type roxOptions struct {
	version                     string
	devModeKey                  string
	fetchInterval               time.Duration
	impressionHandler           model.ImpressionHandler
	configurationFetchedHandler model.ConfigurationFetchedHandler
	roxyURL                     string
	selfManagedOptions          model.SelfManagedOptions
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

	if builder.Logger != nil {
		logging.SetLogger(builder.Logger)
	} else {
		logging.SetLogger(NewServerLogger())
	}

	return &roxOptions{
		version:                     version,
		devModeKey:                  devModeKey,
		fetchInterval:               fetchInterval,
		impressionHandler:           builder.ImpressionHandler,
		configurationFetchedHandler: builder.ConfigurationFetchedHandler,
		roxyURL:                     builder.RoxyURL,
		selfManagedOptions:          builder.SelfManagedOptions,
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
	return ro.roxyURL
}

func (ro *roxOptions) SelfManagedOptions() model.SelfManagedOptions {
	return ro.selfManagedOptions
}
