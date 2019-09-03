package core_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/rollout/rox-go/core"
	"github.com/rollout/rox-go/core/mocks"
)

var validApiKey = "5008ef002000b62ceaaab37b"

func TestCoreWillCheckCoreSetupWhenOptionsWithRoxy(t *testing.T) {
	sdkSettings := &mocks.SdkSettings{}
	sdkSettings.On("DevModeSecret").Return("")
	sdkSettings.On("APIKey").Return(validApiKey)

	deviceProperties := &mocks.DeviceProperties{}
	deviceProperties.On("GetAllProperties").Return(map[string]string{})
	deviceProperties.On("DistinctID").Return("")
	deviceProperties.On("RolloutEnvironment").Return("Test")
	deviceProperties.On("LibVersion").Return("0.0.1-test")

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
	sdkSettings.On("APIKey").Return(validApiKey)

	deviceProperties := &mocks.DeviceProperties{}
	deviceProperties.On("GetAllProperties").Return(map[string]string{})
	deviceProperties.On("DistinctID").Return("")
	deviceProperties.On("RolloutEnvironment").Return("Test")
	deviceProperties.On("LibVersion").Return("0.0.1-test")

	c := core.NewCore()
	<-c.Setup(sdkSettings, deviceProperties, nil)
}

func TestInvalidAPIKey(t *testing.T) {
	c := core.NewCore()

	defer func() {
		if err := recover(); err != nil {
			// Due to the panic() generated but the Setup,
			// we should reach here and not the t.FailNow() underneath
		}
	}()
	sdkSettings := &mocks.SdkSettings{}
	sdkSettings.On("APIKey").Return("invalid api key")

	deviceProperties := &mocks.DeviceProperties{}

	<-c.Setup(sdkSettings, deviceProperties, nil)
	// We should never reach this point because the API key is invalid
	t.FailNow()
}

func TestEmptyAPIKey(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	sdkSettings := &mocks.SdkSettings{}
	sdkSettings.On("APIKey").Return("")

	deviceProperties := &mocks.DeviceProperties{}

	c := core.NewCore()
	<-c.Setup(sdkSettings, deviceProperties, nil)
	// We should never reach this point because the API key is invalid
	t.FailNow()
}
