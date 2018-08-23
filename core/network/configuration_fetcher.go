package network

import (
	"github.com/rollout/rox-go/core/configuration"
	"net/http"
)

type ConfigurationFetcher interface {
	Fetch() *configuration.ConfigurationFetchResult
}

type configurationFetcher struct {
	requestConfigurationBuilder RequestConfigurationBuilder
	request                     Request
	configurationFetcherLogger  configurationFetcherLogger
}

func NewConfigurationFetcher(requestConfigurationBuilder RequestConfigurationBuilder, request Request, configurationFetchedInvoker *configuration.ConfigurationFetchedInvoker) ConfigurationFetcher {
	return &configurationFetcher{
		requestConfigurationBuilder: requestConfigurationBuilder,
		request:                     request,
		configurationFetcherLogger:  configurationFetcherLogger{configurationFetchedInvoker},
	}
}

func (f *configurationFetcher) Fetch() *configuration.ConfigurationFetchResult {
	source := configuration.SourceCDN

	defer func() {
		if r := recover(); r != nil {
			f.configurationFetcherLogger.WriteFetchExceptionToLogAndInvokeFetchHandler(source, r)
		}
	}()

	fetchResult, err := f.fetchFromCDN()
	if err != nil {
		f.configurationFetcherLogger.WriteFetchExceptionToLogAndInvokeFetchHandler(source, err)
		return nil
	}

	if fetchResult.IsSuccessStatusCode() {
		return configuration.NewConfigurationFetchResult(string(fetchResult.Content), source)
	}

	if fetchResult.StatusCode == http.StatusForbidden || fetchResult.StatusCode == http.StatusNotFound {
		f.configurationFetcherLogger.WriteFetchErrorToLog(source, fetchResult, configuration.SourceAPI)
		source = configuration.SourceAPI
		fetchResult, err := f.fetchFromAPI()
		if err != nil {
			f.configurationFetcherLogger.WriteFetchExceptionToLogAndInvokeFetchHandler(source, err)
			return nil
		}

		if fetchResult.IsSuccessStatusCode() {
			return configuration.NewConfigurationFetchResult(string(fetchResult.Content), source)
		}
	}

	f.configurationFetcherLogger.WriteFetchErrorToLogAndInvokeFetchHandler(source, fetchResult)
	return nil
}

func (f *configurationFetcher) fetchFromCDN() (response *Response, err error) {
	return f.request.SendGet(f.requestConfigurationBuilder.BuildForCDN())
}

func (f *configurationFetcher) fetchFromAPI() (response *Response, err error) {
	return f.request.SendGet(f.requestConfigurationBuilder.BuildForAPI())
}
