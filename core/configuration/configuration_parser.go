package configuration

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rollout/rox-go/core/model"
	"github.com/rollout/rox-go/core/reporting"
	"github.com/rollout/rox-go/core/security"
	"time"
)

type ConfigurationParser struct {
	signatureVerifier           security.SignatureVerifier
	errorReporter               reporting.ErrorReporter
	configurationFetchedInvoker *ConfigurationFetchedInvoker
}

func NewConfigurationParser(signatureVerifier security.SignatureVerifier, errorReporter reporting.ErrorReporter, configurationFetchedInvoker *ConfigurationFetchedInvoker) *ConfigurationParser {
	return &ConfigurationParser{
		signatureVerifier:           signatureVerifier,
		errorReporter:               errorReporter,
		configurationFetchedInvoker: configurationFetchedInvoker,
	}
}

func (cp *ConfigurationParser) Parse(fetchResult *ConfigurationFetchResult, sdkSettings model.SdkSettings) (configuration *Configuration) {
	defer func() {
		if r := recover(); r != nil {
			// TODO logger
			fmt.Printf("Failed to parse configurations: %s\n", r)
			cp.configurationFetchedInvoker.InvokeError(model.FetcherErrorUnknown)
			configuration = nil
		}
	}()

	if fetchResult == nil || fetchResult.Data == "" {
		cp.configurationFetchedInvoker.InvokeError(model.FetcherErrorEmptyJson)
		cp.errorReporter.Report("Failed to parse JSON configuration - Null Or Empty", errors.New("null data"))
		return nil
	}

	var jsonConf jsonConfiguration
	err := json.Unmarshal([]byte(fetchResult.Data), &jsonConf)
	if err != nil {
		cp.configurationFetchedInvoker.InvokeError(model.FetcherErrorCorruptedJson)
		cp.errorReporter.Report("Failed to parse JSON configuration", err)
		return nil
	}

	if fetchResult.Source != SourceRoxy && !cp.signatureVerifier.Verify(jsonConf.Data, jsonConf.Signature) {
		cp.configurationFetchedInvoker.InvokeError(model.FetcherErrorSignatureVerification)
		cp.errorReporter.Report("Failed to validate signature", fmt.Errorf("Data : %s Signature : %s", jsonConf.Data, jsonConf.Signature))
		return nil
	}

	signatureDate, err := time.Parse(time.RFC3339, jsonConf.SignedDate)
	if err != nil {
		cp.configurationFetchedInvoker.InvokeError(model.FetcherErrorCorruptedJson)
		cp.errorReporter.Report("Failed to parse signature date", fmt.Errorf("Signature date : %s", jsonConf.SignedDate))
		return nil
	}

	var internalJsonConf jsonInternalConfiguration
	err = json.Unmarshal([]byte(jsonConf.Data), &internalJsonConf)
	if err != nil {
		cp.configurationFetchedInvoker.InvokeError(model.FetcherErrorCorruptedJson)
		cp.errorReporter.Report("Failed to parse JSON configuration", err)
		return nil
	}

	if fetchResult.Source != SourceRoxy && internalJsonConf.Application != sdkSettings.ApiKey() {
		cp.configurationFetchedInvoker.InvokeError(model.FetcherErrorMismatchAppKey)
		cp.errorReporter.Report("Failed to parse JSON configuration - ", fmt.Errorf("Internal Data: %s SdkSettings: %s", internalJsonConf.Application, sdkSettings.ApiKey()))
		return nil
	}

	experiments := cp.parseExperiments(internalJsonConf)
	groups := cp.parseGroups(internalJsonConf)

	return NewConfiguration(experiments, groups, signatureDate)
}

func (cp *ConfigurationParser) parseExperiments(internalJsonConf jsonInternalConfiguration) []*model.ExperimentModel {
	experiments := make([]*model.ExperimentModel, len(internalJsonConf.Experiments))
	for i, e := range internalJsonConf.Experiments {
		flags := make([]string, len(e.Flags))
		for j, f := range e.Flags {
			flags[j] = f.Name
		}

		experiment := model.NewExperimentModel(e.Id, e.Name, e.DeploymentConfiguration.Condition, e.IsArchived, flags, e.Labels)
		experiments[i] = experiment
	}
	return experiments
}

func (cp *ConfigurationParser) parseGroups(internalJsonConf jsonInternalConfiguration) []*model.TargetGroupModel {
	groups := make([]*model.TargetGroupModel, len(internalJsonConf.TargetGroups))
	for i, g := range internalJsonConf.TargetGroups {
		group := model.NewTargetGroupModel(g.Id, g.Condition)
		groups[i] = group
	}
	return groups
}

type jsonConfiguration struct {
	Data       string `json:"data"`
	Signature  string `json:"signature_v0"`
	SignedDate string `json:"signed_date"`
}

type jsonInternalConfiguration struct {
	Application  string            `json:"application"`
	Experiments  []jsonExperiment  `json:"experiments"`
	TargetGroups []jsonTargetGroup `json:"targetGroups"`
}

type jsonExperiment struct {
	Id                      string                       `json:"_id"`
	Name                    string                       `json:"name"`
	IsArchived              bool                         `json:"archived"`
	Labels                  []string                     `json:"labels"`
	Flags                   []jsonExperimentFlag         `json:"featureFlags"`
	DeploymentConfiguration *jsonDeploymentConfiguration `json:"deploymentConfiguration"`
}

type jsonDeploymentConfiguration struct {
	Condition string `json:"condition"`
}

type jsonExperimentFlag struct {
	Name string `json:"name"`
}

type jsonTargetGroup struct {
	Id        string `json:"_id"`
	Condition string `json:"condition"`
}
