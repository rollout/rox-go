package mocks

import (
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/rollout/rox-go/v5/core/model"
)

type Analytics struct {
	mock.Mock
}

func (m *Analytics) CaptureImpressions(impressions []model.Impression) {
	m.Called(impressions)
}

func (m *Analytics) InitiateIntervalReporting(interval time.Duration) {
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
