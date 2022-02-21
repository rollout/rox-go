package context_test

import (
	"github.com/rollout/rox-go/v4/core/context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMergedContextWithNullLocalContext(t *testing.T) {
	globalContext := context.NewContext(map[string]interface{}{
		"a": 1,
	})

	mergedContext := context.NewMergedContext(globalContext, nil)

	assert.Equal(t, 1, mergedContext.Get("a"))
	assert.Equal(t, nil, mergedContext.Get("b"))
}

func TestMergedContextWithNullGlobalContext(t *testing.T) {
	localContext := context.NewContext(map[string]interface{}{
		"a": 1,
	})

	mergedContext := context.NewMergedContext(nil, localContext)

	assert.Equal(t, 1, mergedContext.Get("a"))
	assert.Equal(t, nil, mergedContext.Get("b"))
}

func TestMergedContextWithLocalAndGlobalContext(t *testing.T) {
	globalContext := context.NewContext(map[string]interface{}{
		"a": 1,
		"b": 2,
	})

	localContext := context.NewContext(map[string]interface{}{
		"a": 3,
		"c": 4,
	})

	mergedContext := context.NewMergedContext(globalContext, localContext)

	assert.Equal(t, 3, mergedContext.Get("a"))
	assert.Equal(t, 2, mergedContext.Get("b"))
	assert.Equal(t, 4, mergedContext.Get("c"))
	assert.Equal(t, nil, mergedContext.Get("d"))
}
