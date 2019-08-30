package network

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/rollout/rox-go/core/configuration"
	"github.com/rollout/rox-go/core/consts"
	"github.com/rollout/rox-go/core/logging"
	"github.com/rollout/rox-go/core/model"
	"github.com/rollout/rox-go/core/properties"
	"github.com/rollout/rox-go/core/utils"
)

var relevantAPICallParams = []*consts.PropertyType{
	consts.PropertyTypePlatform,
	consts.PropertyTypeCustomProperties,
	consts.PropertyTypeFeatureFlags,
	consts.PropertyTypeRemoteVariables,
	consts.PropertyTypeDevModeSecret,
}

var stateGenerators = []consts.PropertyType{
	*consts.PropertyTypePlatform,
	*consts.PropertyTypeAppKey,
	*consts.PropertyTypeCustomProperties,
	*consts.PropertyTypeFeatureFlags,
	*consts.PropertyTypeRemoteVariables,
	*consts.PropertyTypeDevModeSecret,
}

type StateSender struct {
	customPropertyRepository model.CustomPropertyRepository
	deviceProperties         model.DeviceProperties
	flagRepository           model.FlagRepository
	request                  model.Request
	stateDebouncer           utils.Debouncer
}

func NewStateSender(r model.Request, deviceProperties model.DeviceProperties, flagRepository model.FlagRepository, customPropertyRepository model.CustomPropertyRepository) *StateSender {
	stateSender := &StateSender{
		customPropertyRepository: customPropertyRepository,
		deviceProperties:         deviceProperties,
		flagRepository:           flagRepository,
		request:                  r,
	}
	stateSender.stateDebouncer = *utils.NewDebouncer(3000, func() {
		stateSender.Send()
	})
	customPropertyRepository.RegisterPropertyAddedHandler(func(p *properties.CustomProperty) {
		stateSender.sendStateDebounce()
	})
	flagRepository.RegisterFlagAddedHandler(func(variant model.Variant) {
		stateSender.sendStateDebounce()
	})
	return stateSender
}

func getStateMd5(properties map[string]string) string {
	return utils.GenerateMD5(properties, stateGenerators)
}

func getPath(properties map[string]string) string {
	return fmt.Sprintf("%s/%s", properties[consts.PropertyTypeAppKey.Name], properties[consts.PropertyTypeStateMD5.Name])
}

func getCDNUrl(properties map[string]string) string {
	return fmt.Sprintf("%s/%s", consts.EnvironmentStateCDNPath(), properties[consts.PropertyTypeCacheMissRelativeURL.Name])
}

func getAPIUrl(properties map[string]string) string {
	return fmt.Sprintf("%s/%s", consts.EnvironmentStateAPIPath(), properties[consts.PropertyTypeCacheMissRelativeURL.Name])
}

func (s *StateSender) serializeFeatureFlags() (string, []jsonFlag) {
	var flags []jsonFlag
	allFlags := s.flagRepository.GetAllFlags()
	sort.Slice(allFlags, func(i, j int) bool {
		return allFlags[i].Name() < allFlags[j].Name()
	})
	for _, f := range allFlags {
		flags = append(flags, jsonFlag{f.Name(), f.DefaultValue(), f.Options()})
	}
	result, _ := json.Marshal(flags)

	return string(result), flags
}

func (s *StateSender) serializeCustomProperties() (string, []jsonProperty) {
	var properties []jsonProperty
	customProperties := s.customPropertyRepository.GetAllCustomProperties()
	sort.Slice(customProperties, func(i, j int) bool {
		return customProperties[i].Name < customProperties[j].Name
	})
	for _, p := range customProperties {
		properties = append(properties, jsonProperty{p.Name, p.Type.Type, p.Type.ExternalType})
	}
	result, _ := json.Marshal(properties)
	return string(result), properties
}

func (s *StateSender) sendStateToCDN(properties map[string]string) (response *model.Response, err error) {
	cdnRequest := model.RequestData{URL: getCDNUrl(properties), QueryParams: nil}
	return s.request.SendGet(cdnRequest)
}

func (s *StateSender) sendStateToAPI(properties map[string]string, featureFlags []jsonFlag, customProperties []jsonProperty) (response *model.Response, err error) {
	queryParams := make(map[string]interface{}, len(relevantAPICallParams))
	for _, prop := range relevantAPICallParams {
		propName := prop.Name

		if propName == consts.PropertyTypeFeatureFlags.Name {
			queryParams[propName] = featureFlags
		} else if propName == consts.PropertyTypeCustomProperties.Name {
			queryParams[propName] = customProperties
		} else {
			queryParams[propName] = properties[propName]
		}
	}

	return s.request.SendPost(getAPIUrl(properties), queryParams)
}

func (s *StateSender) preparePropsFromDeviceProps() (map[string]string, []jsonFlag /* feature flags*/, []jsonProperty /* custom properties */) {
	var featureFlags []jsonFlag
	var customProperties []jsonProperty
	properties := s.deviceProperties.GetAllProperties()
	properties[consts.PropertyTypeFeatureFlags.Name], featureFlags = s.serializeFeatureFlags()
	properties[consts.PropertyTypeRemoteVariables.Name] = ""
	properties[consts.PropertyTypeCustomProperties.Name], customProperties = s.serializeCustomProperties()

	stateMD5 := getStateMd5(properties)
	properties[consts.PropertyTypeStateMD5.Name] = stateMD5
	properties[consts.PropertyTypeCacheMissRelativeURL.Name] = getPath(properties)

	return properties, featureFlags, customProperties
}

func (s *StateSender) sendStateDebounce() {
	s.stateDebouncer.Invoke()
}

func (s *StateSender) Send() {
	properties, featureFlags, customProperties := s.preparePropsFromDeviceProps()
	shouldRetry := false
	source := configuration.SourceCDN

	fetchResult, err := s.sendStateToCDN(properties)

	if err != nil {
		s.logSendStateError(source, err)
		return
	}

	if fetchResult.IsSuccessStatusCode() {
		configurationFetchResult := configuration.NewFetchResult(string(fetchResult.Content), source)
		if configurationFetchResult == nil {
			s.logSendStateError(source, nil)
			return
		}

		if configurationFetchResult.ParsedData.Result == 404 {
			shouldRetry = true
		} else {
			// success from CDN
			return
		}
	}

	if shouldRetry || fetchResult.StatusCode == http.StatusForbidden || fetchResult.StatusCode == http.StatusNotFound {
		s.logSendStateErrorRetry(source, fetchResult, configuration.SourceAPI)
		source = configuration.SourceAPI
		fetchResult, err := s.sendStateToAPI(properties, featureFlags, customProperties)
		if err != nil {
			s.logSendStateError(source, err)
			return
		}

		if fetchResult.IsSuccessStatusCode() {
			// success for api
			return
		}
	}
}

func (s *StateSender) logSendStateErrorRetry(source configuration.Source, response *model.Response, nextSource configuration.Source) {
	retryMsg := fmt.Sprintf("Trying from %s. ", nextSource)
	logging.GetLogger().Debug(fmt.Sprintf("Failed to send state to %s. %shttp error code: %d\n", source, retryMsg, response.StatusCode), nil)
}

func (s *StateSender) logSendStateError(source configuration.Source, err error) {
	logging.GetLogger().Debug(fmt.Sprintf("Failed to send state. Source: %s", err), nil)
}

type jsonFlag struct {
	Name         string   `json:"name"`
	DefaultValue string   `json:"defaultValue"`
	Options      []string `json:"options"`
}

type jsonProperty struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	ExternalType string `json:"externalType"`
}
