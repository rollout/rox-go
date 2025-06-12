package impression_test

import (
	"testing"

	"github.com/rollout/rox-go/v6/core/consts"
	"github.com/rollout/rox-go/v6/core/context"
	"github.com/rollout/rox-go/v6/core/impression"
	"github.com/rollout/rox-go/v6/core/mocks"
	"github.com/rollout/rox-go/v6/core/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const distinctId = "123abc456def789xyz"

func devicePropertiesMock() *mocks.DeviceProperties {
	deviceProperties := &mocks.DeviceProperties{}

	properties := map[string]string{
		consts.PropertyTypeDistinctID.Name: distinctId,
	}
	deviceProperties.On("GetAllProperties").Return(properties)
	deviceProperties.On("DistinctID").Return(distinctId)

	return deviceProperties
}

func TestImpressionInvokerEmptyInvokeNotThrowingException(t *testing.T) {
	internalFlags := &mocks.InternalFlags{}

	deps := &impression.ImpressionsDeps{
		InternalFlags:            internalFlags,
		CustomPropertyRepository: nil,
		DeviceProperties:         devicePropertiesMock(),
		Analytics:                nil,
		IsRoxy:                   false,
	}
	impressionInvoker := impression.NewImpressionInvoker(deps)
	impressionInvoker.Invoke(nil, nil)
}

func TestImpressionInvokerInvokeAndParameters(t *testing.T) {
	internalFlags := &mocks.InternalFlags{}
	internalFlags.On("IsEnabled", "rox.internal.analytics").Return(true)

	analytics := &mocks.Analytics{}
	analytics.On("CaptureImpressions", mock.Anything)

	deps := &impression.ImpressionsDeps{
		InternalFlags:    internalFlags,
		DeviceProperties: devicePropertiesMock(),
		Analytics:        analytics,
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
		DeviceProperties:         devicePropertiesMock(),
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
	internalFlags.On("IsEnabled", "rox.internal.analytics").Return(false)

	analytics := &mocks.Analytics{}
	analytics.On("CaptureImpressions", mock.Anything)

	deps := &impression.ImpressionsDeps{
		InternalFlags:            internalFlags,
		CustomPropertyRepository: nil,
		DeviceProperties:         devicePropertiesMock(),
		Analytics:                analytics,
		IsRoxy:                   false,
	}
	impressionInvoker := impression.NewImpressionInvoker(deps)

	ctx := context.NewContext(map[string]interface{}{"obj1": 1})
	reportingValue := model.NewReportingValue("name", "value", true)

	impressionInvoker.Invoke(reportingValue, ctx)

	// Should not have been called, due to rox.internal.analytics feature flag.
	analytics.AssertNumberOfCalls(t, "CaptureImpressions", 0)
}

func TestImpressionInvokerWillNotInvokeAnalyticsWhenIsRoxy(t *testing.T) {
	internalFlags := &mocks.InternalFlags{}
	analytics := &mocks.Analytics{}

	deps := &impression.ImpressionsDeps{
		InternalFlags:            internalFlags,
		CustomPropertyRepository: nil,
		DeviceProperties:         devicePropertiesMock(),
		Analytics:                analytics,
		IsRoxy:                   true,
	}
	impressionInvoker := impression.NewImpressionInvoker(deps)

	ctx := context.NewContext(map[string]interface{}{"obj1": 1})
	reportingValue := model.NewReportingValue("name", "value", true)

	impressionInvoker.Invoke(reportingValue, ctx)

	// Should not have been called, due to Roxy mode.
	analytics.AssertNumberOfCalls(t, "CaptureImpressions", 0)
}

func TestImpressionInvokerWillInvokeAnalytics(t *testing.T) {
	internalFlags := &mocks.InternalFlags{}
	internalFlags.On("IsEnabled", "rox.internal.analytics").Return(true)

	flagName := "name"
	flagValue := "value"

	analytics := &mocks.Analytics{}
	analytics.On("CaptureImpressions", mock.MatchedBy(impressionMatcher(flagName, flagValue)))

	deps := &impression.ImpressionsDeps{
		InternalFlags:            internalFlags,
		CustomPropertyRepository: nil,
		DeviceProperties:         devicePropertiesMock(),
		Analytics:                analytics,
		IsRoxy:                   false,
	}
	impressionInvoker := impression.NewImpressionInvoker(deps)

	ctx := context.NewContext(map[string]interface{}{"obj1": 1})
	reportingValue := model.NewReportingValue(flagName, flagValue, true)

	impressionInvoker.Invoke(reportingValue, ctx)
	mock.AssertExpectationsForObjects(t, analytics)
}

func impressionMatcher(flagName, flagValue string) func(impressions []model.Impression) bool {

	return func(impressions []model.Impression) bool {
		// No batching implemented atm.
		if len(impressions) != 1 {
			return false
		}
		i := impressions[0]

		// Fixed/static fields (apart from timestamp... but that is hard to validate beyond 'not zero'):
		// (validating unset for the optional fields Count & Type).
		if i.Count != 0 || i.Type != "" || i.Timestamp == 0 {
			return false
		}

		// Dynamic fields relating to the ReportingValue:
		if i.FlagName != flagName || i.Value != flagValue {
			return false
		}

		return true
	}
}
