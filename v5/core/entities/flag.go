package entities

import (
	"github.com/rollout/rox-go/v5/core/consts"
	"github.com/rollout/rox-go/v5/core/context"
	"github.com/rollout/rox-go/v5/core/model"
	"github.com/rollout/rox-go/v5/core/roxx"
)

type flag struct {
	*roxString
}

func NewFlag(defaultValue bool) model.Flag {
	var variantDefaultValue string
	if defaultValue {
		variantDefaultValue = roxx.FlagTrueValue
	} else {
		variantDefaultValue = roxx.FlagFalseValue
	}
	roxString := NewRoxString(variantDefaultValue, []string{roxx.FlagFalseValue, roxx.FlagTrueValue}).(*roxString)
	roxString.flagType = consts.BoolType
	return &flag{
		roxString,
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
