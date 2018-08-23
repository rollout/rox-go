package server

import (
	"fmt"
	"github.com/rollout/rox-go/core"
	"github.com/rollout/rox-go/core/consts"
	"github.com/rollout/rox-go/core/context"
	"github.com/rollout/rox-go/core/custom-properties"
	"github.com/rollout/rox-go/core/model"
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
	r.core.AddCustomPropertyIfNotExists(properties.NewDeviceProperty(consts.PropertyTypeDistinctId.Name, properties.CustomPropertyTypeString, func(ctx context.Context) interface{} {
		return uuid.NewV4().String()
	}))
	r.core.AddCustomPropertyIfNotExists(properties.NewDeviceStringProperty("internal.realPlatform", serverProperties.GetAllProperties()[consts.PropertyTypePlatform.Name]))
	r.core.AddCustomPropertyIfNotExists(properties.NewDeviceStringProperty("internal.customPlatform", serverProperties.GetAllProperties()[consts.PropertyTypePlatform.Name]))
	r.core.AddCustomPropertyIfNotExists(properties.NewDeviceStringProperty("internal.appKey", serverProperties.RolloutKey()))
	r.core.AddCustomPropertyIfNotExists(properties.NewDeviceProperty("internal."+consts.PropertyTypeDistinctId.Name, properties.CustomPropertyTypeString, func(ctx context.Context) interface{} {
		return uuid.NewV4().String()
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
	r.core.AddCustomProperty(properties.NewCustomStringProperty(name, value))
}

func (r *Rox) SetCustomComputedStringProperty(name string, value properties.CustomPropertyGenerator) {
	r.core.AddCustomProperty(properties.NewCustomProperty(name, properties.CustomPropertyTypeString, value))
}

func (r *Rox) SetCustomBooleanProperty(name string, value bool) {
	r.core.AddCustomProperty(properties.NewCustomBooleanProperty(name, value))
}

func (r *Rox) SetCustomComputedBooleanProperty(name string, value properties.CustomPropertyGenerator) {
	r.core.AddCustomProperty(properties.NewCustomProperty(name, properties.CustomPropertyTypeBool, value))
}

func (r *Rox) SetCustomIntegerProperty(name string, value int) {
	r.core.AddCustomProperty(properties.NewCustomIntegerProperty(name, value))
}

func (r *Rox) SetCustomComputedIntegerProperty(name string, value properties.CustomPropertyGenerator) {
	r.core.AddCustomProperty(properties.NewCustomProperty(name, properties.CustomPropertyTypeInt, value))
}

func (r *Rox) SetCustomFloatProperty(name string, value float64) {
	r.core.AddCustomProperty(properties.NewCustomFloatProperty(name, value))
}

func (r *Rox) SetCustomComputedFloatProperty(name string, value properties.CustomPropertyGenerator) {
	r.core.AddCustomProperty(properties.NewCustomProperty(name, properties.CustomPropertyTypeFloat, value))
}

func (r *Rox) SetCustomSemverProperty(name string, value string) {
	r.core.AddCustomProperty(properties.NewCustomSemverProperty(name, value))
}

func (r *Rox) SetCustomComputedSemverProperty(name string, value properties.CustomPropertyGenerator) {
	r.core.AddCustomProperty(properties.NewCustomProperty(name, properties.CustomPropertyTypeSemver, value))
}
