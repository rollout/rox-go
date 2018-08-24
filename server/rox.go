package server

import (
	"fmt"
	"github.com/rollout/rox-go/core"
	"github.com/rollout/rox-go/core/consts"
	"github.com/rollout/rox-go/core/context"
	"github.com/rollout/rox-go/core/model"
	"github.com/rollout/rox-go/core/properties"
	"github.com/satori/go.uuid"
)

type Rox struct {
	core *core.Core
}

func NewRox() *Rox {
	return &Rox{
		core: core.NewCore(),
	}
}

func (r *Rox) Setup(apiKey string, roxOptions model.RoxOptions) <-chan struct{} {
	defer func() {
		if r := recover(); r != nil {
			// TODO logger
			fmt.Printf("Failed in Rox.Setup %s\n", r)
		}
	}()

	if roxOptions == nil {
		roxOptions = NewRoxOptions(RoxOptionsBuilder{})
	}

	sdkSettings := NewSdkSettings(apiKey, roxOptions.DevModeKey())
	serverProperties := NewServerProperties(sdkSettings, roxOptions)

	r.core.AddCustomPropertyIfNotExists(properties.NewDeviceStringProperty(consts.PropertyTypePlatform.Name, serverProperties.GetAllProperties()[consts.PropertyTypePlatform.Name]))
	r.core.AddCustomPropertyIfNotExists(properties.NewDeviceSemverProperty(consts.PropertyTypeAppRelease.Name, serverProperties.GetAllProperties()[consts.PropertyTypeAppRelease.Name]))
	r.core.AddCustomPropertyIfNotExists(properties.NewDeviceComputedStringProperty(consts.PropertyTypeDistinctID.Name, func(ctx context.Context) string {
		value, _ := uuid.NewV4()
		return value.String()
	}))
	r.core.AddCustomPropertyIfNotExists(properties.NewDeviceStringProperty("internal.realPlatform", serverProperties.GetAllProperties()[consts.PropertyTypePlatform.Name]))
	r.core.AddCustomPropertyIfNotExists(properties.NewDeviceStringProperty("internal.customPlatform", serverProperties.GetAllProperties()[consts.PropertyTypePlatform.Name]))
	r.core.AddCustomPropertyIfNotExists(properties.NewDeviceStringProperty("internal.appKey", serverProperties.RolloutKey()))
	r.core.AddCustomPropertyIfNotExists(properties.NewDeviceComputedStringProperty("internal."+consts.PropertyTypeDistinctID.Name, func(ctx context.Context) string {
		value, _ := uuid.NewV4()
		return value.String()
	}))

	done := make(chan struct{})
	go func() {
		defer close(done)

		defer func() {
			if r := recover(); r != nil {
				// TODO logger
				fmt.Printf("Failed in Rox.Setup %s\n", r)
			}
		}()

		<-r.core.Setup(sdkSettings, serverProperties, roxOptions)
	}()
	return done
}

func (r *Rox) Register(namespace string, roxContainer interface{}) {
	r.core.Register(namespace, roxContainer)
}

func (r *Rox) SetContext(ctx context.Context) {
	r.core.SetContext(ctx)
}

func (r *Rox) Fetch() <-chan struct{} {
	done := make(chan struct{})

	go func() {
		defer close(done)

		defer func() {
			if r := recover(); r != nil {
				// TODO logger
				fmt.Printf("Failed in Rox.Fetch %s\n", r)
			}
		}()

		<-r.core.Fetch()
	}()

	return done
}

func (r *Rox) SetCustomStringProperty(name string, value string) {
	r.core.AddCustomProperty(properties.NewStringProperty(name, value))
}

func (r *Rox) SetCustomComputedStringProperty(name string, value properties.CustomStringPropertyGenerator) {
	r.core.AddCustomProperty(properties.NewComputedStringProperty(name, value))
}

func (r *Rox) SetCustomBooleanProperty(name string, value bool) {
	r.core.AddCustomProperty(properties.NewBooleanProperty(name, value))
}

func (r *Rox) SetCustomComputedBooleanProperty(name string, value properties.CustomBooleanPropertyGenerator) {
	r.core.AddCustomProperty(properties.NewComputedBooleanProperty(name, value))
}

func (r *Rox) SetCustomIntegerProperty(name string, value int) {
	r.core.AddCustomProperty(properties.NewIntegerProperty(name, value))
}

func (r *Rox) SetCustomComputedIntegerProperty(name string, value properties.CustomIntegerPropertyGenerator) {
	r.core.AddCustomProperty(properties.NewComputedIntegerProperty(name, value))
}

func (r *Rox) SetCustomFloatProperty(name string, value float64) {
	r.core.AddCustomProperty(properties.NewFloatProperty(name, value))
}

func (r *Rox) SetCustomComputedFloatProperty(name string, value properties.CustomFloatPropertyGenerator) {
	r.core.AddCustomProperty(properties.NewComputedFloatProperty(name, value))
}

func (r *Rox) SetCustomSemverProperty(name string, value string) {
	r.core.AddCustomProperty(properties.NewSemverProperty(name, value))
}

func (r *Rox) SetCustomComputedSemverProperty(name string, value properties.CustomSemverPropertyGenerator) {
	r.core.AddCustomProperty(properties.NewComputedSemverProperty(name, value))
}
