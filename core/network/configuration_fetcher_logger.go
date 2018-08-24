package network

import (
	"fmt"
	"github.com/rollout/rox-go/core/configuration"
	"github.com/rollout/rox-go/core/model"
)

type configurationFetcherLogger struct {
	fetchedInvoker *configuration.FetchedInvoker
}

func (fl *configurationFetcherLogger) WriteFetchErrorToLogAndInvokeFetchHandler(source configuration.Source, response *Response) {
	// TODO logging
	fmt.Printf("Failed to fetch from %s. http error code: %d\n", source, response.StatusCode)
	fl.fetchedInvoker.InvokeError(model.FetcherErrorNetwork)
}

func (fl *configurationFetcherLogger) WriteFetchErrorToLog(source configuration.Source, response *Response, nextSource configuration.Source) {
	retryMsg := fmt.Sprintf("Trying from %s. ", nextSource)
	// TODO logging
	fmt.Printf("Failed to fetch from %s. %shttp error code: %d\n", source, retryMsg, response.StatusCode)
}

func (fl *configurationFetcherLogger) WriteFetchExceptionToLogAndInvokeFetchHandler(source configuration.Source, ex interface{}) {
	// TODO logging
	fmt.Printf("Failed to fetch configuration. Source: %s. Ex: %s\n", source, ex)
	fl.fetchedInvoker.InvokeError(model.FetcherErrorNetwork)
}
