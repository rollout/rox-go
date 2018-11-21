package server

import (
	"github.com/rollout/rox-go/core/context"
	"github.com/rollout/rox-go/core/entities"
	"github.com/rollout/rox-go/core/model"
)

type RoxVariant interface {
	model.Variant
}

type roxVariant struct {
	model.Variant
}

func NewRoxVariant(defaultValue string, options []string) RoxVariant {
	return &roxVariant{
		Variant: entities.NewVariant(defaultValue, options),
	}
}

func (v *roxVariant) InternalGetValue(ctx context.Context) (returnValue string, isDefault bool) {
	return v.Variant.(model.InternalVariant).InternalGetValue(ctx)
}
