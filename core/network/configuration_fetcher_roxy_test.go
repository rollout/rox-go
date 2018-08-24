package network_test

import (
	"errors"
	"github.com/rollout/rox-go/core/configuration"
	"github.com/rollout/rox-go/core/mocks"
	"github.com/rollout/rox-go/core/model"
	"github.com/rollout/rox-go/core/network"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestConfigurationFetcherRoxyWillReturnCDNDataWhenSuccessful(t *testing.T) {
	confFetchInvoker := configuration.NewConfigurationFetchedInvoker()
	numberOfTimesCalled := 0
	confFetchInvoker.RegisterConfigurationFetchedHandler(func(e *model.ConfigurationFetchedArgs) {
		numberOfTimesCalled++
	})

	requestData := network.RequestData{Url: "harta.com"}
	request := &mocks.Request{}
	response := &network.Response{StatusCode: http.StatusOK, Content: []byte("harti")}
	request.On("SendGet", requestData).Return(response, nil)

	requestBuilder := &mocks.RequestConfigurationBuilder{}
	requestBuilder.On("BuildForRoxy").Return(requestData)

	confFetcher := network.NewConfigurationFetcherRoxy(requestBuilder, request, confFetchInvoker)
	result := confFetcher.Fetch()

	assert.Equal(t, "harti", result.Data)
	assert.Equal(t, configuration.SourceRoxy, result.Source)
	assert.Equal(t, 0, numberOfTimesCalled)
}

func TestConfigurationFetcherRoxyWillReturnNullWhenRoxyFailsWithException(t *testing.T) {
	confFetchInvoker := configuration.NewConfigurationFetchedInvoker()
	numberOfTimesCalled := 0
	confFetchInvoker.RegisterConfigurationFetchedHandler(func(e *model.ConfigurationFetchedArgs) {
		numberOfTimesCalled++
	})

	requestData := network.RequestData{Url: "harta.com"}
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
	confFetchInvoker := configuration.NewConfigurationFetchedInvoker()
	numberOfTimesCalled := 0
	confFetchInvoker.RegisterConfigurationFetchedHandler(func(e *model.ConfigurationFetchedArgs) {
		numberOfTimesCalled++
	})

	requestData := network.RequestData{Url: "harta.com"}
	request := &mocks.Request{}
	response := &network.Response{StatusCode: http.StatusNotFound, Content: []byte("harto")}
	request.On("SendGet", requestData).Return(response, nil)

	requestBuilder := &mocks.RequestConfigurationBuilder{}
	requestBuilder.On("BuildForRoxy").Return(requestData)

	confFetcher := network.NewConfigurationFetcherRoxy(requestBuilder, request, confFetchInvoker)
	result := confFetcher.Fetch()

	assert.Nil(t, result)
	assert.Equal(t, 1, numberOfTimesCalled)
}
