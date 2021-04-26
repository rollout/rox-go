package server

import (
	"github.com/rollout/rox-go/core/entities"
	"github.com/rollout/rox-go/core/model"
)

type RoxString = model.RoxString

func NewRoxString(defaultValue string, options []string) RoxString {
	return entities.NewRoxString(defaultValue, options)
}
