package e2e

import (
	"github.com/rollout/rox-go/core/context"
	"github.com/rollout/rox-go/server"
)

func createCustomProps(rox *server.Rox) {
	rox.SetCustomStringProperty("StringProp1", "Hello")
	rox.SetCustomComputedStringProperty("StringProp2", func(context context.Context) string {
		TestVarsIsComputedStringPropCalled = true
		return "World"
	})

	rox.SetCustomBooleanProperty("BoolProp1", true)
	rox.SetCustomComputedBooleanProperty("BoolProp2", func(context context.Context) bool {
		TestVarsIsComputedBooleanPropCalled = true
		return false
	})

	rox.SetCustomIntegerProperty("IntProp1", 6)
	rox.SetCustomComputedIntegerProperty("IntProp2", func(context context.Context) int {
		TestVarsIsComputedIntPropCalled = true
		return 28
	})

	rox.SetCustomFloatProperty("FloatProp1", 3.14)
	rox.SetCustomComputedFloatProperty("FloatProp2", func(context context.Context) float64 {
		TestVarsIsComputedFloatPropCalled = true
		return 1.618
	})

	rox.SetCustomSemverProperty("SmvrProp1", "9.11.2001")
	rox.SetCustomComputedSemverProperty("SmvrProp2", func(context context.Context) string {
		TestVarsIsComputedSemverPropCalled = true
		return "20.7.1969"
	})

	rox.SetCustomComputedBooleanProperty("BoolPropTargetGroupForVariant", func(context context.Context) bool {
		value, _ := context.Get("isDuckAndCover").(bool)
		return value
	})

	rox.SetCustomComputedBooleanProperty("BoolPropTargetGroupOperand1", func(context context.Context) bool {
		return TestVarsTargetGroup1
	})

	rox.SetCustomComputedBooleanProperty("BoolPropTargetGroupOperand2", func(context context.Context) bool {
		return TestVarsTargetGroup2
	})

	rox.SetCustomComputedBooleanProperty("BoolPropTargetGroupForDependency", func(context context.Context) bool {
		return TestVarsIsPropForTargetGroupForDependency
	})

	rox.SetCustomComputedBooleanProperty("BoolPropTargetGroupForVariantDependency", func(context context.Context) bool {
		value, _ := context.Get("isDuckAndCover").(bool)
		return value
	})
}
