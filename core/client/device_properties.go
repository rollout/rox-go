package client

import (
	"github.com/rollout/rox-go/core/consts"
	"github.com/rollout/rox-go/core/model"
	"os"
)

type deviceProperties struct {
	sdkSettings model.SdkSettings
	roxOptions  model.RoxOptions
}

func NewDeviceProperties(sdkSettings model.SdkSettings, roxOptions model.RoxOptions) model.DeviceProperties {
	return &deviceProperties{
		sdkSettings: sdkSettings,
		roxOptions:  roxOptions,
	}
}

func (dp *deviceProperties) GetAllProperties() map[string]string {
	return map[string]string{
		consts.PropertyTypePackageName.Name:  dp.roxOptions.Version(),
		consts.PropertyTypeVersionName.Name:  dp.roxOptions.Version(),
		consts.PropertyTypeLibVersion.Name:   dp.LibVersion(),
		consts.PropertyTypeRolloutBuild.Name: "50",
		consts.PropertyTypeApiVersion.Name:   consts.BuildApiVersion,
		consts.PropertyTypeAppVersion.Name:   dp.roxOptions.Version(),
		consts.PropertyTypeAppRelease.Name:   dp.roxOptions.Version(),
		consts.PropertyTypeDistinctId.Name:   dp.DistinctId(),
		consts.PropertyTypeAppKey.Name:       dp.sdkSettings.ApiKey(),
		consts.PropertyTypePlatform.Name:     consts.BuildPlatform,
	}
}

func (*deviceProperties) RolloutEnvironment() string {
	rolloutMode := os.Getenv("ROLLOUT_MODE")
	if rolloutMode != "QA" && rolloutMode != "LOCAL" {
		return "PRODUCTION"
	}
	return rolloutMode
}

func (*deviceProperties) LibVersion() string {
	return "1.0.0"
}

func (dp *deviceProperties) RolloutKey() string {
	return dp.GetAllProperties()[consts.PropertyTypeAppKey.Name]
}

func (*deviceProperties) DistinctId() string {
	return "stam"
}
