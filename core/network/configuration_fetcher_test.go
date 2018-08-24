package network_test

import (
	"github.com/pkg/errors"
	"github.com/rollout/rox-go/core/configuration"
	"github.com/rollout/rox-go/core/mocks"
	"github.com/rollout/rox-go/core/model"
	"github.com/rollout/rox-go/core/network"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestConfigurationFetcherWillReturnCDNDataWhenSuccessful(t *testing.T) {
	confFetchInvoker := configuration.NewFetchedInvoker()
	numberOfTimesCalled := 0
	confFetchInvoker.RegisterFetchedHandler(func(e *model.ConfigurationFetchedArgs) {
		numberOfTimesCalled++
	})

	requestData := network.RequestData{URL: "harta.com"}
	request := &mocks.Request{}
	response := &network.Response{StatusCode: http.StatusOK, Content: []byte("harti")}
	request.On("SendGet", requestData).Return(response, nil)

	requestBuilder := &mocks.RequestConfigurationBuilder{}
	requestBuilder.On("BuildForCDN").Return(requestData)

	confFetcher := network.NewConfigurationFetcher(requestBuilder, request, confFetchInvoker)
	result := confFetcher.Fetch()

	assert.Equal(t, "harti", result.Data)
	assert.Equal(t, configuration.SourceCDN, result.Source)
	assert.Equal(t, 0, numberOfTimesCalled)
}

func TestConfigurationFetcherWillReturnNullWhenCDNFailsWithException(t *testing.T) {
	confFetchInvoker := configuration.NewFetchedInvoker()
	numberOfTimesCalled := 0
	confFetchInvoker.RegisterFetchedHandler(func(e *model.ConfigurationFetchedArgs) {
		numberOfTimesCalled++
	})

	requestDataCDN := network.RequestData{URL: "harta1.com"}
	requestDataAPI := network.RequestData{URL: "harta2.com"}
	request := &mocks.Request{}
	response := &network.Response{StatusCode: http.StatusOK, Content: []byte("harto")}
	request.On("SendGet", requestDataCDN).Return(nil, errors.New("not found"))
	request.On("SendGet", requestDataAPI).Return(response, nil)

	requestBuilder := &mocks.RequestConfigurationBuilder{}
	requestBuilder.On("BuildForCDN").Return(requestDataCDN)
	requestBuilder.On("BuildForAPI").Return(requestDataAPI)

	confFetcher := network.NewConfigurationFetcher(requestBuilder, request, confFetchInvoker)
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

	requestDataCDN := network.RequestData{URL: "harta1.com"}
	requestDataAPI := network.RequestData{URL: "harta2.com"}
	request := &mocks.Request{}
	response := &network.Response{StatusCode: http.StatusNotFound}
	request.On("SendGet", requestDataCDN).Return(response, nil)
	request.On("SendGet", requestDataAPI).Return(nil, errors.New("exception"))

	requestBuilder := &mocks.RequestConfigurationBuilder{}
	requestBuilder.On("BuildForCDN").Return(requestDataCDN)
	requestBuilder.On("BuildForAPI").Return(requestDataAPI)

	confFetcher := network.NewConfigurationFetcher(requestBuilder, request, confFetchInvoker)
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

	requestDataCDN := network.RequestData{URL: "harta1.com"}
	requestDataAPI := network.RequestData{URL: "harta2.com"}
	request := &mocks.Request{}
	response := &network.Response{StatusCode: http.StatusOK, Content: []byte("harto")}
	responseCDN := &network.Response{StatusCode: http.StatusNotFound}
	request.On("SendGet", requestDataCDN).Return(responseCDN, nil)
	request.On("SendGet", requestDataAPI).Return(response, nil)

	requestBuilder := &mocks.RequestConfigurationBuilder{}
	requestBuilder.On("BuildForCDN").Return(requestDataCDN)
	requestBuilder.On("BuildForAPI").Return(requestDataAPI)

	confFetcher := network.NewConfigurationFetcher(requestBuilder, request, confFetchInvoker)
	result := confFetcher.Fetch()

	assert.Equal(t, "harto", result.Data)
	assert.Equal(t, configuration.SourceAPI, result.Source)
	assert.Equal(t, 0, numberOfTimesCalled)
}

func TestConfigurationFetcherWillReturnNullDataWhenBothNotFound(t *testing.T) {
	confFetchInvoker := configuration.NewFetchedInvoker()
	numberOfTimesCalled := 0
	confFetchInvoker.RegisterFetchedHandler(func(e *model.ConfigurationFetchedArgs) {
		numberOfTimesCalled++
	})

	requestDataCDN := network.RequestData{URL: "harta1.com"}
	requestDataAPI := network.RequestData{URL: "harta2.com"}
	request := &mocks.Request{}
	response := &network.Response{StatusCode: http.StatusNotFound}
	responseCDN := &network.Response{StatusCode: http.StatusNotFound}
	request.On("SendGet", requestDataCDN).Return(responseCDN, nil)
	request.On("SendGet", requestDataAPI).Return(response, nil)

	requestBuilder := &mocks.RequestConfigurationBuilder{}
	requestBuilder.On("BuildForCDN").Return(requestDataCDN)
	requestBuilder.On("BuildForAPI").Return(requestDataAPI)

	confFetcher := network.NewConfigurationFetcher(requestBuilder, request, confFetchInvoker)
	result := confFetcher.Fetch()

	assert.Nil(t, result)
	assert.Equal(t, 1, numberOfTimesCalled)
}
