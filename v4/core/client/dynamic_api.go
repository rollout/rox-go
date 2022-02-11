package client

import (
	"github.com/rollout/rox-go/v4/core/context"
	"github.com/rollout/rox-go/v4/core/model"
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
	variant := api.flagRepository.GetFlag(name)
	if variant == nil {
		variant = api.entitiesProvider.CreateFlag(defaultValue)
		api.flagRepository.AddFlag(variant, name)
	}

	if flag, ok := variant.(model.Flag); !ok {
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

func (api *dynamicAPI) Value(name string, defaultValue string, options []string, ctx context.Context) string {
	variant := api.flagRepository.GetFlag(name)
	if variant == nil {
		variant = api.entitiesProvider.CreateVariant(defaultValue, options)
		api.flagRepository.AddFlag(variant, name)
	}

	value, isDefaultValue := variant.(model.InternalVariant).InternalGetValue(ctx)
	if isDefaultValue {
		return defaultValue
	} else {
		return value
	}
}
