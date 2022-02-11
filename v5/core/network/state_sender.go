package network

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"sort"

	"github.com/rollout/rox-go/v5/core/configuration"
	"github.com/rollout/rox-go/v5/core/consts"
	"github.com/rollout/rox-go/v5/core/logging"
	"github.com/rollout/rox-go/v5/core/model"
	"github.com/rollout/rox-go/v5/core/properties"
	"github.com/rollout/rox-go/v5/core/utils"
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
	environment              model.Environment
}

func NewStateSender(r model.Request, deviceProperties model.DeviceProperties, flagRepository model.FlagRepository, customPropertyRepository model.CustomPropertyRepository, environment model.Environment) *StateSender {
	stateSender := &StateSender{
		customPropertyRepository: customPropertyRepository,
		deviceProperties:         deviceProperties,
		flagRepository:           flagRepository,
		request:                  r,
		environment:              environment,
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

func (s *StateSender) getCDNUrl(properties map[string]string) string {
	return fmt.Sprintf("%s/%s", s.environment.EnvironmentStateCDNPath(), properties[consts.PropertyTypeCacheMissRelativeURL.Name])
}

func (s *StateSender) getAPIUrl(properties map[string]string) string {
	return fmt.Sprintf("%s/%s", s.environment.EnvironmentStateAPIPath(), properties[consts.PropertyTypeCacheMissRelativeURL.Name])
}

func (s *StateSender) serializeFeatureFlags() (string, []jsonFlag) {
	var flags []jsonFlag
	allFlags := s.flagRepository.GetAllFlags()
	sort.Slice(allFlags, func(i, j int) bool {
		return allFlags[i].Name() < allFlags[j].Name()
	})
	for _, f := range allFlags {
		switch f.FlagType() {
		case consts.BoolType:
			options := optionsToInterface(f.(model.Flag).Options(), consts.BoolType)
			flags = append(flags, jsonFlag{Name: f.Name(), DefaultValue: f.(model.Flag).DefaultValue(), Options: options})
		case consts.StringType:
			options := optionsToInterface(f.(model.RoxString).Options(), consts.StringType)
			flags = append(flags, jsonFlag{f.Name(), f.GetDefaultAsString(), options})
		case consts.IntType:
			options := optionsToInterface(f.(model.RoxInt).Options(), consts.IntType)
			flags = append(flags, jsonFlag{f.Name(), f.(model.RoxInt).DefaultValue(), options})
		case consts.DoubleType:
			options := optionsToInterface(f.(model.RoxDouble).Options(), consts.DoubleType)
			flags = append(flags, jsonFlag{f.Name(), f.(model.RoxDouble).DefaultValue(), options})

		}
	}
	result, _ := json.Marshal(flags)

	return string(result), flags
}

func optionsToInterface(options interface{}, flagType int) []interface{} {

	var result []interface{}

	switch reflect.TypeOf(options).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(options)
		result = make([]interface{}, s.Len())
		for i := 0; i < s.Len(); i++ {
			switch flagType {
			case consts.BoolType:
				result[i] = s.Index(i).String()
			case consts.StringType:
				result[i] = s.Index(i).String()
			case consts.IntType:
				result[i] = s.Index(i).Int()
			case consts.DoubleType:
				result[i] = s.Index(i).Float()

			}
		}
	}
	return result
}

func test(t interface{}) {
	switch reflect.TypeOf(t).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(t)

		for i := 0; i < s.Len(); i++ {
			fmt.Println(s.Index(i))
		}
	}
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
	cdnRequest := model.RequestData{URL: s.getCDNUrl(properties), QueryParams: nil}
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

	return s.request.SendPost(s.getAPIUrl(properties), queryParams)
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

	var fetchResult *model.Response = nil
	var err error = nil
	isSelfManaged := s.environment.IsSelfManaged()

	if !isSelfManaged {
		fetchResult, err = s.sendStateToCDN(properties)

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
	}

	if isSelfManaged || shouldRetry || fetchResult.StatusCode == http.StatusForbidden || fetchResult.StatusCode == http.StatusNotFound {
		if !isSelfManaged {
			s.logSendStateErrorRetry(source, fetchResult, configuration.SourceAPI)
		}
		source = configuration.SourceAPI
		fetchResult, err = s.sendStateToAPI(properties, featureFlags, customProperties)
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
	Name         string        `json:"name"`
	DefaultValue interface{}   `json:"defaultValue"`
	Options      []interface{} `json:"options"`
}

type jsonProperty struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	ExternalType string `json:"externalType"`
}
