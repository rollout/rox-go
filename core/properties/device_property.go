package properties

import "github.com/rollout/rox-go/core/context"

func NewDeviceStringProperty(name string, value string) *CustomProperty {
	return NewDeviceComputedStringProperty(name, func(context context.Context) string {
		return value
	})
}

func NewDeviceSemverProperty(name string, value string) *CustomProperty {
	return NewDeviceComputedSemverProperty(name, func(context context.Context) string {
		return value
	})
}

func NewDeviceComputedStringProperty(name string, value CustomStringPropertyGenerator) *CustomProperty {
	return NewCustomComputedStringProperty("rox."+name, value)
}

func NewDeviceComputedSemverProperty(name string, value CustomSemverPropertyGenerator) *CustomProperty {
	return NewCustomComputedStringProperty("rox."+name, value)
}
