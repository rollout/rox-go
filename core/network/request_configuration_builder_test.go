package network_test

import (
	"fmt"
	"github.com/rollout/rox-go/core/consts"
	"github.com/rollout/rox-go/core/network"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRequestConfigurationBuilderCDNRequestDataWillHaveDistinctId(t *testing.T) {
	sdkSettings := &mockedSdkSettings{}
	sdkSettings.On("DevModeSecret").Return("1")

	buid := &mockedBUID{}
	buid.On("GetValue").Return("123")

	deviceProps := &mockedDeviceProperties{}
	deviceProps.On("DistinctId").Return("123")
	deviceProps.On("GetAllProperties").Return(map[string]string{
		"app_key":        "123",
		"api_version":    "4.0.0",
		"cache_miss_url": "harta",
		"distinct_id":    "123",
	})

	requestConfiguraitonBuilder := network.NewRequestConfigurationBuilder(sdkSettings, buid, deviceProps, "")
	result := requestConfiguraitonBuilder.BuildForCDN()

	assert.Equal(t, fmt.Sprintf("%s/123", consts.EnvironmentCDNPath()), result.Url)
	assert.Equal(t, "123", result.QueryParams["distinct_id"])
}

func TestRequestConfigurationBuilderRoxyRequestDataWillHaveServerData(t *testing.T) {
	sdkSettings := &mockedSdkSettings{}
	sdkSettings.On("DevModeSecret").Return("1")

	buid := &mockedBUID{}
	buid.On("GetValue").Return("123")
	buid.On("GetQueryStringParts").Return(map[string]string{
		"buid": "123",
	})

	deviceProps := &mockedDeviceProperties{}
	deviceProps.On("DistinctId").Return("123")
	deviceProps.On("GetAllProperties").Return(map[string]string{
		"app_key":     "123",
		"api_version": "4.0.0",
		"distinct_id": "123",
	})

	requestConfigurationBuilder := network.NewRequestConfigurationBuilder(sdkSettings, buid, deviceProps, "http://bimba.bobi.o.ponpon")
	result := requestConfigurationBuilder.BuildForRoxy()

	assert.Equal(t, "http://bimba.bobi.o.ponpon/device/request_configuration", result.Url)
	assert.Equal(t, "123", result.QueryParams["app_key"])
	assert.Equal(t, "4.0.0", result.QueryParams["api_version"])
	assert.Equal(t, "123", result.QueryParams["distinct_id"])
	assert.Equal(t, "123", result.QueryParams["buid"])
	assert.Equal(t, fmt.Sprintf("%s/123", consts.EnvironmentCDNPath()), result.QueryParams["cache_miss_url"])
	assert.Equal(t, 6, len(result.QueryParams))
}

func TestRequestConfigurationBuilderAPIRequestDataWillHaveServerData(t *testing.T) {
	sdkSettings := &mockedSdkSettings{}
	sdkSettings.On("DevModeSecret").Return("1")

	buid := &mockedBUID{}
	buid.On("GetValue").Return("123")
	buid.On("GetQueryStringParts").Return(map[string]string{
		"buid": "123",
	})

	deviceProps := &mockedDeviceProperties{}
	deviceProps.On("DistinctId").Return("123")
	deviceProps.On("GetAllProperties").Return(map[string]string{
		"app_key":     "123",
		"api_version": "4.0.0",
		"distinct_id": "123",
	})

	requestConfigurationBuilder := network.NewRequestConfigurationBuilder(sdkSettings, buid, deviceProps, "")
	result := requestConfigurationBuilder.BuildForAPI()

	assert.Equal(t, consts.EnvironmentAPIPath(), result.Url)
	assert.Equal(t, "123", result.QueryParams["app_key"])
	assert.Equal(t, "4.0.0", result.QueryParams["api_version"])
	assert.Equal(t, "123", result.QueryParams["distinct_id"])
	assert.Equal(t, "123", result.QueryParams["buid"])
	assert.Equal(t, fmt.Sprintf("%s/123", consts.EnvironmentCDNPath()), result.QueryParams["cache_miss_url"])
	assert.Equal(t, 6, len(result.QueryParams))
}
