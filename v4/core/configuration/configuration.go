package configuration

import (
	"github.com/rollout/rox-go/v4/core/model"
	"time"
)

type Configuration struct {
	Experiments   []*model.ExperimentModel
	TargetGroups  []*model.TargetGroupModel
	SignatureDate time.Time
}

func NewConfiguration(experiments []*model.ExperimentModel, targetGroups []*model.TargetGroupModel, signatureDate time.Time) *Configuration {
	return &Configuration{
		Experiments:   experiments,
		TargetGroups:  targetGroups,
		SignatureDate: signatureDate,
	}
}
