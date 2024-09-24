package context_test

import (
	"testing"

	"github.com/rollout/rox-go/v6/core/context"
	"github.com/stretchr/testify/assert"
)

func TestContextWillReturnValue(t *testing.T) {
	ctx := context.NewContext(map[string]interface{}{"a": "b"})

	assert.Equal(t, "b", ctx.Get("a"))
}

func TestContextWillReturnNull(t *testing.T) {
	ctx := context.NewContext(map[string]interface{}{})

	assert.Equal(t, nil, ctx.Get("a"))
}

func TestContextWithNullMap(t *testing.T) {
	ctx := context.NewContext(nil)

	assert.Equal(t, nil, ctx.Get("a"))
}
