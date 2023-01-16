package register_test

import (
	"github.com/rollout/rox-go/v5/core/entities"
	"github.com/rollout/rox-go/v5/core/model"
	"github.com/rollout/rox-go/v5/core/register"
	"github.com/rollout/rox-go/v5/core/repositories"
	"github.com/stretchr/testify/assert"
	"testing"
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
