package mocks

import (
	"github.com/rollout/rox-go/v6/core/model"
	"github.com/stretchr/testify/mock"
)

type ExperimentRepository struct {
	mock.Mock
}

func (m *ExperimentRepository) SetExperiments(experiments []*model.ExperimentModel) {
	m.Called(experiments)
}

func (m *ExperimentRepository) GetExperimentByFlag(flagName string) *model.ExperimentModel {
	args := m.Called(flagName)
	result := args.Get(0)
	if result == nil {
		return nil
	}
	return result.(*model.ExperimentModel)
}

func (m *ExperimentRepository) GetAllExperiments() []*model.ExperimentModel {
	args := m.Called()
	return args.Get(0).([]*model.ExperimentModel)
}
