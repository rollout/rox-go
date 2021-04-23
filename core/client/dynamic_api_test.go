package client_test

import (
	"github.com/rollout/rox-go/core/client"
	"github.com/rollout/rox-go/core/entities"
	"github.com/rollout/rox-go/core/model"
	"github.com/rollout/rox-go/core/repositories"
	"github.com/rollout/rox-go/core/roxx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDynamicAPIIsEnabled(t *testing.T) {
	parser := roxx.NewParser()
	flagRepo := repositories.NewFlagRepository()
	expRepo := repositories.NewExperimentRepository()
	flagSetter := entities.NewFlagSetter(flagRepo, parser, expRepo, nil)
	dynamicAPI := client.NewDynamicAPI(flagRepo, &entitiesMockProvider{})

	assert.True(t, dynamicAPI.IsEnabled("default.newFlag", true, nil))
	assert.True(t, flagRepo.GetFlag("default.newFlag").(model.Flag).IsEnabled(nil))
	assert.False(t, dynamicAPI.IsEnabled("default.newFlag", false, nil))
	assert.Equal(t, 1, len(flagRepo.GetAllFlags()))

	expRepo.SetExperiments([]*model.ExperimentModel{model.NewExperimentModel("1", "default.newFlag", "and(true, true)", false, []string{"default.newFlag"}, nil)})
	flagSetter.SetExperiments()

	assert.True(t, dynamicAPI.IsEnabled("default.newFlag", false, nil))
}

func TestDynamicAPIIsEnabledAfterSetup(t *testing.T) {
	parser := roxx.NewParser()
	flagRepo := repositories.NewFlagRepository()
	expRepo := repositories.NewExperimentRepository()
	flagSetter := entities.NewFlagSetter(flagRepo, parser, expRepo, nil)
	dynamicAPI := client.NewDynamicAPI(flagRepo, &entitiesMockProvider{})

	expRepo.SetExperiments([]*model.ExperimentModel{model.NewExperimentModel("1", "default.newFlag", "and(true, true)", false, []string{"default.newFlag"}, nil)})
	flagSetter.SetExperiments()

	assert.True(t, dynamicAPI.IsEnabled("default.newFlag", false, nil))
}

func TestDynamicAPIGetStringValue(t *testing.T) {
	parser := roxx.NewParser()
	flagRepo := repositories.NewFlagRepository()
	expRepo := repositories.NewExperimentRepository()
	flagSetter := entities.NewFlagSetter(flagRepo, parser, expRepo, nil)
	dynamicAPI := client.NewDynamicAPI(flagRepo, &entitiesMockProvider{})

	assert.Equal(t, "A", dynamicAPI.StringValue("default.newVariant", "A", []string{"A", "B", "C"}, nil))
	assert.Equal(t, "A", flagRepo.GetFlag("default.newVariant").GetValueAsString(nil))
	assert.Equal(t, "B", dynamicAPI.StringValue("default.newVariant", "B", []string{"A", "B", "C"}, nil))
	assert.Equal(t, 1, len(flagRepo.GetAllFlags()))

	expRepo.SetExperiments([]*model.ExperimentModel{model.NewExperimentModel("1", "default.newVariant", `ifThen(true, "B", "A")`, false, []string{"default.newVariant"}, nil)})
	flagSetter.SetExperiments()

	assert.Equal(t, "B", dynamicAPI.StringValue("default.newVariant", "A", []string{"A", "B", "C"}, nil))
}

func TestDynamicAPIGetIntValue(t *testing.T) {
	parser := roxx.NewParser()
	flagRepo := repositories.NewFlagRepository()
	expRepo := repositories.NewExperimentRepository()
	flagSetter := entities.NewFlagSetter(flagRepo, parser, expRepo, nil)
	dynamicAPI := client.NewDynamicAPI(flagRepo, &entitiesMockProvider{})

	assert.Equal(t, 1, dynamicAPI.IntValue("default.newVariant", 1, []int{1, 2, 3}, nil))
	assert.Equal(t, 1, flagRepo.GetFlag("default.newVariant").(model.RoxInt).GetValue(nil))
	assert.Equal(t, 2, dynamicAPI.IntValue("default.newVariant", 2, []int{3, 4, 5}, nil))
	assert.Equal(t, 1, len(flagRepo.GetAllFlags()))

	expRepo.SetExperiments([]*model.ExperimentModel{model.NewExperimentModel("1", "default.newVariant", `ifThen(true, 2, 1)`, false, []string{"default.newVariant"}, nil)})
	flagSetter.SetExperiments()

	assert.Equal(t, 2, dynamicAPI.IntValue("default.newVariant", 1, []int{1, 2, 3}, nil))
}

func TestDynamicAPIGetDoubleValue(t *testing.T) {
	parser := roxx.NewParser()
	flagRepo := repositories.NewFlagRepository()
	expRepo := repositories.NewExperimentRepository()
	flagSetter := entities.NewFlagSetter(flagRepo, parser, expRepo, nil)
	dynamicAPI := client.NewDynamicAPI(flagRepo, &entitiesMockProvider{})

	assert.Equal(t, 1, dynamicAPI.IntValue("default.newVariant", 1, []int{1, 2, 3}, nil))
	assert.Equal(t, 1, flagRepo.GetFlag("default.newVariant").(model.RoxInt).GetValue(nil))
	assert.Equal(t, 2, dynamicAPI.IntValue("default.newVariant", 2, []int{3, 4, 5}, nil))
	assert.Equal(t, 1, len(flagRepo.GetAllFlags()))

	expRepo.SetExperiments([]*model.ExperimentModel{model.NewExperimentModel("1", "default.newVariant", `ifThen(true, 2, 1)`, false, []string{"default.newVariant"}, nil)})
	flagSetter.SetExperiments()

	assert.Equal(t, 2, dynamicAPI.IntValue("default.newVariant", 1, []int{1, 2, 3}, nil))
}

func TestDynamicAPIGetValueWithoutOptions(t *testing.T) {
	parser := roxx.NewParser()
	flagRepo := repositories.NewFlagRepository()
	expRepo := repositories.NewExperimentRepository()
	flagSetter := entities.NewFlagSetter(flagRepo, parser, expRepo, nil)
	dynamicAPI := client.NewDynamicAPI(flagRepo, &entitiesMockProvider{})

	assert.Equal(t, "A", dynamicAPI.StringValue("default.newVariant", "A", nil, nil))
	assert.Equal(t, "A", flagRepo.GetFlag("default.newVariant").GetValueAsString(nil))
	assert.Equal(t, "B", dynamicAPI.StringValue("default.newVariant", "B", nil, nil))
	assert.Equal(t, 1, len(flagRepo.GetAllFlags()))

	expRepo.SetExperiments([]*model.ExperimentModel{model.NewExperimentModel("1", "default.newVariant", `ifThen(true, "B", "A")`, false, []string{"default.newVariant"}, nil)})
	flagSetter.SetExperiments()

	assert.Equal(t, "B", dynamicAPI.StringValue("default.newVariant", "A", nil, nil))
}

type entitiesMockProvider struct {
}

func (*entitiesMockProvider) CreateFlag(defaultValue bool) model.Flag {
	return entities.NewFlag(defaultValue)
}

func (*entitiesMockProvider) CreateRoxString(defaultValue string, options []string) model.RoxString {
	return entities.NewRoxString(defaultValue, options)
}

func (*entitiesMockProvider) CreateRoxInt(defaultValue int, options []int) model.RoxInt {
	return entities.NewRoxInt(defaultValue, options)
}

func (*entitiesMockProvider) CreateRoxDouble(defaultValue float64, options []float64) model.RoxDouble {
	return entities.NewRoxDouble(defaultValue, options)
}
