package server

import (
	"github.com/rollout/rox-go/core/context"
	"github.com/rollout/rox-go/core/entities"
	"github.com/rollout/rox-go/core/model"
)

type RoxFlag interface {
	model.Variant

	IsEnabled() bool
	Enabled(action func())
	Disabled(action func())

	IsEnabledInContext(ctx context.Context) bool
	EnabledInContext(ctx context.Context, action func())
	DisabledInContext(ctx context.Context, action func())
}

type roxFlag struct {
	model.Flag
}

func NewRoxFlag(defaultValue bool) RoxFlag {
	return &roxFlag{
		Flag: entities.NewFlag(defaultValue),
	}
}

func (f *roxFlag) IsEnabled() bool {
	return f.Flag.IsEnabled(nil)
}

func (f *roxFlag) Enabled(action func()) {
	f.Flag.Enabled(nil, action)
}

func (f *roxFlag) Disabled(action func()) {
	f.Flag.Disabled(nil, action)
}

func (f *roxFlag) IsEnabledInContext(ctx context.Context) bool {
	return f.Flag.IsEnabled(ctx)
}

func (f *roxFlag) EnabledInContext(ctx context.Context, action func()) {
	f.Flag.Enabled(ctx, action)
}

func (f *roxFlag) DisabledInContext(ctx context.Context, action func()) {
	f.Flag.Disabled(ctx, action)
}
