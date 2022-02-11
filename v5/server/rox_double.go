package server

import (
	"github.com/rollout/rox-go/v5/core/entities"
	"github.com/rollout/rox-go/v5/core/model"
)

type RoxDouble = model.RoxDouble

func NewRoxDouble(defaultValue float64, options []float64) RoxDouble {
	return entities.NewRoxDouble(defaultValue, options)
}
