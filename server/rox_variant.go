package server

import (
	"github.com/rollout/rox-go/core/entities"
	"github.com/rollout/rox-go/core/model"
)

type RoxVariant = model.Variant

func NewRoxVariant(defaultValue string, options []string) RoxVariant {
	return entities.NewVariant(defaultValue, options)
}
