package repositories_test

import (
	"github.com/rollout/rox-go/core/model"
	"github.com/rollout/rox-go/core/repositories"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExperimentRepositoryWillReturnNullWhenNotFound(t *testing.T) {
	exp := []*model.ExperimentModel{
		model.NewExperimentModel("1", "1", "1", false, []string{"a"}, nil),
	}

	expRepo := repositories.NewExperimentRepository()
	expRepo.SetExperiments(exp)

	assert.Nil(t, expRepo.GetExperimentByFlag("b"))
}

func TestExperimentRepositoryWillReturnWhenFound(t *testing.T) {
	exp := []*model.ExperimentModel{
		model.NewExperimentModel("1", "1", "1", false, []string{"a"}, nil),
	}

	expRepo := repositories.NewExperimentRepository()
	expRepo.SetExperiments(exp)

	assert.Equal(t, "1", expRepo.GetExperimentByFlag("a").Id)
}
