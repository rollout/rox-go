package mocks

import (
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/rollout/rox-go/v6/core/model"
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

func (m *Analytics) StopIntervalReporting() {
	m.Called()
}
