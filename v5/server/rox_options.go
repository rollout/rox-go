package server

import (
	"time"

	"github.com/rollout/rox-go/v5/core/logging"
	"github.com/rollout/rox-go/v5/core/model"
)

type RoxOptionsBuilder struct {
	Version                      string
	DevModeKey                   string
	FetchInterval                time.Duration
	AnalyticsReportInterval      time.Duration
	DisableAnalyticsReporting    bool
	Logger                       logging.Logger
	ImpressionHandler            model.ImpressionHandler
	ConfigurationFetchedHandler  model.ConfigurationFetchedHandler
	RoxyURL                      string
	SelfManagedOptions           model.SelfManagedOptions
	DynamicPropertyRuleHandler   model.DynamicPropertyRuleHandler
	NetworkConfigurationsOptions model.NetworkConfigurationsOptions
	DisableSignatureVerification bool
}

type roxOptions struct {
	version                      string
	devModeKey                   string
	fetchInterval                time.Duration
	analyticsReportInterval      time.Duration
	disableAnalyticsReporting    bool
	impressionHandler            model.ImpressionHandler
	configurationFetchedHandler  model.ConfigurationFetchedHandler
	roxyURL                      string
	selfManagedOptions           model.SelfManagedOptions
	dynamicPropertyRuleHandler   model.DynamicPropertyRuleHandler
	networkConfigurationsOptions model.NetworkConfigurationsOptions
	disableSignatureVerification bool
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

	if builder.AnalyticsReportInterval == 0 {
		builder.DisableAnalyticsReporting = true
	}

	var dynamicPropertyRuleHandler = builder.DynamicPropertyRuleHandler
	if dynamicPropertyRuleHandler == nil {
		dynamicPropertyRuleHandler = func(args model.DynamicPropertyRuleHandlerArgs) interface{} {
			if args.Context != nil {
				return args.Context.Get(args.PropName)
			}
			return nil
		}
	}

	return &roxOptions{
		version:                      version,
		devModeKey:                   devModeKey,
		fetchInterval:                fetchInterval,
		analyticsReportInterval:      builder.AnalyticsReportInterval,
		disableAnalyticsReporting:    builder.DisableAnalyticsReporting,
		impressionHandler:            builder.ImpressionHandler,
		configurationFetchedHandler:  builder.ConfigurationFetchedHandler,
		roxyURL:                      builder.RoxyURL,
		selfManagedOptions:           builder.SelfManagedOptions,
		dynamicPropertyRuleHandler:   dynamicPropertyRuleHandler,
		networkConfigurationsOptions: builder.NetworkConfigurationsOptions,
		disableSignatureVerification: builder.DisableSignatureVerification,
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

func (ro *roxOptions) DynamicPropertyRuleHandler() model.DynamicPropertyRuleHandler {
	return ro.dynamicPropertyRuleHandler
}

func (ro *roxOptions) NetworkConfigurationsOptions() model.NetworkConfigurationsOptions {
	return ro.networkConfigurationsOptions
}

func (ro *roxOptions) IsSignatureVerificationDisabled() bool {
	return ro.disableSignatureVerification
}

func (ro *roxOptions) AnalyticsReportInterval() time.Duration {
	return ro.analyticsReportInterval
}

func (ro *roxOptions) IsAnalyticsReportingDisabled() bool {
	return ro.disableAnalyticsReporting
}
