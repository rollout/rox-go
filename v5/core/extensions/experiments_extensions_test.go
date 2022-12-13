package extensions_test

import (
	"github.com/rollout/rox-go/v5/core/context"
	"github.com/rollout/rox-go/v5/core/entities"
	"github.com/rollout/rox-go/v5/core/extensions"
	"github.com/rollout/rox-go/v5/core/impression"
	"github.com/rollout/rox-go/v5/core/mocks"
	"github.com/rollout/rox-go/v5/core/model"
	"github.com/rollout/rox-go/v5/core/properties"
	"github.com/rollout/rox-go/v5/core/repositories"
	"github.com/rollout/rox-go/v5/core/roxx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExperimentsExtensionsCustomPropertyWithSimpleValue(t *testing.T) {
	parser := roxx.NewParser()
	targetGroupsRepository := repositories.NewTargetGroupRepository()
	experimentsExtensions := extensions.NewExperimentsExtensions(parser, targetGroupsRepository, nil, nil)
	experimentsExtensions.Extend()

	assert.Equal(t, false, parser.EvaluateExpression(`isInTargetGroup("targetGroup1")`, nil).Value())
}

func TestExperimentsExtensionsIsInPercentageRange(t *testing.T) {
	parser := roxx.NewParser()
	targetGroupsRepository := repositories.NewTargetGroupRepository()
	experimentsExtensions := extensions.NewExperimentsExtensions(parser, targetGroupsRepository, nil, nil)
	experimentsExtensions.Extend()

	assert.Equal(t, true, parser.EvaluateExpression(`isInPercentageRange(0, 0.5, "device2.seed2")`, nil).Value())
}

func TestExperimentsExtensionsNotIsInPercentageRange(t *testing.T) {
	parser := roxx.NewParser()
	targetGroupsRepository := repositories.NewTargetGroupRepository()
	experimentsExtensions := extensions.NewExperimentsExtensions(parser, targetGroupsRepository, nil, nil)
	experimentsExtensions.Extend()

	assert.Equal(t, false, parser.EvaluateExpression(`isInPercentageRange(0.5, 1, "device2.seed2")`, nil).Value())
}

func TestExperimentsExtensionsGetBucket(t *testing.T) {
	parser := roxx.NewParser()
	targetGroupsRepository := repositories.NewTargetGroupRepository()
	experimentsExtensions := extensions.NewExperimentsExtensions(parser, targetGroupsRepository, nil, nil)

	bucket := experimentsExtensions.GetBucket("device2.seed2")
	assert.Equal(t, 0.18721251450181298, bucket)
}

func TestExperimentsExtensionsFlagValueNoFlagNoExperiment(t *testing.T) {
	parser := roxx.NewParser()
	targetGroupsRepository := repositories.NewTargetGroupRepository()
	experimentRepository := repositories.NewExperimentRepository()
	flagRepository := repositories.NewFlagRepository()
	experimentsExtensions := extensions.NewExperimentsExtensions(parser, targetGroupsRepository, flagRepository, experimentRepository)
	experimentsExtensions.Extend()

	assert.Equal(t, "false", parser.EvaluateExpression(`flagValue("f1")`, nil).Value())
}

func TestExperimentsExtensionsFlagValueNoFlagEvaluateExperiment(t *testing.T) {
	parser := roxx.NewParser()
	targetGroupsRepository := repositories.NewTargetGroupRepository()
	experimentRepository := repositories.NewExperimentRepository()
	flagRepository := repositories.NewFlagRepository()
	experimentsExtensions := extensions.NewExperimentsExtensions(parser, targetGroupsRepository, flagRepository, experimentRepository)
	experimentsExtensions.Extend()

	experiments := []*model.ExperimentModel{
		model.NewExperimentModel("id", "name", `"op2"`, false, []string{"f1"}, nil),
	}
	experimentRepository.SetExperiments(experiments)

	assert.Equal(t, "op2", parser.EvaluateExpression(`flagValue("f1")`, nil).Value())
}

func TestExperimentsExtensionsFlagValueFlagEvaluationDefault(t *testing.T) {
	parser := roxx.NewParser()
	targetGroupsRepository := repositories.NewTargetGroupRepository()
	experimentRepository := repositories.NewExperimentRepository()
	flagRepository := repositories.NewFlagRepository()
	experimentsExtensions := extensions.NewExperimentsExtensions(parser, targetGroupsRepository, flagRepository, experimentRepository)
	experimentsExtensions.Extend()

	v := entities.NewRoxString("op1", []string{"op2"})
	flagRepository.AddFlag(v, "f1", "f1")

	assert.Equal(t, "op1", parser.EvaluateExpression(`flagValue("f1")`, nil).Value())
}

func TestExperimentsExtensionsFlagDependencyValue(t *testing.T) {
	parser := roxx.NewParser()
	targetGroupsRepository := repositories.NewTargetGroupRepository()
	experimentRepository := repositories.NewExperimentRepository()
	flagRepository := repositories.NewFlagRepository()
	experimentsExtensions := extensions.NewExperimentsExtensions(parser, targetGroupsRepository, flagRepository, experimentRepository)
	experimentsExtensions.Extend()

	f := entities.NewFlag(false)
	flagRepository.AddFlag(f, "f1", "f1")

	v := entities.NewRoxString("blue", []string{"red", "green"})
	flagRepository.AddFlag(v, "v1", "v1")
	exp := model.NewExperimentModel("id", "name", `ifThen(eq("true", flagValue("f1")), "red", "green")`, false, nil, nil)
	v.(model.InternalVariant).SetForEvaluation(parser, exp, nil)

	assert.Equal(t, "green", v.GetValueAsString(nil))
}

func TestExperimentsExtensionsFlagDependencyImpressionHandler(t *testing.T) {
	parser := roxx.NewParser()
	targetGroupsRepository := repositories.NewTargetGroupRepository()
	experimentRepository := repositories.NewExperimentRepository()
	flagRepository := repositories.NewFlagRepository()
	internalFlags := &mocks.InternalFlags{}
	ii := impression.NewImpressionInvoker(internalFlags, nil, nil, false)
	experimentsExtensions := extensions.NewExperimentsExtensions(parser, targetGroupsRepository, flagRepository, experimentRepository)
	experimentsExtensions.Extend()

	f := entities.NewFlag(false)
	flagRepository.AddFlag(f, "f1", "f1")
	f.(model.InternalVariant).SetForEvaluation(parser, nil, ii)

	v := entities.NewRoxString("blue", []string{"red", "green"})
	flagRepository.AddFlag(v, "v1", "v1")
	exp := model.NewExperimentModel("id", "name", `ifThen(eq("true", flagValue("f1")), "red", "green")`, false, nil, nil)
	v.(model.InternalVariant).SetForEvaluation(parser, exp, ii)

	var impressions []model.ImpressionArgs
	ii.RegisterImpressionHandler(func(args model.ImpressionArgs) {
		impressions = append(impressions, args)
	})

	assert.Equal(t, "green", v.GetValueAsString(nil))

	assert.Equal(t, 1, len(impressions))
	assert.Equal(t, "v1", impressions[0].ReportingValue.Name)
	assert.Equal(t, "green", impressions[0].ReportingValue.Value)
}

func TestExperimentsExtensionsFlagDependency2LevelsBottomNotExists(t *testing.T) {
	parser := roxx.NewParser()
	targetGroupsRepository := repositories.NewTargetGroupRepository()
	experimentRepository := repositories.NewExperimentRepository()
	flagRepository := repositories.NewFlagRepository()
	experimentsExtensions := extensions.NewExperimentsExtensions(parser, targetGroupsRepository, flagRepository, experimentRepository)
	experimentsExtensions.Extend()

	f := entities.NewFlag(false)
	flagRepository.AddFlag(f, "f1", "f1")
	exp1 := model.NewExperimentModel("id1", "name1", `flagValue("someFlag")`, false, nil, nil)
	f.(model.InternalVariant).SetForEvaluation(parser, exp1, nil)

	v := entities.NewRoxString("blue", []string{"red", "green"})
	flagRepository.AddFlag(v, "v1", "f1")
	exp2 := model.NewExperimentModel("id2", "name2", `ifThen(eq("true", flagValue("f1")), "red", "green")`, false, nil, nil)
	v.(model.InternalVariant).SetForEvaluation(parser, exp2, nil)

	assert.Equal(t, "green", v.GetValueAsString(nil))
}

func TestExperimentsExtensionsFlagDependencyUnexistingFlagButExistingExperiment(t *testing.T) {
	parser := roxx.NewParser()
	targetGroupsRepository := repositories.NewTargetGroupRepository()
	experimentRepository := repositories.NewExperimentRepository()
	flagRepository := repositories.NewFlagRepository()

	experimentModels := []*model.ExperimentModel{
		model.NewExperimentModel("exp1id", "exp1name", `ifThen(true, "true", "false")`, false, []string{"someFlag"}, nil),
		model.NewExperimentModel("exp2id", "exp2name", `ifThen(eq("true", flagValue("someFlag")), "blue", "green")`, false, []string{"colorVar"}, nil),
	}

	flagSetter := entities.NewFlagSetter(flagRepository, parser, experimentRepository, nil)
	experimentRepository.SetExperiments(experimentModels)
	flagSetter.SetExperiments()

	experimentsExtensions := extensions.NewExperimentsExtensions(parser, targetGroupsRepository, flagRepository, experimentRepository)
	experimentsExtensions.Extend()

	colorVar := entities.NewRoxString("red", []string{"red", "green", "blue"})
	colorVar.(model.InternalVariant).SetForEvaluation(parser, nil, nil)
	flagRepository.AddFlag(colorVar, "colorVar", "colorVar")

	assert.Equal(t, "blue", colorVar.GetValueAsString(nil))
}

func TestExperimentsExtensionsFlagDependencyUnexistingFlagAndExperimentUndefined(t *testing.T) {
	parser := roxx.NewParser()
	targetGroupsRepository := repositories.NewTargetGroupRepository()
	experimentRepository := repositories.NewExperimentRepository()
	flagRepository := repositories.NewFlagRepository()

	experimentModels := []*model.ExperimentModel{
		model.NewExperimentModel("exp1id", "exp1name", `undefined`, false, []string{"someFlag"}, nil),
		model.NewExperimentModel("exp2id", "exp2name", `ifThen(eq("true", flagValue("someFlag")), "blue", "green")`, false, []string{"colorVar"}, nil),
	}

	flagSetter := entities.NewFlagSetter(flagRepository, parser, experimentRepository, nil)
	experimentRepository.SetExperiments(experimentModels)
	flagSetter.SetExperiments()

	experimentsExtensions := extensions.NewExperimentsExtensions(parser, targetGroupsRepository, flagRepository, experimentRepository)
	experimentsExtensions.Extend()

	colorVar := entities.NewRoxString("red", []string{"red", "green", "blue"})
	colorVar.(model.InternalVariant).SetForEvaluation(parser, nil, nil)
	flagRepository.AddFlag(colorVar, "colorVar", "colorVar")

	assert.Equal(t, "green", colorVar.GetValueAsString(nil))
}

func TestExperimentsExtensionsFlagDependencyWithContext(t *testing.T) {
	parser := roxx.NewParser()
	targetGroupsRepository := repositories.NewTargetGroupRepository()
	experimentRepository := repositories.NewExperimentRepository()
	flagRepository := repositories.NewFlagRepository()
	propertiesRepository := repositories.NewCustomPropertyRepository()

	extensions.NewPropertiesExtensions(parser, propertiesRepository, nil).Extend()
	extensions.NewExperimentsExtensions(parser, targetGroupsRepository, flagRepository, experimentRepository).Extend()

	propertiesRepository.AddCustomProperty(properties.NewComputedBooleanProperty("prop", func(context context.Context) bool {
		return context.Get("isPropOn").(bool)
	}))

	flag1 := entities.NewFlag(false)
	exp1 := model.NewExperimentModel("id1", "name1", `property("prop")`, false, nil, nil)
	flag1.(model.InternalVariant).SetForEvaluation(parser, exp1, nil)
	flagRepository.AddFlag(flag1, "flag1", "flag1")

	flag2 := entities.NewFlag(false)
	exp2 := model.NewExperimentModel("id2", "name2", `flagValue("flag1")`, false, nil, nil)
	flag2.(model.InternalVariant).SetForEvaluation(parser, exp2, nil)
	flagRepository.AddFlag(flag2, "flag2", "flag2")

	flagValue := flag2.GetValueAsString(context.NewContext(map[string]interface{}{"isPropOn": true}))

	assert.Equal(t, "true", flagValue)
}

func TestExperimentsExtensionsFlagDependencyWithContextUsedOnExperimentWithNoFlag(t *testing.T) {
	parser := roxx.NewParser()
	targetGroupsRepository := repositories.NewTargetGroupRepository()
	experimentRepository := repositories.NewExperimentRepository()
	flagRepository := repositories.NewFlagRepository()
	propertiesRepository := repositories.NewCustomPropertyRepository()

	extensions.NewPropertiesExtensions(parser, propertiesRepository, nil).Extend()
	extensions.NewExperimentsExtensions(parser, targetGroupsRepository, flagRepository, experimentRepository).Extend()

	propertiesRepository.AddCustomProperty(properties.NewComputedBooleanProperty("prop", func(context context.Context) bool {
		return context.Get("isPropOn").(bool)
	}))

	flag3 := entities.NewFlag(false)
	exp3 := model.NewExperimentModel("id3", "name3", `flagValue("flag2")`, false, nil, nil)
	flag3.(model.InternalVariant).SetForEvaluation(parser, exp3, nil)
	flagRepository.AddFlag(flag3, "flag3", "flag3")

	experimentModels := []*model.ExperimentModel{
		model.NewExperimentModel("exp1id", "exp1name", `property("prop")`, false, []string{"flag2"}, nil),
	}
	experimentRepository.SetExperiments(experimentModels)

	flagValue := flag3.GetValueAsString(context.NewContext(map[string]interface{}{"isPropOn": true}))

	assert.Equal(t, "true", flagValue)
}

func TestExperimentsExtensionsFlagDependencyWithContext2LevelMidLevelNoFlagEvalExperiment(t *testing.T) {
	parser := roxx.NewParser()
	targetGroupsRepository := repositories.NewTargetGroupRepository()
	experimentRepository := repositories.NewExperimentRepository()
	flagRepository := repositories.NewFlagRepository()
	propertiesRepository := repositories.NewCustomPropertyRepository()

	extensions.NewPropertiesExtensions(parser, propertiesRepository, nil).Extend()
	extensions.NewExperimentsExtensions(parser, targetGroupsRepository, flagRepository, experimentRepository).Extend()

	propertiesRepository.AddCustomProperty(properties.NewComputedBooleanProperty("prop", func(context context.Context) bool {
		return context.Get("isPropOn").(bool)
	}))

	flag1 := entities.NewFlag(false)
	exp1 := model.NewExperimentModel("id1", "name1", `property("prop")`, false, nil, nil)
	flag1.(model.InternalVariant).SetForEvaluation(parser, exp1, nil)
	flagRepository.AddFlag(flag1, "flag1", "flag1")

	flag3 := entities.NewFlag(false)
	exp2 := model.NewExperimentModel("id3", "name3", `flagValue("flag2")`, false, nil, nil)
	flag3.(model.InternalVariant).SetForEvaluation(parser, exp2, nil)
	flagRepository.AddFlag(flag3, "flag3", "flag3")

	experimentModels := []*model.ExperimentModel{
		model.NewExperimentModel("exp1id", "exp1name", `flagValue("flag1")`, false, []string{"flag2"}, nil),
	}
	experimentRepository.SetExperiments(experimentModels)

	flagValue := flag3.GetValueAsString(context.NewContext(map[string]interface{}{"isPropOn": true}))

	assert.Equal(t, "true", flagValue)
}
