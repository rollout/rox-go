package extensions_test

import (
	"github.com/rollout/rox-go/core/extensions"
	"github.com/rollout/rox-go/core/repositories"
	"github.com/rollout/rox-go/core/roxx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExperimentsExtensionsGetBucket(t *testing.T) {
	parser := roxx.NewParser()
	targetGroupsRepository := repositories.NewTargetGroupRepository()
	experimentsExtensions := extensions.NewExperimentsExtensions(parser, targetGroupsRepository, nil, nil)

	bucket := experimentsExtensions.GetBucket("device2.seed2")
	assert.Equal(t, 0.18721251450181298, bucket)
}
