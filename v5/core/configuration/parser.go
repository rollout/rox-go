package configuration

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/rollout/rox-go/v5/core/logging"
	"github.com/rollout/rox-go/v5/core/model"
	"github.com/rollout/rox-go/v5/core/security"
)

type Parser struct {
	signatureVerifier security.SignatureVerifier
	errorReporter     model.ErrorReporter
	fetchedInvoker    *FetchedInvoker
}

func NewParser(signatureVerifier security.SignatureVerifier, errorReporter model.ErrorReporter, fetchedInvoker *FetchedInvoker) *Parser {
	return &Parser{
		signatureVerifier: signatureVerifier,
		errorReporter:     errorReporter,
		fetchedInvoker:    fetchedInvoker,
	}
}

func (cp *Parser) Parse(fetchResult *FetchResult, sdkSettings model.SdkSettings) (configuration *Configuration) {
	defer func() {
		if r := recover(); r != nil {
			logging.GetLogger().Error("Failed to parse configurations", r)
			cp.fetchedInvoker.InvokeError(model.FetcherErrorUnknown)
			configuration = nil
		}
	}()

	jsonConf := fetchResult.ParsedData

	if fetchResult.Source != SourceRoxy && !cp.signatureVerifier.Verify(jsonConf.Data, jsonConf.Signature) {
		cp.fetchedInvoker.InvokeError(model.FetcherErrorSignatureVerification)
		cp.errorReporter.Report("Failed to validate signature", fmt.Errorf("Data : %s Signature : %s", jsonConf.Data, jsonConf.Signature))
		return nil
	}

	signatureDate, err := time.Parse(time.RFC3339, jsonConf.SignedDate)
	if err != nil {
		cp.fetchedInvoker.InvokeError(model.FetcherErrorCorruptedJSON)
		cp.errorReporter.Report("Failed to parse signature date", fmt.Errorf("Signature date : %s", jsonConf.SignedDate))
		return nil
	}

	var internalJSONConf jsonInternalConfiguration
	err = json.Unmarshal([]byte(jsonConf.Data), &internalJSONConf)
	if err != nil {
		cp.fetchedInvoker.InvokeError(model.FetcherErrorCorruptedJSON)
		cp.errorReporter.Report("Failed to parse JSON configuration", err)
		return nil
	}

	if fetchResult.Source != SourceRoxy && internalJSONConf.Application != sdkSettings.APIKey() {
		cp.fetchedInvoker.InvokeError(model.FetcherErrorMismatchAppKey)
		cp.errorReporter.Report("Failed to parse JSON configuration - ", fmt.Errorf("Internal Data: %s SdkSettings: %s", internalJSONConf.Application, sdkSettings.APIKey()))
		return nil
	}

	experiments := cp.parseExperiments(internalJSONConf)
	groups := cp.parseGroups(internalJSONConf)

	return NewConfiguration(experiments, groups, signatureDate)
}

func (cp *Parser) parseExperiments(internalJSONConf jsonInternalConfiguration) []*model.ExperimentModel {
	experiments := make([]*model.ExperimentModel, len(internalJSONConf.Experiments))
	for i, e := range internalJSONConf.Experiments {
		flags := make([]string, len(e.Flags))
		for j, f := range e.Flags {
			flags[j] = f.Name
		}

		experiment := model.NewExperimentModel(e.ID, e.Name, e.DeploymentConfiguration.Condition, e.IsArchived, flags, e.Labels)
		experiments[i] = experiment
	}
	return experiments
}

func (cp *Parser) parseGroups(internalJSONConf jsonInternalConfiguration) []*model.TargetGroupModel {
	groups := make([]*model.TargetGroupModel, len(internalJSONConf.TargetGroups))
	for i, g := range internalJSONConf.TargetGroups {
		group := model.NewTargetGroupModel(g.ID, g.Condition)
		groups[i] = group
	}
	return groups
}

type jsonConfiguration struct {
	Data       string `json:"data"`
	Signature  string `json:"signature_v0"`
	SignedDate string `json:"signed_date"`
	Result     int    `json:"result"`
}

type jsonInternalConfiguration struct {
	Application  string            `json:"application"`
	Experiments  []jsonExperiment  `json:"experiments"`
	TargetGroups []jsonTargetGroup `json:"targetGroups"`
}

type jsonExperiment struct {
	ID                      string                       `json:"_id"`
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
	ID        string `json:"_id"`
	Condition string `json:"condition"`
}
