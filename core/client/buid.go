package client

import (
	"github.com/rollout/rox-go/core/consts"
	"github.com/rollout/rox-go/core/model"
	"github.com/rollout/rox-go/core/utils"
)

var buidGenerators = []consts.PropertyType{
	*consts.PropertyTypePlatform,
	*consts.PropertyTypeAppKey,
	*consts.PropertyTypeLibVersion,
	*consts.PropertyTypeAPIVersion,
}

type buid struct {
	sdkSettings              model.SdkSettings
	deviceProperties         model.DeviceProperties
	flagRepository           model.FlagRepository
	customPropertyRepository model.CustomPropertyRepository
	buid                     string
}

func NewBUID(sdkSettings model.SdkSettings, deviceProperties model.DeviceProperties, flagRepository model.FlagRepository, customPropertyRepository model.CustomPropertyRepository) model.BUID {
	return &buid{
		sdkSettings:              sdkSettings,
		deviceProperties:         deviceProperties,
		flagRepository:           flagRepository,
		customPropertyRepository: customPropertyRepository,
	}
}

func (b *buid) GetQueryStringParts() map[string]string {
	return map[string]string{
		consts.PropertyTypeBuid.Name: b.GetValue(),
	}
}

func (b *buid) GetValue() string {
	properties := b.deviceProperties.GetAllProperties()
	buid := utils.GenerateMD5(properties, buidGenerators, nil)
	return buid
}

func (b *buid) String() string {
	return b.buid
}
