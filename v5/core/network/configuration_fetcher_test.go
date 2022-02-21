package network_test

import (
	"net/http"
	"testing"

	"github.com/pkg/errors"
	"github.com/rollout/rox-go/v5/core/client"
	"github.com/rollout/rox-go/v5/core/configuration"
	"github.com/rollout/rox-go/v5/core/mocks"
	"github.com/rollout/rox-go/v5/core/model"
	"github.com/rollout/rox-go/v5/core/network"
	"github.com/stretchr/testify/assert"
)

func TestConfigurationFetcherWillReturnCDNDataWhenSuccessful(t *testing.T) {
	confFetchInvoker := configuration.NewFetchedInvoker()
	numberOfTimesCalled := 0
	confFetchInvoker.RegisterFetchedHandler(func(e *model.ConfigurationFetchedArgs) {
		numberOfTimesCalled++
	})

	requestData := model.RequestData{URL: "harta.com"}
	request := &mocks.Request{}
	response := &model.Response{StatusCode: http.StatusOK, Content: []byte("{\"data\": \"harti\"}")}
	request.On("SendGet", requestData).Return(response, nil)

	environment := client.NewSaasEnvironment()
	requestBuilder := &mocks.RequestConfigurationBuilder{}
	requestBuilder.On("BuildForCDN").Return(requestData)

	confFetcher := network.NewConfigurationFetcher(environment, requestBuilder, request, confFetchInvoker)
	result := confFetcher.Fetch()

	assert.Equal(t, "harti", result.ParsedData.Data)
	assert.Equal(t, configuration.SourceCDN, result.Source)
	assert.Equal(t, 0, numberOfTimesCalled)
}

func TestConfigurationFetcherWillReturnNullWhenCDNFailsWithException(t *testing.T) {
	confFetchInvoker := configuration.NewFetchedInvoker()
	numberOfTimesCalled := 0
	confFetchInvoker.RegisterFetchedHandler(func(e *model.ConfigurationFetchedArgs) {
		numberOfTimesCalled++
	})

	requestDataCDN := model.RequestData{URL: "harta1.com"}
	requestDataAPI := model.RequestData{URL: "harta2.com"}
	request := &mocks.Request{}
	response := &model.Response{StatusCode: http.StatusOK, Content: []byte("harto")}
	request.On("SendGet", requestDataCDN).Return(nil, errors.New("not found"))
	request.On("SendGet", requestDataAPI).Return(response, nil)

	environment := client.NewSaasEnvironment()
	requestBuilder := &mocks.RequestConfigurationBuilder{}
	requestBuilder.On("BuildForCDN").Return(requestDataCDN)
	requestBuilder.On("BuildForAPI").Return(requestDataAPI)

	confFetcher := network.NewConfigurationFetcher(environment, requestBuilder, request, confFetchInvoker)
	result := confFetcher.Fetch()

	assert.Nil(t, result)
	assert.Equal(t, 1, numberOfTimesCalled)
}

func TestConfigurationFetcherWillReturnNullWhenCDNFails404APIWithException(t *testing.T) {
	confFetchInvoker := configuration.NewFetchedInvoker()
	numberOfTimesCalled := 0
	confFetchInvoker.RegisterFetchedHandler(func(e *model.ConfigurationFetchedArgs) {
		numberOfTimesCalled++
	})

	requestDataCDN := model.RequestData{URL: "harta1.com"}
	requestDataAPI := model.RequestData{URL: "harta2.com"}
	request := &mocks.Request{}
	response := &model.Response{StatusCode: http.StatusNotFound}
	request.On("SendGet", requestDataCDN).Return(response, nil)
	request.On("SendGet", requestDataAPI).Return(nil, errors.New("exception"))

	environment := client.NewSaasEnvironment()
	requestBuilder := &mocks.RequestConfigurationBuilder{}
	requestBuilder.On("BuildForCDN").Return(requestDataCDN)
	requestBuilder.On("BuildForAPI").Return(requestDataAPI)

	confFetcher := network.NewConfigurationFetcher(environment, requestBuilder, request, confFetchInvoker)
	result := confFetcher.Fetch()

	assert.Nil(t, result)
	assert.Equal(t, 1, numberOfTimesCalled)
}

func TestConfigurationFetcherWillReturnAPIDataWhenCDNFails404APIOK(t *testing.T) {
	confFetchInvoker := configuration.NewFetchedInvoker()
	numberOfTimesCalled := 0
	confFetchInvoker.RegisterFetchedHandler(func(e *model.ConfigurationFetchedArgs) {
		numberOfTimesCalled++
	})

	requestDataCDN := model.RequestData{URL: "harta1.com"}
	requestDataAPI := model.RequestData{URL: "harta2.com"}
	request := &mocks.Request{}
	response := &model.Response{StatusCode: http.StatusOK, Content: []byte("{\"data\": \"harto\"}")}
	responseCDN := &model.Response{StatusCode: http.StatusNotFound}
	request.On("SendGet", requestDataCDN).Return(responseCDN, nil)
	request.On("SendPost", requestDataAPI.URL, requestDataAPI.QueryParams).Return(response, nil)

	environment := client.NewSaasEnvironment()
	requestBuilder := &mocks.RequestConfigurationBuilder{}
	requestBuilder.On("BuildForCDN").Return(requestDataCDN)
	requestBuilder.On("BuildForAPI").Return(requestDataAPI)

	confFetcher := network.NewConfigurationFetcher(environment, requestBuilder, request, confFetchInvoker)
	result := confFetcher.Fetch()

	assert.Equal(t, "harto", result.ParsedData.Data)
	assert.Equal(t, configuration.SourceAPI, result.Source)
	assert.Equal(t, 0, numberOfTimesCalled)
}

func TestConfigurationFetcherWillReturnNullDataWhenBothNotFound(t *testing.T) {
	confFetchInvoker := configuration.NewFetchedInvoker()
	numberOfTimesCalled := 0
	confFetchInvoker.RegisterFetchedHandler(func(e *model.ConfigurationFetchedArgs) {
		numberOfTimesCalled++
	})

	requestDataCDN := model.RequestData{URL: "harta1.com"}
	requestDataAPI := model.RequestData{URL: "harta2.com"}
	request := &mocks.Request{}
	response := &model.Response{StatusCode: http.StatusNotFound}
	responseCDN := &model.Response{StatusCode: http.StatusNotFound}
	request.On("SendGet", requestDataCDN).Return(responseCDN, nil)
	request.On("SendGet", requestDataAPI).Return(response, nil)

	environment := client.NewSaasEnvironment()
	requestBuilder := &mocks.RequestConfigurationBuilder{}
	requestBuilder.On("BuildForCDN").Return(requestDataCDN)
	requestBuilder.On("BuildForAPI").Return(requestDataAPI)

	confFetcher := network.NewConfigurationFetcher(environment, requestBuilder, request, confFetchInvoker)
	result := confFetcher.Fetch()

	assert.Nil(t, result)
	assert.Equal(t, 1, numberOfTimesCalled)
}

func TestConfigurationFetcherWillReturnAPIDataWhenSelfManaged(t *testing.T) {
	confFetchInvoker := configuration.NewFetchedInvoker()
	numberOfTimesCalled := 0
	confFetchInvoker.RegisterFetchedHandler(func(e *model.ConfigurationFetchedArgs) {
		numberOfTimesCalled++
	})

	requestDataAPI := model.RequestData{URL: "http://harta2.com"}
	request := &mocks.Request{}
	response := &model.Response{StatusCode: http.StatusNotFound}
	request.On("SendGet", requestDataAPI).Return(response, nil)

	environment := client.NewSelfManagedEnvironment(client.NewSelfManagedOptions(client.SelfManagedOptionsBuilder{
		ServerURL:    "http://harta2.com",
		AnalyticsURL: "http://harta2.com",
	}))
	requestBuilder := &mocks.RequestConfigurationBuilder{}
	requestBuilder.On("BuildForAPI").Return(requestDataAPI)

	confFetcher := network.NewConfigurationFetcher(environment, requestBuilder, request, confFetchInvoker)
	result := confFetcher.Fetch()

	assert.Nil(t, result)
	assert.Equal(t, 1, numberOfTimesCalled)
}
