package core_test

import (
	"github.com/rollout/rox-go/core"
	"github.com/rollout/rox-go/core/mocks"
	"testing"
	"time"
)

func TestCoreWillCheckCoreSetupWhenOptionsWithRoxy(t *testing.T) {
	sdkSettings := &mocks.SdkSettings{}
	sdkSettings.On("DevModeSecret").Return("")
	sdkSettings.On("APIKey").Return("api_key")

	deviceProperties := &mocks.DeviceProperties{}
	deviceProperties.On("GetAllProperties").Return(map[string]string{})
	deviceProperties.On("DistinctID").Return("")

	options := &mocks.RoxOptions{}
	options.On("RoxyURL").Return("http://site.com")
	options.On("FetchInterval").Return(time.Duration(0))
	options.On("ConfigurationFetchedHandler").Return(nil)
	options.On("ImpressionHandler").Return(nil)

	c := core.NewCore()
	<-c.Setup(sdkSettings, deviceProperties, options)
}

func TestCoreWillCheckCoreSetupWhenNoOptions(t *testing.T) {
	sdkSettings := &mocks.SdkSettings{}
	sdkSettings.On("DevModeSecret").Return("")
	sdkSettings.On("APIKey").Return("api_key")

	deviceProperties := &mocks.DeviceProperties{}
	deviceProperties.On("GetAllProperties").Return(map[string]string{})
	deviceProperties.On("DistinctID").Return("")

	c := core.NewCore()
	<-c.Setup(sdkSettings, deviceProperties, nil)
}
