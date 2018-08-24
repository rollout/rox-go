package configuration

import (
	"github.com/rollout/rox-go/core/model"
	"time"
)

type FetchedInvoker struct {
	handlers []model.ConfigurationFetchedHandler
}

func NewFetchedInvoker() *FetchedInvoker {
	return &FetchedInvoker{}
}

func (cfi *FetchedInvoker) Invoke(fetcherStatus model.FetcherStatus, creationDate time.Time, hasChanges bool) {
	cfi.raiseFetchedEvent(model.NewConfigurationFetchedArgs(fetcherStatus, creationDate, hasChanges))
}

func (cfi *FetchedInvoker) InvokeError(errorDetails model.FetcherError) {
	cfi.raiseFetchedEvent(model.NewErrorConfigurationFetchedArgs(errorDetails))
}

func (cfi *FetchedInvoker) RegisterFetchedHandler(handler model.ConfigurationFetchedHandler) {
	cfi.handlers = append(cfi.handlers, handler)
}

func (cfi *FetchedInvoker) raiseFetchedEvent(args model.ConfigurationFetchedArgs) {
	for _, handler := range cfi.handlers {
		handler(&args)
	}
}
