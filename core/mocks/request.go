package mocks

import (
	"github.com/rollout/rox-go/core/network"
	"github.com/stretchr/testify/mock"
)

type Request struct {
	mock.Mock
}

func (m Request) SendGet(requestData network.RequestData) (response *network.Response, err error) {
	args := m.Called(requestData)
	return args.Get(0).(*network.Response), args.Error(1)
}

func (m Request) SendPost(uri string, content interface{}) (response *network.Response, err error) {
	args := m.Called(uri, content)
	return args.Get(0).(*network.Response), args.Error(1)
}
