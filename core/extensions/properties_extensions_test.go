package extensions_test

import (
	"github.com/rollout/rox-go/core/context"
	"github.com/rollout/rox-go/core/extensions"
	"github.com/rollout/rox-go/core/properties"
	"github.com/rollout/rox-go/core/repositories"
	"github.com/rollout/rox-go/core/roxx"
	"github.com/rollout/rox-go/server"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPropertiesExtensionsRoxxPropertiesExtensionsString(t *testing.T) {
	customPropertiesRepository := repositories.NewCustomPropertyRepository()
	parser := roxx.NewParser()
	extensions.NewPropertiesExtensions(parser, customPropertiesRepository, nil).Extend()

	customPropertiesRepository.AddCustomProperty(properties.NewStringProperty("testKey", "test"))

	assert.Equal(t, true, parser.EvaluateExpression(`eq("test", property("testKey"))`, nil).Value())
}

func TestPropertiesExtensionsRoxxPropertiesExtensionsInt(t *testing.T) {
	customPropertiesRepository := repositories.NewCustomPropertyRepository()
	parser := roxx.NewParser()
	extensions.NewPropertiesExtensions(parser, customPropertiesRepository, nil).Extend()

	customPropertiesRepository.AddCustomProperty(properties.NewIntegerProperty("testKey", 3))

	assert.Equal(t, true, parser.EvaluateExpression(`eq(3, property("testKey"))`, nil).Value())
}

func TestPropertiesExtensionsRoxxPropertiesExtensionsFloat(t *testing.T) {
	customPropertiesRepository := repositories.NewCustomPropertyRepository()
	parser := roxx.NewParser()
	extensions.NewPropertiesExtensions(parser, customPropertiesRepository, nil).Extend()

	customPropertiesRepository.AddCustomProperty(properties.NewFloatProperty("testKey", 3.3))

	assert.Equal(t, true, parser.EvaluateExpression(`eq(3.3, property("testKey"))`, nil).Value())
}

func TestPropertiesExtensionsRoxxPropertiesExtensionsWithContextString(t *testing.T) {
	customPropertiesRepository := repositories.NewCustomPropertyRepository()
	parser := roxx.NewParser()
	extensions.NewPropertiesExtensions(parser, customPropertiesRepository, nil).Extend()

	customPropertiesRepository.AddCustomProperty(properties.NewComputedStringProperty("CustomPropertyTestKey", func(context context.Context) string {
		return context.Get("ContextTestKey").(string)
	}))

	ctx := context.NewContext(map[string]interface{}{"ContextTestKey": "test"})
	assert.Equal(t, true, parser.EvaluateExpression(`eq("test", property("CustomPropertyTestKey"))`, ctx).Value())
}

func TestPropertiesExtensionsRoxxPropertiesExtensionsWithContextInt(t *testing.T) {
	customPropertiesRepository := repositories.NewCustomPropertyRepository()
	parser := roxx.NewParser()
	extensions.NewPropertiesExtensions(parser, customPropertiesRepository, nil).Extend()

	customPropertiesRepository.AddCustomProperty(properties.NewComputedIntegerProperty("CustomPropertyTestKey", func(context context.Context) int {
		return context.Get("ContextTestKey").(int)
	}))

	ctx := context.NewContext(map[string]interface{}{"ContextTestKey": 3})
	assert.Equal(t, true, parser.EvaluateExpression(`eq(3, property("CustomPropertyTestKey"))`, ctx).Value())
}

func TestPropertiesExtensionsRoxxPropertiesExtensionsWithContextIntWithString(t *testing.T) {
	customPropertiesRepository := repositories.NewCustomPropertyRepository()
	parser := roxx.NewParser()
	extensions.NewPropertiesExtensions(parser, customPropertiesRepository, nil).Extend()

	customPropertiesRepository.AddCustomProperty(properties.NewComputedIntegerProperty("CustomPropertyTestKey", func(context context.Context) int {
		return context.Get("ContextTestKey").(int)
	}))

	ctx := context.NewContext(map[string]interface{}{"ContextTestKey": 3})
	assert.Equal(t, false, parser.EvaluateExpression(`eq("3", property("CustomPropertyTestKey"))`, ctx).Value())
}

func TestPropertiesExtensionsRoxxPropertiesExtensionsWithContextIntNotEqual(t *testing.T) {
	customPropertiesRepository := repositories.NewCustomPropertyRepository()
	parser := roxx.NewParser()
	extensions.NewPropertiesExtensions(parser, customPropertiesRepository, nil).Extend()

	customPropertiesRepository.AddCustomProperty(properties.NewComputedIntegerProperty("CustomPropertyTestKey", func(context context.Context) int {
		return context.Get("ContextTestKey").(int)
	}))

	ctx := context.NewContext(map[string]interface{}{"ContextTestKey": 3})
	assert.Equal(t, false, parser.EvaluateExpression(`eq(4, property("CustomPropertyTestKey"))`, ctx).Value())
}

func TestPropertiesExtensionsUnknownProperty(t *testing.T) {
	customPropertiesRepository := repositories.NewCustomPropertyRepository()
	parser := roxx.NewParser()
	extensions.NewPropertiesExtensions(parser, customPropertiesRepository, server.NewRoxOptions(server.RoxOptionsBuilder{}).DynamicPropertyRuleHandler()).Extend()

	customPropertiesRepository.AddCustomProperty(properties.NewStringProperty("testKey", "test"))

	assert.Equal(t, false, parser.EvaluateExpression(`eq("test", property("testKey1"))`, nil).Value())
}


func TestPropertiesExtensionsDynamicProperty(t *testing.T) {
	customPropertiesRepository := repositories.NewCustomPropertyRepository()
	parser := roxx.NewParser()
	extensions.NewPropertiesExtensions(parser, customPropertiesRepository, server.NewRoxOptions(server.RoxOptionsBuilder{}).DynamicPropertyRuleHandler()).Extend()

	assert.Equal(t, true, parser.EvaluateExpression(`eq("test", property("testKey1"))`, context.NewContext(map[string]interface{}{"testKey1":"test"})).Value())
}

func TestPropertiesExtensionsDynamicPropertyBool(t *testing.T) {
	customPropertiesRepository := repositories.NewCustomPropertyRepository()
	parser := roxx.NewParser()
	extensions.NewPropertiesExtensions(parser, customPropertiesRepository, server.NewRoxOptions(server.RoxOptionsBuilder{}).DynamicPropertyRuleHandler()).Extend()

	assert.Equal(t, true, parser.EvaluateExpression(`eq(true, property("testKey1"))`, context.NewContext(map[string]interface{}{"testKey1":true})).Value())
}

func TestPropertiesExtensionsDynamicPropertyInt(t *testing.T) {
	customPropertiesRepository := repositories.NewCustomPropertyRepository()
	parser := roxx.NewParser()
	extensions.NewPropertiesExtensions(parser, customPropertiesRepository, server.NewRoxOptions(server.RoxOptionsBuilder{}).DynamicPropertyRuleHandler()).Extend()

	assert.Equal(t, true, parser.EvaluateExpression(`eq(5, property("testKey1"))`, context.NewContext(map[string]interface{}{"testKey1":5})).Value())
}

func TestPropertiesExtensionsDynamicPropertyDouble(t *testing.T) {
	customPropertiesRepository := repositories.NewCustomPropertyRepository()
	parser := roxx.NewParser()
	extensions.NewPropertiesExtensions(parser, customPropertiesRepository, server.NewRoxOptions(server.RoxOptionsBuilder{}).DynamicPropertyRuleHandler()).Extend()

	assert.Equal(t, true, parser.EvaluateExpression(`eq(5.0, property("testKey1"))`, context.NewContext(map[string]interface{}{"testKey1":5.0})).Value())
}

func TestPropertiesExtensionsDynamicPropertyPrecedence(t *testing.T) {
	customPropertiesRepository := repositories.NewCustomPropertyRepository()
	parser := roxx.NewParser()
	extensions.NewPropertiesExtensions(parser, customPropertiesRepository, nil).Extend()

	customPropertiesRepository.AddCustomProperty(properties.NewStringProperty("testKey1", "testCustomProperty"))

	assert.Equal(t, true, parser.EvaluateExpression(`eq("testCustomProperty", property("testKey1"))`, context.NewContext(map[string]interface{}{"testKey1":"testDynamicProperty"})).Value())
}
