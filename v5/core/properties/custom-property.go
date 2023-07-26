package properties

import (
	"time"

	"github.com/rollout/rox-go/v5/core/context"
)

type CustomPropertyGenerator = func(context context.Context) interface{}
type CustomStringPropertyGenerator = func(context context.Context) string
type CustomIntegerPropertyGenerator = func(context context.Context) int
type CustomFloatPropertyGenerator = func(context context.Context) float64
type CustomTimePropertyGenerator = func(context context.Context) time.Time
type CustomBooleanPropertyGenerator = func(context context.Context) bool
type CustomSemverPropertyGenerator = func(context context.Context) string

type CustomProperty struct {
	Name  string
	Type  *CustomPropertyType
	Value CustomPropertyGenerator
}

func NewStringProperty(name string, value string) *CustomProperty {
	return NewComputedStringProperty(name, func(context context.Context) string {
		return value
	})
}

func NewIntegerProperty(name string, value int) *CustomProperty {
	return NewComputedIntegerProperty(name, func(context context.Context) int {
		return value
	})
}

func NewFloatProperty(name string, value float64) *CustomProperty {
	return NewComputedFloatProperty(name, func(context context.Context) float64 {
		return value
	})
}

func NewTimeProperty(name string, value time.Time) *CustomProperty {
	return NewComputedTimeProperty(name, func(context context.Context) time.Time {
		return value
	})
}

func NewBooleanProperty(name string, value bool) *CustomProperty {
	return NewComputedBooleanProperty(name, func(context context.Context) bool {
		return value
	})
}

func NewSemverProperty(name string, value string) *CustomProperty {
	return NewComputedSemverProperty(name, func(context context.Context) string {
		return value
	})
}

func NewComputedStringProperty(name string, value CustomStringPropertyGenerator) *CustomProperty {
	return &CustomProperty{
		Name: name,
		Type: CustomPropertyTypeString,
		Value: func(context context.Context) interface{} {
			return value(context)
		},
	}
}

func NewComputedIntegerProperty(name string, value CustomIntegerPropertyGenerator) *CustomProperty {
	return &CustomProperty{
		Name: name,
		Type: CustomPropertyTypeInt,
		Value: func(context context.Context) interface{} {
			return value(context)
		},
	}
}

func NewComputedFloatProperty(name string, value CustomFloatPropertyGenerator) *CustomProperty {
	return &CustomProperty{
		Name: name,
		Type: CustomPropertyTypeFloat,
		Value: func(context context.Context) interface{} {
			return value(context)
		},
	}
}

func NewComputedTimeProperty(name string, value CustomTimePropertyGenerator) *CustomProperty {
	return &CustomProperty{
		Name: name,
		Type: CustomPropertyTypeTime,
		Value: func(context context.Context) interface{} {
			return value(context)
		},
	}
}

func NewComputedBooleanProperty(name string, value CustomBooleanPropertyGenerator) *CustomProperty {
	return &CustomProperty{
		Name: name,
		Type: CustomPropertyTypeBool,
		Value: func(context context.Context) interface{} {
			return value(context)
		},
	}
}

func NewComputedSemverProperty(name string, value CustomSemverPropertyGenerator) *CustomProperty {
	return &CustomProperty{
		Name: name,
		Type: CustomPropertyTypeSemver,
		Value: func(context context.Context) interface{} {
			return value(context)
		},
	}
}
