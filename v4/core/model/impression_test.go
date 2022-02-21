package model_test

import (
	"github.com/magiconair/properties/assert"
	"github.com/rollout/rox-go/v4/core/model"
	"testing"
)

func TestExperimentConstructor(t *testing.T) {
	originalExperiment := model.NewExperimentModel("id", "name", "cond", true, nil, []string{"name1"})
	experiment := model.NewExperiment(originalExperiment)

	assert.Equal(t, originalExperiment.Name, experiment.Name)
	assert.Equal(t, originalExperiment.ID, experiment.Identifier)
	assert.Equal(t, originalExperiment.IsArchived, experiment.IsArchived)
	assert.Equal(t, originalExperiment.Labels, experiment.Labels)
}

func TestReportingValueConstructor(t *testing.T) {
	reportingValue := model.NewReportingValue("pi", "ka")

	assert.Equal(t, "pi", reportingValue.Name)
	assert.Equal(t, "ka", reportingValue.Value)
}
