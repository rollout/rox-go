package server

import "github.com/rollout/rox-go/core/model"

type ServerEntitiesProvider struct {
}

func (*ServerEntitiesProvider) CreateFlag(defaultValue bool) model.Flag {
	return NewRoxFlag(defaultValue)
}

func (*ServerEntitiesProvider) CreateVariant(defaultValue string, options []string) model.Variant {
	return NewRoxVariant(defaultValue, options)
}
