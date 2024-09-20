package model

import (
	"github.com/rollout/rox-go/v6/core/context"
	"github.com/rollout/rox-go/v6/core/roxx"
)

type Variant interface {
	Name() string
	FlagType() int
	GetValueAsString(context context.Context) string
	GetDefaultAsString() string
	GetOptionsAsString() []string
}

type RoxString interface {
	Variant
	DefaultValue() string
	Options() []string
	GetValue(context context.Context) string
}

type RoxInt interface {
	Variant
	DefaultValue() int
	Options() []int
	GetValue(context context.Context) int
}

type RoxDouble interface {
	Variant
	DefaultValue() float64
	Options() []float64
	GetValue(context context.Context) float64
}

type Flag interface {
	RoxString
	IsEnabled(ctx context.Context) bool
	Enabled(ctx context.Context, action func())
	Disabled(ctx context.Context, action func())
}

type EntitiesProvider interface {
	CreateFlag(defaultValue bool) Flag
	CreateRoxString(defaultValue string, options []string) RoxString
	CreateRoxInt(defaultValue int, options []int) RoxInt
	CreateRoxDouble(defaultValue float64, options []float64) RoxDouble
}

type InternalVariant interface {
	SetName(name string)
	SetContext(globalContext context.Context)
	SetForEvaluation(parser roxx.Parser, experiment *ExperimentModel, impressionInvoker ImpressionInvoker)
}

type InternalRoxString interface {
	InternalGetValue(ctx context.Context) (returnValue string, isDefault bool)
}

type InternalRoxInt interface {
	InternalGetValue(ctx context.Context) (returnValue int, isDefault bool)
}

type InternalRoxDouble interface {
	InternalGetValue(ctx context.Context) (returnValue float64, isDefault bool)
}

type InternalFlag interface {
	InternalIsEnabled(ctx context.Context) (isEnabled bool, isDefault bool)
}
