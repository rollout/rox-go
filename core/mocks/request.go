package mocks

import (
	"github.com/rollout/rox-go/core/model"
	"github.com/stretchr/testify/mock"
)

type Request struct {
	mock.Mock
}

func (m *Request) SendGet(requestData model.RequestData) (response *model.Response, err error) {
	args := m.Called(requestData)
	return args.Get(0).(*model.Response), args.Error(1)
}

func (m *Request) SendPost(uri string, content interface{}) (response *model.Response, err error) {
	args := m.Called(uri, content)
	return args.Get(0).(*model.Response), args.Error(1)
}
