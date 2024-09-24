package network_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/rollout/rox-go/v6/core/configuration"
	"github.com/rollout/rox-go/v6/core/mocks"
	"github.com/rollout/rox-go/v6/core/model"
	"github.com/rollout/rox-go/v6/core/network"
	"github.com/stretchr/testify/assert"
)

func TestConfigurationFetcherRoxyWillReturnCDNDataWhenSuccessful(t *testing.T) {
	confFetchInvoker := configuration.NewFetchedInvoker()
	numberOfTimesCalled := 0
	confFetchInvoker.RegisterFetchedHandler(func(e *model.ConfigurationFetchedArgs) {
		numberOfTimesCalled++
	})

	requestData := model.RequestData{URL: "harta.com"}
	request := &mocks.Request{}
	response := &model.Response{StatusCode: http.StatusOK, Content: []byte("{\"data\": \"harti\"}")}
	request.On("SendGet", requestData).Return(response, nil)

	requestBuilder := &mocks.RequestConfigurationBuilder{}
	requestBuilder.On("BuildForRoxy").Return(requestData)

	confFetcher := network.NewConfigurationFetcherRoxy(requestBuilder, request, confFetchInvoker)
	result := confFetcher.Fetch()

	assert.Equal(t, "harti", result.ParsedData.Data)
	assert.Equal(t, configuration.SourceRoxy, result.Source)
	assert.Equal(t, 0, numberOfTimesCalled)
}

func TestConfigurationFetcherRoxyWillReturnNullWhenRoxyFailsWithException(t *testing.T) {
	confFetchInvoker := configuration.NewFetchedInvoker()
	numberOfTimesCalled := 0
	confFetchInvoker.RegisterFetchedHandler(func(e *model.ConfigurationFetchedArgs) {
		numberOfTimesCalled++
	})

	requestData := model.RequestData{URL: "harta.com"}
	request := &mocks.Request{}
	request.On("SendGet", requestData).Return(nil, errors.New("not found"))

	requestBuilder := &mocks.RequestConfigurationBuilder{}
	requestBuilder.On("BuildForRoxy").Return(requestData)

	confFetcher := network.NewConfigurationFetcherRoxy(requestBuilder, request, confFetchInvoker)
	result := confFetcher.Fetch()

	assert.Nil(t, result)
	assert.Equal(t, 1, numberOfTimesCalled)
}

func TestConfigurationFetcherRoxyWillReturnNullWhenRoxyFailsWithHttpStatus(t *testing.T) {
	confFetchInvoker := configuration.NewFetchedInvoker()
	numberOfTimesCalled := 0
	confFetchInvoker.RegisterFetchedHandler(func(e *model.ConfigurationFetchedArgs) {
		numberOfTimesCalled++
	})

	requestData := model.RequestData{URL: "harta.com"}
	request := &mocks.Request{}
	response := &model.Response{StatusCode: http.StatusNotFound, Content: []byte("harto")}
	request.On("SendGet", requestData).Return(response, nil)

	requestBuilder := &mocks.RequestConfigurationBuilder{}
	requestBuilder.On("BuildForRoxy").Return(requestData)

	confFetcher := network.NewConfigurationFetcherRoxy(requestBuilder, request, confFetchInvoker)
	result := confFetcher.Fetch()

	assert.Nil(t, result)
	assert.Equal(t, 1, numberOfTimesCalled)
}
