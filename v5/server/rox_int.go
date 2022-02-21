package server

import (
	"github.com/rollout/rox-go/v5/core/entities"
	"github.com/rollout/rox-go/v5/core/model"
)

type RoxInt = model.RoxInt

func NewRoxInt(defaultValue int, options []int) RoxInt {
	return entities.NewRoxInt(defaultValue, options)
}
