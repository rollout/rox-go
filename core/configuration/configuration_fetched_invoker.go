package configuration

import (
	"github.com/rollout/rox-go/core/model"
	"time"
)

type ConfigurationFetchedInvoker struct {
	handlers []model.ConfigurationFetchedHandler
}

func NewConfigurationFetchedInvoker() *ConfigurationFetchedInvoker {
	return &ConfigurationFetchedInvoker{}
}

func (cfi *ConfigurationFetchedInvoker) Invoke(fetcherStatus model.FetcherStatus, creationDate time.Time, hasChanges bool) {
	cfi.raiseConfigurationFetchedEvent(model.NewConfigurationFetchedArgs(fetcherStatus, creationDate, hasChanges))
}

func (cfi *ConfigurationFetchedInvoker) InvokeError(errorDetails model.FetcherError) {
	cfi.raiseConfigurationFetchedEvent(model.NewErrorConfigurationFetchedArgs(errorDetails))
}

func (cfi *ConfigurationFetchedInvoker) RegisterConfigurationFetchedHandler(handler model.ConfigurationFetchedHandler) {
	cfi.handlers = append(cfi.handlers, handler)
}

func (cfi *ConfigurationFetchedInvoker) raiseConfigurationFetchedEvent(args model.ConfigurationFetchedArgs) {
	for _, handler := range cfi.handlers {
		handler(&args)
	}
}
