package properties_test

import (
	"github.com/rollout/rox-go/core/context"
	"github.com/rollout/rox-go/core/properties"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCustomPropertyWillCreatePropertyWithConstValue(t *testing.T) {
	propString := properties.NewCustomStringProperty("prop1", "123")

	assert.Equal(t, "prop1", propString.Name)
	assert.Equal(t, properties.CustomPropertyTypeString, propString.Type)
	assert.Equal(t, "123", propString.Value(nil))

	propFloat := properties.NewCustomFloatProperty("prop1", 123.12)

	assert.Equal(t, "prop1", propFloat.Name)
	assert.Equal(t, properties.CustomPropertyTypeFloat, propFloat.Type)
	assert.Equal(t, 123.12, propFloat.Value(nil))

	propInt := properties.NewCustomIntegerProperty("prop1", 123)

	assert.Equal(t, "prop1", propInt.Name)
	assert.Equal(t, properties.CustomPropertyTypeInt, propInt.Type)
	assert.Equal(t, 123, propInt.Value(nil))

	propBool := properties.NewCustomBooleanProperty("prop1", true)

	assert.Equal(t, "prop1", propBool.Name)
	assert.Equal(t, properties.CustomPropertyTypeBool, propBool.Type)
	assert.Equal(t, true, propBool.Value(nil))

	propSemver := properties.NewCustomSemverProperty("prop1", "1.2.3")

	assert.Equal(t, "prop1", propSemver.Name)
	assert.Equal(t, properties.CustomPropertyTypeSemver, propSemver.Type)
	assert.Equal(t, "1.2.3", propSemver.Value(nil))
}

func TestCustomPropertyWillCreatePropertyWithFuncValue(t *testing.T) {
	propString := properties.NewCustomComputedStringProperty("prop1", func(context context.Context) string {
		return "123"
	})

	assert.Equal(t, "prop1", propString.Name)
	assert.Equal(t, properties.CustomPropertyTypeString, propString.Type)
	assert.Equal(t, "123", propString.Value(nil))

	propFloat := properties.NewCustomComputedFloatProperty("prop1", func(context context.Context) float64 {
		return 123.12
	})

	assert.Equal(t, "prop1", propFloat.Name)
	assert.Equal(t, properties.CustomPropertyTypeFloat, propFloat.Type)
	assert.Equal(t, 123.12, propFloat.Value(nil))

	propInt := properties.NewCustomComputedIntegerProperty("prop1", func(context context.Context) int {
		return 123
	})

	assert.Equal(t, "prop1", propInt.Name)
	assert.Equal(t, properties.CustomPropertyTypeInt, propInt.Type)
	assert.Equal(t, 123, propInt.Value(nil))

	propBool := properties.NewCustomComputedBooleanProperty("prop1", func(context context.Context) bool {
		return true
	})

	assert.Equal(t, "prop1", propBool.Name)
	assert.Equal(t, properties.CustomPropertyTypeBool, propBool.Type)
	assert.Equal(t, true, propBool.Value(nil))

	propSemver := properties.NewCustomComputedSemverProperty("prop1", func(context context.Context) string {
		return "1.2.3"
	})

	assert.Equal(t, "prop1", propSemver.Name)
	assert.Equal(t, properties.CustomPropertyTypeSemver, propSemver.Type)
	assert.Equal(t, "1.2.3", propSemver.Value(nil))
}

func TestCustomPropertyWillPassContext(t *testing.T) {
	ctx := context.NewContext(map[string]interface{}{"a": 1})
	var contextFromFunc context.Context

	propString := properties.NewCustomComputedStringProperty("prop1", func(context context.Context) string {
		contextFromFunc = context
		return "123"
	})

	assert.Equal(t, "123", propString.Value(ctx))
	assert.Equal(t, 1, contextFromFunc.Get("a"))
}

func TestDevicePropertyWilAddRoxToTheName(t *testing.T) {
	prop := properties.NewDeviceStringProperty("prop1", "123")

	assert.Equal(t, "rox.prop1", prop.Name)
}
