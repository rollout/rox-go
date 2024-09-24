package network

import (
	"net/http"

	"github.com/rollout/rox-go/v6/core/configuration"
	"github.com/rollout/rox-go/v6/core/model"
)

type ConfigurationFetcher interface {
	Fetch() *configuration.FetchResult
}

type configurationFetcher struct {
	environment                 model.Environment
	requestConfigurationBuilder RequestConfigurationBuilder
	request                     model.Request
	fetcherLogger               configurationFetcherLogger
}

func NewConfigurationFetcher(environment model.Environment, requestConfigurationBuilder RequestConfigurationBuilder, request model.Request, fetchedInvoker *configuration.FetchedInvoker) ConfigurationFetcher {
	return &configurationFetcher{
		environment:                 environment,
		requestConfigurationBuilder: requestConfigurationBuilder,
		request:                     request,
		fetcherLogger:               configurationFetcherLogger{fetchedInvoker},
	}
}

func (f *configurationFetcher) Fetch() *configuration.FetchResult {
	shouldRetry := false
	source := configuration.SourceCDN

	defer func() {
		if r := recover(); r != nil {
			f.fetcherLogger.WriteFetchExceptionToLogAndInvokeFetchHandler(source, r)
		}
	}()

	var fetchResult *model.Response = nil
	var err error = nil
	isSelfManaged := f.environment.IsSelfManaged()

	if !isSelfManaged {
		fetchResult, err = f.fetchFromCDN()
		if err != nil {
			f.fetcherLogger.WriteFetchExceptionToLogAndInvokeFetchHandler(source, err)
			return nil
		}

		if fetchResult.IsSuccessStatusCode() {
			configurationFetchResult := configuration.NewFetchResult(string(fetchResult.Content), source)
			if configurationFetchResult == nil {
				return nil
			}

			if configurationFetchResult.ParsedData.Result == 404 {
				shouldRetry = true
			} else {
				return configurationFetchResult
			}
		}
	}

	if isSelfManaged || shouldRetry || fetchResult.StatusCode == http.StatusForbidden || fetchResult.StatusCode == http.StatusNotFound {
		if !isSelfManaged {
			f.fetcherLogger.WriteFetchErrorToLog(source, fetchResult, configuration.SourceAPI)
		}
		source = configuration.SourceAPI
		fetchResult, err = f.fetchFromAPI()
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
	requestData := f.requestConfigurationBuilder.BuildForAPI()
	return f.request.SendPost(requestData.URL, requestData.QueryParams)
}
