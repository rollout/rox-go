package model

import (
	"github.com/rollout/rox-go/v4/core/context"
	"github.com/rollout/rox-go/v4/core/roxx"
)

type Variant interface {
	DefaultValue() string
	Options() []string
	Name() string
	GetValue(context context.Context) string
}

type Flag interface {
	Variant
	IsEnabled(ctx context.Context) bool
	Enabled(ctx context.Context, action func())
	Disabled(ctx context.Context, action func())
}

type EntitiesProvider interface {
	CreateFlag(defaultValue bool) Flag
	CreateVariant(defaultValue string, options []string) Variant
}

type InternalVariant interface {
	SetName(name string)
	SetContext(globalContext context.Context)
	SetForEvaluation(parser roxx.Parser, experiment *ExperimentModel, impressionInvoker ImpressionInvoker)
	InternalGetValue(ctx context.Context) (returnValue string, isDefault bool)
}

type InternalFlag interface {
	InternalIsEnabled(ctx context.Context) (isEnabled bool, isDefault bool)
}
