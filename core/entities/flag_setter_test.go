package entities

import (
	"github.com/rollout/rox-go/core/impression"
	"github.com/rollout/rox-go/core/mocks"
	"github.com/rollout/rox-go/core/model"
	"github.com/rollout/rox-go/core/repositories"
	"github.com/rollout/rox-go/core/roxx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestFlagSetterWillSetFlagData(t *testing.T) {
	flagRepo := repositories.NewFlagRepository()
	expRepo := repositories.NewExperimentRepository()

	parser := &mocks.Parser{}
	parser.On("EvaluateExpression", mock.Anything, mock.Anything).Return(roxx.NewEvaluationResult(nil))

	internalFlags := &mocks.InternalFlags{}
	impInvoker := impression.NewImpressionInvoker(internalFlags, nil, nil, false)

	flagRepo.AddFlag(NewFlag(false), "f1")
	expRepo.SetExperiments([]*model.ExperimentModel{
		model.NewExperimentModel("33", "1", "1", false, []string{"f1"}, nil),
	})

	flagSetter := NewFlagSetter(flagRepo, parser, expRepo, impInvoker)
	flagSetter.SetExperiments()

	assert.Equal(t, "1", flagRepo.GetFlag("f1").(internalVariant).Condition())
	assert.Equal(t, parser, flagRepo.GetFlag("f1").(internalVariant).Parser())
	assert.Equal(t, impInvoker, flagRepo.GetFlag("f1").(internalVariant).ImpressionInvoker())
	assert.Equal(t, "33", flagRepo.GetFlag("f1").(internalVariant).ClientExperiment().Identifier)
}

func TestFlagSetterWillNotSetForOtherFlag(t *testing.T) {
	flagRepo := repositories.NewFlagRepository()
	expRepo := repositories.NewExperimentRepository()

	parser := &mocks.Parser{}
	parser.On("EvaluateExpression", mock.Anything, mock.Anything).Return(roxx.NewEvaluationResult(nil))

	internalFlags := &mocks.InternalFlags{}
	impInvoker := impression.NewImpressionInvoker(internalFlags, nil, nil, false)

	flagRepo.AddFlag(NewFlag(false), "f1")
	flagRepo.AddFlag(NewFlag(false), "f2")
	expRepo.SetExperiments([]*model.ExperimentModel{
		model.NewExperimentModel("1", "1", "1", false, []string{"f1"}, nil),
	})

	flagSetter := NewFlagSetter(flagRepo, parser, expRepo, impInvoker)
	flagSetter.SetExperiments()

	assert.Equal(t, "", flagRepo.GetFlag("f2").(internalVariant).Condition())
	assert.Equal(t, parser, flagRepo.GetFlag("f2").(internalVariant).Parser())
	assert.Equal(t, impInvoker, flagRepo.GetFlag("f2").(internalVariant).ImpressionInvoker())
	assert.Nil(t, flagRepo.GetFlag("f2").(internalVariant).ClientExperiment())
}

func TestFlagSetterWillSetExperimentForFlagAndWillRemoveIt(t *testing.T) {
	flagRepo := repositories.NewFlagRepository()
	expRepo := repositories.NewExperimentRepository()

	parser := &mocks.Parser{}
	parser.On("EvaluateExpression", mock.Anything, mock.Anything).Return(roxx.NewEvaluationResult(nil))

	internalFlags := &mocks.InternalFlags{}
	impInvoker := impression.NewImpressionInvoker(internalFlags, nil, nil, false)

	flagSetter := NewFlagSetter(flagRepo, parser, expRepo, impInvoker)
	flagRepo.AddFlag(NewFlag(false), "f2")

	expRepo.SetExperiments([]*model.ExperimentModel{
		model.NewExperimentModel("id1", "1", "con", false, []string{"f2"}, nil),
	})
	flagSetter.SetExperiments()

	assert.Equal(t, "con", flagRepo.GetFlag("f2").(internalVariant).Condition())
	assert.Equal(t, parser, flagRepo.GetFlag("f2").(internalVariant).Parser())
	assert.Equal(t, impInvoker, flagRepo.GetFlag("f2").(internalVariant).ImpressionInvoker())
	assert.Equal(t, "id1", flagRepo.GetFlag("f2").(internalVariant).ClientExperiment().Identifier)

	expRepo.SetExperiments([]*model.ExperimentModel{})
	flagSetter.SetExperiments()

	assert.Equal(t, "", flagRepo.GetFlag("f2").(internalVariant).Condition())
	assert.Equal(t, parser, flagRepo.GetFlag("f2").(internalVariant).Parser())
	assert.Equal(t, impInvoker, flagRepo.GetFlag("f2").(internalVariant).ImpressionInvoker())
	assert.Nil(t, flagRepo.GetFlag("f2").(internalVariant).ClientExperiment())
}

func TestFlagSetterWillSetFlagWithoutExperimentAndThenAddExperiment(t *testing.T) {
	flagRepo := repositories.NewFlagRepository()
	expRepo := repositories.NewExperimentRepository()

	parser := &mocks.Parser{}
	parser.On("EvaluateExpression", mock.Anything, mock.Anything).Return(roxx.NewEvaluationResult(nil))

	internalFlags := &mocks.InternalFlags{}
	impInvoker := impression.NewImpressionInvoker(internalFlags, nil, nil, false)

	flagSetter := NewFlagSetter(flagRepo, parser, expRepo, impInvoker)
	flagRepo.AddFlag(NewFlag(false), "f2")
	flagSetter.SetExperiments()

	assert.Equal(t, "", flagRepo.GetFlag("f2").(internalVariant).Condition())
	assert.Equal(t, parser, flagRepo.GetFlag("f2").(internalVariant).Parser())
	assert.Equal(t, impInvoker, flagRepo.GetFlag("f2").(internalVariant).ImpressionInvoker())
	assert.Nil(t, flagRepo.GetFlag("f2").(internalVariant).ClientExperiment())

	expRepo.SetExperiments([]*model.ExperimentModel{
		model.NewExperimentModel("id1", "1", "con", false, []string{"f2"}, nil),
	})
	flagSetter.SetExperiments()

	assert.Equal(t, "con", flagRepo.GetFlag("f2").(internalVariant).Condition())
	assert.Equal(t, parser, flagRepo.GetFlag("f2").(internalVariant).Parser())
	assert.Equal(t, impInvoker, flagRepo.GetFlag("f2").(internalVariant).ImpressionInvoker())
	assert.Equal(t, "id1", flagRepo.GetFlag("f2").(internalVariant).ClientExperiment().Identifier)
}

func TestFlagSetterWillSetDataForAddedFlag(t *testing.T) {
	flagRepo := repositories.NewFlagRepository()
	expRepo := repositories.NewExperimentRepository()

	parser := &mocks.Parser{}
	parser.On("EvaluateExpression", mock.Anything, mock.Anything).Return(roxx.NewEvaluationResult(nil))

	internalFlags := &mocks.InternalFlags{}
	impInvoker := impression.NewImpressionInvoker(internalFlags, nil, nil, false)

	expRepo.SetExperiments([]*model.ExperimentModel{
		model.NewExperimentModel("1", "1", "1", false, []string{"f1"}, nil),
	})
	flagSetter := NewFlagSetter(flagRepo, parser, expRepo, impInvoker)
	flagSetter.SetExperiments()

	flagRepo.AddFlag(NewFlag(false), "f1")
	flagRepo.AddFlag(NewFlag(false), "f2")

	assert.Equal(t, "1", flagRepo.GetFlag("f1").(internalVariant).Condition())
	assert.Equal(t, "", flagRepo.GetFlag("f2").(internalVariant).Condition())
	assert.Equal(t, parser, flagRepo.GetFlag("f1").(internalVariant).Parser())
	assert.Equal(t, parser, flagRepo.GetFlag("f2").(internalVariant).Parser())
	assert.Equal(t, impInvoker, flagRepo.GetFlag("f1").(internalVariant).ImpressionInvoker())
	assert.Equal(t, impInvoker, flagRepo.GetFlag("f2").(internalVariant).ImpressionInvoker())
	assert.Equal(t, "1", flagRepo.GetFlag("f1").(internalVariant).ClientExperiment().Identifier)
	assert.Nil(t, flagRepo.GetFlag("f2").(internalVariant).ClientExperiment())
}
