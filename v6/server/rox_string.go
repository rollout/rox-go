package server

import (
	"github.com/rollout/rox-go/v6/core/entities"
	"github.com/rollout/rox-go/v6/core/model"
)

type RoxString = model.RoxString

func NewRoxString(defaultValue string, options []string) RoxString {
	return entities.NewRoxString(defaultValue, options)
}
