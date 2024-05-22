package mocks

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type Analytics struct {
	mock.Mock
}

func (m *Analytics) Enqueue(timeStamp float64, name string, value interface{}) {
	m.Called(timeStamp, name, value)
}

func (m *Analytics) InitiateReporting(interval time.Duration) {
	m.Called()
}

func (m *Analytics) IsAnalyticsReportingDisabled() bool {
	ret := m.Called()

	r0 := ret.Get(0)
	if r0 == nil {
		return false
	}
	return r0.(bool)
}
