package entities

import (
	"github.com/rollout/rox-go/v5/core/impression"
	"github.com/rollout/rox-go/v5/core/mocks"
	"github.com/rollout/rox-go/v5/core/model"
	"github.com/rollout/rox-go/v5/core/roxx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFlagWithDefaultValue(t *testing.T) {
	flag1 := NewFlag(false)
	assert.False(t, flag1.IsEnabled(nil))

	flag2 := NewFlag(true)
	assert.True(t, flag2.IsEnabled(nil))
}

func TestFlagWillInvokeEnabledAction(t *testing.T) {
	flag := NewFlag(true)
	isCalled := false
	flag.Enabled(nil, func() {
		isCalled = true
	})

	assert.True(t, isCalled)
}

func TestFlagWillInvokeDisabledAction(t *testing.T) {
	flag := NewFlag(false)
	isCalled := false
	flag.Disabled(nil, func() {
		isCalled = true
	})

	assert.True(t, isCalled)
}

func TestFlagForConsistencyWithString(t *testing.T) {
	parser := roxx.NewParser()

	isImpressionRaised := false
	internalFlags := &mocks.InternalFlags{}
	impInvoker := impression.NewImpressionInvoker(&impression.ImpressionsDeps{
		InternalFlags:            internalFlags,
		CustomPropertyRepository: nil,
		DeviceProperties:         nil,
		Analytics:                nil,
		IsRoxy:                   false,
	})
	impInvoker.RegisterImpressionHandler(func(e model.ImpressionArgs) {
		isImpressionRaised = true
	})

	flag := NewFlag(false)
	flag.(model.InternalVariant).SetForEvaluation(parser, model.NewExperimentModel("id", "name", `ifThen(true, "hey", "yo")`, false, []string{"1"}, nil), impInvoker)

	assert.Equal(t, false, flag.IsEnabled(nil))
	assert.False(t, isImpressionRaised)
}

func TestFlagForConsistencyWithInt(t *testing.T) {
	parser := roxx.NewParser()

	isImpressionRaised := false
	internalFlags := &mocks.InternalFlags{}
	impInvoker := impression.NewImpressionInvoker(&impression.ImpressionsDeps{
		InternalFlags:            internalFlags,
		CustomPropertyRepository: nil,
		DeviceProperties:         nil,
		Analytics:                nil,
		IsRoxy:                   false,
	})
	impInvoker.RegisterImpressionHandler(func(e model.ImpressionArgs) {
		isImpressionRaised = true
	})

	flag := NewFlag(true)
	flag.(model.InternalVariant).SetForEvaluation(parser, model.NewExperimentModel("id", "name", `ifThen(true, 2, 3)`, false, []string{"1"}, nil), impInvoker)

	assert.Equal(t, true, flag.IsEnabled(nil))
	assert.False(t, isImpressionRaised)
}

func TestFlagForConsistencyWithDouble(t *testing.T) {
	parser := roxx.NewParser()

	isImpressionRaised := false
	internalFlags := &mocks.InternalFlags{}
	impInvoker := impression.NewImpressionInvoker(&impression.ImpressionsDeps{
		InternalFlags:            internalFlags,
		CustomPropertyRepository: nil,
		DeviceProperties:         nil,
		Analytics:                nil,
		IsRoxy:                   false,
	})
	impInvoker.RegisterImpressionHandler(func(e model.ImpressionArgs) {
		isImpressionRaised = true
	})

	flag := NewFlag(true)
	flag.(model.InternalVariant).SetForEvaluation(parser, model.NewExperimentModel("id", "name", `ifThen(true, 2.5, 3.5)`, false, []string{"1"}, nil), impInvoker)

	assert.Equal(t, true, flag.IsEnabled(nil))
	assert.False(t, isImpressionRaised)
}
