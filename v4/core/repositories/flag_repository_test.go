package repositories_test

import (
	"github.com/rollout/rox-go/v4/core/entities"
	"github.com/rollout/rox-go/v4/core/model"
	"github.com/rollout/rox-go/v4/core/repositories"
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
