package register_test

import (
	"testing"

	"github.com/rollout/rox-go/v6/core/entities"
	"github.com/rollout/rox-go/v6/core/logging"
	"github.com/rollout/rox-go/v6/core/model"
	"github.com/rollout/rox-go/v6/core/register"
	"github.com/rollout/rox-go/v6/core/repositories"
	"github.com/stretchr/testify/assert"
)

type Container1 struct {
	Variant1 model.Variant
	Flag1    model.Flag
	Flag2    model.Flag
	Flag3    model.Flag

	SomethingElse interface{}
}

type Container2 struct {
	Variant1 model.Variant `flagName:"variant-1"`
	Flag1    model.Flag    `flagName:"flag.long.name"`
	Flag2    model.Flag    `otherTag:"flag-2"`
	Flag3    model.Flag    `flagName:"flag_3"`

	SomethingElse interface{} `flagName:"something_else"`
}

func TestRegistererWillThrowWhenNSRegisteredTwice(t *testing.T) {
	flagRepo := repositories.NewFlagRepository()
	container := &Container1{
		Variant1:      entities.NewRoxString("1", []string{"1", "2", "3"}),
		Flag1:         entities.NewFlag(false),
		Flag2:         entities.NewFlag(false),
		Flag3:         entities.NewFlag(false),
		SomethingElse: struct{}{},
	}
	registerer := register.NewRegisterer(flagRepo)

	registerer.RegisterInstance(container, "ns1")

	assert.Panics(t, func() {
		registerer.RegisterInstance(container, "ns1")
	})
}

func TestRegisterWillRegisterVariantAndFlag(t *testing.T) {
	flagRepo := repositories.NewFlagRepository()
	container := &Container1{
		Variant1:      entities.NewRoxString("1", []string{"1", "2", "3"}),
		Flag1:         entities.NewFlag(false),
		Flag2:         entities.NewFlag(false),
		Flag3:         entities.NewFlag(false),
		SomethingElse: struct{}{},
	}
	registerer := register.NewRegisterer(flagRepo)
	registerer.RegisterInstance(container, "ns1")

	assert.Equal(t, 4, len(flagRepo.GetAllFlags()))
	assert.Equal(t, "1", flagRepo.GetFlag("ns1.Variant1").GetDefaultAsString())
	assert.Equal(t, "false", flagRepo.GetFlag("ns1.Flag1").GetDefaultAsString())
}

func TestRegistererWillRegisterWithEmptyNS(t *testing.T) {
	flagRepo := repositories.NewFlagRepository()
	container := &Container1{
		Variant1:      entities.NewRoxString("1", []string{"1", "2", "3"}),
		Flag1:         entities.NewFlag(false),
		Flag2:         entities.NewFlag(false),
		Flag3:         entities.NewFlag(false),
		SomethingElse: struct{}{},
	}
	registerer := register.NewRegisterer(flagRepo)
	registerer.RegisterInstance(container, "")

	assert.Equal(t, 4, len(flagRepo.GetAllFlags()))
	assert.Equal(t, "1", flagRepo.GetFlag("Variant1").GetDefaultAsString())
	assert.Equal(t, "false", flagRepo.GetFlag("Flag1").GetDefaultAsString())
}

func TestRegisterWillSetNameFromTagIfAvailable(t *testing.T) {
	flagRepo := repositories.NewFlagRepository()
	container2 := &Container2{
		Variant1:      entities.NewRoxString("1", []string{"1", "2", "3"}),
		Flag1:         entities.NewFlag(false),
		Flag2:         entities.NewFlag(false),
		Flag3:         entities.NewFlag(false),
		SomethingElse: struct{}{},
	}
	registerer := register.NewRegisterer(flagRepo)
	registerer.RegisterInstance(container2, "ns1")

	assert.Equal(t, "ns1.variant-1", flagRepo.GetFlag("ns1.variant-1").Name())
}

func TestRegisterWillUseVariableNameIgnoringOtherTags(t *testing.T) {
	flagRepo := repositories.NewFlagRepository()

	container2 := &Container2{
		Variant1:      entities.NewRoxString("1", []string{"1", "2", "3"}),
		Flag1:         entities.NewFlag(false),
		Flag2:         entities.NewFlag(false),
		Flag3:         entities.NewFlag(false),
		SomethingElse: struct{}{},
	}
	registerer := register.NewRegisterer(flagRepo)
	registerer.RegisterInstance(container2, "ns1")

	flag := flagRepo.GetFlag("ns1.flag-2")
	assert.Equal(t, nil, flag, "flag was wrongly registered with unrecognized tag")

	assert.Equal(t, "ns1.Flag2", flagRepo.GetFlag("ns1.Flag2").Name())
}

type warnCapturingLogger struct {
	warnings []string
}

func (l *warnCapturingLogger) Debug(message string, err interface{}) {}
func (l *warnCapturingLogger) Error(message string, err interface{}) {}
func (l *warnCapturingLogger) Warn(message string, err interface{}) {
	l.warnings = append(l.warnings, message)
}

func captureWarnings(t *testing.T) *warnCapturingLogger {
	logger := &warnCapturingLogger{}
	logging.SetLogger(logger)
	t.Cleanup(func() { logging.SetLogger(nil) })
	return logger
}

type ContainerWithUnexported struct {
	Flag1   model.Flag
	flag2   model.Flag
	string1 model.RoxString `flagName:"string-1"`

	somethingElse interface{}
	counter       int
}

func TestRegisterWillWarnAboutUnexportedFlags(t *testing.T) {
	logger := captureWarnings(t)
	flagRepo := repositories.NewFlagRepository()
	container := &ContainerWithUnexported{
		Flag1:         entities.NewFlag(false),
		flag2:         entities.NewFlag(false),
		string1:       entities.NewRoxString("1", []string{"1", "2"}),
		somethingElse: struct{}{},
	}
	registerer := register.NewRegisterer(flagRepo)
	registerer.RegisterInstance(container, "")

	assert.Equal(t, 1, len(flagRepo.GetAllFlags()))
	assert.NotNil(t, flagRepo.GetFlag("Flag1"))
	assert.Nil(t, flagRepo.GetFlag("flag2"))
	assert.Nil(t, flagRepo.GetFlag("string-1"))

	assert.Equal(t, []string{
		"Flag field \"flag2\" is unexported, so it cannot be registered. Rename it to \"Flag2\" and add the tag `flagName:\"flag2\"` to keep the same flag name.",
		"Flag field \"string1\" is unexported, so it cannot be registered. Rename it to \"String1\" and add the tag `flagName:\"string-1\"` to keep the same flag name.",
	}, logger.warnings)
}

func TestRegisterWillIncludeNamespaceInUnexportedFlagWarning(t *testing.T) {
	logger := captureWarnings(t)
	flagRepo := repositories.NewFlagRepository()
	container := &ContainerWithUnexported{
		Flag1: entities.NewFlag(false),
		flag2: entities.NewFlag(false),
	}
	registerer := register.NewRegisterer(flagRepo)
	registerer.RegisterInstance(container, "ns1")

	assert.Contains(t, logger.warnings, "Flag field \"flag2\" in namespace \"ns1\" is unexported, so it cannot be registered. Rename it to \"Flag2\" and add the tag `flagName:\"flag2\"` to keep the same flag name.")
}

func TestRegisterWillNotWarnWhenAllFlagsAreExported(t *testing.T) {
	logger := captureWarnings(t)
	flagRepo := repositories.NewFlagRepository()
	container := &Container2{
		Variant1:      entities.NewRoxString("1", []string{"1", "2", "3"}),
		Flag1:         entities.NewFlag(false),
		Flag2:         entities.NewFlag(false),
		Flag3:         entities.NewFlag(false),
		SomethingElse: struct{}{},
	}
	registerer := register.NewRegisterer(flagRepo)
	registerer.RegisterInstance(container, "")

	assert.Equal(t, 4, len(flagRepo.GetAllFlags()))
	assert.Empty(t, logger.warnings)
}

func TestRegisterWillWarnAboutUnexportedFlagsWhenContainerPassedByValue(t *testing.T) {
	logger := captureWarnings(t)
	flagRepo := repositories.NewFlagRepository()
	container := ContainerWithUnexported{
		Flag1: entities.NewFlag(false),
		flag2: entities.NewFlag(false),
	}
	registerer := register.NewRegisterer(flagRepo)
	registerer.RegisterInstance(container, "")

	assert.Equal(t, 1, len(flagRepo.GetAllFlags()))
	assert.NotNil(t, flagRepo.GetFlag("Flag1"))
	assert.Equal(t, 2, len(logger.warnings))
}

func TestRegisterWillLogUnexportedFlagWarningsOnFlushWhenNoLoggerConfigured(t *testing.T) {
	registerer := register.NewRegisterer(repositories.NewFlagRepository())
	registerer.RegisterInstance(&ContainerWithUnexported{flag2: entities.NewFlag(false)}, "")

	logger := captureWarnings(t)
	assert.Empty(t, logger.warnings)

	registerer.FlushWarnings()
	assert.Equal(t, 2, len(logger.warnings))

	registerer.FlushWarnings()
	assert.Equal(t, 2, len(logger.warnings))
}
