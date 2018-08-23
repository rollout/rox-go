package server

import (
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
