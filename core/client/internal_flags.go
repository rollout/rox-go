package client

import (
	"github.com/rollout/rox-go/core/model"
	"github.com/rollout/rox-go/core/roxx"
)

type internalFlags struct {
	experimentRepository model.ExperimentRepository
	parser               roxx.Parser
}

func NewInternalFlags(experimentRepository model.ExperimentRepository, parser roxx.Parser) model.InternalFlags {
	return &internalFlags{
		experimentRepository: experimentRepository,
		parser:               parser,
	}
}

func (f *internalFlags) IsEnabled(flagName string) bool {
	internalExperiment := f.experimentRepository.GetExperimentByFlag(flagName)
	if internalExperiment == nil {
		return false
	}

	value := f.parser.EvaluateExpression(internalExperiment.Condition, nil).StringValue()
	return value == roxx.FlagTrueValue
}
