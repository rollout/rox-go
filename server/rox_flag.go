package server

import (
	"github.com/rollout/rox-go/core/context"
	"github.com/rollout/rox-go/core/entities"
	"github.com/rollout/rox-go/core/model"
)

type RoxFlag interface {
	model.Flag
}

type roxFlag struct {
	model.Flag
}

func NewRoxFlag(defaultValue bool) RoxFlag {
	return &roxFlag{
		Flag: entities.NewFlag(defaultValue),
	}
}

func (f *roxFlag) InternalIsEnabled(ctx context.Context) (isEnabled bool, isDefault bool) {
	return f.Flag.(model.InternalFlag).InternalIsEnabled(ctx)
}
