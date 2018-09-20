package network

import (
	"fmt"
	"github.com/rollout/rox-go/core/configuration"
	"github.com/rollout/rox-go/core/logging"
	"github.com/rollout/rox-go/core/model"
)

type configurationFetcherLogger struct {
	fetchedInvoker *configuration.FetchedInvoker
}

func (fl *configurationFetcherLogger) WriteFetchErrorToLogAndInvokeFetchHandler(source configuration.Source, response *model.Response) {
	logging.GetLogger().Debug(fmt.Sprintf("Failed to fetch from %s. http error code: %d\n", source, response.StatusCode), nil)
	fl.fetchedInvoker.InvokeError(model.FetcherErrorNetwork)
}

func (fl *configurationFetcherLogger) WriteFetchErrorToLog(source configuration.Source, response *model.Response, nextSource configuration.Source) {
	retryMsg := fmt.Sprintf("Trying from %s. ", nextSource)
	logging.GetLogger().Error(fmt.Sprintf("Failed to fetch from %s. %shttp error code: %d\n", source, retryMsg, response.StatusCode), nil)
}

func (fl *configurationFetcherLogger) WriteFetchExceptionToLogAndInvokeFetchHandler(source configuration.Source, ex interface{}) {
	logging.GetLogger().Error(fmt.Sprintf("Failed to fetch configuration. Source: %s. Ex: %s\n", source, ex), nil)
	fl.fetchedInvoker.InvokeError(model.FetcherErrorNetwork)
}
