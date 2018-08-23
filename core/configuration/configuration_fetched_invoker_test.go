package configuration

import (
	"github.com/rollout/rox-go/core/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConfigurationFetchedInvokerWithNoSubscriberNoException(t *testing.T) {
	configurationFetchedInvoker := NewConfigurationFetchedInvoker()
	configurationFetchedInvoker.InvokeError(model.FetcherErrorUnknown)

	configurationFetchedInvoker2 := NewConfigurationFetchedInvoker()
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
	configurationFetchedInvoker := NewConfigurationFetchedInvoker()

	configurationFetchedInvoker.RegisterConfigurationFetchedHandler(func(e *model.ConfigurationFetchedArgs) {
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
	configurationFetchedInvoker := NewConfigurationFetchedInvoker()

	now := time.Now()
	status := model.FetcherStatusAppliedFromNetwork
	hasChanges := true

	configurationFetchedInvoker.RegisterConfigurationFetchedHandler(func(e *model.ConfigurationFetchedArgs) {
		assert.Equal(t, status, e.FetcherStatus)
		assert.Equal(t, now, e.CreationDate)
		assert.Equal(t, hasChanges, e.HasChanges)
		assert.Equal(t, model.FetcherErrorNoError, e.ErrorDetails)

		isConfigurationHandlerInvokerRaised = true
	})

	configurationFetchedInvoker.Invoke(status, now, hasChanges)

	assert.True(t, isConfigurationHandlerInvokerRaised)
}
