package repositories

import (
	"github.com/rollout/rox-go/v5/core/model"
)

type experimentRepository struct {
	experiments []*model.ExperimentModel
}

func NewExperimentRepository() model.ExperimentRepository {
	return &experimentRepository{}
}

func (r *experimentRepository) SetExperiments(experiments []*model.ExperimentModel) {
	r.experiments = experiments
}

func (r *experimentRepository) GetExperimentByFlag(flagName string) *model.ExperimentModel {
	for _, e := range r.experiments {
		for _, f := range e.Flags {
			if f == flagName {
				return e
			}
		}
	}
	return nil
}

func (r *experimentRepository) GetAllExperiments() []*model.ExperimentModel {
	return r.experiments
}
