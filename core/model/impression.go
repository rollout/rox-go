package model

import "github.com/rollout/rox-go/core/context"

type ImpressionError int

const (
	ImpressionErrorNoError ImpressionError = iota
	ImpressionErrorCausedByHandler
)


type ImpressionArgs struct {
	ReportingValue *ReportingValue
	Experiment     *Experiment
	ErrorDetails ImpressionError
	Context        context.Context
}

type ImpressionHandler = func(args ImpressionArgs)

type ImpressionInvoker interface {
	Invoke(value *ReportingValue, experiment *Experiment, context context.Context)
	RegisterImpressionHandler(handler ImpressionHandler)
}

type ExperimentModel struct {
	ID         string
	Name       string
	Condition  string
	IsArchived bool
	Flags      []string
	Labels     []string
}

func NewExperimentModel(id, name, condition string, isArchived bool, flags, labels []string) *ExperimentModel {
	return &ExperimentModel{
		ID:         id,
		Name:       name,
		Condition:  condition,
		IsArchived: isArchived,
		Flags:      flags,
		Labels:     labels,
	}
}

type ReportingValue struct {
	Name      string
	Value     string
	Targeting bool
}

func NewReportingValue(name, value string) *ReportingValue {
	return &ReportingValue{
		Name:      name,
		Value:     value,
	}
}
