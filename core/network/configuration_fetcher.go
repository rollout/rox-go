package network

import (
	"github.com/rollout/rox-go/core/configuration"
	"github.com/rollout/rox-go/core/model"
	"net/http"
)

type ConfigurationFetcher interface {
	Fetch() *configuration.FetchResult
}

type configurationFetcher struct {
	requestConfigurationBuilder RequestConfigurationBuilder
	request                     model.Request
	fetcherLogger               configurationFetcherLogger
}

func NewConfigurationFetcher(requestConfigurationBuilder RequestConfigurationBuilder, request model.Request, fetchedInvoker *configuration.FetchedInvoker) ConfigurationFetcher {
	return &configurationFetcher{
		requestConfigurationBuilder: requestConfigurationBuilder,
		request:                     request,
		fetcherLogger:               configurationFetcherLogger{fetchedInvoker},
	}
}

func (f *configurationFetcher) Fetch() *configuration.FetchResult {
	source := configuration.SourceCDN

	defer func() {
		if r := recover(); r != nil {
			f.fetcherLogger.WriteFetchExceptionToLogAndInvokeFetchHandler(source, r)
		}
	}()

	fetchResult, err := f.fetchFromCDN()
	if err != nil {
		f.fetcherLogger.WriteFetchExceptionToLogAndInvokeFetchHandler(source, err)
		return nil
	}

	if fetchResult.IsSuccessStatusCode() {
		return configuration.NewFetchResult(string(fetchResult.Content), source)
	}

	if fetchResult.StatusCode == http.StatusForbidden || fetchResult.StatusCode == http.StatusNotFound {
		f.fetcherLogger.WriteFetchErrorToLog(source, fetchResult, configuration.SourceAPI)
		source = configuration.SourceAPI
		fetchResult, err := f.fetchFromAPI()
		if err != nil {
			f.fetcherLogger.WriteFetchExceptionToLogAndInvokeFetchHandler(source, err)
			return nil
		}

		if fetchResult.IsSuccessStatusCode() {
			return configuration.NewFetchResult(string(fetchResult.Content), source)
		}
	}

	f.fetcherLogger.WriteFetchErrorToLogAndInvokeFetchHandler(source, fetchResult)
	return nil
}

func (f *configurationFetcher) fetchFromCDN() (response *model.Response, err error) {
	return f.request.SendGet(f.requestConfigurationBuilder.BuildForCDN())
}

func (f *configurationFetcher) fetchFromAPI() (response *model.Response, err error) {
	return f.request.SendGet(f.requestConfigurationBuilder.BuildForAPI())
}
