package client_test

import (
	"github.com/rollout/rox-go/v4/core/client"
	"github.com/rollout/rox-go/v4/core/entities"
	"github.com/rollout/rox-go/v4/core/model"
	"github.com/rollout/rox-go/v4/core/repositories"
	"github.com/rollout/rox-go/v4/core/roxx"
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

func TestDynamicAPIGetValue(t *testing.T) {
	parser := roxx.NewParser()
	flagRepo := repositories.NewFlagRepository()
	expRepo := repositories.NewExperimentRepository()
	flagSetter := entities.NewFlagSetter(flagRepo, parser, expRepo, nil)
	dynamicAPI := client.NewDynamicAPI(flagRepo, &entitiesMockProvider{})

	assert.Equal(t, "A", dynamicAPI.Value("default.newVariant", "A", []string{"A", "B", "C"}, nil))
	assert.Equal(t, "A", flagRepo.GetFlag("default.newVariant").GetValue(nil))
	assert.Equal(t, "B", dynamicAPI.Value("default.newVariant", "B", []string{"A", "B", "C"}, nil))
	assert.Equal(t, 1, len(flagRepo.GetAllFlags()))

	expRepo.SetExperiments([]*model.ExperimentModel{model.NewExperimentModel("1", "default.newVariant", `ifThen(true, "B", "A")`, false, []string{"default.newVariant"}, nil)})
	flagSetter.SetExperiments()

	assert.Equal(t, "B", dynamicAPI.Value("default.newVariant", "A", []string{"A", "B", "C"}, nil))
}

func TestDynamicAPIGetValueWithoutOptions(t *testing.T) {
	parser := roxx.NewParser()
	flagRepo := repositories.NewFlagRepository()
	expRepo := repositories.NewExperimentRepository()
	flagSetter := entities.NewFlagSetter(flagRepo, parser, expRepo, nil)
	dynamicAPI := client.NewDynamicAPI(flagRepo, &entitiesMockProvider{})

	assert.Equal(t, "A", dynamicAPI.Value("default.newVariant", "A", nil, nil))
	assert.Equal(t, "A", flagRepo.GetFlag("default.newVariant").GetValue(nil))
	assert.Equal(t, "B", dynamicAPI.Value("default.newVariant", "B", nil, nil))
	assert.Equal(t, 1, len(flagRepo.GetAllFlags()))

	expRepo.SetExperiments([]*model.ExperimentModel{model.NewExperimentModel("1", "default.newVariant", `ifThen(true, "B", "A")`, false, []string{"default.newVariant"}, nil)})
	flagSetter.SetExperiments()

	assert.Equal(t, "B", dynamicAPI.Value("default.newVariant", "A", nil, nil))
}

type entitiesMockProvider struct {
}

func (*entitiesMockProvider) CreateFlag(defaultValue bool) model.Flag {
	return entities.NewFlag(defaultValue)
}

func (*entitiesMockProvider) CreateVariant(defaultValue string, options []string) model.Variant {
	return entities.NewVariant(defaultValue, options)
}
