package e2e

import (
	"fmt"
	"github.com/rollout/rox-go/core/context"
	"github.com/rollout/rox-go/core/model"
	"github.com/rollout/rox-go/server"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

type logger struct {
}

func (*logger) Debug(message string, err interface{}) {
	fmt.Println("Before Rox.Setup", message)
}

func (*logger) Warn(message string, err interface{}) {
	fmt.Println("Before Rox.Setup", message)
}

func (*logger) Error(message string, err interface{}) {
	fmt.Println("Before Rox.Setup", message)
}

func TestMain(m *testing.M) {
	setUp()
	retCode := m.Run()
	os.Exit(retCode)
}

var rox *server.Rox

func setUp() {
	os.Setenv("ROLLOUT_MODE", "QA")

	options := server.NewRoxOptions(server.RoxOptionsBuilder{
		DevModeKey: "c44d343131b9231f5a001a47",
		Logger:     &logger{},
		ConfigurationFetchedHandler: func(e *model.ConfigurationFetchedArgs) {
			if e != nil && e.FetcherStatus == model.FetcherStatusAppliedFromNetwork {
				TestVarsConfigurationFetchedCount++
			}
		},
		ImpressionHandler: func(e model.ImpressionArgs) {
			if e.ReportingValue != nil {
				if e.ReportingValue.Name == "FlagForImpression" {
					TestVarsIsImpressionRaised = true
				}
			}
			TestVarsImpressionReturnedArgs = &e
		},
	})

	rox = server.NewRox()
	rox.Register("", container)
	createCustomProps(rox)

	<-rox.Setup("5b8268e3de813f34648debd1", options)
}

func TestSimpleFlag(t *testing.T) {
	assert.True(t, container.SimpleFlag.IsEnabled(nil))
}

func TestSimpleFlagOverwritten(t *testing.T) {
	assert.False(t, container.SimpleFlagOverwritten.IsEnabled(nil))
}

func TestVariant(t *testing.T) {
	assert.Equal(t, "red", container.Variant.GetValue(nil))
}

func TestVariantOverwritten(t *testing.T) {
	assert.Equal(t, "green", container.VariantOverwritten.GetValue(nil))
}

func TestAllCustomProperties(t *testing.T) {
	assert.True(t, container.FlagCustomProperties.IsEnabled(nil))

	assert.True(t, TestVarsIsComputedBooleanPropCalled)
	assert.True(t, TestVarsIsComputedFloatPropCalled)
	assert.True(t, TestVarsIsComputedIntPropCalled)
	assert.True(t, TestVarsIsComputedSemverPropCalled)
	assert.True(t, TestVarsIsComputedStringPropCalled)
}

func TestFetchWithinTimeout(t *testing.T) {
	numberOfConfigFetches := TestVarsConfigurationFetchedCount
	timer := time.NewTimer(time.Second * 5)
	select {
	case <-rox.Fetch():
	case <-timer.C:
		assert.Fail(t, "timeout")
	}

	assert.True(t, numberOfConfigFetches < TestVarsConfigurationFetchedCount)
}

func TestVariantWithContext(t *testing.T) {
	somePositiveContext := context.NewContext(map[string]interface{}{"isDuckAndCover": true})
	someNegativeContext := context.NewContext(map[string]interface{}{"isDuckAndCover": false})

	assert.Equal(t, "red", container.VariantWithContext.GetValue(nil))
	assert.Equal(t, "blue", container.VariantWithContext.GetValue(somePositiveContext))
	assert.Equal(t, "red", container.VariantWithContext.GetValue(someNegativeContext))
}

func TestTargetGroupsAllAnyNone(t *testing.T) {
	TestVarsTargetGroup1 = true
	TestVarsTargetGroup2 = true

	assert.True(t, container.FlagTargetGroupsAll.IsEnabled(nil))
	assert.True(t, container.FlagTargetGroupsAny.IsEnabled(nil))
	assert.False(t, container.FlagTargetGroupsNone.IsEnabled(nil))

	TestVarsTargetGroup1 = false

	assert.False(t, container.FlagTargetGroupsAll.IsEnabled(nil))
	assert.True(t, container.FlagTargetGroupsAny.IsEnabled(nil))
	assert.False(t, container.FlagTargetGroupsNone.IsEnabled(nil))

	TestVarsTargetGroup2 = false

	assert.False(t, container.FlagTargetGroupsAll.IsEnabled(nil))
	assert.False(t, container.FlagTargetGroupsAny.IsEnabled(nil))
	assert.True(t, container.FlagTargetGroupsNone.IsEnabled(nil))
}

func TestImpressionHandler(t *testing.T) {
	container.FlagForImpression.IsEnabled(nil)
	assert.True(t, TestVarsIsImpressionRaised)
	TestVarsIsImpressionRaised = false

	ctx := context.NewContext(map[string]interface{}{"var": "val"})
	flagImpressionValue := container.FlagForImpressionWithExperimentAndContext.IsEnabled(ctx)
	assert.NotNil(t, TestVarsImpressionReturnedArgs)
	assert.NotNil(t, TestVarsImpressionReturnedArgs.ReportingValue)
	assert.Equal(t, "true", TestVarsImpressionReturnedArgs.ReportingValue.Value)
	assert.True(t, flagImpressionValue)
	assert.Equal(t, "FlagForImpressionWithExperimentAndContext", TestVarsImpressionReturnedArgs.ReportingValue.Name)

	assert.NotNil(t, TestVarsImpressionReturnedArgs)
	assert.NotNil(t, TestVarsImpressionReturnedArgs.Experiment)
	assert.Equal(t, "5b828888de813f34648ded70", TestVarsImpressionReturnedArgs.Experiment.Identifier)
	assert.Equal(t, "flag for impression with experiment and context", TestVarsImpressionReturnedArgs.Experiment.Name)

	assert.Equal(t, "val", TestVarsImpressionReturnedArgs.Context.Get("var"))
}

func TestFlagDependency(t *testing.T) {
	TestVarsIsPropForTargetGroupForDependency = true
	assert.True(t, container.FlagForDependency.IsEnabled(nil))
	assert.False(t, container.FlagDependent.IsEnabled(nil))

	TestVarsIsPropForTargetGroupForDependency = false
	assert.True(t, container.FlagDependent.IsEnabled(nil))
	assert.False(t, container.FlagForDependency.IsEnabled(nil))
}

func TestVariantDependencyWithContext(t *testing.T) {
	somePositiveContext := context.NewContext(map[string]interface{}{"isDuckAndCover": true})
	someNegativeContext := context.NewContext(map[string]interface{}{"isDuckAndCover": false})

	assert.Equal(t, "White", container.FlagColorDependentWithContext.GetValue(nil))
	assert.Equal(t, "White", container.FlagColorDependentWithContext.GetValue(someNegativeContext))
	assert.Equal(t, "Yellow", container.FlagColorDependentWithContext.GetValue(somePositiveContext))
}
