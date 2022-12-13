package entities

import (
	"github.com/rollout/rox-go/v5/core/consts"
	"github.com/rollout/rox-go/v5/core/context"
	"github.com/rollout/rox-go/v5/core/model"
	"github.com/rollout/rox-go/v5/core/roxx"
	"github.com/rollout/rox-go/v5/core/utils"
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
			flagType: consts.IntType,
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

func (v *roxInt) SetTag(tag string) {
	v.tag = tag
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
	sendImpression := false

	if v.parser != nil && v.condition != "" {
		evaluationResult := v.parser.EvaluateExpression(v.condition, mergedContext)
		value, err := evaluationResult.IntValue()
		if err == nil {
			returnValue, isDefault = value, false
			sendImpression = true
		}
	}

	if v.impressionInvoker != nil && sendImpression {

		targeting := false
		if v.clientExperiment != nil {
			targeting = true
		}

		v.impressionInvoker.Invoke(model.NewReportingValue(v.name, strconv.Itoa(returnValue), targeting), mergedContext)
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
