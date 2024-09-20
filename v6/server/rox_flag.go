package server

import (
	"github.com/rollout/rox-go/v6/core/entities"
	"github.com/rollout/rox-go/v6/core/model"
)

type RoxFlag = model.Flag

func NewRoxFlag(defaultValue bool) RoxFlag {
	return entities.NewFlag(defaultValue)
}
