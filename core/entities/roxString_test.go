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

func TestRoxStringWithoutOptions(t *testing.T) {
	roxString := NewRoxString("1", nil)

	assert.Equal(t, 1, len(roxString.Options()))
}

func TestRoxStringWillNotAddDefaultToOptionsIfExists(t *testing.T) {
	roxString := NewRoxString("1", []string{"1", "2", "3"})

	assert.Equal(t, 3, len(roxString.Options()))
}

func TestRoxStringWillAddDefaultToOptionsIfNotExists(t *testing.T) {
	roxString := NewRoxString("1", []string{"2", "3"})

	assert.Equal(t, 3, len(roxString.Options()))
	assert.Contains(t, roxString.Options(), "1")
}

func TestRoxStringWillSetName(t *testing.T) {
	roxString := NewRoxString("1", []string{"2", "3"})

	assert.Equal(t, "", roxString.Name())

	roxString.(model.InternalVariant).SetName("bop")

	assert.Equal(t, "bop", roxString.Name())
}

func TestRoxStringWillReturnDefaultValueWhenNoParserOrCondition(t *testing.T) {
	roxString := NewRoxString("1", []string{"2", "3"})

	assert.Equal(t, "1", roxString.GetValueAsString(nil))

	roxString.(model.InternalVariant).SetForEvaluation(&mocks.Parser{}, nil, nil)

	assert.Equal(t, "1", roxString.GetValueAsString(nil))

	parser := &mocks.Parser{}
	parser.On("EvaluateExpression", mock.Anything, mock.Anything).Return(roxx.NewEvaluationResult(nil))
	roxString.(model.InternalVariant).SetForEvaluation(parser, model.NewExperimentModel("id", "name", "123", false, []string{"1"}, nil), nil)

	assert.Equal(t, "1", roxString.GetValueAsString(nil))
}

func TestRoxStringWillReturnDefaultValueWhenResultNotInOptions(t *testing.T) {
	parser := &mocks.Parser{}
	parser.On("EvaluateExpression", mock.Anything, mock.Anything).Return(roxx.NewEvaluationResult("xxx"))

	roxString := NewRoxString("1", []string{"2", "3"})
	roxString.(model.InternalVariant).SetForEvaluation(parser, model.NewExperimentModel("id", "name", "123", false, []string{"1"}, nil), nil)

	assert.Equal(t, "xxx", roxString.GetValueAsString(nil))
}

func TestRoxStringWillReturnValueWhenOnEvaluation(t *testing.T) {
	parser := &mocks.Parser{}
	parser.On("EvaluateExpression", mock.Anything, mock.Anything).Return(roxx.NewEvaluationResult("2"))

	roxString := NewRoxString("1", []string{"2", "3"})
	roxString.(model.InternalVariant).SetForEvaluation(parser, model.NewExperimentModel("id", "name", "123", false, []string{"1"}, nil), nil)

	assert.Equal(t, "2", roxString.GetValue(nil))
}

func TestRoxStringWillRaiseImpression(t *testing.T) {
	parser := &mocks.Parser{}
	parser.On("EvaluateExpression", mock.Anything, mock.Anything).Return(roxx.NewEvaluationResult("2"))

	isImpressionRaised := false
	internalFlags := &mocks.InternalFlags{}
	impInvoker := impression.NewImpressionInvoker(internalFlags, nil, nil, false)
	impInvoker.RegisterImpressionHandler(func(e model.ImpressionArgs) {
		isImpressionRaised = true
	})

	roxString := NewRoxString("1", []string{"2", "3"})
	roxString.(model.InternalVariant).SetForEvaluation(parser, model.NewExperimentModel("id", "name", `ifThen(true, "hi","hey")`, false, []string{"1"}, nil), impInvoker)

	assert.Equal(t, "2", roxString.GetValueAsString(nil))
	assert.True(t, isImpressionRaised)
}


func TestRoxStringForConsistencyWithInt(t *testing.T) {
	parser := roxx.NewParser()

	isImpressionRaised := false
	internalFlags := &mocks.InternalFlags{}
	impInvoker := impression.NewImpressionInvoker(internalFlags, nil, nil, false)
	impInvoker.RegisterImpressionHandler(func(e model.ImpressionArgs) {
		isImpressionRaised = true
	})

	roxString := NewRoxString("a", []string{"b", "c"})
	roxString.(model.InternalVariant).SetForEvaluation(parser, model.NewExperimentModel("id", "name", `ifThen(true, 2, 3)`, false, []string{"1"}, nil), impInvoker)

	assert.Equal(t, "a", roxString.GetValue(nil))
	assert.False(t, isImpressionRaised)
}

func TestRoxStringForConsistencyWithDouble(t *testing.T) {
	parser := roxx.NewParser()

	isImpressionRaised := false
	internalFlags := &mocks.InternalFlags{}
	impInvoker := impression.NewImpressionInvoker(internalFlags, nil, nil, false)
	impInvoker.RegisterImpressionHandler(func(e model.ImpressionArgs) {
		isImpressionRaised = true
	})

	roxString := NewRoxString("a", []string{"b", "c"})
	roxString.(model.InternalVariant).SetForEvaluation(parser, model.NewExperimentModel("id", "name", `ifThen(true, 2.5, 3.5)`, false, []string{"1"}, nil), impInvoker)

	assert.Equal(t, "a", roxString.GetValue(nil))
	assert.False(t, isImpressionRaised)
}

func TestRoxStringForConsistencyWithBool(t *testing.T) {
	parser := roxx.NewParser()

	isImpressionRaised := false
	internalFlags := &mocks.InternalFlags{}
	impInvoker := impression.NewImpressionInvoker(internalFlags, nil, nil, false)
	impInvoker.RegisterImpressionHandler(func(e model.ImpressionArgs) {
		isImpressionRaised = true
	})

	roxString := NewRoxString("a", []string{"b", "c"})
	roxString.(model.InternalVariant).SetForEvaluation(parser, model.NewExperimentModel("id", "name", `ifThen(true, false, false)`, false, []string{"1"}, nil), impInvoker)

	assert.Equal(t, "a", roxString.GetValue(nil))
	assert.False(t, isImpressionRaised)
}
