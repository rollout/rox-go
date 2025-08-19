package core

import (
	"net/http"
	"regexp"

	uuid "github.com/google/uuid"

	"github.com/rollout/rox-go/v5/core/analytics"
	"github.com/rollout/rox-go/v5/core/consts"
	"github.com/rollout/rox-go/v5/core/security"

	"github.com/rollout/rox-go/v5/core/client"
	"github.com/rollout/rox-go/v5/core/configuration"
	"github.com/rollout/rox-go/v5/core/context"
	"github.com/rollout/rox-go/v5/core/entities"
	"github.com/rollout/rox-go/v5/core/extensions"
	"github.com/rollout/rox-go/v5/core/impression"
	"github.com/rollout/rox-go/v5/core/model"
	"github.com/rollout/rox-go/v5/core/network"
	"github.com/rollout/rox-go/v5/core/notifications"
	"github.com/rollout/rox-go/v5/core/properties"
	"github.com/rollout/rox-go/v5/core/register"
	"github.com/rollout/rox-go/v5/core/reporting"
	"github.com/rollout/rox-go/v5/core/repositories"
	"github.com/rollout/rox-go/v5/core/roxx"
	"github.com/rollout/rox-go/v5/core/utils"
)

type Core struct {
	registerer                   *register.Registerer
	flagRepository               model.FlagRepository
	customPropertyRepository     model.CustomPropertyRepository
	experimentRepository         model.ExperimentRepository
	targetGroupRepository        model.TargetGroupRepository
	flagSetter                   *entities.FlagSetter
	parser                       roxx.Parser
	impressionInvoker            model.ImpressionInvoker
	analyticsHandler             model.Analytics
	configurationFetchedInvoker  *configuration.FetchedInvoker
	stateSender                  *network.StateSender
	sdkSettings                  model.SdkSettings
	configurationFetcher         network.ConfigurationFetcher
	errorReporter                model.ErrorReporter
	lastConfigurations           *configuration.FetchResult
	internalFlags                model.InternalFlags
	pushUpdatesListener          *notifications.NotificationListener
	environment                  model.Environment
	disableSignatureVerification bool
	quit                         chan struct{}
}

const invalidAPIKeyErrorMessage = "Invalid rollout apikey"

func NewCore() *Core {
	parser := roxx.NewParser()
	flagRepository := repositories.NewFlagRepository()
	targetGroupRepository := repositories.NewTargetGroupRepository()
	experimentRepository := repositories.NewExperimentRepository()
	customPropertyRepository := repositories.NewCustomPropertyRepository()

	return &Core{
		flagRepository:              flagRepository,
		customPropertyRepository:    customPropertyRepository,
		targetGroupRepository:       targetGroupRepository,
		experimentRepository:        experimentRepository,
		parser:                      parser,
		configurationFetchedInvoker: configuration.NewFetchedInvoker(),
		registerer:                  register.NewRegisterer(flagRepository),
		quit:                        make(chan struct{}),
	}
}

func (core *Core) Setup(sdkSettings model.SdkSettings, deviceProperties model.DeviceProperties, roxOptions model.RoxOptions) <-chan struct{} {
	core.sdkSettings = sdkSettings

	roxyPath := ""
	if roxOptions != nil && roxOptions.RoxyURL() != "" {
		roxyPath = roxOptions.RoxyURL()
	}

	if roxOptions != nil {
		core.disableSignatureVerification = roxOptions.IsSignatureVerificationDisabled()
	}
	envApi := consts.ROLLOUT_API
	if roxyPath == "" {
		validMongoIdPattern := "^[a-f\\d]{24}$"
		// Try to parse it as a mongo ID (rollout.io)
		matched, err := regexp.Match(validMongoIdPattern, []byte(sdkSettings.APIKey()))

		if err != nil || !matched {
			// try to parse it as a UUID (platform)
			_, err = uuid.Parse(sdkSettings.APIKey())
			if err != nil {
				panic(invalidAPIKeyErrorMessage)
			}
			envApi = consts.PLATFORM_API
			core.disableSignatureVerification = true
		}
	}

	if roxOptions != nil && roxOptions.SelfManagedOptions() != nil {
		core.environment = client.NewSelfManagedEnvironment(roxOptions.SelfManagedOptions())
	} else if roxOptions != nil && roxOptions.NetworkConfigurationsOptions() != nil {
		core.environment = client.NewCustomEnvironment(roxOptions.NetworkConfigurationsOptions())
	} else {
		core.environment = client.NewSaasEnvironment(envApi)
	}

	core.internalFlags = client.NewInternalFlags(core.experimentRepository, core.parser, core.environment)
	impressionDeps := &impression.ImpressionsDeps{
		InternalFlags:            core.internalFlags,
		CustomPropertyRepository: core.customPropertyRepository,
		DeviceProperties:         deviceProperties,
		IsRoxy:                   roxyPath != "",
	}
	analyticsEnabled := roxOptions != nil && !roxOptions.IsAnalyticsReportingDisabled() && roxyPath != ""
	if analyticsEnabled {
		analyticsHandler := analytics.NewAnalyticsHandler(&analytics.AnalyticsDeps{
			UriPath:          core.environment.EnvironmentAnalyticsPath(),
			Request:          network.NewRequest(http.DefaultClient),
			DeviceProperties: deviceProperties,
			FlushAtSize:      roxOptions.AnalyticsQueueSize(),
		})
		impressionDeps.Analytics = analyticsHandler
		core.analyticsHandler = analyticsHandler
		analyticsHandler.InitiateIntervalReporting(roxOptions.AnalyticsReportInterval())
	}
	core.impressionInvoker = impression.NewImpressionInvoker(impressionDeps)

	core.flagSetter = entities.NewFlagSetter(core.flagRepository, core.parser, core.experimentRepository, core.impressionInvoker)
	buid := client.NewBUID(sdkSettings, deviceProperties, core.flagRepository, core.customPropertyRepository)

	experimentsExtensions := extensions.NewExperimentsExtensions(core.parser, core.targetGroupRepository, core.flagRepository, core.experimentRepository)
	var dynamicPropertyRuleHandler model.DynamicPropertyRuleHandler
	if roxOptions != nil {
		dynamicPropertyRuleHandler = roxOptions.DynamicPropertyRuleHandler()
	}
	propertiesExtensions := extensions.NewPropertiesExtensions(core.parser, core.customPropertyRepository, dynamicPropertyRuleHandler)
	experimentsExtensions.Extend()
	propertiesExtensions.Extend()

	requestConfigBuilder := network.NewRequestConfigurationBuilder(sdkSettings, buid, deviceProperties, roxyPath, core.environment)

	// TODO http client
	clientRequest := network.NewRequest(http.DefaultClient)

	// TODO http client
	errReporterRequest := network.NewRequest(http.DefaultClient)
	core.errorReporter = reporting.NewErrorReporter(core.environment, errReporterRequest, deviceProperties, buid)

	if roxyPath != "" {
		core.configurationFetcher = network.NewConfigurationFetcherRoxy(requestConfigBuilder, clientRequest, core.configurationFetchedInvoker)
	} else {
		core.stateSender = network.NewStateSender(clientRequest, deviceProperties, core.flagRepository, core.customPropertyRepository, core.environment, core.disableSignatureVerification)
		core.configurationFetcher = network.NewConfigurationFetcher(core.environment, requestConfigBuilder, clientRequest, core.configurationFetchedInvoker)
	}

	var configurationFetchedHandler model.ConfigurationFetchedHandler
	if roxOptions != nil {
		configurationFetchedHandler = roxOptions.ConfigurationFetchedHandler()
	}
	core.configurationFetchedInvoker.RegisterFetchedHandler(core.wrapConfigurationFetchedHandler(configurationFetchedHandler))

	done := make(chan struct{})
	go func() {
		defer close(done)
		<-core.Fetch()

		if roxOptions != nil && roxOptions.ImpressionHandler() != nil {
			core.impressionInvoker.RegisterImpressionHandler(roxOptions.ImpressionHandler())
		}

		if roxOptions != nil && roxOptions.FetchInterval() != 0 {
			utils.RunPeriodicTask(func() {
				<-core.Fetch()
			}, roxOptions.FetchInterval(), core.quit)
		}
		if core.stateSender != nil {
			core.stateSender.Send()
		}
	}()
	return done
}

func (core *Core) Fetch() <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		select {
		default:
			if core.configurationFetcher == nil {
				return
			}

			result := core.configurationFetcher.Fetch()
			if result == nil {
				return
			}

			var signatureVerifier security.SignatureVerifier
			if core.disableSignatureVerification {
				signatureVerifier = security.NewDisabledSignatureVerifier()
			} else {
				signatureVerifier = security.NewSignatureVerifier(core.environment)
			}
			configurationParser := configuration.NewParser(signatureVerifier, core.errorReporter, core.configurationFetchedInvoker)
			config := configurationParser.Parse(result, core.sdkSettings)
			if config != nil {
				core.experimentRepository.SetExperiments(config.Experiments)
				core.targetGroupRepository.SetTargetGroups(config.TargetGroups)
				core.flagSetter.SetExperiments()

				hasChanges := core.lastConfigurations == nil || *core.lastConfigurations != *result
				core.lastConfigurations = result
				core.configurationFetchedInvoker.Invoke(model.FetcherStatusAppliedFromNetwork, config.SignatureDate, hasChanges)
			}
			return
		case <-core.quit:
			return
		}
	}()
	return done
}

func (core *Core) Register(ns string, roxContainer interface{}) {
	core.registerer.RegisterInstance(roxContainer, ns)
}

func (core *Core) SetContext(ctx context.Context) {
	for _, flag := range core.flagRepository.GetAllFlags() {
		flag.(model.InternalVariant).SetContext(ctx)
	}
}

func (core *Core) AddCustomProperty(property *properties.CustomProperty) {
	core.customPropertyRepository.AddCustomProperty(property)
}

func (core *Core) AddCustomPropertyIfNotExists(property *properties.CustomProperty) {
	core.customPropertyRepository.AddCustomPropertyIfNotExists(property)
}

func (core *Core) wrapConfigurationFetchedHandler(handler model.ConfigurationFetchedHandler) model.ConfigurationFetchedHandler {
	return func(args *model.ConfigurationFetchedArgs) {
		if args.FetcherStatus != model.FetcherStatusErrorFetchedFailed {
			core.startOrStopPushUpdatesListener()
		}

		if handler != nil {
			handler(args)
		}
	}
}

func (core *Core) startOrStopPushUpdatesListener() {

	if core.pushUpdatesListener == nil {
		core.pushUpdatesListener = notifications.NewNotificationListener(core.environment.EnvironmentNotificationsPath(), core.sdkSettings.APIKey())
		core.pushUpdatesListener.On("changed", func(event notifications.Event) {
			<-core.Fetch()
		})
		core.pushUpdatesListener.Start()
	}
}

func (core *Core) DynamicAPI(entitiesProvider model.EntitiesProvider) model.DynamicAPI {
	return client.NewDynamicAPI(core.flagRepository, entitiesProvider)
}

func (core *Core) Shutdown() <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)

		if core.pushUpdatesListener != nil {
			core.pushUpdatesListener.Stop()
			core.pushUpdatesListener = nil
		}
		if core.analyticsHandler != nil {
			core.analyticsHandler.StopIntervalReporting()
			core.analyticsHandler = nil
		}
		close(core.quit)
	}()

	return done
}
