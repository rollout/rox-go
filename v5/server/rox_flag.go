package server

import (
	"github.com/rollout/rox-go/v5/core/entities"
	"github.com/rollout/rox-go/v5/core/model"
)

type RoxFlag = model.Flag

func NewRoxFlag(defaultValue bool) RoxFlag {
	return entities.NewFlag(defaultValue)
}
