package network

import (
	"github.com/rollout/rox-go/core/configuration"
)

type configurationFetcherRoxy struct {
	requestConfigurationBuilder RequestConfigurationBuilder
	request                     Request
	fetcherLogger               configurationFetcherLogger
}

func NewConfigurationFetcherRoxy(requestConfigurationBuilder RequestConfigurationBuilder, request Request, fetchedInvoker *configuration.FetchedInvoker) ConfigurationFetcher {
	return &configurationFetcherRoxy{
		requestConfigurationBuilder: requestConfigurationBuilder,
		request:                     request,
		fetcherLogger:               configurationFetcherLogger{fetchedInvoker},
	}
}

func (f *configurationFetcherRoxy) Fetch() *configuration.FetchResult {
	source := configuration.SourceRoxy

	defer func() {
		if r := recover(); r != nil {
			f.fetcherLogger.WriteFetchExceptionToLogAndInvokeFetchHandler(source, r)
		}
	}()

	fetchResult, err := f.fetchFromRoxy()
	if err != nil {
		f.fetcherLogger.WriteFetchExceptionToLogAndInvokeFetchHandler(source, err)
		return nil
	}

	if fetchResult.IsSuccessStatusCode() {
		return configuration.NewFetchResult(string(fetchResult.Content), source)
	}

	f.fetcherLogger.WriteFetchErrorToLogAndInvokeFetchHandler(source, fetchResult)
	return nil
}

func (f *configurationFetcherRoxy) fetchFromRoxy() (response *Response, err error) {
	return f.request.SendGet(f.requestConfigurationBuilder.BuildForRoxy())
}
