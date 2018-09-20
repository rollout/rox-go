package client_test

import (
	"github.com/rollout/rox-go/core/client"
	"github.com/rollout/rox-go/core/mocks"
	"github.com/rollout/rox-go/core/model"
	"github.com/rollout/rox-go/core/roxx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestInternalFlagsWillReturnFalseWhenNoExperiment(t *testing.T) {
	parser := &mocks.Parser{}
	expRepo := &mocks.ExperimentRepository{}
	expRepo.On("GetExperimentByFlag", mock.Anything).Return(nil)
	internalFlags := client.NewInternalFlags(expRepo, parser)

	assert.False(t, internalFlags.IsEnabled("stam"))
}

func TestInternalFlagsWillReturnFalseWhenExpressionIsNull(t *testing.T) {
	parser := &mocks.Parser{}
	parser.On("EvaluateExpression", mock.Anything, mock.Anything).Return(roxx.NewEvaluationResult(nil))
	expRepo := &mocks.ExperimentRepository{}
	expRepo.On("GetExperimentByFlag", mock.Anything).Return(model.NewExperimentModel("id", "name", "stam", false, nil, nil))
	internalFlags := client.NewInternalFlags(expRepo, parser)

	assert.False(t, internalFlags.IsEnabled("stam"))
}

func TestInternalFlagsWillReturnFalseWhenExpressionIsFalse(t *testing.T) {
	parser := &mocks.Parser{}
	parser.On("EvaluateExpression", mock.Anything, mock.Anything).Return(roxx.NewEvaluationResult(roxx.FlagFalseValue))
	expRepo := &mocks.ExperimentRepository{}
	expRepo.On("GetExperimentByFlag", mock.Anything).Return(model.NewExperimentModel("id", "name", "stam", false, nil, nil))
	internalFlags := client.NewInternalFlags(expRepo, parser)

	assert.False(t, internalFlags.IsEnabled("stam"))
}

func TestInternalFlagsWillReturnTrueWhenExpressionIsTrue(t *testing.T) {
	parser := &mocks.Parser{}
	parser.On("EvaluateExpression", mock.Anything, mock.Anything).Return(roxx.NewEvaluationResult(roxx.FlagTrueValue))
	expRepo := &mocks.ExperimentRepository{}
	expRepo.On("GetExperimentByFlag", mock.Anything).Return(model.NewExperimentModel("id", "name", "stam", false, nil, nil))
	internalFlags := client.NewInternalFlags(expRepo, parser)

	assert.True(t, internalFlags.IsEnabled("stam"))
}
