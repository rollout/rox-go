package server

import (
	"github.com/rollout/rox-go/v6/core/entities"
	"github.com/rollout/rox-go/v6/core/model"
)

type RoxDouble = model.RoxDouble

func NewRoxDouble(defaultValue float64, options []float64) RoxDouble {
	return entities.NewRoxDouble(defaultValue, options)
}
