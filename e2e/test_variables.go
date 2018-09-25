package e2e

import "github.com/rollout/rox-go/core/model"

var (
	TestVarsIsComputedBooleanPropCalled bool
	TestVarsIsComputedStringPropCalled  bool
	TestVarsIsComputedIntPropCalled     bool
	TestVarsIsComputedFloatPropCalled   bool
	TestVarsIsComputedSemverPropCalled  bool

	TestVarsTargetGroup1 bool
	TestVarsTargetGroup2 bool

	TestVarsIsImpressionRaised     bool
	TestVarsImpressionReturnedArgs *model.ImpressionArgs

	TestVarsIsPropForTargetGroupForDependency bool

	TestVarsConfigurationFetchedCount int
)
