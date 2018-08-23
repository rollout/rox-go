package entities

import (
	"github.com/rollout/rox-go/core/context"
	"github.com/rollout/rox-go/core/model"
	"github.com/rollout/rox-go/core/roxx"
	"github.com/rollout/rox-go/core/utils"
)

type variant struct {
	defaultValue      string
	options           []string
	condition         string
	parser            roxx.Parser
	globalContext     context.Context
	impressionInvoker model.ImpressionInvoker
	clientExperiment  *model.Experiment
	name              string
}

func NewVariant(defaultValue string, options []string) model.Variant {
	allOptions := make([]string, len(options))
	copy(allOptions, options)
	if !utils.ContainsString(allOptions, defaultValue) {
		allOptions = append(allOptions, defaultValue)
	}

	return &variant{
		defaultValue: defaultValue,
		options:      allOptions,
	}
}

func (v *variant) DefaultValue() string {
	return v.defaultValue
}

func (v *variant) Options() []string {
	return v.options
}

func (v *variant) Name() string {
	return v.name
}

func (v *variant) SetForEvaluation(parser roxx.Parser, experiment *model.ExperimentModel, impressionInvoker model.ImpressionInvoker) {
	if experiment != nil {
		v.clientExperiment = model.NewExperiment(experiment)
		v.condition = experiment.Condition
	} else {
		v.clientExperiment = nil
		v.condition = ""
	}

	v.parser = parser
	v.impressionInvoker = impressionInvoker
}

func (v *variant) SetContext(globalContext context.Context) {
	v.globalContext = globalContext
}

func (v *variant) SetName(name string) {
	v.name = name
}

func (v *variant) GetValue(ctx context.Context) string {
	returnValue := v.defaultValue
	mergedContext := context.NewMergedContext(v.globalContext, ctx)

	if v.parser != nil && v.condition != "" {
		evaluationResult := v.parser.EvaluateExpression(v.condition, mergedContext)
		value := evaluationResult.StringValue()
		if value != "" {
			if utils.ContainsString(v.options, value) {
				returnValue = value
			}
		}
	}

	if v.impressionInvoker != nil {
		v.impressionInvoker.Invoke(model.NewReportingValue(v.name, returnValue), v.clientExperiment, mergedContext)
	}

	return returnValue
}
