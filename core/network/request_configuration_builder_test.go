package network_test

import (
	"fmt"
	"testing"

	"github.com/rollout/rox-go/core/consts"
	"github.com/rollout/rox-go/core/mocks"
	"github.com/rollout/rox-go/core/network"
	"github.com/stretchr/testify/assert"
)

func TestRequestConfigurationBuilderCDNRequestDataWillHaveDistinctID(t *testing.T) {
	sdkSettings := &mocks.SdkSettings{}
	sdkSettings.On("DevModeSecret").Return("1")

	buid := &mocks.BUID{}
	buid.On("GetValue").Return("123")

	deviceProps := &mocks.DeviceProperties{}
	deviceProps.On("DistinctID").Return("123")
	deviceProps.On("GetAllProperties").Return(map[string]string{
		"app_key":        "123",
		"api_version":    "4.0.0",
		"cache_miss_url": "harta",
		"distinct_id":    "123",
	})
	appKey := "ABCD"
	deviceProps.On("RolloutKey").Return(appKey)

	requestConfigurationBuilder := network.NewRequestConfigurationBuilder(sdkSettings, buid, deviceProps, "")
	result := requestConfigurationBuilder.BuildForCDN()

	assert.Equal(t, fmt.Sprintf("%s/%s/123", consts.EnvironmentCDNPath(), appKey), result.URL)
	assert.Equal(t, "123", result.QueryParams["distinct_id"])
}

func TestRequestConfigurationBuilderRoxyRequestDataWillHaveServerData(t *testing.T) {
	sdkSettings := &mocks.SdkSettings{}
	sdkSettings.On("DevModeSecret").Return("1")

	buid := &mocks.BUID{}
	buid.On("GetValue").Return("123")
	buid.On("GetQueryStringParts").Return(map[string]string{
		"buid": "123",
	})

	deviceProps := &mocks.DeviceProperties{}
	deviceProps.On("DistinctID").Return("123")
	deviceProps.On("GetAllProperties").Return(map[string]string{
		"app_key":     "123",
		"api_version": "4.0.0",
		"distinct_id": "123",
	})
	appKey := "ABCD"
	deviceProps.On("RolloutKey").Return(appKey)

	requestConfigurationBuilder := network.NewRequestConfigurationBuilder(sdkSettings, buid, deviceProps, "http://bimba.bobi.o.ponpon")
	result := requestConfigurationBuilder.BuildForRoxy()

	assert.Equal(t, "http://bimba.bobi.o.ponpon/device/request_configuration", result.URL)
	assert.Equal(t, "123", result.QueryParams["app_key"])
	assert.Equal(t, "4.0.0", result.QueryParams["api_version"])
	assert.Equal(t, "123", result.QueryParams["distinct_id"])
	assert.Equal(t, "123", result.QueryParams["buid"])
	assert.Equal(t, fmt.Sprintf("%s/123", appKey), result.QueryParams["cache_miss_relative_url"])
	assert.Equal(t, 6, len(result.QueryParams))
}

func TestRequestConfigurationBuilderAPIRequestDataWillHaveServerData(t *testing.T) {
	sdkSettings := &mocks.SdkSettings{}
	sdkSettings.On("DevModeSecret").Return("1")

	buid := &mocks.BUID{}
	buid.On("GetValue").Return("123")
	buid.On("GetQueryStringParts").Return(map[string]string{
		"buid": "123",
	})

	deviceProps := &mocks.DeviceProperties{}
	deviceProps.On("DistinctID").Return("123")
	deviceProps.On("GetAllProperties").Return(map[string]string{
		"app_key":     "123",
		"api_version": "4.0.0",
		"distinct_id": "123",
	})
	appKey := "ABCD"
	deviceProps.On("RolloutKey").Return(appKey)

	requestConfigurationBuilder := network.NewRequestConfigurationBuilder(sdkSettings, buid, deviceProps, "")
	result := requestConfigurationBuilder.BuildForAPI()

	assert.Equal(t, fmt.Sprintf("%s/%s/123", consts.EnvironmentAPIPath(), appKey), result.URL)
	assert.Equal(t, "123", result.QueryParams["app_key"])
	assert.Equal(t, "4.0.0", result.QueryParams["api_version"])
	assert.Equal(t, "123", result.QueryParams["distinct_id"])
	assert.Equal(t, "123", result.QueryParams["buid"])
	assert.Equal(t, fmt.Sprintf("%s/123", appKey), result.QueryParams["cache_miss_relative_url"])
	assert.Equal(t, 6, len(result.QueryParams))
}
