package entities

import (
	"github.com/rollout/rox-go/v5/core/impression"
	"github.com/rollout/rox-go/v5/core/mocks"
	"github.com/rollout/rox-go/v5/core/model"
	"github.com/rollout/rox-go/v5/core/roxx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestRoxIntWithoutOptions(t *testing.T) {
	roxInt := NewRoxInt(1, nil)

	assert.Equal(t, 1, len(roxInt.Options()))
}

func TestRoxIntWillNotAddDefaultToOptionsIfExists(t *testing.T) {
	roxInt := NewRoxInt(1, []int{1, 2, 3})

	assert.Equal(t, 3, len(roxInt.Options()))
}

func TestRoxIntWillAddDefaultToOptionsIfNotExists(t *testing.T) {
	roxInt := NewRoxInt(1, []int{1, 2, 3})

	assert.Equal(t, 3, len(roxInt.Options()))
	assert.Contains(t, roxInt.Options(), 1)
}

func TestRoxIntWillSetName(t *testing.T) {
	roxInt := NewRoxInt(1, []int{2, 3})

	assert.Equal(t, "", roxInt.Name())

	roxInt.(model.InternalVariant).SetName("bop")

	assert.Equal(t, "bop", roxInt.Name())
}

func TestRoxIntWillReturnDefaultValueWhenNoParserOrCondition(t *testing.T) {
	roxInt := NewRoxInt(1, []int{2, 3})

	assert.Equal(t, 1, roxInt.GetValue(nil))

	roxInt.(model.InternalVariant).SetForEvaluation(&mocks.Parser{}, nil, nil)

	assert.Equal(t, 1, roxInt.GetValue(nil))

	parser := &mocks.Parser{}
	parser.On("EvaluateExpression", mock.Anything, mock.Anything).Return(roxx.NewEvaluationResult(nil))
	roxInt.(model.InternalVariant).SetForEvaluation(parser, model.NewExperimentModel("id", "name", "123", false, []string{"1"}, nil), nil)

	assert.Equal(t, 1, roxInt.GetValue(nil))
}

func TestRoxIntWillReturnDefaultValueWhenResultNotInOptions(t *testing.T) {
	parser := &mocks.Parser{}
	parser.On("EvaluateExpression", mock.Anything, mock.Anything).Return(roxx.NewEvaluationResult(666))

	roxInt := NewRoxInt(1, []int{2, 3})
	roxInt.(model.InternalVariant).SetForEvaluation(parser, model.NewExperimentModel("id", "name", "123", false, []string{"1"}, nil), nil)

	assert.Equal(t, 666, roxInt.GetValue(nil))
}

func TestRoxIntWillReturnValueWhenOnEvaluation(t *testing.T) {
	parser := &mocks.Parser{}
	parser.On("EvaluateExpression", mock.Anything, mock.Anything).Return(roxx.NewEvaluationResult(2))

	roxInt := NewRoxInt(1, []int{2, 3})
	roxInt.(model.InternalVariant).SetForEvaluation(parser, model.NewExperimentModel("id", "name", "123", false, []string{"1"}, nil), nil)

	assert.Equal(t, 2, roxInt.GetValue(nil))
}

func TestRoxIntWillRaiseImpressionInvoker(t *testing.T) {
	parser := &mocks.Parser{}
	parser.On("EvaluateExpression", mock.Anything, mock.Anything).Return(roxx.NewEvaluationResult(2))
	analytics := &mocks.Analytics{}
	analytics.On("IsAnalyticsReportingDisabled").Return(true)

	isImpressionRaised := false
	internalFlags := &mocks.InternalFlags{}
	impInvoker := impression.NewImpressionInvoker(&impression.ImpressionsDeps{
		InternalFlags:            internalFlags,
		CustomPropertyRepository: nil,
		DeviceProperties:         nil,
		Analytics:                analytics,
		IsRoxy:                   false,
	})
	impInvoker.RegisterImpressionHandler(func(e model.ImpressionArgs) {
		isImpressionRaised = true
	})

	roxInt := NewRoxInt(1, []int{2, 3})
	roxInt.(model.InternalVariant).SetForEvaluation(parser, model.NewExperimentModel("id", "name", "123", false, []string{"1"}, nil), impInvoker)

	assert.Equal(t, 2, roxInt.GetValue(nil))
	assert.True(t, isImpressionRaised)
}

func TestRoxIntForConsistencyWithString(t *testing.T) {
	parser := roxx.NewParser()

	isImpressionRaised := false
	internalFlags := &mocks.InternalFlags{}
	impInvoker := impression.NewImpressionInvoker(&impression.ImpressionsDeps{
		InternalFlags:            internalFlags,
		CustomPropertyRepository: nil,
		DeviceProperties:         nil,
		Analytics:                nil,
		IsRoxy:                   false,
	})
	impInvoker.RegisterImpressionHandler(func(e model.ImpressionArgs) {
		isImpressionRaised = true
	})

	roxInt := NewRoxInt(1, []int{2, 3})
	roxInt.(model.InternalVariant).SetForEvaluation(parser, model.NewExperimentModel("id", "name", `ifThen(true, "hi","hey")`, false, []string{"1"}, nil), impInvoker)

	assert.Equal(t, 1, roxInt.GetValue(nil))
	assert.False(t, isImpressionRaised)
}

func TestRoxIntForConsistencyWithDouble(t *testing.T) {
	parser := roxx.NewParser()

	isImpressionRaised := false
	internalFlags := &mocks.InternalFlags{}
	impInvoker := impression.NewImpressionInvoker(&impression.ImpressionsDeps{
		InternalFlags:            internalFlags,
		CustomPropertyRepository: nil,
		DeviceProperties:         nil,
		Analytics:                nil,
		IsRoxy:                   false,
	})
	impInvoker.RegisterImpressionHandler(func(e model.ImpressionArgs) {
		isImpressionRaised = true
	})

	roxInt := NewRoxInt(1, []int{2, 3})
	roxInt.(model.InternalVariant).SetForEvaluation(parser, model.NewExperimentModel("id", "name", `ifThen(true, 2.5,3.5)`, false, []string{"1"}, nil), impInvoker)

	assert.Equal(t, 1, roxInt.GetValue(nil))
	assert.False(t, isImpressionRaised)
}

func TestRoxIntForConsistencyWithBoolean(t *testing.T) {
	parser := roxx.NewParser()

	isImpressionRaised := false
	internalFlags := &mocks.InternalFlags{}
	impInvoker := impression.NewImpressionInvoker(&impression.ImpressionsDeps{
		InternalFlags:            internalFlags,
		CustomPropertyRepository: nil,
		DeviceProperties:         nil,
		Analytics:                nil,
		IsRoxy:                   false,
	})
	impInvoker.RegisterImpressionHandler(func(e model.ImpressionArgs) {
		isImpressionRaised = true
	})

	roxInt := NewRoxInt(1, []int{2, 3})
	roxInt.(model.InternalVariant).SetForEvaluation(parser, model.NewExperimentModel("id", "name", `ifThen(true, false, false)`, false, []string{"1"}, nil), impInvoker)

	assert.Equal(t, 1, roxInt.GetValue(nil))
	assert.False(t, isImpressionRaised)
}
