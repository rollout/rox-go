package server

import (
	"github.com/rollout/rox-go/core/client"
	"github.com/rollout/rox-go/core/model"
	"github.com/satori/go.uuid"
)

type serverProperties struct {
	model.DeviceProperties
	distinctId string
}

func NewServerProperties(sdkSettings model.SdkSettings, roxOptions model.RoxOptions) model.DeviceProperties {
	return &serverProperties{
		DeviceProperties: client.NewDeviceProperties(sdkSettings, roxOptions),
		distinctId:       uuid.NewV4().String(),
	}
}

func (sp *serverProperties) DistinctId() string {
	return sp.distinctId
}

// TODO
//func (sp *serverProperties) LibVersion() string {
//	return
//}
