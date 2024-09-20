package server

import (
	"github.com/rollout/rox-go/v6/core/entities"
	"github.com/rollout/rox-go/v6/core/model"
)

type RoxInt = model.RoxInt

func NewRoxInt(defaultValue int, options []int) RoxInt {
	return entities.NewRoxInt(defaultValue, options)
}
