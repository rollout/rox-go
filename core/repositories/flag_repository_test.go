package repositories_test

import (
	"github.com/rollout/rox-go/core/entities"
	"github.com/rollout/rox-go/core/model"
	"github.com/rollout/rox-go/core/repositories"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFlagRepositoryWillReturnNullWhenFlagNotFound(t *testing.T) {
	repo := repositories.NewFlagRepository()

	assert.Nil(t, repo.GetFlag("harti"))
}

func TestFlagRepositoryWillAddFlagAndSetName(t *testing.T) {
	repo := repositories.NewFlagRepository()
	flag := entities.NewFlag(false)
	repo.AddFlag(flag, "harti")

	assert.Equal(t, "harti", repo.GetFlag("harti").Name())
}

func TestFlagRepositoryWillAddRoxStringAndSetName(t *testing.T) {
	repo := repositories.NewFlagRepository()
	flag := entities.NewRoxString("1", []string{"2", "3"})
	repo.AddFlag(flag, "harti")

	assert.Equal(t, "harti", repo.GetFlag("harti").Name())
}

func TestFlagRepositoryWillAddRoxIntAndSetName(t *testing.T) {
	repo := repositories.NewFlagRepository()
	flag := entities.NewRoxInt(1, []int{2, 3})
	repo.AddFlag(flag, "harti")

	assert.Equal(t, "harti", repo.GetFlag("harti").Name())
}

func TestFlagRepositoryWillRaiseFlagAddedEvent(t *testing.T) {
	repo := repositories.NewFlagRepository()
	flag := entities.NewFlag(false)

	var variantFromEvent model.Variant
	repo.RegisterFlagAddedHandler(func(variant model.Variant) {
		variantFromEvent = variant
	})

	repo.AddFlag(flag, "harti")

	assert.Equal(t, "harti", variantFromEvent.Name())
}
