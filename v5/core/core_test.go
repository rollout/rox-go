package core_test

import (
	"fmt"
	"github.com/rollout/rox-go/v5/core"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/rollout/rox-go/v5/core/mocks"
)

var validApiKey = "5008ef002000b62ceaaab37b"

func TestCoreWillCheckCoreSetupWhenOptionsWithRoxy(t *testing.T) {
	sdkSettings := &mocks.SdkSettings{}
	sdkSettings.On("DevModeSecret").Return("")
	sdkSettings.On("APIKey").Return(validApiKey)

	deviceProperties := &mocks.DeviceProperties{}
	deviceProperties.On("GetAllProperties").Return(map[string]string{})
	deviceProperties.On("DistinctID").Return("")

	options := &mocks.RoxOptions{}
	options.On("RoxyURL").Return("http://site.com")
	options.On("FetchInterval").Return(time.Duration(0))
	options.On("ConfigurationFetchedHandler").Return(nil)
	options.On("ImpressionHandler").Return(nil)
	options.On("SelfManagedOptions").Return(nil)
	options.On("DynamicPropertyRuleHandler").Return(nil)

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

	defer func() {
		assert.Equal(t, "Invalid rollout apikey", recover())
	}()
	<-c.Setup(sdkSettings, deviceProperties, nil)
	assert.Fail(t, "We should never reach this point because the API key is invalid")
}

func TestValidAPIKey_MongoId(t *testing.T) {
	c := core.NewCore()

	defer func() {
		if err := recover(); err != nil {
			// Due to the panic() generated but the Setup,
			// we should reach here and not the t.FailNow() underneath
		}
	}()
	sdkSettings := &mocks.SdkSettings{}
	sdkSettings.On("APIKey").Return("12345678901234567890abcd") // Valid Mongo ID

	deviceProperties := &mocks.DeviceProperties{}
	deviceProperties.On("GetAllProperties").Return(map[string]string{})

	defer func() {
		assert.Nil(t, recover(), "we should not have panicked as the API key was valid")
	}()
	<-c.Setup(sdkSettings, deviceProperties, nil)
	// Success
}

func TestValidAPIKey_Uuid(t *testing.T) {
	c := core.NewCore()

	defer func() {
		if err := recover(); err != nil {
			// Due to the panic() generated but the Setup,
			// we should reach here and not the t.FailNow() underneath
		}
	}()
	sdkSettings := &mocks.SdkSettings{}
	sdkSettings.On("APIKey").Return("a9faa59e-f005-11ed-abc0-00155d7746b3") // Valid Mongo ID

	deviceProperties := &mocks.DeviceProperties{}
	deviceProperties.On("GetAllProperties").Return(map[string]string{})

	defer func() {
		assert.Nil(t, recover(), "we should not have panicked as the API key was valid")
	}()
	<-c.Setup(sdkSettings, deviceProperties, nil)
	// Success
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
	defer func() {
		assert.Equal(t, "Invalid rollout apikey", recover())
	}()
	<-c.Setup(sdkSettings, deviceProperties, nil)
	assert.Fail(t, "We should never reach this point because the API key is invalid")
}
