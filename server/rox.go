package server

import (
	"github.com/rollout/rox-go/core"
	"github.com/rollout/rox-go/core/consts"
	"github.com/rollout/rox-go/core/context"
	"github.com/rollout/rox-go/core/logging"
	"github.com/rollout/rox-go/core/model"
	"github.com/rollout/rox-go/core/properties"
	uuid "github.com/satori/go.uuid"
	"sync"
)

type RoxState int

const (
	Idle RoxState = iota
	SettingUp
	Set
	ShuttingDown
	Corrupted
)

type Rox struct {
	core               *core.Core
	state              RoxState
	setupShutdownMutex sync.RWMutex
}

func NewRox() *Rox {
	return &Rox{
		core:  core.NewCore(),
		state: Idle,
	}
}

func (r *Rox) Shutdown() <-chan struct{} {
	done := make(chan struct{})
	go func() {
		r.setupShutdownMutex.Lock()
		defer r.setupShutdownMutex.Unlock()
		defer close(done)
		if r.state != Set && r.state != Corrupted {
			logging.GetLogger().Warn("Rox can only be shutdown when it is already in Set or Corrupted state.", nil)
			return
		} else {
			reset(r)
		}
	}()
	return done
}

func (r *Rox) Setup(apiKey string, roxOptions model.RoxOptions) <-chan struct{} {
	r.setupShutdownMutex.Lock()
	defer r.setupShutdownMutex.Unlock()
	defer func() {
		if r := recover(); r != nil {
			logging.GetLogger().Error("Failed in Rox.Setup", r)
		}
	}()

	if r.state != Idle && r.state != Corrupted {
		logging.GetLogger().Warn("Rox has already been initialised, skipping setup", nil)
		done := make(chan struct{})
		defer close(done)
		return done
	}

	if r.state == Corrupted {
		reset(r)
	}

	r.state = SettingUp

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
	r.core.AddCustomPropertyIfNotExists(properties.NewSemverProperty("internal."+consts.PropertyTypeAPIVersion.Name, serverProperties.GetAllProperties()[consts.PropertyTypeAPIVersion.Name]))
	r.core.AddCustomPropertyIfNotExists(properties.NewSemverProperty("internal."+consts.PropertyTypeLibVersion.Name, serverProperties.GetAllProperties()[consts.PropertyTypeLibVersion.Name]))
	r.core.AddCustomPropertyIfNotExists(properties.NewDeviceComputedStringProperty("internal."+consts.PropertyTypeDistinctID.Name, func(ctx context.Context) string {
		value, _ := uuid.NewV4()
		return value.String()
	}))

	done := make(chan struct{})
	go func() {
		defer close(done)

		defer func() {
			if err := recover(); err != nil {
				logging.GetLogger().Error("Failed in Rox.Setup", err)
				r.state = Corrupted
			}
		}()
		<-r.core.Setup(sdkSettings, serverProperties, roxOptions)
		r.state = Set
	}()
	return done
}

func (r *Rox) RegisterWithEmptyNamespace(roxContainer interface{}) {
	r.Register("", roxContainer)
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
				logging.GetLogger().Error("Failed in Rox.Fetch", r)
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

func (r *Rox) DynamicAPI() model.DynamicAPI {
	return r.core.DynamicAPI(&ServerEntitiesProvider{})
}

func reset(r *Rox) {
	r.state = ShuttingDown
	r.core.Shutdown()
	r.core = core.NewCore()
	r.state = Idle
}
