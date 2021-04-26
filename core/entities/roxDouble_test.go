package entities

import (
	"github.com/rollout/rox-go/core/impression"
	"github.com/rollout/rox-go/core/mocks"
	"github.com/rollout/rox-go/core/model"
	"github.com/rollout/rox-go/core/roxx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestRoxDoubleWithoutOptions(t *testing.T) {
	roxDouble := NewRoxDouble(1.0, nil)

	assert.Equal(t, 1, len(roxDouble.Options()))
}

func TestRoxDoubleWillNotAddDefaultToOptionsIfExists(t *testing.T) {
	roxDouble := NewRoxDouble(1.0, []float64{1.0, 2.0, 3.0})

	assert.Equal(t, 3, len(roxDouble.Options()))
}

func TestRoxDoubleWillAddDefaultToOptionsIfNotExists(t *testing.T) {
	roxDouble := NewRoxDouble(1.0, []float64{1.0, 2.0, 3.0})

	assert.Equal(t, 3, len(roxDouble.Options()))
	assert.Contains(t, roxDouble.Options(), 1.0)
}

func TestRoxDoubleWillSetName(t *testing.T) {
	roxDouble := NewRoxDouble(1.0, []float64{2.0, 3.0})

	assert.Equal(t, "", roxDouble.Name())

	roxDouble.(model.InternalVariant).SetName("bop")

	assert.Equal(t, "bop", roxDouble.Name())
}

func TestRoxDoubleWillReturnDefaultValueWhenNoParserOrCondition(t *testing.T) {
	roxDouble := NewRoxDouble(1.0, []float64{2.0, 3.0})

	assert.Equal(t, 1.0, roxDouble.GetValue(nil))

	roxDouble.(model.InternalVariant).SetForEvaluation(&mocks.Parser{}, nil, nil)

	assert.Equal(t, 1.0, roxDouble.GetValue(nil))

	parser := &mocks.Parser{}
	parser.On("EvaluateExpression", mock.Anything, mock.Anything).Return(roxx.NewEvaluationResult(nil))
	roxDouble.(model.InternalVariant).SetForEvaluation(parser, model.NewExperimentModel("id", "name", "123", false, []string{"1"}, nil), nil)

	assert.Equal(t, 1.0, roxDouble.GetValue(nil))
}

func TestRoxDoubleWillReturnDefaultValueWhenResultNotInOptions(t *testing.T) {
	parser := &mocks.Parser{}
	parser.On("EvaluateExpression", mock.Anything, mock.Anything).Return(roxx.NewEvaluationResult(666.0))

	roxDouble := NewRoxDouble(1.0, []float64{2.0, 3.0})
	roxDouble.(model.InternalVariant).SetForEvaluation(parser, model.NewExperimentModel("id", "name", "123", false, []string{"1"}, nil), nil)

	assert.Equal(t, 666.0, roxDouble.GetValue(nil))
}

func TestRoxDoubleWillReturnValueWhenOnEvaluation(t *testing.T) {
	parser := &mocks.Parser{}
	parser.On("EvaluateExpression", mock.Anything, mock.Anything).Return(roxx.NewEvaluationResult(2.0))

	roxDouble := NewRoxDouble(1.0, []float64{2.0, 3.0})
	roxDouble.(model.InternalVariant).SetForEvaluation(parser, model.NewExperimentModel("id", "name", "123", false, []string{"1"}, nil), nil)

	assert.Equal(t, 2.0, roxDouble.GetValue(nil))
}

func TestRoxDoubleWillRaiseImpressionInvoker(t *testing.T) {
	parser := &mocks.Parser{}
	parser.On("EvaluateExpression", mock.Anything, mock.Anything).Return(roxx.NewEvaluationResult(2.0))

	isImpressionRaised := false
	internalFlags := &mocks.InternalFlags{}
	impInvoker := impression.NewImpressionInvoker(internalFlags, nil, nil, false)
	impInvoker.RegisterImpressionHandler(func(e model.ImpressionArgs) {
		isImpressionRaised = true
	})

	roxDouble := NewRoxDouble(1.0, []float64{2.0, 3.0})
	roxDouble.(model.InternalVariant).SetForEvaluation(parser, model.NewExperimentModel("id", "name", "123", false, []string{"1"}, nil), impInvoker)

	assert.Equal(t, 2.0, roxDouble.GetValue(nil))
	assert.True(t, isImpressionRaised)
}
