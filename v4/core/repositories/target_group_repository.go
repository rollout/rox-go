package repositories

import (
	"github.com/rollout/rox-go/v4/core/model"
)

type targetGroupRepository struct {
	targetGroups []*model.TargetGroupModel
}

func NewTargetGroupRepository() model.TargetGroupRepository {
	return &targetGroupRepository{}
}

func (r *targetGroupRepository) SetTargetGroups(targetGroups []*model.TargetGroupModel) {
	r.targetGroups = targetGroups
}

func (r *targetGroupRepository) GetTargetGroup(id string) *model.TargetGroupModel {
	for _, g := range r.targetGroups {
		if g.ID == id {
			return g
		}
	}
	return nil
}
