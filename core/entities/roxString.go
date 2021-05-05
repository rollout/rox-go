package entities

import (
	"github.com/rollout/rox-go/core/consts"
	"github.com/rollout/rox-go/core/context"
	"github.com/rollout/rox-go/core/model"
	"github.com/rollout/rox-go/core/roxx"
	"github.com/rollout/rox-go/core/utils"
)

type internalVariant interface {
	Condition() string
	Parser() roxx.Parser
	ImpressionInvoker() model.ImpressionInvoker
	ClientExperiment() *model.Experiment
}

type roxString struct {
	roxVariant
	defaultValue      string
	options           []string
	condition         string
	parser            roxx.Parser
	globalContext     context.Context
	impressionInvoker model.ImpressionInvoker
	clientExperiment  *model.Experiment
}

func NewRoxString(defaultValue string, options []string) model.RoxString {
	if options == nil {
		options = []string{}
	}
	allOptions := make([]string, len(options))
	copy(allOptions, options)
	if !utils.ContainsString(allOptions, defaultValue) {
		allOptions = append(allOptions, defaultValue)
	}

	roxString := &roxString{
		roxVariant: roxVariant{
			flagType: consts.StringType,
		},
		defaultValue: defaultValue,
		options:      allOptions,
	}

	return roxString
}

func (v *roxString) GetDefaultAsString() string {
	return v.DefaultValue()
}

func (v *roxString) DefaultValue() string {
	return v.defaultValue
}

func (v *roxString) GetOptionsAsString() []string {
	return v.Options()
}

func (v *roxString) Options() []string {
	return v.options
}

func (v *roxString) SetForEvaluation(parser roxx.Parser, experiment *model.ExperimentModel, impressionInvoker model.ImpressionInvoker) {
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

func (v *roxString) SetContext(globalContext context.Context) {
	v.globalContext = globalContext
}

func (v *roxString) SetName(name string) {
	v.name = name
}

func (v *roxString) GetValueAsString(ctx context.Context) string {
	return v.GetValue(ctx)
}

func (v *roxString) GetValue(ctx context.Context) string {
	returnValue, _ := v.InternalGetValue(ctx)
	return returnValue
}

func (v *roxString) InternalGetValue(ctx context.Context) (returnValue string, isDefault bool) {
	returnValue, isDefault = v.defaultValue, true
	mergedContext := context.NewMergedContext(v.globalContext, ctx)
	sendImpression := false

	if v.parser != nil && v.condition != "" {
		evaluationResult := v.parser.EvaluateExpression(v.condition, mergedContext)
		value := evaluationResult.StringValue()
		if value != "" {
			switch v.FlagType() {
			case consts.StringType:
				if _, ok := evaluationResult.Value().(string); ok {
					returnValue, isDefault = value, false
					sendImpression = true
				}
			case consts.BoolType:
				if value == roxx.FlagFalseValue || value == roxx.FlagTrueValue {
					returnValue, isDefault = value, false
					sendImpression = true
				}
			}
		}
	}

	if v.impressionInvoker != nil && sendImpression {
		v.impressionInvoker.Invoke(model.NewReportingValue(v.name, returnValue), v.clientExperiment, mergedContext)
	}

	return returnValue, isDefault
}

func (v *roxString) Condition() string {
	return v.condition
}

func (v *roxString) Parser() roxx.Parser {
	return v.parser
}

func (v *roxString) ImpressionInvoker() model.ImpressionInvoker {
	return v.impressionInvoker
}

func (v *roxString) ClientExperiment() *model.Experiment {
	return v.clientExperiment
}
