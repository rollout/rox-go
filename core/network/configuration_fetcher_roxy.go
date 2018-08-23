package network

import (
	"github.com/rollout/rox-go/core/configuration"
)

type configurationFetcherRoxy struct {
	requestConfigurationBuilder RequestConfigurationBuilder
	request                     Request
	configurationFetcherLogger  configurationFetcherLogger
}

func NewConfigurationFetcherRoxy(requestConfigurationBuilder RequestConfigurationBuilder, request Request, configurationFetchedInvoker *configuration.ConfigurationFetchedInvoker) ConfigurationFetcher {
	return &configurationFetcherRoxy{
		requestConfigurationBuilder: requestConfigurationBuilder,
		request:                     request,
		configurationFetcherLogger:  configurationFetcherLogger{configurationFetchedInvoker},
	}
}

func (f *configurationFetcherRoxy) Fetch() *configuration.ConfigurationFetchResult {
	source := configuration.SourceRoxy

	defer func() {
		if r := recover(); r != nil {
			f.configurationFetcherLogger.WriteFetchExceptionToLogAndInvokeFetchHandler(source, r)
		}
	}()

	fetchResult, err := f.fetchFromRoxy()
	if err != nil {
		f.configurationFetcherLogger.WriteFetchExceptionToLogAndInvokeFetchHandler(source, err)
		return nil
	}

	if fetchResult.IsSuccessStatusCode() {
		return configuration.NewConfigurationFetchResult(string(fetchResult.Content), source)
	}

	f.configurationFetcherLogger.WriteFetchErrorToLogAndInvokeFetchHandler(source, fetchResult)
	return nil
}

func (f *configurationFetcherRoxy) fetchFromRoxy() (response *Response, err error) {
	return f.request.SendGet(f.requestConfigurationBuilder.BuildForRoxy())
}
