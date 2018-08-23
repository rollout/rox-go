package properties

import "github.com/rollout/rox-go/core/context"

type CustomPropertyGenerator = func(context context.Context) interface{}
type CustomStringPropertyGenerator = func(context context.Context) string
type CustomIntegerPropertyGenerator = func(context context.Context) int
type CustomFloatPropertyGenerator = func(context context.Context) float64
type CustomBooleanPropertyGenerator = func(context context.Context) bool
type CustomSemverPropertyGenerator = func(context context.Context) string

type CustomProperty struct {
	Name  string
	Type  *CustomPropertyType
	Value CustomPropertyGenerator
}

func NewCustomStringProperty(name string, value string) *CustomProperty {
	return NewCustomComputedStringProperty(name, func(context context.Context) string {
		return value
	})
}

func NewCustomIntegerProperty(name string, value int) *CustomProperty {
	return NewCustomComputedIntegerProperty(name, func(context context.Context) int {
		return value
	})
}

func NewCustomFloatProperty(name string, value float64) *CustomProperty {
	return NewCustomComputedFloatProperty(name, func(context context.Context) float64 {
		return value
	})
}

func NewCustomBooleanProperty(name string, value bool) *CustomProperty {
	return NewCustomComputedBooleanProperty(name, func(context context.Context) bool {
		return value
	})
}

func NewCustomSemverProperty(name string, value string) *CustomProperty {
	return NewCustomComputedSemverProperty(name, func(context context.Context) string {
		return value
	})
}

func NewCustomComputedStringProperty(name string, value CustomStringPropertyGenerator) *CustomProperty {
	return &CustomProperty{
		Name: name,
		Type: CustomPropertyTypeString,
		Value: func(context context.Context) interface{} {
			return value(context)
		},
	}
}

func NewCustomComputedIntegerProperty(name string, value CustomIntegerPropertyGenerator) *CustomProperty {
	return &CustomProperty{
		Name: name,
		Type: CustomPropertyTypeInt,
		Value: func(context context.Context) interface{} {
			return value(context)
		},
	}
}

func NewCustomComputedFloatProperty(name string, value CustomFloatPropertyGenerator) *CustomProperty {
	return &CustomProperty{
		Name: name,
		Type: CustomPropertyTypeFloat,
		Value: func(context context.Context) interface{} {
			return value(context)
		},
	}
}

func NewCustomComputedBooleanProperty(name string, value CustomBooleanPropertyGenerator) *CustomProperty {
	return &CustomProperty{
		Name: name,
		Type: CustomPropertyTypeBool,
		Value: func(context context.Context) interface{} {
			return value(context)
		},
	}
}

func NewCustomComputedSemverProperty(name string, value CustomSemverPropertyGenerator) *CustomProperty {
	return &CustomProperty{
		Name: name,
		Type: CustomPropertyTypeSemver,
		Value: func(context context.Context) interface{} {
			return value(context)
		},
	}
}
