package server

import "github.com/rollout/rox-go/core/model"

type ServerEntitiesProvider struct {
}

func (*ServerEntitiesProvider) CreateFlag(defaultValue bool) model.Flag {
	return NewRoxFlag(defaultValue)
}

func (*ServerEntitiesProvider) CreateRoxString(defaultValue string, options []string) model.RoxString {
	return NewRoxString(defaultValue, options)
}

func (*ServerEntitiesProvider) CreateRoxInt(defaultValue int, options []int) model.RoxInt {
	return NewRoxInt(defaultValue, options)
}

func (*ServerEntitiesProvider) CreateRoxDouble(defaultValue float64, options []float64) model.RoxDouble {
	return NewRoxDouble(defaultValue, options)
}
