package model

import (
	"github.com/rollout/rox-go/core/context"
	"github.com/rollout/rox-go/core/roxx"
)

type Variant interface {
	DefaultValue() string
	Options() []string
	Name() string
	SetName(name string)
	SetContext(globalContext context.Context)
	GetValue(context context.Context) string
	SetForEvaluation(parser roxx.Parser, experiment *ExperimentModel, impressionInvoker ImpressionInvoker)
}

type Flag interface {
	Variant
	IsEnabled(ctx context.Context) bool
	Enabled(ctx context.Context, action func())
	Disabled(ctx context.Context, action func())
}
