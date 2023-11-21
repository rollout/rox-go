package server

import (
	"github.com/google/uuid"
	"github.com/rollout/rox-go/v4/core/client"
	"github.com/rollout/rox-go/v4/core/model"
)

type serverProperties struct {
	model.DeviceProperties
	distinctID string
}

func NewServerProperties(sdkSettings model.SdkSettings, roxOptions model.RoxOptions) model.DeviceProperties {
	distinctID := uuid.New()
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
