package entities

import (
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
