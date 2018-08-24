package core

import (
	"github.com/rollout/rox-go/core/client"
	"github.com/rollout/rox-go/core/configuration"
	"github.com/rollout/rox-go/core/context"
	"github.com/rollout/rox-go/core/entities"
	"github.com/rollout/rox-go/core/extensions"
	"github.com/rollout/rox-go/core/impression"
	"github.com/rollout/rox-go/core/model"
	"github.com/rollout/rox-go/core/network"
	"github.com/rollout/rox-go/core/properties"
	"github.com/rollout/rox-go/core/register"
	"github.com/rollout/rox-go/core/reporting"
	"github.com/rollout/rox-go/core/repositories"
	"github.com/rollout/rox-go/core/roxx"
	"github.com/rollout/rox-go/core/security"
	"github.com/rollout/rox-go/core/utils"
	"net/http"
)

type Core struct {
	registerer                  *register.Registerer
	flagRepository              model.FlagRepository
	customPropertyRepository    model.CustomPropertyRepository
	experimentRepository        model.ExperimentRepository
	targetGroupRepository       model.TargetGroupRepository
	flagSetter                  *entities.FlagSetter
	parser                      roxx.Parser
	impressionInvoker           model.ImpressionInvoker
	configurationFetchedInvoker *configuration.FetchedInvoker
	sdkSettings                 model.SdkSettings
	configurationFetcher        network.ConfigurationFetcher
	errorReporter               reporting.ErrorReporter
	lastConfigurations          *configuration.FetchResult
}

func NewCore() *Core {
	parser := roxx.NewParser()
	flagRepository := repositories.NewFlagRepository()
	targetGroupRepository := repositories.NewTargetGroupRepository()
	experimentRepository := repositories.NewExperimentRepository()
	customPropertyRepository := repositories.NewCustomPropertyRepository()

	experimentsExtensions := extensions.NewExperimentsExtensions(parser, targetGroupRepository, flagRepository, experimentRepository)
	propertiesExtensions := extensions.NewPropertiesExtensions(parser, customPropertyRepository)
	experimentsExtensions.Extend()
	propertiesExtensions.Extend()

	return &Core{
		flagRepository:           flagRepository,
		customPropertyRepository: customPropertyRepository,
		targetGroupRepository:    targetGroupRepository,
		experimentRepository:     experimentRepository,
		parser:                   parser,
		configurationFetchedInvoker: configuration.NewFetchedInvoker(),
		registerer:                  register.NewRegisterer(flagRepository),
	}
}

func (core *Core) Setup(sdkSettings model.SdkSettings, deviceProperties model.DeviceProperties, roxOptions model.RoxOptions) <-chan struct{} {
	core.sdkSettings = sdkSettings

	roxyPath := ""
	if roxOptions != nil && roxOptions.RoxyURL() != "" {
		roxyPath = roxOptions.RoxyURL()
	}

	// TODO Analytics.Analytics.Initialize(deviceProperties.RolloutKey, deviceProperties)

	internalFlags := client.NewInternalFlags(core.experimentRepository, core.parser)
	core.impressionInvoker = impression.NewImpressionInvoker(internalFlags, core.customPropertyRepository, deviceProperties /* TODO Analytics.Analytics.Client, */, roxyPath != "")
	core.flagSetter = entities.NewFlagSetter(core.flagRepository, core.parser, core.experimentRepository, core.impressionInvoker)
	buid := client.NewBUID(sdkSettings, deviceProperties, core.flagRepository, core.customPropertyRepository)

	requestConfigBuilder := network.NewRequestConfigurationBuilder(sdkSettings, buid, deviceProperties, roxyPath)

	// TODO gzip

	// TODO http client
	clientRequest := network.NewRequest(http.DefaultClient)

	// TODO errReporterRequest := network.NewRequest(http.DefaultClient)
	core.errorReporter = reporting.NewErrorReporter()

	if roxyPath != "" {
		core.configurationFetcher = network.NewConfigurationFetcherRoxy(requestConfigBuilder, clientRequest, core.configurationFetchedInvoker)
	} else {
		core.configurationFetcher = network.NewConfigurationFetcher(requestConfigBuilder, clientRequest, core.configurationFetchedInvoker)
	}

	if roxOptions != nil && roxOptions.ConfigurationFetchedHandler() != nil {
		core.configurationFetchedInvoker.RegisterFetchedHandler(roxOptions.ConfigurationFetchedHandler())
	}

	done := make(chan struct{})
	go func() {
		defer close(done)

		<-core.Fetch()

		if roxOptions != nil && roxOptions.ImpressionHandler() != nil {
			core.impressionInvoker.RegisterImpressionHandler(roxOptions.ImpressionHandler())
		}

		if roxOptions != nil && roxOptions.FetchInterval() != 0 {
			go utils.RunPeriodicTask(func() {
				<-core.Fetch()
			}, roxOptions.FetchInterval())
		}
	}()
	return done
}

func (core *Core) Fetch() <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)

		if core.configurationFetcher == nil {
			return
		}

		result := core.configurationFetcher.Fetch()
		if result == nil {
			return
		}

		configurationParser := configuration.NewParser(security.NewSignatureVerifier(), core.errorReporter, core.configurationFetchedInvoker)
		config := configurationParser.Parse(result, core.sdkSettings)
		if config != nil {
			core.experimentRepository.SetExperiments(config.Experiments)
			core.targetGroupRepository.SetTargetGroups(config.TargetGroups)
			core.flagSetter.SetExperiments()

			hasChanges := core.lastConfigurations == nil || core.lastConfigurations.Data != result.Data
			core.lastConfigurations = result
			core.configurationFetchedInvoker.Invoke(model.FetcherStatusAppliedFromNetwork, config.SignatureDate, hasChanges)
		}
	}()
	return done
}

func (core *Core) Register(ns string, roxContainer interface{}) {
	core.registerer.RegisterInstance(roxContainer, ns)
}

func (core *Core) SetContext(ctx context.Context) {
	for _, flag := range core.flagRepository.GetAllFlags() {
		flag.SetContext(ctx)
	}
}

func (core *Core) AddCustomProperty(property *properties.CustomProperty) {
	core.customPropertyRepository.AddCustomProperty(property)
}

func (core *Core) AddCustomPropertyIfNotExists(property *properties.CustomProperty) {
	core.customPropertyRepository.AddCustomPropertyIfNotExists(property)
}
