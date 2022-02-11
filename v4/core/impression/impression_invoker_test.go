package impression_test

import (
	"github.com/rollout/rox-go/v4/core/context"
	"github.com/rollout/rox-go/v4/core/impression"
	"github.com/rollout/rox-go/v4/core/mocks"
	"github.com/rollout/rox-go/v4/core/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestImpressionInvokerEmptyInvokeNotThrowingException(t *testing.T) {
	internalFlags := &mocks.InternalFlags{}
	impressionInvoker := impression.NewImpressionInvoker(internalFlags, nil, nil, false)
	impressionInvoker.Invoke(nil, nil, nil)
}

func TestImpressionInvokerInvokeAndParameters(t *testing.T) {
	internalFlags := &mocks.InternalFlags{}
	impressionInvoker := impression.NewImpressionInvoker(internalFlags, nil, nil, false)

	ctx := context.NewContext(map[string]interface{}{"obj1": 1})
	reportingValue := model.NewReportingValue("name", "value")
	originalExperiment := model.NewExperimentModel("id", "name", "cond", true, nil, nil)
	experiment := model.NewExperiment(originalExperiment)

	isImpressionRaised := false
	impressionInvoker.RegisterImpressionHandler(func(e model.ImpressionArgs) {
		assert.Equal(t, reportingValue, e.ReportingValue)
		assert.Equal(t, experiment, e.Experiment)
		assert.Equal(t, ctx, e.Context)

		isImpressionRaised = true
	})

	impressionInvoker.Invoke(reportingValue, experiment, ctx)

	assert.True(t, isImpressionRaised)
}

func TestImpressionInvokerWillNotInvokeAnalyticsWhenFlagIsOff(t *testing.T) {
	// TODO
}

func TestImpressionInvokerWillNotInvokeAnalyticsWhenIsRoxy(t *testing.T) {
	// TODO
}

func TestImpressionInvokerWillInvokeAnalytics(t *testing.T) {
	// TODO
}

func TestImpressionInvokerWillInvokeAnalyticsWithBadDistinctID(t *testing.T) {
	// TODO
}
