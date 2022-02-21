package server

import (
	"github.com/rollout/rox-go/v4/core/client"
	"github.com/rollout/rox-go/v4/core/model"
	"github.com/satori/go.uuid"
)

type serverProperties struct {
	model.DeviceProperties
	distinctID string
}

func NewServerProperties(sdkSettings model.SdkSettings, roxOptions model.RoxOptions) model.DeviceProperties {
	distinctID, _ := uuid.NewV4()
	return &serverProperties{
		DeviceProperties: client.NewDeviceProperties(sdkSettings, roxOptions),
		distinctID:       distinctID.String(),
	}
}

func (sp *serverProperties) DistinctID() string {
	return sp.distinctID
}

// TODO
//func (sp *serverProperties) LibVersion() string {
//	return
//}
