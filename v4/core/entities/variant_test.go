package entities

import (
	"github.com/rollout/rox-go/v4/core/impression"
	"github.com/rollout/rox-go/v4/core/mocks"
	"github.com/rollout/rox-go/v4/core/model"
	"github.com/rollout/rox-go/v4/core/roxx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestVariantWithoutOptions(t *testing.T) {
	variant := NewVariant("1", nil)

	assert.Equal(t, 1, len(variant.Options()))
}

func TestVariantWillNotAddDefaultToOptionsIfExists(t *testing.T) {
	variant := NewVariant("1", []string{"1", "2", "3"})

	assert.Equal(t, 3, len(variant.Options()))
}

func TestVariantWillAddDefaultToOptionsIfNotExists(t *testing.T) {
	variant := NewVariant("1", []string{"2", "3"})

	assert.Equal(t, 3, len(variant.Options()))
	assert.Contains(t, variant.Options(), "1")
}

func TestVariantWillSetName(t *testing.T) {
	variant := NewVariant("1", []string{"2", "3"})

	assert.Equal(t, "", variant.Name())

	variant.(model.InternalVariant).SetName("bop")

	assert.Equal(t, "bop", variant.Name())
}

func TestVariantWillReturnDefaultValueWhenNoParserOrCondition(t *testing.T) {
	variant := NewVariant("1", []string{"2", "3"})

	assert.Equal(t, "1", variant.GetValue(nil))

	variant.(model.InternalVariant).SetForEvaluation(&mocks.Parser{}, nil, nil)

	assert.Equal(t, "1", variant.GetValue(nil))

	parser := &mocks.Parser{}
	parser.On("EvaluateExpression", mock.Anything, mock.Anything).Return(roxx.NewEvaluationResult(nil))
	variant.(model.InternalVariant).SetForEvaluation(parser, model.NewExperimentModel("id", "name", "123", false, []string{"1"}, nil), nil)

	assert.Equal(t, "1", variant.GetValue(nil))
}

func TestVariantWillReturnDefaultValueWhenResultNotInOptions(t *testing.T) {
	parser := &mocks.Parser{}
	parser.On("EvaluateExpression", mock.Anything, mock.Anything).Return(roxx.NewEvaluationResult("xxx"))

	variant := NewVariant("1", []string{"2", "3"})
	variant.(model.InternalVariant).SetForEvaluation(parser, model.NewExperimentModel("id", "name", "123", false, []string{"1"}, nil), nil)

	assert.Equal(t, "xxx", variant.GetValue(nil))
}

func TestVariantWillReturnValueWhenOnEvaluation(t *testing.T) {
	parser := &mocks.Parser{}
	parser.On("EvaluateExpression", mock.Anything, mock.Anything).Return(roxx.NewEvaluationResult("2"))

	variant := NewVariant("1", []string{"2", "3"})
	variant.(model.InternalVariant).SetForEvaluation(parser, model.NewExperimentModel("id", "name", "123", false, []string{"1"}, nil), nil)

	assert.Equal(t, "2", variant.GetValue(nil))
}

func TestVariantWillRaiseImpression(t *testing.T) {
	parser := &mocks.Parser{}
	parser.On("EvaluateExpression", mock.Anything, mock.Anything).Return(roxx.NewEvaluationResult("2"))

	isImpressionRaised := false
	internalFlags := &mocks.InternalFlags{}
	impInvoker := impression.NewImpressionInvoker(internalFlags, nil, nil, false)
	impInvoker.RegisterImpressionHandler(func(e model.ImpressionArgs) {
		isImpressionRaised = true
	})

	variant := NewVariant("1", []string{"2", "3"})
	variant.(model.InternalVariant).SetForEvaluation(parser, model.NewExperimentModel("id", "name", "123", false, []string{"1"}, nil), impInvoker)

	assert.Equal(t, "2", variant.GetValue(nil))
	assert.True(t, isImpressionRaised)
}
