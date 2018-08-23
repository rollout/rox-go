package network

import (
	"fmt"
	"github.com/rollout/rox-go/core/consts"
	"github.com/rollout/rox-go/core/model"
	"net/url"
)

type RequestConfigurationBuilder interface {
	BuildForRoxy() RequestData
	BuildForCDN() RequestData
	BuildForAPI() RequestData
}

type requestConfigurationBuilder struct {
	sdkSettings      model.SdkSettings
	buid             model.BUID
	deviceProperties model.DeviceProperties
	roxyUrl          string
}

func NewRequestConfigurationBuilder(sdkSettings model.SdkSettings, buid model.BUID, deviceProperties model.DeviceProperties, roxyUrl string) RequestConfigurationBuilder {
	return &requestConfigurationBuilder{
		sdkSettings:      sdkSettings,
		buid:             buid,
		deviceProperties: deviceProperties,
		roxyUrl:          roxyUrl,
	}
}

func (b *requestConfigurationBuilder) BuildForRoxy() RequestData {
	uri, _ := url.Parse(b.roxyUrl)
	internalUri, _ := url.Parse(consts.EnvironmentRoxyInternalPath())
	uri = uri.ResolveReference(internalUri)
	return b.buildRequestWithFullParams(uri.String())
}

func (b *requestConfigurationBuilder) BuildForCDN() RequestData {
	return RequestData{
		fmt.Sprintf("%s/%s", consts.EnvironmentCDNPath(), b.buid.GetValue()),
		map[string]string{consts.PropertyTypeDistinctId.Name: b.deviceProperties.DistinctId()},
	}
}

func (b *requestConfigurationBuilder) BuildForAPI() RequestData {
	return b.buildRequestWithFullParams(consts.EnvironmentAPIPath())
}

func (b *requestConfigurationBuilder) buildRequestWithFullParams(uri string) RequestData {
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
	queryParams[consts.PropertyTypeCacheMissUrl.Name] = cdnData.Url
	queryParams["devModeSecret"] = b.sdkSettings.DevModeSecret()

	return RequestData{uri, queryParams}
}
