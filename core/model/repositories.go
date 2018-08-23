package model

import "github.com/rollout/rox-go/core/properties"

type CustomPropertyAddedHandler = func(property *properties.CustomProperty)

type CustomPropertyRepository interface {
	AddCustomProperty(customProperty *properties.CustomProperty)
	AddCustomPropertyIfNotExists(customProperty *properties.CustomProperty)
	GetCustomProperty(name string) *properties.CustomProperty
	GetAllCustomProperties() []*properties.CustomProperty
	RegisterCustomPropertyAddedHandler(handler CustomPropertyAddedHandler)
}

type ExperimentRepository interface {
	SetExperiments(experiments []*ExperimentModel)
	GetExperimentByFlag(flagName string) *ExperimentModel
	GetAllExperiments() []*ExperimentModel
}

type FlagAddedHandler = func(variant Variant)

type FlagRepository interface {
	AddFlag(variant Variant, name string)
	GetFlag(name string) Variant
	GetAllFlags() []Variant

	RegisterFlagAddedHandler(handler FlagAddedHandler)
}

type TargetGroupRepository interface {
	SetTargetGroups(targetGroups []*TargetGroupModel)
	GetTargetGroup(id string) *TargetGroupModel
}
