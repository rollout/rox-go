package network

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/rollout/rox-go/core/configuration"
	"github.com/rollout/rox-go/core/consts"
	"github.com/rollout/rox-go/core/model"
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
	*consts.PropertyTypeCustomPropertiesString,
	*consts.PropertyTypeFeatureFlagsString,
	*consts.PropertyTypeRemoteVariablesString,
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
	return stateSender
}

func getStateMd5(properties map[string]string) string {
	return utils.GenerateMD5(properties, stateGenerators, nil)
}

func getPath(properties map[string]string) string {
	return fmt.Sprintf("%s/%s", properties[consts.PropertyTypePlatform.Name], properties[consts.PropertyTypeStateMD5.Name])
}

func getCDNUrl(properties map[string]string) string {
	return fmt.Sprintf("%s/%s", consts.EnvironmentStateCDNPath(), properties[consts.PropertyTypeCacheMissRelativeURL.Name])
}

func getAPIUrl(properties map[string]string) string {
	return fmt.Sprintf("%s/%s", consts.EnvironmentStateAPIPath(), properties[consts.PropertyTypeCacheMissRelativeURL.Name])
}

func (s *StateSender) serializeFeatureFlags() string {
	var flags []jsonFlag
	allFlags := s.flagRepository.GetAllFlags()
	sort.Slice(allFlags, func(i, j int) bool {
		return allFlags[i].Name() < allFlags[j].Name()
	})
	for _, f := range allFlags {
		flags = append(flags, jsonFlag{f.Name(), f.DefaultValue(), f.Options()})
	}
	result, _ := json.Marshal(flags)
	return string(result)
}

func (s *StateSender) serializeCustomProperties() string {
	var properties []jsonProperty
	customProperties := s.customPropertyRepository.GetAllCustomProperties()
	sort.Slice(customProperties, func(i, j int) bool {
		return customProperties[i].Name < customProperties[j].Name
	})
	for _, p := range customProperties {
		properties = append(properties, jsonProperty{p.Name, p.Type.Type, p.Type.ExternalType})
	}
	result, _ := json.Marshal(properties)
	return string(result)
}

func (s *StateSender) sendStateToCDN(properties map[string]string) (response *model.Response, err error) {
	cdnRequest := model.RequestData{URL: getCDNUrl(properties), QueryParams: nil}
	return s.request.SendGet(cdnRequest)
}

func (s *StateSender) sendStateToAPI(properties map[string]string) (response *model.Response, err error) {
	queryParams := make(map[string]string)
	for _, prop := range relevantAPICallParams {
		propName := prop.Name
		propValue := properties[propName]
		if propValue != "" {
			queryParams[propName] = propValue
		}
	}

	return s.request.SendPost(getAPIUrl(properties), queryParams)
}

func (s *StateSender) preparePropsFromDeviceProps() map[string]string {
	properties := s.deviceProperties.GetAllProperties()
	properties[consts.PropertyTypeFeatureFlags.Name] = s.serializeFeatureFlags()
	properties[consts.PropertyTypeRemoteVariables.Name] = ""
	properties[consts.PropertyTypeCustomProperties.Name] = s.serializeCustomProperties()

	properties[consts.PropertyTypeFeatureFlagsString.Name] = properties[consts.PropertyTypeFeatureFlags.Name]
	properties[consts.PropertyTypeRemoteVariablesString.Name] = properties[consts.PropertyTypeRemoteVariables.Name]
	properties[consts.PropertyTypeCustomPropertiesString.Name] = properties[consts.PropertyTypeCustomProperties.Name]

	stateMD5 := getStateMd5(properties)
	properties[consts.PropertyTypeStateMD5.Name] = stateMD5
	properties[consts.PropertyTypeCacheMissRelativeURL.Name] = getPath(properties)

	return properties
}

func (s *StateSender) Send() {
	properties := s.preparePropsFromDeviceProps()
	shouldRetry := false
	source := configuration.SourceCDN

	fetchResult, err := s.sendStateToCDN(properties)

	if err != nil {
		// TODO: log
		return
	}

	if fetchResult.IsSuccessStatusCode() {
		configurationFetchResult := configuration.NewFetchResult(string(fetchResult.Content), source)
		if configurationFetchResult == nil {
			// TODO: log
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
		// TODO: log
		source = configuration.SourceAPI
		fetchResult, err := s.sendStateToAPI(properties)
		if err != nil {
			// TODO log
			return
		}

		if fetchResult.IsSuccessStatusCode() {
			// success for api
			return
		}
	}
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
