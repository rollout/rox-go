package configuration

import (
	"github.com/rollout/rox-go/core/model"
	"sync"
	"time"
)

type FetchedInvoker struct {
	fetchedHandlers []model.ConfigurationFetchedHandler
	handlersMutex   sync.RWMutex
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
	cfi.handlersMutex.Lock()
	cfi.fetchedHandlers = append(cfi.fetchedHandlers, handler)
	cfi.handlersMutex.Unlock()
}

func (cfi *FetchedInvoker) raiseFetchedEvent(args model.ConfigurationFetchedArgs) {
	cfi.handlersMutex.RLock()
	handlers := make([]model.ConfigurationFetchedHandler, len(cfi.fetchedHandlers))
	copy(handlers, cfi.fetchedHandlers)
	cfi.handlersMutex.RUnlock()

	for _, handler := range handlers {
		handler(&args)
	}
}
