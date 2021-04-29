package client

import (
	"github.com/rollout/rox-go/core/consts"
	"github.com/rollout/rox-go/core/context"
	"github.com/rollout/rox-go/core/model"
)

type dynamicAPI struct {
	flagRepository   model.FlagRepository
	entitiesProvider model.EntitiesProvider
}

func NewDynamicAPI(flagRepository model.FlagRepository, entitiesProvider model.EntitiesProvider) model.DynamicAPI {
	return &dynamicAPI{
		flagRepository:   flagRepository,
		entitiesProvider: entitiesProvider,
	}
}

func (api *dynamicAPI) IsEnabled(name string, defaultValue bool, ctx context.Context) bool {
	flag := api.flagRepository.GetFlag(name)
	if flag == nil {
		flag = api.entitiesProvider.CreateFlag(defaultValue)
		api.flagRepository.AddFlag(flag, name)
	}

	if flag, ok := flag.(model.Flag); !ok {
		return defaultValue
	} else {
		isEnabled, isDefaultValue := flag.(model.InternalFlag).InternalIsEnabled(ctx)
		if isDefaultValue {
			return defaultValue
		} else {
			return isEnabled
		}
	}
}

func (api *dynamicAPI) StringValue(name string, defaultValue string, options []string, ctx context.Context) string {
	variant := api.flagRepository.GetFlag(name)
	if variant == nil {
		variant = api.entitiesProvider.CreateRoxString(defaultValue, options)
		api.flagRepository.AddFlag(variant, name)
	}

	switch variant.FlagType() {
	case consts.StringType:
		value, isDefaultValue := variant.(model.InternalRoxString).InternalGetValue(ctx)
		if isDefaultValue {
			return defaultValue
		} else {
			return value
		}
	default:
		return defaultValue
	}
	return defaultValue
}

func (api *dynamicAPI) IntValue(name string, defaultValue int, options []int, ctx context.Context) int {
	variant := api.flagRepository.GetFlag(name)
	if variant == nil {
		variant = api.entitiesProvider.CreateRoxInt(defaultValue, options)
		api.flagRepository.AddFlag(variant, name)
	}

	switch variant.FlagType() {
	case consts.IntType:
		value, isDefaultValue := variant.(model.InternalRoxInt).InternalGetValue(ctx)
		if isDefaultValue {
			return defaultValue
		} else {
			return value
		}
	default:
		return defaultValue
	}
	return defaultValue
}

func (api *dynamicAPI) DoubleValue(name string, defaultValue float64, options []float64, ctx context.Context) float64 {
	variant := api.flagRepository.GetFlag(name)
	if variant == nil {
		variant = api.entitiesProvider.CreateRoxDouble(defaultValue, options)
		api.flagRepository.AddFlag(variant, name)
	}

	switch variant.FlagType() {
	case consts.DoubleType:
		value, isDefaultValue := variant.(model.InternalRoxDouble).InternalGetValue(ctx)
		if isDefaultValue {
			return defaultValue
		} else {
			return value
		}
	default:
		return defaultValue
	}
	return defaultValue
}
