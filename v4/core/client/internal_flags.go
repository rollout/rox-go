package client

import (
	"github.com/rollout/rox-go/v4/core/model"
	"github.com/rollout/rox-go/v4/core/roxx"
)

type internalFlags struct {
	experimentRepository model.ExperimentRepository
	parser               roxx.Parser
	environment          model.Environment
}

var FLAG_DEFAULTS = map[string]string{
	"rox.internal.pushUpdates": "false",
}

func NewInternalFlags(experimentRepository model.ExperimentRepository, parser roxx.Parser, environment model.Environment) model.InternalFlags {
	return &internalFlags{
		experimentRepository: experimentRepository,
		parser:               parser,
		environment:          environment,
	}
}

func (f *internalFlags) IsEnabled(flagName string) bool {
	value := ""
	if f.environment.IsSelfManaged() {
		value = FLAG_DEFAULTS[flagName]
	} else {
		internalExperiment := f.experimentRepository.GetExperimentByFlag(flagName)
		if internalExperiment == nil {
			return false
		}

		value = f.parser.EvaluateExpression(internalExperiment.Condition, nil).StringValue()
	}
	return value == roxx.FlagTrueValue
}
