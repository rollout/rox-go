package entities

import (
	"github.com/rollout/rox-go/v5/core/consts"
	"github.com/rollout/rox-go/v5/core/context"
	"github.com/rollout/rox-go/v5/core/model"
	"github.com/rollout/rox-go/v5/core/roxx"
	"github.com/rollout/rox-go/v5/core/utils"
	"strconv"
)

type roxDouble struct {
	roxVariant
	defaultValue      float64
	options           []float64
	condition         string
	parser            roxx.Parser
	globalContext     context.Context
	impressionInvoker model.ImpressionInvoker
	clientExperiment  *model.Experiment
}

func NewRoxDouble(defaultValue float64, options []float64) model.RoxDouble {
	if options == nil {
		options = []float64{}
	}
	allOptions := make([]float64, len(options))
	copy(allOptions, options)
	if !utils.ContainsDouble(allOptions, defaultValue) {
		allOptions = append(allOptions, defaultValue)
	}

	roxDouble := &roxDouble{
		roxVariant: roxVariant{
			flagType: consts.DoubleType,
		},
		defaultValue: defaultValue,
		options:      allOptions,
	}

	return roxDouble
}

func (v *roxDouble) GetDefaultAsString() string {
	return strconv.FormatFloat(v.DefaultValue(), 'f', -1, 64)
}

func (v *roxDouble) DefaultValue() float64 {
	return v.defaultValue
}

func (v *roxDouble) GetOptionsAsString() []string {
	options := make([]string, len(v.options))

	for _, option := range v.options {
		options = append(options, strconv.FormatFloat(option, 'f', -1, 64))
	}

	return options
}

func (v *roxDouble) Options() []float64 {
	return v.options
}

func (v *roxDouble) SetForEvaluation(parser roxx.Parser, experiment *model.ExperimentModel, impressionInvoker model.ImpressionInvoker) {
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

func (v *roxDouble) SetContext(globalContext context.Context) {
	v.globalContext = globalContext
}

func (v *roxDouble) SetName(name string) {
	v.name = name
}

func (v *roxDouble) SetTag(tag string) {
	v.tag = tag
}

func (v *roxDouble) GetValueAsString(ctx context.Context) string {
	return strconv.FormatFloat(v.GetValue(ctx), 'f', -1, 64)
}

func (v *roxDouble) GetValue(ctx context.Context) float64 {
	returnValue, _ := v.InternalGetValue(ctx)
	return returnValue
}

func (v *roxDouble) InternalGetValue(ctx context.Context) (returnValue float64, isDefault bool) {
	returnValue, isDefault = v.defaultValue, true
	mergedContext := context.NewMergedContext(v.globalContext, ctx)
	sendImpression := false

	if v.parser != nil && v.condition != "" {
		evaluationResult := v.parser.EvaluateExpression(v.condition, mergedContext)
		value, err := evaluationResult.DoubleValue()
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
		v.impressionInvoker.Invoke(model.NewReportingValue(v.name, strconv.FormatFloat(returnValue, 'f', -1, 64), targeting), mergedContext)
	}

	return returnValue, isDefault
}

func (v *roxDouble) Condition() string {
	return v.condition
}

func (v *roxDouble) Parser() roxx.Parser {
	return v.parser
}

func (v *roxDouble) ImpressionInvoker() model.ImpressionInvoker {
	return v.impressionInvoker
}

func (v *roxDouble) ClientExperiment() *model.Experiment {
	return v.clientExperiment
}
