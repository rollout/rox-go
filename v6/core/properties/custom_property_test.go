package properties_test

import (
	"testing"

	"github.com/rollout/rox-go/v6/core/context"
	"github.com/rollout/rox-go/v6/core/properties"
	"github.com/stretchr/testify/assert"
)

func TestCustomPropertyWillCreatePropertyWithConstValue(t *testing.T) {
	propString := properties.NewStringProperty("prop1", "123")

	assert.Equal(t, "prop1", propString.Name)
	assert.Equal(t, properties.CustomPropertyTypeString, propString.Type)
	assert.Equal(t, "123", propString.Value(nil))

	propFloat := properties.NewFloatProperty("prop1", 123.12)

	assert.Equal(t, "prop1", propFloat.Name)
	assert.Equal(t, properties.CustomPropertyTypeFloat, propFloat.Type)
	assert.Equal(t, 123.12, propFloat.Value(nil))

	propInt := properties.NewIntegerProperty("prop1", 123)

	assert.Equal(t, "prop1", propInt.Name)
	assert.Equal(t, properties.CustomPropertyTypeInt, propInt.Type)
	assert.Equal(t, 123, propInt.Value(nil))

	propBool := properties.NewBooleanProperty("prop1", true)

	assert.Equal(t, "prop1", propBool.Name)
	assert.Equal(t, properties.CustomPropertyTypeBool, propBool.Type)
	assert.Equal(t, true, propBool.Value(nil))

	propSemver := properties.NewSemverProperty("prop1", "1.2.3")

	assert.Equal(t, "prop1", propSemver.Name)
	assert.Equal(t, properties.CustomPropertyTypeSemver, propSemver.Type)
	assert.Equal(t, "1.2.3", propSemver.Value(nil))
}

func TestCustomPropertyWillCreatePropertyWithFuncValue(t *testing.T) {
	propString := properties.NewComputedStringProperty("prop1", func(context context.Context) string {
		return "123"
	})

	assert.Equal(t, "prop1", propString.Name)
	assert.Equal(t, properties.CustomPropertyTypeString, propString.Type)
	assert.Equal(t, "123", propString.Value(nil))

	propFloat := properties.NewComputedFloatProperty("prop1", func(context context.Context) float64 {
		return 123.12
	})

	assert.Equal(t, "prop1", propFloat.Name)
	assert.Equal(t, properties.CustomPropertyTypeFloat, propFloat.Type)
	assert.Equal(t, 123.12, propFloat.Value(nil))

	propInt := properties.NewComputedIntegerProperty("prop1", func(context context.Context) int {
		return 123
	})

	assert.Equal(t, "prop1", propInt.Name)
	assert.Equal(t, properties.CustomPropertyTypeInt, propInt.Type)
	assert.Equal(t, 123, propInt.Value(nil))

	propBool := properties.NewComputedBooleanProperty("prop1", func(context context.Context) bool {
		return true
	})

	assert.Equal(t, "prop1", propBool.Name)
	assert.Equal(t, properties.CustomPropertyTypeBool, propBool.Type)
	assert.Equal(t, true, propBool.Value(nil))

	propSemver := properties.NewComputedSemverProperty("prop1", func(context context.Context) string {
		return "1.2.3"
	})

	assert.Equal(t, "prop1", propSemver.Name)
	assert.Equal(t, properties.CustomPropertyTypeSemver, propSemver.Type)
	assert.Equal(t, "1.2.3", propSemver.Value(nil))
}

func TestCustomPropertyWillPassContext(t *testing.T) {
	ctx := context.NewContext(map[string]interface{}{"a": 1})
	var contextFromFunc context.Context

	propString := properties.NewComputedStringProperty("prop1", func(context context.Context) string {
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
