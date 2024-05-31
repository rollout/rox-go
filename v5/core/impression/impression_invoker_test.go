package impression_test

import (
	"testing"

	"github.com/rollout/rox-go/v5/core/context"
	"github.com/rollout/rox-go/v5/core/impression"
	"github.com/rollout/rox-go/v5/core/mocks"
	"github.com/rollout/rox-go/v5/core/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestImpressionInvokerEmptyInvokeNotThrowingException(t *testing.T) {
	internalFlags := &mocks.InternalFlags{}
	deps := &impression.ImpressionsDeps{
		InternalFlags:            internalFlags,
		CustomPropertyRepository: nil,
		DeviceProperties:         nil,
		Analytics:                nil,
		IsRoxy:                   false,
	}
	impressionInvoker := impression.NewImpressionInvoker(deps)
	impressionInvoker.Invoke(nil, nil)
}

func TestImpressionInvokerInvokeAndParameters(t *testing.T) {
	internalFlags := &mocks.InternalFlags{}
	deps := &impression.ImpressionsDeps{
		InternalFlags:    internalFlags,
		DeviceProperties: nil,
		IsRoxy:           false,
	}
	impressionInvoker := impression.NewImpressionInvoker(deps)

	ctx := context.NewContext(map[string]interface{}{"obj1": 1})
	reportingValue := model.NewReportingValue("name", "value", true)

	isImpressionRaised := false
	impressionInvoker.RegisterImpressionHandler(func(e model.ImpressionArgs) {
		assert.Equal(t, reportingValue, e.ReportingValue)
		assert.Equal(t, ctx, e.Context)

		isImpressionRaised = true
	})

	impressionInvoker.Invoke(reportingValue, ctx)

	assert.True(t, isImpressionRaised)
}

func TestImpressionInvokerInvokeHandleUserCodePanic(t *testing.T) {
	internalFlags := &mocks.InternalFlags{}

	deps := &impression.ImpressionsDeps{
		InternalFlags:            internalFlags,
		CustomPropertyRepository: nil,
		DeviceProperties:         nil,
		IsRoxy:                   false,
	}
	impressionInvoker := impression.NewImpressionInvoker(deps)

	ctx := context.NewContext(map[string]interface{}{"obj1": 1})
	reportingValue := model.NewReportingValue("name", "value", true)

	isImpressionRaised := false
	impressionInvoker.RegisterImpressionHandler(func(e model.ImpressionArgs) {
		panic("mwahahahahaEvilUser")
	})

	impressionInvoker.Invoke(reportingValue, ctx)

	assert.False(t, isImpressionRaised)
}

func TestImpressionInvokerWillNotInvokeAnalyticsWhenFlagIsOff(t *testing.T) {
	internalFlags := &mocks.InternalFlags{}

	deps := &impression.ImpressionsDeps{
		InternalFlags:            internalFlags,
		CustomPropertyRepository: nil,
		DeviceProperties:         nil,
		IsRoxy:                   false,
	}
	impressionInvoker := impression.NewImpressionInvoker(deps)

	ctx := context.NewContext(map[string]interface{}{"obj1": 1})
	reportingValue := model.NewReportingValue("name", "value", true)

	impressionInvoker.Invoke(reportingValue, ctx)
}

func TestImpressionInvokerWillNotInvokeAnalyticsWhenIsRoxy(t *testing.T) {
	internalFlags := &mocks.InternalFlags{}

	deps := &impression.ImpressionsDeps{
		InternalFlags:            internalFlags,
		CustomPropertyRepository: nil,
		DeviceProperties:         nil,
		IsRoxy:                   true,
	}
	impressionInvoker := impression.NewImpressionInvoker(deps)

	ctx := context.NewContext(map[string]interface{}{"obj1": 1})
	reportingValue := model.NewReportingValue("name", "value", true)

	impressionInvoker.Invoke(reportingValue, ctx)
}

func TestImpressionInvokerWillInvokeAnalytics(t *testing.T) {
	internalFlags := &mocks.InternalFlags{}
	analytics := &mocks.Analytics{}
	analytics.On("CaptureImpressions", mock.Anything).Return().Times(1)

	deps := &impression.ImpressionsDeps{
		InternalFlags:            internalFlags,
		CustomPropertyRepository: nil,
		DeviceProperties:         nil,
		Analytics:                analytics,
		IsRoxy:                   false,
	}
	impressionInvoker := impression.NewImpressionInvoker(deps)

	ctx := context.NewContext(map[string]interface{}{"obj1": 1})
	reportingValue := model.NewReportingValue("name", "value", true)

	impressionInvoker.Invoke(reportingValue, ctx)
	mock.AssertExpectationsForObjects(t, analytics)
}
