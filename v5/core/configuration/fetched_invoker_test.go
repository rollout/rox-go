package configuration

import (
	"github.com/rollout/rox-go/v5/core/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConfigurationFetchedInvokerWithNoSubscriberNoException(t *testing.T) {
	configurationFetchedInvoker := NewFetchedInvoker()
	configurationFetchedInvoker.InvokeError(model.FetcherErrorUnknown)

	configurationFetchedInvoker2 := NewFetchedInvoker()
	configurationFetchedInvoker2.Invoke(model.FetcherStatusAppliedFromEmbedded, time.Now(), true)
}

func TestConfigurationFetchedInvokerArgsConstructor(t *testing.T) {
	status := model.FetcherStatusAppliedFromEmbedded
	now := time.Now()
	hasChanges := true

	confFetchedArgs := model.NewConfigurationFetchedArgs(status, now, hasChanges)

	assert.Equal(t, status, confFetchedArgs.FetcherStatus)
	assert.Equal(t, now, confFetchedArgs.CreationDate)
	assert.Equal(t, hasChanges, confFetchedArgs.HasChanges)
	assert.Equal(t, model.FetcherErrorNoError, confFetchedArgs.ErrorDetails)

	confFetchedArgs2 := model.NewErrorConfigurationFetchedArgs(model.FetcherErrorSignatureVerification)

	assert.Equal(t, model.FetcherStatusErrorFetchedFailed, confFetchedArgs2.FetcherStatus)
	assert.True(t, confFetchedArgs2.CreationDate.IsZero())
	assert.Equal(t, false, confFetchedArgs2.HasChanges)
	assert.Equal(t, model.FetcherErrorSignatureVerification, confFetchedArgs2.ErrorDetails)
}

func TestConfigurationFetchedInvokerInvokeWithError(t *testing.T) {
	isConfigurationHandlerInvokerRaised := false
	configurationFetchedInvoker := NewFetchedInvoker()

	configurationFetchedInvoker.RegisterFetchedHandler(func(e *model.ConfigurationFetchedArgs) {
		assert.Equal(t, model.FetcherStatusErrorFetchedFailed, e.FetcherStatus)
		assert.True(t, e.CreationDate.IsZero())
		assert.Equal(t, false, e.HasChanges)
		assert.Equal(t, model.FetcherErrorUnknown, e.ErrorDetails)

		isConfigurationHandlerInvokerRaised = true
	})

	configurationFetchedInvoker.InvokeError(model.FetcherErrorUnknown)

	assert.True(t, isConfigurationHandlerInvokerRaised)
}

func TestConfigurationFetchedInvokerInvokeOK(t *testing.T) {
	isConfigurationHandlerInvokerRaised := false
	configurationFetchedInvoker := NewFetchedInvoker()

	now := time.Now()
	status := model.FetcherStatusAppliedFromNetwork
	hasChanges := true

	configurationFetchedInvoker.RegisterFetchedHandler(func(e *model.ConfigurationFetchedArgs) {
		assert.Equal(t, status, e.FetcherStatus)
		assert.Equal(t, now, e.CreationDate)
		assert.Equal(t, hasChanges, e.HasChanges)
		assert.Equal(t, model.FetcherErrorNoError, e.ErrorDetails)

		isConfigurationHandlerInvokerRaised = true
	})

	configurationFetchedInvoker.Invoke(status, now, hasChanges)

	assert.True(t, isConfigurationHandlerInvokerRaised)
}

func TestConfigurationFetchedInvokerPanicHandler(t *testing.T) {
	isConfigurationHandlerInvokerRaised := false
	configurationFetchedInvoker := NewFetchedInvoker()

	now := time.Now()
	status := model.FetcherStatusAppliedFromNetwork
	hasChanges := true

	configurationFetchedInvoker.RegisterFetchedHandler(func(e *model.ConfigurationFetchedArgs) {
		panic(nil)
	})

	configurationFetchedInvoker.Invoke(status, now, hasChanges)

	assert.False(t, isConfigurationHandlerInvokerRaised)
}
