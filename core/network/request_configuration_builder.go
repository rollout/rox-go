package network

import (
	"fmt"
	"github.com/rollout/rox-go/core/consts"
	"github.com/rollout/rox-go/core/model"
	"net/url"
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
}

func NewRequestConfigurationBuilder(sdkSettings model.SdkSettings, buid model.BUID, deviceProperties model.DeviceProperties, roxyURL string) RequestConfigurationBuilder {
	return &requestConfigurationBuilder{
		sdkSettings:      sdkSettings,
		buid:             buid,
		deviceProperties: deviceProperties,
		roxyURL:          roxyURL,
	}
}

func (b *requestConfigurationBuilder) BuildForRoxy() model.RequestData {
	uri, _ := url.Parse(b.roxyURL)
	internalURI, _ := url.Parse(consts.EnvironmentRoxyInternalPath())
	uri = uri.ResolveReference(internalURI)
	return b.buildRequestWithFullParams(uri.String())
}

func (b *requestConfigurationBuilder) BuildForCDN() model.RequestData {
	return model.RequestData{
		fmt.Sprintf("%s/%s", consts.EnvironmentCDNPath(), b.buid.GetValue()),
		map[string]string{consts.PropertyTypeDistinctID.Name: b.deviceProperties.DistinctID()},
	}
}

func (b *requestConfigurationBuilder) BuildForAPI() model.RequestData {
	return b.buildRequestWithFullParams(consts.EnvironmentAPIPath())
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

	cdnData := b.BuildForCDN()
	queryParams[consts.PropertyTypeCacheMissURL.Name] = cdnData.URL
	queryParams["devModeSecret"] = b.sdkSettings.DevModeSecret()

	return model.RequestData{uri, queryParams}
}
