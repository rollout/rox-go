package repositories_test

import (
	"testing"

	"github.com/rollout/rox-go/v6/core/model"
	"github.com/rollout/rox-go/v6/core/repositories"
	"github.com/stretchr/testify/assert"
)

func TestTargetGroupRepositoryWillReturnNullWhenNotFound(t *testing.T) {
	tgs := []*model.TargetGroupModel{
		model.NewTargetGroupModel("1", "x"),
	}
	tgRepo := repositories.NewTargetGroupRepository()
	tgRepo.SetTargetGroups(tgs)

	assert.Nil(t, tgRepo.GetTargetGroup("2"))
}

func TestTargetGroupRepositoryWillReturnWhenFound(t *testing.T) {
	tgs := []*model.TargetGroupModel{
		model.NewTargetGroupModel("1", "x"),
	}
	tgRepo := repositories.NewTargetGroupRepository()
	tgRepo.SetTargetGroups(tgs)

	assert.Equal(t, "x", tgRepo.GetTargetGroup("1").Condition)
}
