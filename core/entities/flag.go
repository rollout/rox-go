package entities

import (
	"github.com/rollout/rox-go/core/context"
	"github.com/rollout/rox-go/core/model"
	"github.com/rollout/rox-go/core/roxx"
)

type flag struct {
	*variant
}

func NewFlag(defaultValue bool) model.Flag {
	var variantDefaultValue string
	if defaultValue {
		variantDefaultValue = roxx.FlagTrueValue
	} else {
		variantDefaultValue = roxx.FlagFalseValue
	}
	return &flag{
		variant: NewVariant(variantDefaultValue, []string{roxx.FlagFalseValue, roxx.FlagTrueValue}).(*variant),
	}
}

func (f *flag) IsEnabled(ctx context.Context) bool {
	isEnabled, _ := f.InternalIsEnabled(ctx)
	return isEnabled
}

func (f *flag) InternalIsEnabled(ctx context.Context) (isEnabled bool, isDefault bool) {
	value, isDefault := f.InternalGetValue(ctx)
	return value == roxx.FlagTrueValue, isDefault
}

func (f *flag) Enabled(ctx context.Context, action func()) {
	if f.IsEnabled(ctx) {
		action()
	}
}

func (f *flag) Disabled(ctx context.Context, action func()) {
	if !f.IsEnabled(ctx) {
		action()
	}
}
