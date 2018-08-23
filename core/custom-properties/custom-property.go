package properties

import "github.com/rollout/rox-go/core/context"

type CustomPropertyGenerator = func(context context.Context) interface{}

type CustomProperty struct {
	Name  string
	Type  *CustomPropertyType
	Value CustomPropertyGenerator
}

func NewCustomStringProperty(name string, value string) *CustomProperty {
	return NewCustomProperty(name, CustomPropertyTypeString, func(context context.Context) interface{} {
		return value
	})
}

func NewCustomIntegerProperty(name string, value int) *CustomProperty {
	return NewCustomProperty(name, CustomPropertyTypeInt, func(context context.Context) interface{} {
		return value
	})
}

func NewCustomFloatProperty(name string, value float64) *CustomProperty {
	return NewCustomProperty(name, CustomPropertyTypeFloat, func(context context.Context) interface{} {
		return value
	})
}

func NewCustomBooleanProperty(name string, value bool) *CustomProperty {
	return NewCustomProperty(name, CustomPropertyTypeBool, func(context context.Context) interface{} {
		return value
	})
}

func NewCustomSemverProperty(name string, value string) *CustomProperty {
	return NewCustomProperty(name, CustomPropertyTypeSemver, func(context context.Context) interface{} {
		return value
	})
}

func NewCustomProperty(name string, propertyType *CustomPropertyType, value CustomPropertyGenerator) *CustomProperty {
	return &CustomProperty{
		Name:  name,
		Type:  propertyType,
		Value: value,
	}
}
