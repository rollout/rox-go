package client

import (
	"os"

	"github.com/rollout/rox-go/v5/core/consts"
	"github.com/rollout/rox-go/v5/core/model"
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
		consts.PropertyTypeLibVersion.Name:    dp.LibVersion(),
		consts.PropertyTypeRolloutBuild.Name:  "50", // TODO: fix the build number
		consts.PropertyTypeAPIVersion.Name:    consts.BuildAPIVersion,
		consts.PropertyTypeAppRelease.Name:    dp.roxOptions.Version(), // used for the version filter
		consts.PropertyTypeDistinctID.Name:    dp.DistinctID(),
		consts.PropertyTypeAppKey.Name:        dp.sdkSettings.APIKey(),
		consts.PropertyTypePlatform.Name:      consts.BuildPlatform,
		consts.PropertyTypeDevModeSecret.Name: dp.roxOptions.DevModeKey(),
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
	return "5.0.6"
}

func (dp *deviceProperties) RolloutKey() string {
	return dp.GetAllProperties()[consts.PropertyTypeAppKey.Name]
}

func (*deviceProperties) DistinctID() string {
	return "stam"
}
