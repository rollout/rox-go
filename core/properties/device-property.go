package properties

import "github.com/rollout/rox-go/core/context"

func NewDeviceStringProperty(name string, value string) *CustomProperty {
	return NewDeviceProperty("rox."+name, CustomPropertyTypeString, func(context context.Context) interface{} {
		return value
	})
}

func NewDeviceSemverProperty(name string, value string) *CustomProperty {
	return NewDeviceProperty("rox."+name, CustomPropertyTypeSemver, func(context context.Context) interface{} {
		return value
	})
}

func NewDeviceProperty(name string, propertyType *CustomPropertyType, value CustomPropertyGenerator) *CustomProperty {
	return NewCustomProperty("rox."+name, propertyType, value)
}
