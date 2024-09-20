package network

import (
	"fmt"
	"net/url"

	"github.com/rollout/rox-go/v6/core/consts"
	"github.com/rollout/rox-go/v6/core/model"
)

type RequestConfigurationBuilder interface {
	BuildForRoxy() model.RequestData
	BuildForCDN() model.RequestData
	BuildForAPI() model.RequestData
}

type requestConfigurationBuilder struct {
	sdkSettings      model.SdkSettings
	buid             model.BUID
	deviceProperties model.DeviceProperties
	roxyURL          string
	environment      model.Environment
}

func NewRequestConfigurationBuilder(sdkSettings model.SdkSettings, buid model.BUID, deviceProperties model.DeviceProperties, roxyURL string, environment model.Environment) RequestConfigurationBuilder {
	return &requestConfigurationBuilder{
		sdkSettings:      sdkSettings,
		buid:             buid,
		deviceProperties: deviceProperties,
		roxyURL:          roxyURL,
		environment:      environment,
	}
}

func (b *requestConfigurationBuilder) BuildForRoxy() model.RequestData {
	uri, _ := url.Parse(b.roxyURL)
	internalURI, _ := url.Parse(b.environment.EnvironmentRoxyInternalPath())
	uri = uri.ResolveReference(internalURI)
	return b.buildRequestWithFullParams(uri.String())
}

func (b *requestConfigurationBuilder) GetPath() string {
	return fmt.Sprintf("%s/%s", b.deviceProperties.RolloutKey(), b.buid.GetValue())
}

func (b *requestConfigurationBuilder) BuildForCDN() model.RequestData {
	return model.RequestData{
		fmt.Sprintf("%s/%s", b.environment.EnvironmentCDNPath(), b.GetPath()),
		map[string]string{consts.PropertyTypeDistinctID.Name: b.deviceProperties.DistinctID()},
	}
}

func (b *requestConfigurationBuilder) BuildForAPI() model.RequestData {
	return b.buildRequestWithFullParams(fmt.Sprintf("%s/%s", b.environment.EnvironmentAPIPath(), b.GetPath()))
}

func (b *requestConfigurationBuilder) buildRequestWithFullParams(uri string) model.RequestData {
	queryParams := make(map[string]string)

	for k, v := range b.buid.GetQueryStringParts() {
		if _, ok := queryParams[k]; !ok {
			queryParams[k] = v
		}
	}

	for k, v := range b.deviceProperties.GetAllProperties() {
		if _, ok := queryParams[k]; !ok {
			queryParams[k] = v
		}
	}

	queryParams[consts.PropertyTypeCacheMissRelativeURL.Name] = b.GetPath()
	queryParams[consts.PropertyTypeDevModeSecret.Name] = b.sdkSettings.DevModeSecret()

	return model.RequestData{uri, queryParams}
}
