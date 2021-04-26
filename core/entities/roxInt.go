package entities

import (
	"github.com/rollout/rox-go/core/context"
	"github.com/rollout/rox-go/core/model"
	"github.com/rollout/rox-go/core/roxx"
	"github.com/rollout/rox-go/core/utils"
	"strconv"
)

type roxInt struct {
	roxVariant
	defaultValue      int
	options           []int
	condition         string
	parser            roxx.Parser
	globalContext     context.Context
	impressionInvoker model.ImpressionInvoker
	clientExperiment  *model.Experiment
}

func NewRoxInt(defaultValue int, options []int) model.RoxInt {
	if options == nil {
		options = []int{}
	}
	allOptions := make([]int, len(options))
	copy(allOptions, options)
	if !utils.ContainsInt(allOptions, defaultValue) {
		allOptions = append(allOptions, defaultValue)
	}

	roxInt := &roxInt{
		roxVariant: roxVariant{
			flagType: "intType",
		},
		defaultValue: defaultValue,
		options:      allOptions,
	}

	return roxInt
}

func (v *roxInt) GetDefaultAsString() string {
	return strconv.Itoa(v.DefaultValue())
}

func (v *roxInt) DefaultValue() int {
	return v.defaultValue
}

func (v *roxInt) GetOptionsAsString() []string {
	options := make([]string, len(v.options))

	for _, option := range v.options {
		options = append(options, strconv.Itoa(option))
	}

	return options
}

func (v *roxInt) Options() []int {
	return v.options
}

func (v *roxInt) SetForEvaluation(parser roxx.Parser, experiment *model.ExperimentModel, impressionInvoker model.ImpressionInvoker) {
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

func (v *roxInt) SetContext(globalContext context.Context) {
	v.globalContext = globalContext
}

func (v *roxInt) SetName(name string) {
	v.name = name
}

func (v *roxInt) GetValueAsString(ctx context.Context) string {
	return strconv.Itoa(v.GetValue(ctx))
}

func (v *roxInt) GetValue(ctx context.Context) int {
	returnValue, _ := v.InternalGetValue(ctx)
	return returnValue
}

func (v *roxInt) InternalGetValue(ctx context.Context) (returnValue int, isDefault bool) {
	returnValue, isDefault = v.defaultValue, true
	mergedContext := context.NewMergedContext(v.globalContext, ctx)

	if v.parser != nil && v.condition != "" {
		evaluationResult := v.parser.EvaluateExpression(v.condition, mergedContext)
		value := evaluationResult.IntValue()
		if value != 0 {
			returnValue, isDefault = value, false
		}
	}

	if v.impressionInvoker != nil {
		v.impressionInvoker.Invoke(model.NewReportingValue(v.name, strconv.Itoa(returnValue)), v.clientExperiment, mergedContext)
	}

	return returnValue, isDefault
}

func (v *roxInt) Condition() string {
	return v.condition
}

func (v *roxInt) Parser() roxx.Parser {
	return v.parser
}

func (v *roxInt) ImpressionInvoker() model.ImpressionInvoker {
	return v.impressionInvoker
}

func (v *roxInt) ClientExperiment() *model.Experiment {
	return v.clientExperiment
}
